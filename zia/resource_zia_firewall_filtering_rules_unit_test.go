package zia

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func newFakeFirewallRules(f *unitFixture, n int) *fakeRuleList {
	return newFakeRuleList(f, n, common.ZIAPath("firewallFilteringRules"), common.ZIAPath("firewallFilteringRules", "*"), "", nil)
}

func firewallRuleConfig(name string, order int) map[string]any {
	return map[string]any{"name": name, "order": order, "rank": 7, "action": "ALLOW", "state": "DISABLED"}
}

// SUP-4231: 150 rules created in one apply, 10 at a time (Terraform's default),
// with declared orders in random creation order, below 10 existing rules. Each
// new rule is created at the end of the list and the rules are placed relative
// to each other until the list is complete, so the run needs about one move per
// rule instead of thousands (4.8.13 moved every placed rule again after each
// batch: ~3,500 moves for 125 rules).
func TestFirewallFilteringRules_BatchCreatePlacesEachRuleAboutOnce(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeFirewallRules(f, 10)
	firewallFilteringStartingOrder = 0
	t.Cleanup(func() { firewallFilteringStartingOrder = 0 })
	r := resourceFirewallFilteringRules()

	const n, batch = 150, 10
	orders := rand.New(rand.NewSource(1)).Perm(n)
	ids := make([]int, n)
	for b := 0; b < n; b += batch {
		var wg sync.WaitGroup
		for i := b; i < b+batch; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				d := f.resourceData(r, "", firewallRuleConfig("New "+strconv.Itoa(i), 11+orders[i]))
				if e := f.run(resourceFirewallFilteringRulesCreate, d); e != "" {
					t.Errorf("create %d: %s", i, e)
				}
				ids[i], _ = strconv.Atoi(d.Id())
			}()
		}
		wg.Wait()
	}

	for i, id := range ids {
		if got := fake.orderOf(id); got != 11+orders[i] {
			t.Errorf("rule %d: order %d, want %d", id, got, 11+orders[i])
		}
	}
	if moves := len(fake.puts); moves > 2*n {
		t.Errorf("%d moves for %d rules, want at most %d", moves, n, 2*n)
	}
	t.Logf("%d moves for %d rules", len(fake.puts), n)
}

// The reorder outlives the request of the rule that started it: cancelling
// that rule's context (Terraform's create timeout) must not stop it.
func TestFirewallFilteringRules_ReorderSurvivesCancelledCreate(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeFirewallRules(f, 5)
	firewallFilteringStartingOrder = 0
	t.Cleanup(func() { firewallFilteringStartingOrder = 0 })
	r := resourceFirewallFilteringRules()

	ctx, cancel := context.WithCancel(context.Background())
	d := f.resourceData(r, "", firewallRuleConfig("New", 2))
	// The create reads the list once for its starting order; the second read
	// is the reorder's. Cancel the create's context from then on, as
	// Terraform's create timeout would while the rules are being placed.
	var reads int
	var mu sync.Mutex
	fake.onGetAll = func() {
		mu.Lock()
		defer mu.Unlock()
		if reads++; reads == 2 {
			cancel()
		}
	}
	_ = r.CreateContext(ctx, d, f.client) // the read after the reorder fails on the cancelled context
	if got := fake.orderOf(1006); got != 2 {
		t.Errorf("order %d, want 2: the reorder stopped with the cancelled request", got)
	}
}

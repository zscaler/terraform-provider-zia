package zia

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func newFakeFirewallRules(f *unitFixture, n int) *fakeRuleList {
	return newFakeRuleList(f, n, common.ZIAPath("firewallFilteringRules"), common.ZIAPath("firewallFilteringRules", "*"), "", nil)
}

func firewallRuleConfig(name string, order int) map[string]any {
	return map[string]any{"name": name, "order": order, "rank": 7, "action": "ALLOW", "state": "DISABLED"}
}

// SUP-4231: create 150 rules in batches of ten, in a different sequence
// from their declared orders. Verify the final API list and refreshed plan;
// request counts alone cannot establish correctness.
func TestFirewallFilteringRules_BatchCreatePreservesDeclaredOrders(t *testing.T) {
	t.Run("after_existing_custom_rules", func(t *testing.T) {
		testFirewallBatchCreate(t, 10, 11, false)
	})
	t.Run("unmanaged_predefined_rules_move_automatically", func(t *testing.T) {
		testFirewallBatchCreate(t, 3, 1, true)
	})
	t.Run("no_predefined_rules", func(t *testing.T) {
		testFirewallBatchCreate(t, 0, 1, false)
	})
}

func testFirewallBatchCreate(t *testing.T, existing, firstOrder int, predefined bool) {
	t.Helper()
	f := newUnitFixture(t)
	var fields map[string]any
	if predefined {
		fields = map[string]any{"predefined": true, "accessControl": "READ_WRITE"}
	}
	fake := newFakeRuleList(f, existing, common.ZIAPath("firewallFilteringRules"), common.ZIAPath("firewallFilteringRules", "*"), "", fields)
	firewallFilteringStartingOrder = 0
	t.Cleanup(func() { firewallFilteringStartingOrder = 0 })
	r := resourceFirewallFilteringRules()

	const n, batch = 150, 10
	orders := rand.New(rand.NewSource(1)).Perm(n)
	ids := make([]int, n)
	states := make([]*schema.ResourceData, n)
	for b := 0; b < n; b += batch {
		var wg sync.WaitGroup
		for i := b; i < b+batch; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				d := f.resourceData(r, "", firewallRuleConfig("New "+strconv.Itoa(i), firstOrder+orders[i]))
				if e := f.run(resourceFirewallFilteringRulesCreate, d); e != "" {
					t.Errorf("create %d: %s", i, e)
				}
				ids[i], _ = strconv.Atoi(d.Id())
				states[i] = d
			}()
		}
		wg.Wait()
	}

	for i, id := range ids {
		if got := fake.orderOf(id); got != firstOrder+orders[i] {
			t.Errorf("rule %d: order %d, want %d", id, got, firstOrder+orders[i])
		}
		// The state Create returned is what Terraform persists. Batches are
		// released before their rules are placed, and nothing rewrites state
		// once a later cycle moves the rule, so Create must record the declared
		// order rather than the position the rule held when it returned.
		if got := states[i].Get("order").(int); got != firstOrder+orders[i] {
			t.Errorf("rule %d: state after create has order %d, want declared %d", id, got, firstOrder+orders[i])
		}
		// A subsequent plan refreshes state from the final upstream list.
		// Temporary positions from an earlier batch must not leave order drift.
		if e := f.run(resourceFirewallFilteringRulesRead, states[i]); e != "" {
			t.Fatal(e)
		}
		config := firewallRuleConfig("New "+strconv.Itoa(i), firstOrder+orders[i])
		diff, err := r.Diff(context.Background(), states[i].State(), terraform.NewResourceConfigRaw(config), f.client)
		if err != nil {
			t.Fatal(err)
		}
		if diff != nil && !diff.Empty() {
			t.Errorf("rule %d: follow-up plan has drift: %+v", id, diff.Attributes)
		}
	}
	if moves := len(fake.puts); moves > 3*n {
		t.Errorf("%d order PUTs for %d rules; repeated batch placement has regressed", moves, n)
	}
	for _, put := range fake.puts {
		if id := asInt(put["id"]); id <= 1000+existing {
			t.Errorf("sent an order PUT for unmanaged rule %d", id)
		}
		// No temporary or normalized targets: each reorder PUT must carry
		// exactly the order declared for this managed rule.
		name, _ := put["name"].(string)
		i, err := strconv.Atoi(strings.TrimPrefix(name, "New "))
		if err != nil || i < 0 || i >= len(orders) {
			t.Fatalf("unexpected reordered rule %q", name)
		}
		if got, want := asInt(put["order"]), firstOrder+orders[i]; got != want {
			t.Errorf("PUT for %s used order %d, want declared order %d", name, got, want)
		}
	}
	if fake.rejected != 0 {
		t.Errorf("sent %d requests at unavailable positions", fake.rejected)
	}
	if predefined {
		for i := 1; i <= existing; i++ {
			if got := fake.orderOf(1000 + i); got != n+i {
				t.Errorf("predefined rule %d: API position %d, want %d after managed rules", 1000+i, got, n+i)
			}
		}
	}
	t.Logf("%d moves for %d rules", len(fake.puts), n)
}

// Importing a predefined rule and declaring it in HCL makes its ID eligible
// for the same reorder path as a custom rule. Other predefined IDs stay out.
func TestFirewallFilteringRules_ImportedPredefinedRuleIsManaged(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeRuleList(f, 3, common.ZIAPath("firewallFilteringRules"), common.ZIAPath("firewallFilteringRules", "*"), "", map[string]any{
		"predefined": true, "accessControl": "READ_WRITE", "action": "ALLOW", "state": "DISABLED",
	})
	r := resourceFirewallFilteringRules()
	state := firewallRuleConfig(ruleName(1), 1)
	state["rule_id"] = 1001
	d := f.updateData(r, "1001", state, firewallRuleConfig(ruleName(1), 3))
	if e := f.run(resourceFirewallFilteringRulesUpdate, d); e != "" {
		t.Fatal(e)
	}
	if got := fake.orderOf(1001); got != 3 {
		t.Errorf("imported predefined rule at order %d, want 3", got)
	}
	if fake.putCount(1001) == 0 {
		t.Fatal("imported predefined rule was excluded from reordering")
	}
	for _, put := range fake.puts {
		if id := asInt(put["id"]); id != 1001 {
			t.Errorf("sent an order PUT for unmanaged predefined rule %d", id)
		}
	}
}

// Position selection and POST must be one operation: advancing the cached
// position after releasing the lock lets concurrent creates reuse it and push
// rules that were already placed down the list.
func TestFirewallFilteringRules_ConcurrentCreatesAppend(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeFirewallRules(f, 10)
	fake.createDelay = 5 * time.Millisecond
	firewallFilteringStartingOrder = 0
	t.Cleanup(func() { firewallFilteringStartingOrder = 0 })
	r := resourceFirewallFilteringRules()

	const n = 10
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			d := f.resourceData(r, "", firewallRuleConfig("Concurrent "+strconv.Itoa(i), 11+i))
			if e := f.run(resourceFirewallFilteringRulesCreate, d); e != "" {
				t.Errorf("create %d: %s", i, e)
			}
		}(i)
	}
	close(start)
	wg.Wait()

	if len(fake.createOrders) != n {
		t.Fatalf("accepted %d creates, want %d", len(fake.createOrders), n)
	}
	for i, order := range fake.createOrders {
		if want := 11 + i; order != want {
			t.Errorf("POST %d used order %d, want append position %d", i, order, want)
		}
	}
}

func TestFirewallFilteringRules_CreateListErrorStopsBeforePost(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeFirewallRules(f, 10)
	firewallFilteringStartingOrder = 0
	t.Cleanup(func() { firewallFilteringStartingOrder = 0 })
	f.api.OnFunc(http.MethodGet, common.ZIAPath("firewallFilteringRules"), func(*http.Request, []byte) common.MockResponse {
		return common.MockResponse{StatusCode: http.StatusForbidden, Body: map[string]any{
			"code": "ACCESS_DENIED", "message": "Insufficient permissions"}}
	})
	r := resourceFirewallFilteringRules()
	d := f.resourceData(r, "", firewallRuleConfig("New", 11))
	if e := f.run(resourceFirewallFilteringRulesCreate, d); e == "" {
		t.Fatal("create succeeded after the initial rule list could not be read")
	}
	if len(fake.createOrders) != 0 {
		t.Errorf("sent %d POSTs without knowing the append position", len(fake.createOrders))
	}
}

func TestFirewallFilteringRules_CreateRefreshesAppendPositionAfterDelete(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeFirewallRules(f, 5)
	firewallFilteringStartingOrder = 0
	t.Cleanup(func() { firewallFilteringStartingOrder = 0 })
	var deleted sync.Once
	fake.onGetAll = func() {
		// Return the original five-rule snapshot, but shrink the policy
		// before the POST uses its append position, as an external delete can.
		deleted.Do(func() {
			fake.deleteRule(1005)
			fake.deleteRule(1004)
		})
	}
	r := resourceFirewallFilteringRules()
	d := f.resourceData(r, "", firewallRuleConfig("New", 2))
	if e := f.run(resourceFirewallFilteringRulesCreate, d); e != "" {
		t.Fatal(e)
	}
	if fake.rejected != 1 {
		t.Errorf("rejected %d stale positions, want 1", fake.rejected)
	}
	if len(fake.createOrders) != 1 || fake.createOrders[0] != 4 {
		t.Errorf("accepted POST positions %v, want [4]", fake.createOrders)
	}
	if got := fake.orderOf(1006); got != 2 {
		t.Errorf("final order %d, want 2", got)
	}
}

// Move 150 existing managed rules from API positions 4..153 into declared
// positions 1..150. The three unmanaged rules must shift automatically, without
// reserving a prefix or receiving PUTs. Check the final API and refreshed plan.
func TestFirewallFilteringRules_BatchUpdatesConvergeWithoutPlanDrift(t *testing.T) {
	f := newUnitFixture(t)
	const existing, n, batch = 3, 150, 10
	fake := newFakeRuleList(f, existing+n, common.ZIAPath("firewallFilteringRules"), common.ZIAPath("firewallFilteringRules", "*"), "", map[string]any{"action": "ALLOW", "state": "DISABLED"})
	r := resourceFirewallFilteringRules()
	orders := rand.New(rand.NewSource(2)).Perm(n)
	states := make([]*schema.ResourceData, n)
	for b := 0; b < n; b += batch {
		var wg sync.WaitGroup
		for i := b; i < b+batch; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				id := 1000 + existing + 1 + i
				state := firewallRuleConfig(ruleName(existing+1+i), existing+1+i)
				state["rule_id"] = id
				config := firewallRuleConfig(ruleName(existing+1+i), 1+orders[i])
				d := f.updateData(r, strconv.Itoa(id), state, config)
				if e := f.run(resourceFirewallFilteringRulesUpdate, d); e != "" {
					t.Errorf("update %d: %s", id, e)
				}
				states[i] = d
			}(i)
		}
		wg.Wait()
	}
	for i, d := range states {
		id := 1000 + existing + 1 + i
		want := 1 + orders[i]
		if got := fake.orderOf(id); got != want {
			t.Errorf("rule %d: upstream order %d, want %d", id, got, want)
		}
		// The state Update returned must carry the declared order even when
		// its batch was released before the rule was moved.
		if got := d.Get("order").(int); got != want {
			t.Errorf("rule %d: state after update has order %d, want declared %d", id, got, want)
		}
		if e := f.run(resourceFirewallFilteringRulesRead, d); e != "" {
			t.Fatal(e)
		}
		config := firewallRuleConfig(ruleName(existing+1+i), want)
		diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), f.client)
		if err != nil {
			t.Fatal(err)
		}
		if diff != nil && !diff.Empty() {
			t.Errorf("rule %d: follow-up plan has drift: %+v", id, diff.Attributes)
		}
	}
	for _, put := range fake.puts {
		if id := asInt(put["id"]); id <= 1000+existing {
			t.Errorf("sent an order PUT for unmanaged rule %d", id)
		}
	}
	for i := 1; i <= existing; i++ {
		if got := fake.orderOf(1000 + i); got != n+i {
			t.Errorf("unmanaged rule %d: order %d, want automatic API position %d", 1000+i, got, n+i)
		}
	}
	t.Logf("%d updates reached their declared positions; refreshed plan has no changes", n)
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

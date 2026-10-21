package zia

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func newFakeURLRules(f *unitFixture, n int) *fakeRuleList {
	return newFakeRuleList(f, n, common.ZIAPath("urlFilteringRules"), common.ZIAPath("urlFilteringRules", "*"), "", nil)
}

func urlRuleConfig(name string, order, rank int) map[string]any {
	return map[string]any{
		"name": name, "description": name, "order": order, "rank": rank,
		"action": "ALLOW", "state": "ENABLED",
		"protocols": []any{"ANY_RULE"}, "url_categories": []any{"ANY"},
	}
}

func urlRuleState(id int, name string, order, rank int) map[string]any {
	s := urlRuleConfig(name, order, rank)
	s["rule_id"] = id
	return s
}

// The customer's failure: a rule in the middle is deleted and every rule below
// it is renumbered, all in one apply. The delete is slow and the updates read
// the list while it is in progress, as in the customer's log. The updates only
// change the order, so none of them writes the rule; the reorder engine finds
// them in place once the delete has closed the gap.
func TestURLFilteringRules_DeleteWithRenumberingInOneApply(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeURLRules(f, 20)
	fake.deleteDelay = 200 * time.Millisecond
	r := resourceURLFilteringRules()

	var wg sync.WaitGroup
	errs := make(chan string, 13)
	wg.Add(1)
	go func() {
		defer wg.Done()
		d := f.resourceData(r, "1008", urlRuleState(1008, ruleName(8), 8, 7))
		errs <- f.run(resourceURLFilteringRulesDelete, d)
	}()
	for i := 9; i <= 20; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := 1000 + i
			d := f.updateData(r, strconv.Itoa(id), urlRuleState(id, ruleName(i), i, 7), urlRuleConfig(ruleName(i), i-1, 7))
			errs <- f.run(resourceURLFilteringRulesUpdate, d)
		}()
	}
	wg.Wait()
	close(errs)

	for e := range errs {
		if e != "" {
			t.Errorf("unexpected error: %s", e)
		}
	}
	if fake.rejected != 0 {
		t.Errorf("%d writes were rejected for an invalid order", fake.rejected)
	}
	for i := 9; i <= 20; i++ {
		if got := fake.orderOf(1000 + i); got != i-1 {
			t.Errorf("%s: order %d, want %d", ruleName(i), got, i-1)
		}
	}
}

// After the delete has landed the rule is already where the configuration puts
// it: an order-only update sends no write at all.
func TestURLFilteringRules_OrderOnlyUpdateSendsNoWrite(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeURLRules(f, 20)
	fake.deleteRule(1008)
	r := resourceURLFilteringRules()

	d := f.updateData(r, "1009", urlRuleState(1009, ruleName(9), 9, 7), urlRuleConfig(ruleName(9), 8, 7))
	if e := f.run(resourceURLFilteringRulesUpdate, d); e != "" {
		t.Fatalf("unexpected error: %s", e)
	}
	if n := fake.putCount(1009); n != 0 {
		t.Errorf("sent %d PUTs for an order-only update of a rule already in place, want 0", n)
	}
}

// An update of other attributes writes the rule at its current position. When a
// delete lands between reading that position and writing (here the last rule,
// order 20, becomes order 19), the rejected write is retried at the position
// read again.
func TestURLFilteringRules_InPlaceUpdateRetriesAfterConcurrentDelete(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeURLRules(f, 20)
	var once sync.Once
	fake.onGetRule = func(id int) {
		if id == 1020 {
			once.Do(func() { fake.deleteRule(1008) })
		}
	}
	r := resourceURLFilteringRules()

	config := urlRuleConfig(ruleName(20), 19, 7)
	config["description"] = "updated"
	d := f.updateData(r, "1020", urlRuleState(1020, ruleName(20), 20, 7), config)
	if e := f.run(resourceURLFilteringRulesUpdate, d); e != "" {
		t.Fatalf("unexpected error: %s", e)
	}
	if fake.rejected != 1 {
		t.Errorf("rejected writes = %d, want 1 (the write at the position read before the delete)", fake.rejected)
	}
	if got := fake.orderOf(1020); got != 19 {
		t.Errorf("order %d, want 19", got)
	}
	if got := fake.puts[0]; got["description"] != "updated" || asInt(got["order"]) != 19 {
		t.Errorf("first accepted PUT: description %v order %v, want \"updated\" at order 19", got["description"], got["order"])
	}
}

// A rank change still writes the rule at the end of the list with rank 7, as the
// API requires, before the reorder engine places it with its new rank.
func TestURLFilteringRules_RankChangeWritesAtEndOfList(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeURLRules(f, 5)
	r := resourceURLFilteringRules()

	d := f.updateData(r, "1001", urlRuleState(1001, ruleName(1), 1, 7), urlRuleConfig(ruleName(1), 1, 6))
	if e := f.run(resourceURLFilteringRulesUpdate, d); e != "" {
		t.Fatalf("unexpected error: %s", e)
	}
	if len(fake.puts) == 0 || asInt(fake.puts[0]["order"]) != 5 || asInt(fake.puts[0]["rank"]) != 7 {
		t.Fatalf("first PUT = %v, want order 5 rank 7", fake.puts)
	}
	if got := fake.orderOf(1001); got != 1 {
		t.Errorf("order %d, want 1", got)
	}
}

// Create caches its starting order for the apply. When rules were deleted since
// it was read, the rejected create is retried at the end of the list read again.
func TestURLFilteringRules_CreateRetriesStaleStartingOrder(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeURLRules(f, 18)
	urlFilteringStartingOrder = 20 // read before two rules were deleted
	t.Cleanup(func() { urlFilteringStartingOrder = 0 })
	r := resourceURLFilteringRules()

	d := f.resourceData(r, "", urlRuleConfig("New Rule", 19, 7))
	if e := f.run(resourceURLFilteringRulesCreate, d); e != "" {
		t.Fatalf("unexpected error: %s", e)
	}
	if fake.rejected != 1 {
		t.Errorf("rejected writes = %d, want 1", fake.rejected)
	}
	id, _ := strconv.Atoi(d.Id())
	if got := fake.orderOf(id); got != 19 {
		t.Errorf("order %d, want 19", got)
	}
}

func TestIsStaleOrderError(t *testing.T) {
	for msg, want := range map[string]bool{
		`{"code":"INVALID_INPUT_ARGUMENT","message":"Rule is not allowed at order 20"}`:             true,
		`{"code":"INVALID_INPUT_ARGUMENT","message":"Rule with rank 7 is not allowed at order 14"}`: true,
		`{"code":"INVALID_INPUT_ARGUMENT","message":"Invalid url category"}`:                        false,
		`{"code":"DUPLICATE_ITEM","message":"Rule is not allowed at order 3"}`:                      false,
	} {
		if got := isStaleOrderError(fmt.Errorf("%s", msg)); got != want {
			t.Errorf("%s: got %v, want %v", strings.TrimSpace(msg), got, want)
		}
	}
}

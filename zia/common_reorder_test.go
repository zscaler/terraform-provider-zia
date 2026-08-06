package zia

import (
	"sort"
	"sync"
	"testing"
	"time"
)

// resetReorderState clears the global reorder state between tests.
func resetReorderState() {
	rules.Lock()
	defer rules.Unlock()
	rules.orders = make(map[string]map[int]orderWithState)
	rules.orderer = nil
	rules.reorderDone = nil
}

// =====================================================
// fakeAPI: a minimal in-memory stand-in for the ZIA
// rule API used by reorderAll's getCurrent/updateOrder.
// =====================================================
//
// Tests preSeed each rule with its "current" (API-side)
// state. updateOrder mutates that state so subsequent
// passes see the rule at its new position. This mirrors
// production: a rule that's already at its desired Order
// produces ZERO PUTs.
type fakeAPI struct {
	mu       sync.Mutex
	state    map[int]OrderRule
	putsByID map[int]int
}

func newFakeAPI() *fakeAPI {
	return &fakeAPI{
		state:    map[int]OrderRule{},
		putsByID: map[int]int{},
	}
}

// preSeed inserts/overwrites rule id with the supplied current state.
// Use cur != desired to simulate a rule that needs to be reordered;
// use cur == desired to simulate a rule already at target (no PUT
// expected).
func (f *fakeAPI) preSeed(id, curOrder, curRank int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state[id] = OrderRule{Order: curOrder, Rank: curRank}
}

func (f *fakeAPI) getCurrent() (map[int]OrderRule, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[int]OrderRule, len(f.state))
	for k, v := range f.state {
		out[k] = v
	}
	return out, nil
}

func (f *fakeAPI) updateOrder(id int, ord OrderRule) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state[id] = ord
	f.putsByID[id]++
	return nil
}

func (f *fakeAPI) putsTotal() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, v := range f.putsByID {
		n += v
	}
	return n
}

func (f *fakeAPI) putsFor(id int) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.putsByID[id]
}

// =====================================================
// Sort Logic Tests
// =====================================================

func TestSortOrders_SingleRank(t *testing.T) {
	input := map[int]orderWithState{
		100: {order: OrderRule{Order: 3, Rank: 7}},
		200: {order: OrderRule{Order: 1, Rank: 7}},
		300: {order: OrderRule{Order: 5, Rank: 7}},
		400: {order: OrderRule{Order: 2, Rank: 7}},
		500: {order: OrderRule{Order: 4, Rank: 7}},
	}

	sorted := sortOrders(input)

	if len(sorted) != 5 {
		t.Fatalf("expected 5 sorted rules, got %d", len(sorted))
	}
	expectedIDs := []int{200, 400, 100, 500, 300}
	for i, expected := range expectedIDs {
		if sorted[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, sorted[i].ID)
		}
	}
}

func TestSortOrders_MixedRanks(t *testing.T) {
	input := map[int]orderWithState{
		100: {order: OrderRule{Order: 1, Rank: 7}},
		200: {order: OrderRule{Order: 1, Rank: 1}},
		300: {order: OrderRule{Order: 2, Rank: 7}},
		400: {order: OrderRule{Order: 2, Rank: 1}},
	}

	sorted := sortOrders(input)

	if len(sorted) != 4 {
		t.Fatalf("expected 4 sorted rules, got %d", len(sorted))
	}
	expectedIDs := []int{200, 400, 100, 300}
	for i, expected := range expectedIDs {
		if sorted[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, sorted[i].ID)
		}
	}
}

func TestRuleIDOrderPairList_Sort(t *testing.T) {
	pairs := RuleIDOrderPairList{
		{ID: 1, Order: OrderRule{Order: 5, Rank: 7}},
		{ID: 2, Order: OrderRule{Order: 2, Rank: 7}},
		{ID: 3, Order: OrderRule{Order: 1, Rank: 7}},
		{ID: 4, Order: OrderRule{Order: 3, Rank: 7}},
	}

	sort.Sort(pairs)

	expectedIDs := []int{3, 2, 4, 1}
	for i, expected := range expectedIDs {
		if pairs[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, pairs[i].ID)
		}
	}
}

func TestRuleIDOrderPairList_SameOrder_TiebreakByID(t *testing.T) {
	pairs := RuleIDOrderPairList{
		{ID: 300, Order: OrderRule{Order: 1, Rank: 7}},
		{ID: 100, Order: OrderRule{Order: 1, Rank: 7}},
		{ID: 200, Order: OrderRule{Order: 1, Rank: 7}},
	}

	sort.Sort(pairs)

	expectedIDs := []int{100, 200, 300}
	for i, expected := range expectedIDs {
		if pairs[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, pairs[i].ID)
		}
	}
}

// =====================================================
// Registration & Done Logic Tests
// =====================================================

func TestMarkOrderRuleAsDone(t *testing.T) {
	resetReorderState()

	rules.Lock()
	rules.orders["test_type"] = map[int]orderWithState{
		100: {order: OrderRule{Order: 1, Rank: 7}, done: false},
	}
	rules.Unlock()

	markOrderRuleAsDone(100, "test_type")

	rules.Lock()
	defer rules.Unlock()
	if !rules.orders["test_type"][100].done {
		t.Error("expected rule 100 to be marked done")
	}
}

func TestReorderWithBeforeReorder_FirstRule_RegistersWithDoneFalse(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	api.preSeed(100, 999, 7)

	reorderWithBeforeReorder(
		OrderRule{Order: 1, Rank: 7}, 100, "test_reg",
		api.getCurrent, api.updateOrder, nil,
	)

	rules.Lock()
	state, exists := rules.orders["test_reg"][100]
	ch := rules.reorderDone["test_reg"]
	rules.Unlock()

	if !exists {
		t.Fatal("expected rule 100 to be registered")
	}
	if state.done {
		t.Error("expected rule 100 to be registered with done=false")
	}
	if ch == nil {
		t.Error("expected reorderDone channel to be created")
	}

	markOrderRuleAsDone(100, "test_reg")
	waitForReorder("test_reg")
}

// =====================================================
// Full Lifecycle Tests
// =====================================================

func TestReorder_AllRulesRegisteredBeforeTick(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	// Each rule already exists in API at a "wrong" current position
	// (Order=999) so reorderAll will need to PUT it to its desired
	// 1..5 position.
	for i := 1; i <= 5; i++ {
		api.preSeed(100+i, 999, 7)
	}

	for i := 1; i <= 5; i++ {
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 100+i, "test_all_before",
			api.getCurrent, api.updateOrder, nil,
		)
	}
	for i := 1; i <= 5; i++ {
		markOrderRuleAsDone(100+i, "test_all_before")
	}

	waitForReorder("test_all_before")

	for i := 1; i <= 5; i++ {
		got := api.putsFor(100 + i)
		if got != 1 {
			t.Errorf("rule %d: expected exactly 1 PUT (diff-based), got %d", 100+i, got)
		}
	}
}

func TestReorder_LateArrivingRules_NewCycleStarted(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	api.preSeed(101, 999, 7)
	api.preSeed(102, 999, 7)

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 101, "test_late", api.getCurrent, api.updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 102, "test_late", api.getCurrent, api.updateOrder, nil)
	markOrderRuleAsDone(101, "test_late")
	markOrderRuleAsDone(102, "test_late")
	waitForReorder("test_late")

	if api.putsFor(101) != 1 || api.putsFor(102) != 1 {
		t.Fatalf("first cycle: each rule should PUT exactly once; got 101=%d 102=%d", api.putsFor(101), api.putsFor(102))
	}

	// Late arrivals after first cycle finished.
	api.preSeed(103, 999, 7)
	api.preSeed(104, 999, 7)
	api.preSeed(105, 999, 7)

	reorderWithBeforeReorder(OrderRule{Order: 3, Rank: 7}, 103, "test_late", api.getCurrent, api.updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 4, Rank: 7}, 104, "test_late", api.getCurrent, api.updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 5, Rank: 7}, 105, "test_late", api.getCurrent, api.updateOrder, nil)
	markOrderRuleAsDone(103, "test_late")
	markOrderRuleAsDone(104, "test_late")
	markOrderRuleAsDone(105, "test_late")
	waitForReorder("test_late")

	for _, id := range []int{103, 104, 105} {
		if api.putsFor(id) != 1 {
			t.Errorf("late rule %d: expected 1 PUT, got %d", id, api.putsFor(id))
		}
	}
	// Rules from the first cycle must NOT have been re-PUT — they
	// are already at their target order.
	if api.putsFor(101) != 1 {
		t.Errorf("rule 101 should still have exactly 1 PUT, got %d (we are PUTing rules already at target!)", api.putsFor(101))
	}
	if api.putsFor(102) != 1 {
		t.Errorf("rule 102 should still have exactly 1 PUT, got %d (we are PUTing rules already at target!)", api.putsFor(102))
	}
}

func TestReorder_MultipleResourceTypes_Independent(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	dnsAPI := newFakeAPI()
	sslAPI := newFakeAPI()
	dnsAPI.preSeed(201, 999, 7)
	dnsAPI.preSeed(202, 999, 7)
	sslAPI.preSeed(301, 999, 7)
	sslAPI.preSeed(302, 999, 7)

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 201, "dns_test", dnsAPI.getCurrent, dnsAPI.updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 202, "dns_test", dnsAPI.getCurrent, dnsAPI.updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 301, "ssl_test", sslAPI.getCurrent, sslAPI.updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 302, "ssl_test", sslAPI.getCurrent, sslAPI.updateOrder, nil)

	markOrderRuleAsDone(201, "dns_test")
	markOrderRuleAsDone(202, "dns_test")
	markOrderRuleAsDone(301, "ssl_test")
	markOrderRuleAsDone(302, "ssl_test")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { waitForReorder("dns_test"); wg.Done() }()
	go func() { waitForReorder("ssl_test"); wg.Done() }()
	wg.Wait()

	if dnsAPI.putsFor(201) != 1 || dnsAPI.putsFor(202) != 1 {
		t.Errorf("DNS rules: expected 1 PUT each, got 201=%d 202=%d", dnsAPI.putsFor(201), dnsAPI.putsFor(202))
	}
	if sslAPI.putsFor(301) != 1 || sslAPI.putsFor(302) != 1 {
		t.Errorf("SSL rules: expected 1 PUT each, got 301=%d 302=%d", sslAPI.putsFor(301), sslAPI.putsFor(302))
	}
}

func TestReorder_ReadMustHappenAfterWaitForReorder(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	api.preSeed(100, 999, 7)

	var mu sync.Mutex
	reorderCompleted := false
	readCalledAfterReorder := false

	wrappedUpdate := func(id int, order OrderRule) error {
		err := api.updateOrder(id, order)
		mu.Lock()
		reorderCompleted = true
		mu.Unlock()
		return err
	}

	simulateRead := func() {
		mu.Lock()
		defer mu.Unlock()
		if reorderCompleted {
			readCalledAfterReorder = true
		}
	}

	reorderWithBeforeReorder(
		OrderRule{Order: 1, Rank: 7}, 100, "test_read_order",
		api.getCurrent, wrappedUpdate, nil,
	)

	markOrderRuleAsDone(100, "test_read_order")
	waitForReorder("test_read_order")
	simulateRead()

	mu.Lock()
	defer mu.Unlock()
	if !reorderCompleted {
		t.Fatal("expected reorder to have completed")
	}
	if !readCalledAfterReorder {
		t.Fatal("Read was called before reorder completed — this would store stale state")
	}
}

func TestReorder_SequentialRulesSimulateUpdate(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()

	for i := 1; i <= 5; i++ {
		ruleID := 500 + i
		api.preSeed(ruleID, 999, 7)

		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, ruleID, "test_sequential_update",
			api.getCurrent, api.updateOrder, nil,
		)

		markOrderRuleAsDone(ruleID, "test_sequential_update")
		waitForReorder("test_sequential_update")

		if api.putsFor(ruleID) != 1 {
			t.Errorf("rule %d: expected 1 PUT, got %d", ruleID, api.putsFor(ruleID))
		}
	}
}

func TestReorder_ConcurrentRegistration(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	for i := 1; i <= 5; i++ {
		api.preSeed(400+i, 999, 7)
	}

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			reorderWithBeforeReorder(
				OrderRule{Order: idx, Rank: 7}, 400+idx, "test_concurrent",
				api.getCurrent, api.updateOrder, nil,
			)
			time.Sleep(50 * time.Millisecond)
			markOrderRuleAsDone(400+idx, "test_concurrent")
			waitForReorder("test_concurrent")
		}(i)
	}
	wg.Wait()

	for i := 1; i <= 5; i++ {
		if api.putsFor(400+i) != 1 {
			t.Errorf("rule %d: expected 1 PUT, got %d", 400+i, api.putsFor(400+i))
		}
	}
}

// =====================================================
// Diff-based skip — the core SUP-3988 optimization.
// =====================================================

// TestReorder_SkipsRulesAlreadyAtTargetOrder is the "no work to do"
// case the user explicitly asked for. Every registered rule's API
// state already matches its desired Order/Rank, so reorderAll MUST
// issue zero PUTs — no GET-then-PUT, no waste.
func TestReorder_SkipsRulesAlreadyAtTargetOrder(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	// Pre-seed each rule already at its target order.
	for i := 1; i <= 10; i++ {
		api.preSeed(800+i, i, 7)
	}

	for i := 1; i <= 10; i++ {
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 800+i, "test_no_op",
			api.getCurrent, api.updateOrder, nil,
		)
	}
	for i := 1; i <= 10; i++ {
		markOrderRuleAsDone(800+i, "test_no_op")
	}
	waitForReorder("test_no_op")

	if api.putsTotal() != 0 {
		t.Fatalf("expected ZERO PUTs when every rule is already at its target order; got %d", api.putsTotal())
	}
}

// TestReorder_OnlyDriftedRulesArePut: a mixed scenario where 8 of 10
// rules are already at the right position and 2 have drifted. Only
// the 2 drifted rules should be PUT.
func TestReorder_OnlyDriftedRulesArePut(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	for i := 1; i <= 10; i++ {
		api.preSeed(900+i, i, 7) // matches target by default
	}
	// Now drift two rules.
	api.preSeed(900+3, 99, 7) // rule 903: drifted
	api.preSeed(900+7, 99, 7) // rule 907: drifted

	for i := 1; i <= 10; i++ {
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 900+i, "test_partial_drift",
			api.getCurrent, api.updateOrder, nil,
		)
	}
	for i := 1; i <= 10; i++ {
		markOrderRuleAsDone(900+i, "test_partial_drift")
	}
	waitForReorder("test_partial_drift")

	if api.putsFor(903) != 1 {
		t.Errorf("drifted rule 903: expected 1 PUT, got %d", api.putsFor(903))
	}
	if api.putsFor(907) != 1 {
		t.Errorf("drifted rule 907: expected 1 PUT, got %d", api.putsFor(907))
	}
	if api.putsTotal() != 2 {
		t.Errorf("expected exactly 2 PUTs (only the drifted rules); got %d total", api.putsTotal())
	}
}

// =====================================================
// Out-of-range PUT prevention (SUP-3988)
// =====================================================

// TestReorder_DefersOutOfRangePut simulates the SUP-3988 / parallelism
// deadlock scenario:
//
//   - A rule declares Order=10 but the API only has 9 orderable
//     positions when the rule's reorder pass runs.
//   - reorderAll must NOT issue a PUT for order=10 (API would reject
//     with INVALID_INPUT_ARGUMENT).
//   - Instead, reorderAll bails so waitForReorder unblocks and the
//     terraform engine can schedule the next batch of POSTs.
//   - The next rule's registration (which extends apiOrderable) starts
//     a fresh reorder cycle that PUTs the previously deferred rule.
//
// This mirrors the production flow with terraform parallelism: each
// "batch" of in-flight rules' Create calls returns once reorder bails;
// the next batch's first Create kicks off a new cycle.
func TestReorder_DefersOutOfRangePut(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	// 9 filler rules already in API at orders 1..9, plus our test
	// rule (id 600) currently at "orderless" position 0 — included
	// in the map (so "exists") but not counted as orderable. count=9.
	for i := 1; i <= 9; i++ {
		api.preSeed(50000+i, i, 7)
	}
	api.preSeed(600, 0, 7)

	reorderWithBeforeReorder(
		OrderRule{Order: 10, Rank: 7}, 600, "test_out_of_range",
		api.getCurrent, api.updateOrder, nil,
	)
	markOrderRuleAsDone(600, "test_out_of_range")

	// Cycle 1: rule 600 is deferred (desired=10 > count=9). After
	// maxStuckOnSkippedTicks (2 ticks * 50ms ≈ 100ms), reorderAll
	// bails so terraform can schedule the next batch.
	waitForReorder("test_out_of_range")

	if api.putsFor(600) != 0 {
		t.Fatalf("cycle 1: expected ZERO PUTs while desired (10) > count (9); got %d", api.putsFor(600))
	}

	// Simulate the next batch arriving: a new rule registers and
	// the API now has 10 orderable rules. The new registration kicks
	// off cycle 2, which sees rule 600 now in range and PUTs it.
	api.preSeed(50000+10, 10, 7)
	api.preSeed(601, 0, 7)
	reorderWithBeforeReorder(
		OrderRule{Order: 11, Rank: 7}, 601, "test_out_of_range",
		api.getCurrent, api.updateOrder, nil,
	)
	// Bump count again so 601's desired=11 is also in range.
	api.preSeed(50000+11, 11, 7)
	markOrderRuleAsDone(601, "test_out_of_range")
	waitForReorder("test_out_of_range")

	if api.putsFor(600) < 1 {
		t.Fatalf("cycle 2: expected at least one PUT for the previously deferred rule after count caught up; got %d", api.putsFor(600))
	}
}

// TestReorder_NoProgress_GivesUp ensures the reorder loop bails out
// instead of looping forever when the configuration can never
// converge (e.g. declared Order=100 but only 1 orderable rule will
// ever exist).
func TestReorder_NoProgress_GivesUp(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 20 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	api := newFakeAPI()
	api.preSeed(70000, 1, 7) // single orderable rule, count stays at 1.
	api.preSeed(700, 0, 7)   // our test rule, exists but unordered.

	reorderWithBeforeReorder(
		OrderRule{Order: 100, Rank: 7}, 700, "test_no_progress",
		api.getCurrent, api.updateOrder, nil,
	)
	markOrderRuleAsDone(700, "test_no_progress")

	done := make(chan struct{})
	go func() {
		waitForReorder("test_no_progress")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("reorderAll did not give up on no-progress within 2s — it would loop forever")
	}

	if api.putsFor(700) != 0 {
		t.Errorf("expected zero PUTs (order=100 always out of range); got %d", api.putsFor(700))
	}
}

// =====================================================
// Cloud App Control — per-rule-type reorder isolation
// =====================================================
//
// The Cloud App Control API orders rules independently per rule type
// (each type has its own 1..N sequence), and the resource's
// getCurrent/updateOrder callbacks are scoped to a single type's
// endpoint. The resource therefore registers rules under a
// type-scoped registry key (cloudAppRuleResourceType) so each type
// gets its own reorder cycle.

func TestCloudAppRuleResourceType_DistinctKeysPerType(t *testing.T) {
	a := cloudAppRuleResourceType("ENTERPRISE_COLLABORATION")
	b := cloudAppRuleResourceType("HOSTING_PROVIDER")
	if a == b {
		t.Fatalf("expected distinct registry keys per rule type, both were %q", a)
	}
	if a != cloudAppRuleResourceType("ENTERPRISE_COLLABORATION") {
		t.Error("expected the key to be stable for the same rule type")
	}
}

// TestReorder_CloudAppControl_PerTypeIsolation mirrors the customer
// scenario from the 4.8.x ordering report: two rule types applied in
// one configuration, 5 rules each, created concurrently. Each type
// has its own type-scoped API view (getCurrent only returns that
// type's rules). With per-type registry keys, both types must
// converge to their declared 1..5 order, with exactly one PUT per
// drifted rule and no cross-type interference.
func TestReorder_CloudAppControl_PerTypeIsolation(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	collabAPI := newFakeAPI()  // ENTERPRISE_COLLABORATION endpoint
	hostingAPI := newFakeAPI() // HOSTING_PROVIDER endpoint
	collabKey := cloudAppRuleResourceType("ENTERPRISE_COLLABORATION")
	hostingKey := cloudAppRuleResourceType("HOSTING_PROVIDER")

	// Every rule exists in its own type's collection at a wrong
	// position, so each needs exactly one PUT to reach its target.
	for i := 1; i <= 5; i++ {
		collabAPI.preSeed(1000+i, 999, 7)
		hostingAPI.preSeed(2000+i, 999, 7)
	}

	// Register both types concurrently, as terraform parallelism does.
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(2)
		go func(idx int) {
			defer wg.Done()
			reorderWithBeforeReorder(
				OrderRule{Order: idx, Rank: 7}, 1000+idx, collabKey,
				collabAPI.getCurrent, collabAPI.updateOrder, nil,
			)
			markOrderRuleAsDone(1000+idx, collabKey)
			waitForReorder(collabKey)
		}(i)
		go func(idx int) {
			defer wg.Done()
			reorderWithBeforeReorder(
				OrderRule{Order: idx, Rank: 7}, 2000+idx, hostingKey,
				hostingAPI.getCurrent, hostingAPI.updateOrder, nil,
			)
			markOrderRuleAsDone(2000+idx, hostingKey)
			waitForReorder(hostingKey)
		}(i)
	}
	wg.Wait()

	// Both types converged to the declared order — this is what the
	// Read path stores in state, so it also covers the state-file
	// expectation.
	for i := 1; i <= 5; i++ {
		if got := collabAPI.state[1000+i]; got.Order != i {
			t.Errorf("ENTERPRISE_COLLABORATION rule %d: expected final order %d, got %d", 1000+i, i, got.Order)
		}
		if got := hostingAPI.state[2000+i]; got.Order != i {
			t.Errorf("HOSTING_PROVIDER rule %d: expected final order %d, got %d", 2000+i, i, got.Order)
		}
		if got := collabAPI.putsFor(1000 + i); got != 1 {
			t.Errorf("ENTERPRISE_COLLABORATION rule %d: expected exactly 1 PUT, got %d", 1000+i, got)
		}
		if got := hostingAPI.putsFor(2000 + i); got != 1 {
			t.Errorf("HOSTING_PROVIDER rule %d: expected exactly 1 PUT, got %d", 2000+i, got)
		}
	}
	// No cross-type traffic: each endpoint saw only its own rules.
	if collabAPI.putsTotal() != 5 || hostingAPI.putsTotal() != 5 {
		t.Errorf("expected 5 PUTs per type endpoint, got collab=%d hosting=%d", collabAPI.putsTotal(), hostingAPI.putsTotal())
	}
}

// TestReorder_CloudAppControl_SharedKeyMixesTypes_Characterization
// documents the bug the per-type key fixes. Rules of two types are
// registered under ONE shared key, but the cycle's getCurrent is
// scoped to a single type's endpoint (as the resource's closures
// are). The other type's rules are invisible to getCurrent, so they
// are skipped every pass and never reordered; the loop bails via its
// deadlock-breaker instead of hanging. If cloudAppRuleResourceType is
// ever collapsed back to a shared key, this is the behavior users
// get — do not wire the resource this way.
func TestReorder_CloudAppControl_SharedKeyMixesTypes_Characterization(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	hostingAPI := newFakeAPI()
	sharedKey := "test_shared_cac_key"

	// 5 hosting rules visible to getCurrent, drifted.
	for i := 1; i <= 5; i++ {
		hostingAPI.preSeed(2000+i, 999, 7)
	}

	// Register the hosting rules AND 5 collaboration rules (ids
	// 3001..3005) under the same key. The collaboration rules do not
	// exist in hostingAPI — GetByRuleType on the hosting endpoint
	// can't see them.
	for i := 1; i <= 5; i++ {
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 2000+i, sharedKey,
			hostingAPI.getCurrent, hostingAPI.updateOrder, nil,
		)
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 3000+i, sharedKey,
			hostingAPI.getCurrent, hostingAPI.updateOrder, nil,
		)
	}
	for i := 1; i <= 5; i++ {
		markOrderRuleAsDone(2000+i, sharedKey)
		markOrderRuleAsDone(3000+i, sharedKey)
	}

	done := make(chan struct{})
	go func() {
		waitForReorder(sharedKey)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("reorder cycle with permanently invisible rules did not bail out within 5s")
	}

	// The visible type converged; the invisible type never got a PUT.
	for i := 1; i <= 5; i++ {
		if got := hostingAPI.state[2000+i]; got.Order != i {
			t.Errorf("visible rule %d: expected final order %d, got %d", 2000+i, i, got.Order)
		}
		if got := hostingAPI.putsFor(3000 + i); got != 0 {
			t.Errorf("invisible rule %d: expected 0 PUTs (not in getCurrent view), got %d", 3000+i, got)
		}
	}
}

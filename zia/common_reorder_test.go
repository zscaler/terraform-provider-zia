package zia

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
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

// snapshotFromOrders builds a SnapshotProvider that mirrors the registered
// reorder targets back as the "current server state". Tests that don't care
// about skip-if-at-target behaviour can use this; the snapshot will never
// match the target order because the engine compares snapshot against target,
// and the target itself is what we want to apply.
//
// To force PUTs in tests, return a snapshot whose Order/Rank differ from the
// target — easiest is to use Order:0,Rank:0 for every snapshot row. We do so
// by snapshotting at "everything is at order 0", which guarantees the engine
// will issue PUTs for every registered rule.
func snapshotFromOrders(resourceType string) SnapshotProvider {
	return func() ([]RuleSnapshot, error) {
		rules.Lock()
		defer rules.Unlock()
		out := make([]RuleSnapshot, 0, len(rules.orders[resourceType]))
		for id := range rules.orders[resourceType] {
			out = append(out, RuleSnapshot{
				ID:    id,
				Order: 0,
				Rank:  0,
				Body:  nil,
			})
		}
		return out, nil
	}
}

// snapshotFromOrdersMatching returns a SnapshotProvider where each snapshot's
// Order/Rank matches the registered target — exercises the skip-if-at-target
// path: the engine should skip every PUT.
func snapshotFromOrdersMatching(resourceType string) SnapshotProvider {
	return func() ([]RuleSnapshot, error) {
		rules.Lock()
		defer rules.Unlock()
		out := make([]RuleSnapshot, 0, len(rules.orders[resourceType]))
		for id, st := range rules.orders[resourceType] {
			out = append(out, RuleSnapshot{
				ID:    id,
				Order: st.order.Order,
				Rank:  st.order.Rank,
				Body:  nil,
			})
		}
		return out, nil
	}
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
	// Rank 1 rules come first, then rank 7
	expectedIDs := []int{200, 400, 100, 300}
	for i, expected := range expectedIDs {
		if sorted[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, sorted[i].ID)
		}
	}
}

func TestSortOrders_EmptyInput(t *testing.T) {
	sorted := sortOrders(map[int]orderWithState{})
	if len(sorted) != 0 {
		t.Errorf("expected empty result for empty input, got %d items", len(sorted))
	}
}

func TestSortOrders_StableForEqualValues(t *testing.T) {
	input := map[int]orderWithState{
		100: {order: OrderRule{Order: 1, Rank: 7}},
		200: {order: OrderRule{Order: 1, Rank: 7}},
		300: {order: OrderRule{Order: 1, Rank: 7}},
	}
	sorted := sortOrders(input)
	if len(sorted) != 3 {
		t.Fatalf("expected 3 results, got %d", len(sorted))
	}
	// Verify sort.Sort behaviour (not necessarily stable, but well-defined)
	ids := []int{}
	for _, p := range sorted {
		ids = append(ids, p.ID)
	}
	sort.Ints(ids)
	if ids[0] != 100 || ids[1] != 200 || ids[2] != 300 {
		t.Errorf("expected IDs to round-trip through sort, got %v", ids)
	}
}

// =====================================================
// markOrderRuleAsDone Tests
// =====================================================

func TestMarkOrderRuleAsDone_SetsDoneTrue(t *testing.T) {
	resetReorderState()
	rules.Lock()
	rules.orders["test_mark"] = map[int]orderWithState{
		100: {order: OrderRule{Order: 1, Rank: 7}, done: false},
	}
	rules.Unlock()

	markOrderRuleAsDone(100, "test_mark")

	rules.Lock()
	state := rules.orders["test_mark"][100]
	rules.Unlock()
	if !state.done {
		t.Error("expected rule 100 to be marked done")
	}
}

func TestReorderWithBeforeReorder_FirstRule_RegistersWithDoneFalse(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	reorderWithBeforeReorder(
		OrderRule{Order: 1, Rank: 7}, 100, "test_reg",
		snapshotFromOrders("test_reg"),
		func(id int, order OrderRule, body interface{}) error { return nil },
		nil,
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

	// Cleanup: mark done so goroutine can finish
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

	var mu sync.Mutex
	reorderedIDs := map[int]int{}

	snap := snapshotFromOrders("test_all_before")
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		reorderedIDs[id] = order.Order
		mu.Unlock()
		return nil
	}

	for i := 1; i <= 5; i++ {
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 100+i, "test_all_before",
			snap, updateOrder, nil,
		)
	}

	for i := 1; i <= 5; i++ {
		markOrderRuleAsDone(100+i, "test_all_before")
	}

	waitForReorder("test_all_before")

	mu.Lock()
	defer mu.Unlock()

	if len(reorderedIDs) != 5 {
		t.Fatalf("expected 5 reordered rules, got %d", len(reorderedIDs))
	}
	for i := 1; i <= 5; i++ {
		if reorderedIDs[100+i] != i {
			t.Errorf("rule %d: expected order %d, got %d", 100+i, i, reorderedIDs[100+i])
		}
	}
}

func TestReorder_LateArrivingRules_NewCycleStarted(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	// shrink the late-arrival debounce so the test runs quickly
	prevDebounce := lateArrivalDebounce
	lateArrivalDebounce = 100 * time.Millisecond
	defer func() {
		reorderTickInterval = 30 * time.Second
		lateArrivalDebounce = prevDebounce
	}()

	var mu sync.Mutex
	reorderedIDs := map[int]int{}

	snap := snapshotFromOrders("test_late")
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		reorderedIDs[id] = order.Order
		mu.Unlock()
		return nil
	}

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 101, "test_late", snap, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 102, "test_late", snap, updateOrder, nil)
	markOrderRuleAsDone(101, "test_late")
	markOrderRuleAsDone(102, "test_late")

	waitForReorder("test_late")

	mu.Lock()
	firstCycleCount := len(reorderedIDs)
	mu.Unlock()
	if firstCycleCount != 2 {
		t.Fatalf("first cycle: expected 2 reordered rules, got %d", firstCycleCount)
	}

	reorderWithBeforeReorder(OrderRule{Order: 3, Rank: 7}, 103, "test_late", snap, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 4, Rank: 7}, 104, "test_late", snap, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 5, Rank: 7}, 105, "test_late", snap, updateOrder, nil)
	markOrderRuleAsDone(103, "test_late")
	markOrderRuleAsDone(104, "test_late")
	markOrderRuleAsDone(105, "test_late")

	waitForReorder("test_late")

	mu.Lock()
	defer mu.Unlock()

	if reorderedIDs[103] != 3 {
		t.Errorf("rule 103: expected order 3, got %d", reorderedIDs[103])
	}
	if reorderedIDs[104] != 4 {
		t.Errorf("rule 104: expected order 4, got %d", reorderedIDs[104])
	}
	if reorderedIDs[105] != 5 {
		t.Errorf("rule 105: expected order 5, got %d", reorderedIDs[105])
	}
}

func TestReorder_MultipleResourceTypes_Independent(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var mu sync.Mutex
	dnsOrders := map[int]int{}
	sslOrders := map[int]int{}

	dnsSnap := snapshotFromOrders("dns_test")
	sslSnap := snapshotFromOrders("ssl_test")
	dnsUpdate := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		dnsOrders[id] = order.Order
		mu.Unlock()
		return nil
	}
	sslUpdate := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		sslOrders[id] = order.Order
		mu.Unlock()
		return nil
	}

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 201, "dns_test", dnsSnap, dnsUpdate, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 202, "dns_test", dnsSnap, dnsUpdate, nil)

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 301, "ssl_test", sslSnap, sslUpdate, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 302, "ssl_test", sslSnap, sslUpdate, nil)

	markOrderRuleAsDone(201, "dns_test")
	markOrderRuleAsDone(202, "dns_test")
	markOrderRuleAsDone(301, "ssl_test")
	markOrderRuleAsDone(302, "ssl_test")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { waitForReorder("dns_test"); wg.Done() }()
	go func() { waitForReorder("ssl_test"); wg.Done() }()
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if len(dnsOrders) != 2 {
		t.Fatalf("expected 2 DNS reordered rules, got %d", len(dnsOrders))
	}
	if len(sslOrders) != 2 {
		t.Fatalf("expected 2 SSL reordered rules, got %d", len(sslOrders))
	}
	if dnsOrders[201] != 1 {
		t.Errorf("DNS rule 201: expected order 1, got %d", dnsOrders[201])
	}
	if dnsOrders[202] != 2 {
		t.Errorf("DNS rule 202: expected order 2, got %d", dnsOrders[202])
	}
	if sslOrders[301] != 1 {
		t.Errorf("SSL rule 301: expected order 1, got %d", sslOrders[301])
	}
	if sslOrders[302] != 2 {
		t.Errorf("SSL rule 302: expected order 2, got %d", sslOrders[302])
	}
}

func TestReorder_ReadMustHappenAfterWaitForReorder(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var mu sync.Mutex
	reorderCompleted := false
	readCalledAfterReorder := false

	snap := snapshotFromOrders("test_read_order")
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		reorderCompleted = true
		mu.Unlock()
		return nil
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
		snap, updateOrder, nil,
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
	prevDebounce := lateArrivalDebounce
	lateArrivalDebounce = 50 * time.Millisecond
	defer func() {
		reorderTickInterval = 30 * time.Second
		lateArrivalDebounce = prevDebounce
	}()

	var mu sync.Mutex
	finalOrders := map[int]int{}

	snap := snapshotFromOrders("test_sequential_update")
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		finalOrders[id] = order.Order
		mu.Unlock()
		return nil
	}

	for i := 1; i <= 5; i++ {
		ruleID := 500 + i

		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, ruleID, "test_sequential_update",
			snap, updateOrder, nil,
		)

		markOrderRuleAsDone(ruleID, "test_sequential_update")
		waitForReorder("test_sequential_update")

		mu.Lock()
		order, exists := finalOrders[ruleID]
		mu.Unlock()
		if !exists {
			t.Fatalf("rule %d was not reordered", ruleID)
		}
		if order != i {
			t.Errorf("rule %d: expected order %d, got %d", ruleID, i, order)
		}
	}
}

func TestReorder_ConcurrentRegistration(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var mu sync.Mutex
	reorderedIDs := map[int]int{}

	snap := snapshotFromOrders("test_concurrent")
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		reorderedIDs[id] = order.Order
		mu.Unlock()
		return nil
	}

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			reorderWithBeforeReorder(
				OrderRule{Order: idx, Rank: 7}, 400+idx, "test_concurrent",
				snap, updateOrder, nil,
			)
			time.Sleep(50 * time.Millisecond)
			markOrderRuleAsDone(400+idx, "test_concurrent")
			waitForReorder("test_concurrent")
		}(i)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if len(reorderedIDs) != 5 {
		t.Fatalf("expected 5 reordered rules, got %d: %v", len(reorderedIDs), reorderedIDs)
	}
	for i := 1; i <= 5; i++ {
		if reorderedIDs[400+i] != i {
			t.Errorf("rule %d: expected order %d, got %d", 400+i, i, reorderedIDs[400+i])
		}
	}
}

// =====================================================
// New behaviour: skip-if-at-target and late-arrival debounce
// =====================================================

// TestReorder_SkipsPUTWhenAlreadyAtTarget validates that the engine does not
// call updateOrder for a rule whose snapshot Order/Rank already matches the
// requested target. This is the single biggest win on bulk-reorder applies
// (SUP-3988-class slowness).
func TestReorder_SkipsPUTWhenAlreadyAtTarget(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var puts int32
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		atomic.AddInt32(&puts, 1)
		return nil
	}

	for i := 1; i <= 5; i++ {
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 700+i, "test_skip",
			snapshotFromOrdersMatching("test_skip"), updateOrder, nil,
		)
	}
	for i := 1; i <= 5; i++ {
		markOrderRuleAsDone(700+i, "test_skip")
	}
	waitForReorder("test_skip")

	if got := atomic.LoadInt32(&puts); got != 0 {
		t.Errorf("expected 0 PUTs when snapshot already matches target, got %d", got)
	}
}

// TestReorder_PUTsRulesNotYetAtTarget validates that the engine PUTs only
// rules whose snapshot Order/Rank differs from the target — even when other
// rules in the same pass are already at target.
func TestReorder_PUTsOnlyMisalignedRules(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	// Snapshot says rule 800 is at order 1 already, rule 801 is at order 0
	// (misaligned). Engine should only PUT 801.
	mixed := func() ([]RuleSnapshot, error) {
		return []RuleSnapshot{
			{ID: 800, Order: 1, Rank: 7},
			{ID: 801, Order: 0, Rank: 7},
			{ID: 802, Order: 3, Rank: 7},
		}, nil
	}

	var mu sync.Mutex
	putIDs := []int{}
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		putIDs = append(putIDs, id)
		mu.Unlock()
		return nil
	}

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 800, "test_mixed", mixed, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 801, "test_mixed", mixed, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 3, Rank: 7}, 802, "test_mixed", mixed, updateOrder, nil)
	markOrderRuleAsDone(800, "test_mixed")
	markOrderRuleAsDone(801, "test_mixed")
	markOrderRuleAsDone(802, "test_mixed")
	waitForReorder("test_mixed")

	mu.Lock()
	defer mu.Unlock()
	if len(putIDs) != 1 || putIDs[0] != 801 {
		t.Errorf("expected single PUT for rule 801, got %v", putIDs)
	}
}

// TestReorder_LateArrivalDebounceCoalesces validates that multiple late-
// arriving registrations within the lateArrivalDebounce window collapse into
// a single follow-up reorder cycle, instead of each spawning its own.
//
// This is the second key fix for SUP-3988-class slowness with Terraform
// parallelism=N: with N=10 and 40+ in-flight Updates, the previous engine
// would spawn 5+ overlapping cycles, each performing a full GET+PUT sweep.
func TestReorder_LateArrivalDebounceCoalesces(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	prevDebounce := lateArrivalDebounce
	lateArrivalDebounce = 200 * time.Millisecond
	defer func() {
		reorderTickInterval = 30 * time.Second
		lateArrivalDebounce = prevDebounce
	}()

	var passes int32
	snap := func() ([]RuleSnapshot, error) {
		atomic.AddInt32(&passes, 1)
		rules.Lock()
		defer rules.Unlock()
		out := make([]RuleSnapshot, 0, len(rules.orders["test_debounce"]))
		for id := range rules.orders["test_debounce"] {
			out = append(out, RuleSnapshot{ID: id, Order: 0, Rank: 0, Body: nil})
		}
		return out, nil
	}
	updateOrder := func(id int, order OrderRule, body interface{}) error { return nil }

	// First cycle: 2 rules.
	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 901, "test_debounce", snap, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 902, "test_debounce", snap, updateOrder, nil)
	markOrderRuleAsDone(901, "test_debounce")
	markOrderRuleAsDone(902, "test_debounce")
	waitForReorder("test_debounce")

	cyclesAfterFirst := atomic.LoadInt32(&passes)
	if cyclesAfterFirst < 1 {
		t.Fatalf("expected at least 1 reorder pass after first cycle, got %d", cyclesAfterFirst)
	}

	// 4 late arrivals, all registered within the debounce window. They
	// MUST share a single follow-up cycle, not spawn 4 separate ones.
	for i := 0; i < 4; i++ {
		reorderWithBeforeReorder(OrderRule{Order: 3 + i, Rank: 7}, 910+i, "test_debounce", snap, updateOrder, nil)
		markOrderRuleAsDone(910+i, "test_debounce")
		// Stagger registrations within the debounce window
		time.Sleep(20 * time.Millisecond)
	}
	waitForReorder("test_debounce")

	cyclesAfterLate := atomic.LoadInt32(&passes)
	delta := cyclesAfterLate - cyclesAfterFirst
	// Each cycle calls snap at least once; some retries may bump it. We
	// allow up to 3 (initial pass + 2 stability ticks). What we MUST NOT
	// see is one cycle per late arrival (>=4 deltas).
	if delta >= 4 {
		t.Errorf("expected late arrivals to coalesce into <=3 reorder passes, got %d", delta)
	}
}

// TestReorder_SnapshotErrorRetries validates that a transient snapshot error
// (e.g. API 5xx) is logged but does NOT advance lastReorderedSize, so the
// engine retries on the next tick instead of silently giving up.
func TestReorder_SnapshotErrorRetries(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var snapCalls int32
	snap := func() ([]RuleSnapshot, error) {
		n := atomic.AddInt32(&snapCalls, 1)
		if n <= 2 {
			return nil, fmt.Errorf("transient")
		}
		rules.Lock()
		defer rules.Unlock()
		out := make([]RuleSnapshot, 0, len(rules.orders["test_retry"]))
		for id := range rules.orders["test_retry"] {
			out = append(out, RuleSnapshot{ID: id, Order: 0, Rank: 0})
		}
		return out, nil
	}

	var puts int32
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		atomic.AddInt32(&puts, 1)
		return nil
	}

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 950, "test_retry", snap, updateOrder, nil)
	markOrderRuleAsDone(950, "test_retry")
	waitForReorder("test_retry")

	if got := atomic.LoadInt32(&puts); got != 1 {
		t.Errorf("expected 1 PUT after retried snapshot, got %d", got)
	}
	if got := atomic.LoadInt32(&snapCalls); got < 3 {
		t.Errorf("expected at least 3 snapshot attempts (2 fail + 1 succeed), got %d", got)
	}
}

// =====================================================
// Transient "order not allowed" handling
// =====================================================

// TestIsTransientReorderError verifies the regex catches both wordings the
// ZIA API has been observed to return for "this order is not yet placeable".
func TestIsTransientReorderError(t *testing.T) {
	cases := map[string]bool{
		`{"code":"INVALID_INPUT_ARGUMENT","message":"Rule is not allowed at order 34"}`: true,
		`Rule is not allowed at order 5`:                                                true,
		`INVALID_INPUT_ARGUMENT: order out of range`:                                    true,
		`STALE_CONFIGURATION_ERROR`:                                                     false,
		`DUPLICATE_ITEM`:                                                                false,
		``:                                                                              false,
	}
	for msg, want := range cases {
		got := isTransientReorderError(fmt.Errorf("%s", msg))
		if got != want {
			t.Errorf("isTransientReorderError(%q) = %v, want %v", msg, got, want)
		}
	}
	if isTransientReorderError(nil) {
		t.Error("isTransientReorderError(nil) should be false")
	}
}

// TestReorder_TransientErrorDeferredAndRetriedInPass simulates the SUP-3988
// scenario: a rule's target order is rejected on first attempt because its
// lower-order predecessor hasn't been written *yet*, but the engine writes
// the predecessor later in the same pass and the deferred rule then succeeds
// on the in-pass retry.
func TestReorder_TransientErrorDeferredAndRetriedInPass(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	// Snapshot reports both rules at order 0 (so both will be PUT).
	snap := func() ([]RuleSnapshot, error) {
		return []RuleSnapshot{
			{ID: 1001, Order: 0, Rank: 7},
			{ID: 1002, Order: 0, Rank: 7},
		}, nil
	}

	var (
		mu         sync.Mutex
		attempts   = map[int]int{}
		writtenLog []int
	)
	// Rule 1002 (target order 2) gets rejected on first attempt — predecessor
	// rule 1001 (target order 1) hasn't been written yet — and accepted on
	// the in-pass retry once 1001 has been placed.
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		attempts[id]++
		n := attempts[id]
		mu.Unlock()
		if id == 1002 && n == 1 {
			return fmt.Errorf(`{"code":"INVALID_INPUT_ARGUMENT","message":"Rule is not allowed at order 2"}`)
		}
		mu.Lock()
		writtenLog = append(writtenLog, id)
		mu.Unlock()
		return nil
	}

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 1001, "test_transient", snap, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 1002, "test_transient", snap, updateOrder, nil)
	markOrderRuleAsDone(1001, "test_transient")
	markOrderRuleAsDone(1002, "test_transient")
	waitForReorder("test_transient")

	mu.Lock()
	defer mu.Unlock()
	if attempts[1001] != 1 {
		t.Errorf("rule 1001 should be attempted exactly once, got %d", attempts[1001])
	}
	if attempts[1002] != 2 {
		t.Errorf("rule 1002 should be attempted twice (initial + in-pass retry), got %d", attempts[1002])
	}
	// Both should ultimately be in writtenLog.
	if len(writtenLog) != 2 {
		t.Errorf("expected 2 successful writes, got %d (%v)", len(writtenLog), writtenLog)
	}
}

// TestReorder_TransientErrorPersistsAcrossPasses ensures that a rule whose
// transient error doesn't clear in the same pass does NOT block the engine —
// the pass completes, lastReorderedSize advances, stable ticks fire, and
// waitForReorder() unblocks. The next opportunity to fix the rule comes from
// a later size-changed pass or a subsequent Update (matching the engine's
// behaviour for any other write failure).
//
// CRITICAL: this is the property that prevents the deadlock observed in the
// SUP-3988 50-rule reproduction — without it, transient failures on one
// rule would prevent the engine's goroutine from ever returning, blocking
// every Create's waitForReorder() and starving the next Create batch.
func TestReorder_TransientErrorDoesNotBlockEngine(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var (
		mu         sync.Mutex
		writeCalls = map[int]int{}
	)

	snap := func() ([]RuleSnapshot, error) {
		return []RuleSnapshot{
			{ID: 1100, Order: 0, Rank: 7},
		}, nil
	}
	// Rule 1100 always fails with a transient error. The engine must still
	// finish (reach 3 stable ticks and return).
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		mu.Lock()
		writeCalls[id]++
		mu.Unlock()
		return fmt.Errorf(`{"code":"INVALID_INPUT_ARGUMENT","message":"Rule is not allowed at order 1"}`)
	}

	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 1100, "test_no_block", snap, updateOrder, nil)
	markOrderRuleAsDone(1100, "test_no_block")

	// Bound the wait — if waitForReorder hangs, the test fails fast instead
	// of timing out the suite.
	done := make(chan struct{})
	go func() {
		waitForReorder("test_no_block")
		close(done)
	}()
	select {
	case <-done:
		// pass
	case <-time.After(10 * time.Second):
		t.Fatal("reorder goroutine did not return — engine deadlocked on persistent transient errors")
	}

	mu.Lock()
	defer mu.Unlock()
	// Rule 1100 should have been attempted twice: initial + in-pass retry.
	// Both fail, but the engine still completes the pass.
	if writeCalls[1100] < 2 {
		t.Errorf("expected at least 2 write attempts (initial + retry), got %d", writeCalls[1100])
	}
}

// TestReorder_OutOfRangeTargetDoesNotBlockEngine validates the deadlock fix:
// when registered rules have target orders > current snapshot count (because
// their lower-order siblings haven't been Created yet by Terraform), the
// engine must NOT spin forever — it must skip the out-of-range targets,
// finish the pass, hit stability, and return. Otherwise waitForReorder()
// blocks the Create goroutines and the next Create batch never starts.
func TestReorder_OutOfRangeTargetDoesNotBlockEngine(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 50 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	// Snapshot returns only 5 rules even though we'll register 10. The 5
	// missing rules represent the not-yet-Created lower-order siblings.
	// Rules with target order 6..10 have no slot in the policy yet, so
	// their PUTs would be rejected — the engine must skip them.
	snap := func() ([]RuleSnapshot, error) {
		out := make([]RuleSnapshot, 0, 5)
		for id := 2000; id < 2005; id++ {
			out = append(out, RuleSnapshot{ID: id, Order: 0, Rank: 7})
		}
		return out, nil
	}
	var puts int32
	updateOrder := func(id int, order OrderRule, body interface{}) error {
		atomic.AddInt32(&puts, 1)
		return nil
	}

	for i := 0; i < 10; i++ {
		reorderWithBeforeReorder(OrderRule{Order: i + 1, Rank: 7}, 2000+i, "test_oob", snap, updateOrder, nil)
	}
	for i := 0; i < 10; i++ {
		markOrderRuleAsDone(2000+i, "test_oob")
	}

	done := make(chan struct{})
	go func() {
		waitForReorder("test_oob")
		close(done)
	}()
	select {
	case <-done:
		// pass
	case <-time.After(10 * time.Second):
		t.Fatal("reorder goroutine did not return when targets exceeded snapshot size — deadlock")
	}

	// Only rules 2000..2004 (target orders 1..5, all within snapshot count=5)
	// are PUT-able. Rules with target orders 6..10 are skipped.
	if got := atomic.LoadInt32(&puts); got != 5 {
		t.Errorf("expected 5 PUTs (only in-range targets), got %d", got)
	}
}

// TestFamilyWriteLock_SerializesPostAgainstPut directly exercises the
// invariant that motivated the DUPLICATE_ITEM fix: for a single resource
// family, a Create's POST and the engine's reorder PUT must NEVER overlap.
//
// Without the family writer lock, the engine's PUT (which runs without
// rules.Lock held during I/O) and a parallel Terraform Create POST can
// collide on ZIA's edit lock, producing 409 → SDK retry → DUPLICATE_ITEM if
// the originally-retried POST committed server-side.
//
// This test runs N goroutines that each repeatedly try to "POST" and N
// goroutines that try to "PUT". Both call paths route through
// withFamilyWriteLock(family). A counter tracks how many writers are
// "in-flight" (between lock acquisition and release). The test fails if
// that counter ever exceeds 1 — i.e., if any two writes for the same family
// were ever observed running concurrently.
//
// A separate "other_family" PUT path uses a different family key and MUST
// run concurrently with the first family's writes; the test asserts that
// concurrency is preserved across families (the lock is not global).
func TestFamilyWriteLock_SerializesPostAgainstPut(t *testing.T) {
	const family = "test_family_serialization"
	const otherFamily = "other_test_family"

	var (
		inFlightFamily atomic.Int32
		maxInFlight    atomic.Int32

		inFlightOther atomic.Int32
		maxOtherSeen  atomic.Int32

		crossFamilyConcurrent atomic.Int32
	)

	bumpMax := func(cur int32, max *atomic.Int32) {
		for {
			old := max.Load()
			if cur <= old {
				return
			}
			if max.CompareAndSwap(old, cur) {
				return
			}
		}
	}

	doFamilyWrite := func() {
		withFamilyWriteLock(family, func() {
			cur := inFlightFamily.Add(1)
			bumpMax(cur, &maxInFlight)
			// Observe whether the *other* family has a concurrent in-flight
			// writer. We expect this to happen at least once across the run
			// — that proves the lock is per-family, not global.
			if inFlightOther.Load() > 0 {
				crossFamilyConcurrent.Add(1)
			}
			// Simulate a short network call. Sleep is chosen large enough
			// that goroutines collide if no lock were held, but short enough
			// to keep the test fast.
			time.Sleep(2 * time.Millisecond)
			inFlightFamily.Add(-1)
		})
	}

	doOtherFamilyWrite := func() {
		withFamilyWriteLock(otherFamily, func() {
			cur := inFlightOther.Add(1)
			bumpMax(cur, &maxOtherSeen)
			time.Sleep(2 * time.Millisecond)
			inFlightOther.Add(-1)
		})
	}

	const writers = 16
	const opsPerWriter = 25
	var wg sync.WaitGroup
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerWriter; j++ {
				doFamilyWrite()
			}
		}()
	}
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerWriter; j++ {
				doOtherFamilyWrite()
			}
		}()
	}
	wg.Wait()

	if got := maxInFlight.Load(); got > 1 {
		t.Fatalf("family %q: expected at most 1 concurrent writer, observed %d — lock failed to serialize POST vs PUT", family, got)
	}
	if got := maxOtherSeen.Load(); got > 1 {
		t.Fatalf("family %q: expected at most 1 concurrent writer, observed %d — lock failed to serialize", otherFamily, got)
	}
	// Cross-family concurrency must have been observed at least once. If
	// not, the family lock has degenerated into a global lock, which would
	// re-introduce the SUP-3988 starvation we explicitly want to avoid.
	if crossFamilyConcurrent.Load() == 0 {
		t.Fatalf("expected concurrent writes across different families at least once, observed 0 — family lock has degenerated into a global lock")
	}
}

// TestFamilyWriteLock_DistinctFamiliesGetDistinctMutexes is a structural
// check on the lock-creation helper: two different family names must
// resolve to different *sync.Mutex instances, and the same family name must
// always resolve to the same instance (otherwise serialization would be
// broken).
func TestFamilyWriteLock_DistinctFamiliesGetDistinctMutexes(t *testing.T) {
	a1 := familyWriteLock("alpha")
	a2 := familyWriteLock("alpha")
	b := familyWriteLock("beta")

	if a1 != a2 {
		t.Fatalf("familyWriteLock(\"alpha\") returned different *sync.Mutex on second call — registration is not stable")
	}
	if a1 == b {
		t.Fatalf("familyWriteLock(\"alpha\") and familyWriteLock(\"beta\") returned the same *sync.Mutex — distinct families must have distinct locks")
	}
}

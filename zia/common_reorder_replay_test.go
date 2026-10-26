package zia

import (
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"sync"
	"testing"
	"time"
)

// This is the registration sequence from the SUP-4231 update log. IDs are
// anonymized to 1000+declaredOrder. The API initially has three unmanaged rules
// followed by managed rules in positions 4..153; HCL requests positions 1..150.
// Releasing each batch only after its registered rules reach their exact
// targets reproduces the provider's existing completion contract.
func TestReorder_LoggedUpdateBatches(t *testing.T) {
	registration := []int{
		53, 90, 20, 39, 125, 88, 83, 118, 99, 63,
		143, 69, 44, 100, 95, 140, 49, 38, 89, 25,
		24, 136, 54, 10, 124, 64, 50, 134, 84, 5,
		139, 47, 102, 94, 7, 13, 145, 81, 19, 80,
		30, 11, 111, 1, 93, 133, 3, 15, 16, 26,
		29, 46, 42, 121, 98, 149, 18, 79, 59, 110,
		115, 135, 120, 40, 12, 91, 123, 28, 33, 45,
		34, 122, 129, 96, 72, 150, 103, 127, 66, 148,
		51, 128, 87, 101, 104, 48, 113, 77, 108, 92,
		62, 126, 17, 58, 71, 27, 131, 114, 36, 55,
		73, 137, 106, 21, 31, 22, 147, 60, 132, 138,
		142, 78, 74, 37, 41, 68, 4, 56, 141, 43,
		119, 109, 117, 75, 9, 52, 144, 86, 70, 85,
		61, 76, 107, 105, 2, 112, 32, 6, 146, 97,
		35, 23, 116, 8, 65, 130, 82, 57, 67, 14,
	}
	resetReorderState()
	tick := reorderTickInterval
	reorderTickInterval = 2 * time.Millisecond
	t.Cleanup(func() { reorderTickInterval = tick; resetReorderState() })

	var mu sync.Mutex
	list := []int{1, 2, 3}
	for i := 1; i <= 150; i++ {
		list = append(list, 1000+i)
	}
	registered := map[int]bool{}
	puts := 0
	get := func() (map[int]OrderRule, error) {
		mu.Lock()
		defer mu.Unlock()
		out := make(map[int]OrderRule, len(list))
		for i, id := range list {
			out[id] = OrderRule{Order: i + 1, Rank: 7}
		}
		return out, nil
	}
	put := func(id int, target OrderRule) error {
		mu.Lock()
		defer mu.Unlock()
		if !registered[id] || target.Order != id-1000 || target.Rank != 7 {
			return fmt.Errorf("PUT for unregistered rule or changed HCL target: %d %+v", id, target)
		}
		if target.Order < 1 || target.Order > len(list) {
			return fmt.Errorf("PUT at unavailable position %d", target.Order)
		}
		// Model insertion independently of the provider's applyMove helper.
		from := slices.Index(list, id)
		if from < 0 || from+1 == target.Order {
			return fmt.Errorf("missing rule or redundant PUT for %d", id)
		}
		list = slices.Delete(list, from, from+1)
		list = slices.Insert(list, target.Order-1, id)
		puts++
		return nil
	}

	for b := 0; b < len(registration); b += 10 {
		before := slices.Clone(list)
		for _, target := range registration[b : b+10] {
			id := 1000 + target
			mu.Lock()
			registered[id] = true
			mu.Unlock()
			reorderWithBeforeReorder(OrderRule{Order: target, Rank: 7}, id, "logged_updates", get, put, nil)
		}
		for _, target := range registration[b : b+10] {
			markOrderRuleAsDone(1000+target, "logged_updates")
		}
		waitForReorder("logged_updates")
		mu.Lock()
		// Every batch must produce the same complete API list as ordinary
		// stable insertion: registered rules at exact HCL positions, with
		// every other rule retaining its relative order.
		want := make([]int, len(list))
		for id := range registered {
			want[id-1000-1] = id
			if got := slices.Index(list, id) + 1; got != id-1000 {
				t.Errorf("batch %d: rule %d at %d, want HCL target %d", b/10+1, id, got, id-1000)
			}
		}
		free := 0
		for _, id := range before {
			if registered[id] {
				continue
			}
			for want[free] != 0 {
				free++
			}
			want[free] = id
		}
		if !slices.Equal(list, want) {
			t.Errorf("batch %d: changed API placement beyond registered targets", b/10+1)
		}
		mu.Unlock()
	}
	if !slices.Equal(list[150:], []int{1, 2, 3}) {
		t.Errorf("unmanaged rules were not shifted automatically: %v", list[150:])
	}
	if puts > 688 {
		t.Errorf("logged update replay used %d PUTs, want at most 688 (previously 725)", puts)
	}
	t.Logf("logged update replay: %d order PUTs; exact HCL targets verified after every batch", puts)
}

// Exercise both choices against independent slice insertion. The optimization
// may bypass the relative phase only when it completes all targets and saves
// requests; inspecting a trial must never modify the real snapshot.
func TestPreferDirectPlacement_IndependentInsertions(t *testing.T) {
	rng := rand.New(rand.NewSource(4231))
	preferred, fallback := 0, 0
	for sample := 0; sample < 2000; sample++ {
		n := 2 + rng.Intn(18)
		initial := rng.Perm(n)
		current := make(map[int]OrderRule, n)
		for i, id := range initial {
			current[id] = OrderRule{Order: i + 1, Rank: 7}
		}
		m := 1 + rng.Intn(n)
		ids, positions := rng.Perm(n)[:m], rng.Perm(n)[:m]
		slices.Sort(positions)
		candidates := make([]placement, m)
		for i, id := range ids {
			candidates[i] = placement{id: id, target: OrderRule{Order: positions[i] + 1, Rank: 7}}
		}
		keep := inOrderRun(candidates, current)
		before := maps.Clone(current)
		direct := preferDirectPlacement(candidates, current, keep)
		if !maps.Equal(before, current) {
			t.Fatal("evaluating placement changed the API snapshot")
		}
		if !direct {
			fallback++
			continue
		}
		preferred++
		run := func(relative bool) ([]int, int) {
			list := slices.Clone(initial)
			puts := 0
			move := func(c placement) {
				from := slices.Index(list, c.id)
				if from+1 != c.target.Order {
					list = slices.Delete(list, from, from+1)
					list = slices.Insert(list, c.target.Order-1, c.id)
					puts++
				}
			}
			if relative {
				for _, c := range candidates {
					if !keep[c.id] {
						move(c)
					}
				}
			}
			for _, c := range candidates {
				move(c)
			}
			return list, puts
		}
		list, directPuts := run(false)
		_, previousPuts := run(true)
		if directPuts >= previousPuts {
			t.Fatalf("sample %d: selected %d PUTs over previous %d", sample, directPuts, previousPuts)
		}
		for _, c := range candidates {
			if list[c.target.Order-1] != c.id {
				t.Fatalf("sample %d: direct placement did not reach HCL target %+v", sample, c)
			}
		}
	}
	if preferred == 0 || fallback == 0 {
		t.Fatalf("did not exercise both paths: preferred=%d fallback=%d", preferred, fallback)
	}
	t.Logf("validated %d direct selections and %d fallbacks", preferred, fallback)
}

func TestPreferDirectPlacement_PreservesSpecialCases(t *testing.T) {
	base := map[int]OrderRule{1: {Order: 1, Rank: 7}, 2: {Order: 2, Rank: 7}, 3: {Order: 3, Rank: 7}}
	for _, tc := range []struct {
		name    string
		current map[int]OrderRule
		targets []placement
	}{
		{"rank_change", base, []placement{{1, OrderRule{2, 6}}, {2, OrderRule{1, 7}}}},
		{"duplicate_target", base, []placement{{1, OrderRule{2, 7}}, {2, OrderRule{2, 7}}}},
		{"unavailable_target", base, []placement{{1, OrderRule{2, 7}}, {2, OrderRule{4, 7}}}},
		{"missing_rule", base, []placement{{4, OrderRule{1, 7}}, {2, OrderRule{2, 7}}}},
		{"mixed_ranks", map[int]OrderRule{1: {1, 0}, 2: {2, 7}}, []placement{{2, OrderRule{1, 7}}, {1, OrderRule{2, 0}}}},
		{"duplicate_api_order", map[int]OrderRule{1: {1, 7}, 2: {1, 7}}, []placement{{2, OrderRule{1, 7}}, {1, OrderRule{2, 7}}}},
		{"noncontiguous_api", map[int]OrderRule{1: {1, 7}, 2: {3, 7}}, []placement{{2, OrderRule{1, 7}}, {1, OrderRule{2, 7}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if preferDirectPlacement(tc.targets, tc.current, inOrderRun(tc.targets, tc.current)) {
				t.Fatal("special case bypassed the existing placement path")
			}
		})
	}
}

func TestReorder_DirectPlacementRetriesFailedPut(t *testing.T) {
	resetReorderState()
	tick := reorderTickInterval
	reorderTickInterval = time.Millisecond
	t.Cleanup(func() { reorderTickInterval = tick; resetReorderState() })
	api := newFakeAPI()
	for i, id := range []int{1, 6, 4, 3, 2, 5, 8, 7} {
		api.preSeed(id, i+1, 7)
	}
	targets := map[int]int{1: 1, 5: 2, 7: 4, 8: 6}
	candidates := []placement{{1, OrderRule{1, 7}}, {5, OrderRule{2, 7}}, {7, OrderRule{4, 7}}, {8, OrderRule{6, 7}}}
	current, _ := api.getCurrent()
	if !preferDirectPlacement(candidates, current, inOrderRun(candidates, current)) {
		t.Fatal("fixture did not select the optimized path")
	}
	failures := 0
	put := func(id int, target OrderRule) error {
		if failures == 0 {
			failures++
			return errors.New("temporary PUT failure")
		}
		if target.Order != targets[id] {
			return fmt.Errorf("changed HCL target for %d: %+v", id, target)
		}
		return api.updateOrder(id, target)
	}
	for _, c := range candidates {
		reorderWithBeforeReorder(c.target, c.id, "direct_retry", api.getCurrent, put, nil)
	}
	for _, c := range candidates {
		markOrderRuleAsDone(c.id, "direct_retry")
	}
	waitForReorder("direct_retry")
	current, _ = api.getCurrent()
	for id, order := range targets {
		if current[id] != (OrderRule{Order: order, Rank: 7}) {
			t.Errorf("rule %d remained off target after retry: %+v", id, current[id])
		}
	}
}

package zia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"sort"
	gostrings "strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

// fakeURLCategory serves one custom URL category the way the ZIA API handles
// updates: ADD_TO_LIST / REMOVE_FROM_LIST change only the custom URLs, and a
// REMOVE_FROM_LIST that would leave no custom URL, keyword or IP range is
// rejected even when URLs retaining the parent category remain. An update
// without an action replaces the category, and an empty custom URL list is
// accepted there (as in the UI).
type fakeURLCategory struct {
	mu      sync.Mutex
	cat     map[string]any
	updates []urlCategoryUpdate // accepted and rejected PUTs, in order
}

type urlCategoryUpdate struct {
	action   string
	urls     []string
	dbURLs   []string
	rejected bool
}

func newFakeURLCategory(f *unitFixture, id string, urls, dbURLs []string) *fakeURLCategory {
	fake := &fakeURLCategory{cat: map[string]any{
		"id": id, "configuredName": "Category " + id, "superCategory": "USER_DEFINED",
		"customCategory": true, "type": "URL_CATEGORY", "editable": true,
		"urls": urls, "dbCategorizedUrls": dbURLs,
		"keywords": []string{}, "ipRanges": []string{},
	}}
	f.api.OnFunc(http.MethodGet, common.ZIAPath("urlCategories"), func(*http.Request, []byte) common.MockResponse {
		fake.mu.Lock()
		defer fake.mu.Unlock()
		return common.MockResponse{StatusCode: http.StatusOK, Body: []map[string]any{fake.cat}}
	})
	f.api.OnFunc(http.MethodGet, common.ZIAPath("urlCategories", "*"), func(req *http.Request, _ []byte) common.MockResponse {
		fake.mu.Lock()
		defer fake.mu.Unlock()
		if path.Base(req.URL.Path) != id {
			return common.MockResponse{StatusCode: http.StatusNotFound, Body: common.ZIANotFoundBody()}
		}
		return common.MockResponse{StatusCode: http.StatusOK, Body: fake.cat}
	})
	f.api.OnFunc(http.MethodPut, common.ZIAPath("urlCategories", "*"), fake.update)
	return fake
}

func (f *fakeURLCategory) update(req *http.Request, body []byte) common.MockResponse {
	in := map[string]any{}
	_ = json.Unmarshal(body, &in)
	action := req.URL.Query().Get("action")
	sent := urlCategoryUpdate{action: action, urls: urlStrs(in["urls"]), dbURLs: urlStrs(in["dbCategorizedUrls"])}

	f.mu.Lock()
	defer f.mu.Unlock()
	current := urlStrs(f.cat["urls"])
	switch action {
	case "ADD_TO_LIST":
		f.cat["urls"] = urlUnion(current, sent.urls)
	case "REMOVE_FROM_LIST":
		remaining := urlMinus(current, sent.urls)
		if len(remaining) == 0 && len(urlStrs(f.cat["keywords"])) == 0 && len(urlStrs(f.cat["ipRanges"])) == 0 {
			sent.rejected = true
			f.updates = append(f.updates, sent)
			return common.MockResponse{StatusCode: http.StatusBadRequest, Body: map[string]any{
				"code": "INVALID_INPUT_ARGUMENT", "message": "At least 1 URL or keyword or IP range should be entered"}}
		}
		f.cat["urls"] = remaining
	default:
		f.cat["urls"] = sent.urls
		f.cat["dbCategorizedUrls"] = sent.dbURLs
	}
	f.updates = append(f.updates, sent)
	return common.MockResponse{StatusCode: http.StatusOK, Body: f.cat}
}

func (f *fakeURLCategory) state() (urls, dbURLs []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return sortedURLs(urlStrs(f.cat["urls"])), sortedURLs(urlStrs(f.cat["dbCategorizedUrls"]))
}

func urlStrs(v any) []string {
	out := []string{}
	switch l := v.(type) {
	case []string:
		out = append(out, l...)
	case []any:
		for _, s := range l {
			out = append(out, s.(string))
		}
	}
	return out
}

func urlUnion(a, b []string) []string {
	return append(append([]string{}, a...), urlMinus(b, a)...)
}

func urlMinus(a, b []string) []string {
	drop := stringSliceToMap(b)
	out := []string{}
	for _, s := range a {
		if !drop[s] {
			out = append(out, s)
		}
	}
	return out
}

func sortedURLs(s []string) []string {
	out := append([]string{}, s...)
	sort.Strings(out)
	return out
}

func urlAnyList(s ...string) []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

// applyURLCategoryUpdate plans urls from "from" to "to" (db_categorized_urls
// unchanged) and runs the resource Update against the fake category.
func applyURLCategoryUpdate(t *testing.T, id string, from, to, dbURLs []string) (*fakeURLCategory, string) {
	t.Helper()
	return runURLCategoryUpdate(t, id, from, map[string]any{"urls": urlAnyList(to...)}, dbURLs)
}

// urlCategoryPlanConfig is the configuration of the category, with urls set to
// urls, or without the urls attribute when urls is nil.
func urlCategoryPlanConfig(id string, urls []any, dbURLs []string) map[string]any {
	m := map[string]any{
		"category_id": id, "configured_name": "Category " + id, "super_category": "USER_DEFINED",
		"custom_category": true, "type": "URL_CATEGORY", "db_categorized_urls": urlAnyList(dbURLs...),
	}
	if urls != nil {
		m["urls"] = urls
	}
	return m
}

// runURLCategoryUpdate plans from a state holding urls "from" to config (whose
// "urls" entry, if any, is the configured list) and runs the resource Update.
func runURLCategoryUpdate(t *testing.T, id string, from []string, config map[string]any, dbURLs []string) (*fakeURLCategory, string) {
	t.Helper()
	f := newUnitFixture(t)
	fake := newFakeURLCategory(f, id, from, dbURLs)
	r := resourceURLCategories()
	urls, _ := config["urls"].([]any)
	d := f.updateData(r, id, urlCategoryPlanConfig(id, urlAnyList(from...), dbURLs), urlCategoryPlanConfig(id, urls, dbURLs))
	return fake, f.run(r.UpdateContext, d)
}

func assertURLCategoryUpdates(t *testing.T, fake *fakeURLCategory, want ...urlCategoryUpdate) {
	t.Helper()
	fake.mu.Lock()
	got := append([]urlCategoryUpdate{}, fake.updates...)
	fake.mu.Unlock()
	if len(got) != len(want) {
		t.Fatalf("updates = %+v, want %+v", got, want)
	}
	for i := range want {
		if got[i].action != want[i].action || !reflect.DeepEqual(sortedURLs(got[i].urls), sortedURLs(want[i].urls)) || got[i].rejected {
			t.Fatalf("update %d = %+v, want action %q urls %v", i, got[i], want[i].action, want[i].urls)
		}
	}
}

// Issue: emptying urls while db_categorized_urls has entries failed with
// "At least 1 URL or keyword or IP range should be entered".
func TestURLCategories_ClearingURLsUsesFullUpdate(t *testing.T) {
	db := []string{".creditkarma8.com", ".youku8.com"}
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_91", []string{"blockme.example.com"}, nil, db)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "", urls: []string{}})
	if got := fake.updates[0].dbURLs; !reflect.DeepEqual(sortedURLs(got), sortedURLs(db)) {
		t.Fatalf("full update sent db_categorized_urls %v, want %v", got, db)
	}
	urls, dbURLs := fake.state()
	if len(urls) != 0 || !reflect.DeepEqual(dbURLs, sortedURLs(db)) {
		t.Fatalf("category urls = %v, db_categorized_urls = %v; want none and %v", urls, dbURLs, db)
	}
}

// Replacing every URL used to remove them all first, which the API rejects
// in the same way.
func TestURLCategories_ReplacingAllURLsUsesFullUpdate(t *testing.T) {
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_92", []string{"a.example.com", "b.example.com"},
		[]string{"c.example.com"}, []string{".retained.example.com"})
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "", urls: []string{"c.example.com"}})
	if urls, _ := fake.state(); !reflect.DeepEqual(urls, []string{"c.example.com"}) {
		t.Fatalf("category urls = %v, want [c.example.com]", urls)
	}
}

func TestURLCategories_PartialRemovalKeepsRemoveAction(t *testing.T) {
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_93", []string{"a.example.com", "b.example.com"},
		[]string{"a.example.com"}, []string{".retained.example.com"})
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "REMOVE_FROM_LIST", urls: []string{"b.example.com"}})
}

func TestURLCategories_AdditionKeepsAddAction(t *testing.T) {
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_94", []string{"a.example.com"},
		[]string{"a.example.com", "b.example.com"}, nil)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "ADD_TO_LIST", urls: []string{"b.example.com"}})
}

func TestURLCategories_SwapWithOverlapKeepsRemoveThenAdd(t *testing.T) {
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_95", []string{"a.example.com", "b.example.com"},
		[]string{"a.example.com", "c.example.com"}, nil)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake,
		urlCategoryUpdate{action: "REMOVE_FROM_LIST", urls: []string{"b.example.com"}},
		urlCategoryUpdate{action: "ADD_TO_LIST", urls: []string{"c.example.com"}})
	if urls, _ := fake.state(); !reflect.DeepEqual(urls, []string{"a.example.com", "c.example.com"}) {
		t.Fatalf("category urls = %v, want [a.example.com c.example.com]", urls)
	}
}

// planURLs plans the category from a state holding urls "from" to the given
// configured urls (nil: attribute removed from the configuration) and returns
// the planned urls, or nil when the plan has no urls change.
func planURLs(t *testing.T, id string, from []string, configured []any) []string {
	t.Helper()
	r := resourceURLCategories()
	instance := schema.TestResourceDataRaw(t, r.Schema, urlCategoryPlanConfig(id, urlAnyList(from...), nil))
	instance.SetId(id)
	diff, err := r.Diff(context.Background(), instance.State(), terraform.NewResourceConfigRaw(urlCategoryPlanConfig(id, configured, nil)), nil)
	if err != nil {
		t.Fatalf("plan failed: %v", err)
	}
	if diff == nil {
		return nil
	}
	changed := false
	for k := range diff.Attributes {
		if k == "urls.#" || gostrings.HasPrefix(k, "urls.") {
			changed = true
		}
	}
	if !changed {
		return nil
	}
	d, err := schema.InternalMap(r.Schema).Data(instance.State(), diff)
	if err != nil {
		t.Fatalf("reading plan: %v", err)
	}
	return ListToStringList(d, "urls")
}

// Removing the urls attribute from the configuration must plan the removal of
// every URL. The reorder suppression compared the prior state with itself and
// planned no change, so the URLs stayed in ZIA and no drift was reported.
func TestURLCategories_PlanRemovingURLs(t *testing.T) {
	cases := []struct {
		name       string
		from       []string
		configured []any
		want       []string
	}{
		{"attribute removed, one url", []string{"a.example.com"}, nil, []string{}},
		{"attribute removed, two urls", []string{"a.example.com", "b.example.com"}, nil, []string{}},
		{"empty list", []string{"a.example.com"}, []any{}, []string{}},
		{"one of two removed", []string{"a.example.com", "b.example.com"}, urlAnyList("b.example.com"), []string{"b.example.com"}},
		{"one added", []string{"a.example.com"}, urlAnyList("a.example.com", "b.example.com"), []string{"a.example.com", "b.example.com"}},
		{"same length, one replaced", []string{"a.example.com", "b.example.com"}, urlAnyList("a.example.com", "c.example.com"), []string{"a.example.com", "c.example.com"}},
	}
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := planURLs(t, fmt.Sprintf("PLAN_%d", i), c.from, c.configured)
			if got == nil {
				t.Fatalf("no urls change planned, want %v", c.want)
			}
			if !reflect.DeepEqual(sortedURLs(got), sortedURLs(c.want)) {
				t.Fatalf("planned urls = %v, want %v", got, c.want)
			}
		})
	}
}

// Only a reorder of the same URLs is suppressed.
func TestURLCategories_PlanReorderStillSuppressed(t *testing.T) {
	got := planURLs(t, "PLAN_REORDER", []string{"a.example.com", "b.example.com"}, urlAnyList("b.example.com", "a.example.com"))
	if got != nil {
		t.Fatalf("reorder planned a urls change: %v", got)
	}
}

// End to end: urls removed from the configuration while db_categorized_urls
// remains, as in the report. The plan carries the removal and the update
// clears the URLs with a full update.
func TestURLCategories_RemovingURLsAttributeClearsURLs(t *testing.T) {
	db := []string{".creditkarma8.com", ".youku8.com"}
	fake, errMsg := runURLCategoryUpdate(t, "CUSTOM_96", []string{"blockme.example.com"}, map[string]any{}, db)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "", urls: []string{}})
	urls, dbURLs := fake.state()
	if len(urls) != 0 || !reflect.DeepEqual(dbURLs, sortedURLs(db)) {
		t.Fatalf("category urls = %v, db_categorized_urls = %v; want none and %v", urls, dbURLs, db)
	}
}

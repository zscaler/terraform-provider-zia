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
// updates, as observed on a tenant (issue #615): ADD_TO_LIST /
// REMOVE_FROM_LIST apply to urls, dbCategorizedUrls and ipRanges alike (a
// REMOVE_FROM_LIST carrying the full dbCategorizedUrls or ipRanges list
// removed them all), while keywords present in the request replace the
// existing ones. A REMOVE_FROM_LIST that would leave no custom URL, keyword or
// IP range is rejected. An update without an action replaces every list, and
// an empty custom URL list is accepted there (as in the UI).
type fakeURLCategory struct {
	mu      sync.Mutex
	cat     map[string]any
	updates []urlCategoryUpdate // accepted and rejected PUTs, in order
}

type urlCategoryUpdate struct {
	action   string
	urls     []string
	dbURLs   []string
	ipRanges []string
	rejected bool
}

// fullUpdateLists are the lists an update without an action replaces.
var fullUpdateLists = []string{"urls", "dbCategorizedUrls", "keywords", "ipRanges",
	"keywordsRetainingParentCategory", "ipRangesRetainingParentCategory"}

// setList sets a list of the category, e.g. "ipRanges".
func (f *fakeURLCategory) setList(field string, values ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cat[field] = values
}

// list returns a list of the category, sorted.
func (f *fakeURLCategory) list(field string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return sortedURLs(urlStrs(f.cat[field]))
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
	sent := urlCategoryUpdate{action: action, urls: urlStrs(in["urls"]), dbURLs: urlStrs(in["dbCategorizedUrls"]), ipRanges: urlStrs(in["ipRanges"])}

	f.mu.Lock()
	defer f.mu.Unlock()
	current := urlStrs(f.cat["urls"])
	if action != "" {
		if kw, ok := in["keywords"]; ok {
			f.cat["keywords"] = urlStrs(kw)
		}
	}
	switch action {
	case "ADD_TO_LIST":
		f.cat["urls"] = urlUnion(current, sent.urls)
		f.cat["dbCategorizedUrls"] = urlUnion(urlStrs(f.cat["dbCategorizedUrls"]), sent.dbURLs)
		f.cat["ipRanges"] = urlUnion(urlStrs(f.cat["ipRanges"]), sent.ipRanges)
	case "REMOVE_FROM_LIST":
		remaining := urlMinus(current, sent.urls)
		remainingIPs := urlMinus(urlStrs(f.cat["ipRanges"]), sent.ipRanges)
		if len(remaining) == 0 && len(urlStrs(f.cat["keywords"])) == 0 && len(remainingIPs) == 0 {
			sent.rejected = true
			f.updates = append(f.updates, sent)
			return common.MockResponse{StatusCode: http.StatusBadRequest, Body: map[string]any{
				"code": "INVALID_INPUT_ARGUMENT", "message": "At least 1 URL or keyword or IP range should be entered"}}
		}
		f.cat["urls"] = remaining
		f.cat["dbCategorizedUrls"] = urlMinus(urlStrs(f.cat["dbCategorizedUrls"]), sent.dbURLs)
		f.cat["ipRanges"] = remainingIPs
	default:
		for _, field := range fullUpdateLists {
			f.cat[field] = urlStrs(in[field])
		}
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

// assertRetainedURLs checks that every incremental request carried no URLs
// retaining the parent category, and that the category still holds want.
func assertRetainedURLs(t *testing.T, fake *fakeURLCategory, want []string) {
	t.Helper()
	fake.mu.Lock()
	updates := append([]urlCategoryUpdate{}, fake.updates...)
	fake.mu.Unlock()
	for i, u := range updates {
		if u.action != "" && len(u.dbURLs) != 0 {
			t.Fatalf("update %d (%s) sent dbCategorizedUrls %v, want none", i, u.action, u.dbURLs)
		}
	}
	if _, dbURLs := fake.state(); !reflect.DeepEqual(dbURLs, sortedURLs(want)) {
		t.Fatalf("category db_categorized_urls = %v, want %v", dbURLs, sortedURLs(want))
	}
}

// Issue #615: removing some URLs from a category with URLs retaining the
// parent category removed those retained URLs as well.
func TestURLCategories_PartialRemovalKeepsRemoveAction(t *testing.T) {
	db := []string{"keep-c.example.net", "keep-d.example.net"}
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_93", []string{"partial-a.example.com", "partial-b.example.com"},
		[]string{"partial-a.example.com"}, db)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "REMOVE_FROM_LIST", urls: []string{"partial-b.example.com"}})
	assertRetainedURLs(t, fake, db)
	if urls, _ := fake.state(); !reflect.DeepEqual(urls, []string{"partial-a.example.com"}) {
		t.Fatalf("category urls = %v, want [partial-a.example.com]", urls)
	}
}

func TestURLCategories_AdditionKeepsAddAction(t *testing.T) {
	db := []string{"keep-c.example.net"}
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_94", []string{"a.example.com"},
		[]string{"a.example.com", "b.example.com"}, db)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "ADD_TO_LIST", urls: []string{"b.example.com"}})
	assertRetainedURLs(t, fake, db)
}

// URLs and retained URLs changed in one apply: a single full update makes both
// changes, since the incremental requests no longer carry retained URLs.
func TestURLCategories_URLsAndRetainedURLsChangedUsesFullUpdate(t *testing.T) {
	f := newUnitFixture(t)
	id := "CUSTOM_97"
	fake := newFakeURLCategory(f, id, []string{"a.example.com", "b.example.com"}, []string{"keep-c.example.net", "keep-d.example.net"})
	r := resourceURLCategories()
	state := urlCategoryPlanConfig(id, urlAnyList("a.example.com", "b.example.com"), []string{"keep-c.example.net", "keep-d.example.net"})
	config := urlCategoryPlanConfig(id, urlAnyList("a.example.com", "e.example.com"), []string{"keep-c.example.net"})
	if errMsg := f.run(r.UpdateContext, f.updateData(r, id, state, config)); errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "", urls: []string{"a.example.com", "e.example.com"}})
	urls, dbURLs := fake.state()
	if !reflect.DeepEqual(urls, []string{"a.example.com", "e.example.com"}) || !reflect.DeepEqual(dbURLs, []string{"keep-c.example.net"}) {
		t.Fatalf("category urls = %v, db_categorized_urls = %v", urls, dbURLs)
	}
}

func TestURLCategories_SwapWithOverlapKeepsRemoveThenAdd(t *testing.T) {
	fake, errMsg := applyURLCategoryUpdate(t, "CUSTOM_95", []string{"a.example.com", "b.example.com"},
		[]string{"a.example.com", "c.example.com"}, []string{"keep-c.example.net"})
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake,
		urlCategoryUpdate{action: "REMOVE_FROM_LIST", urls: []string{"b.example.com"}},
		urlCategoryUpdate{action: "ADD_TO_LIST", urls: []string{"c.example.com"}})
	assertRetainedURLs(t, fake, []string{"keep-c.example.net"})
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

// runURLCategoryUpdateWithLists runs Update for a category whose state and
// configuration hold urls from -> to plus the given extra list attributes
// (Terraform name -> API name -> values) in both, unless overridden in
// configLists. db_categorized_urls is db in both.
func runURLCategoryUpdateWithLists(t *testing.T, id string, from, to, db []string,
	lists map[string][]string, configLists map[string][]string) (*fakeURLCategory, string) {
	t.Helper()
	apiName := map[string]string{
		"keywords": "keywords", "ip_ranges": "ipRanges",
		"keywords_retaining_parent_category":  "keywordsRetainingParentCategory",
		"ip_ranges_retaining_parent_category": "ipRangesRetainingParentCategory",
	}
	f := newUnitFixture(t)
	fake := newFakeURLCategory(f, id, from, db)
	state := urlCategoryPlanConfig(id, urlAnyList(from...), db)
	config := urlCategoryPlanConfig(id, urlAnyList(to...), db)
	for attr, values := range lists {
		fake.setList(apiName[attr], values...)
		state[attr] = urlAnyList(values...)
		config[attr] = urlAnyList(values...)
	}
	for attr, values := range configLists {
		config[attr] = urlAnyList(values...)
	}
	r := resourceURLCategories()
	return fake, f.run(r.UpdateContext, f.updateData(r, id, state, config))
}

// Step 2 of the tenant test for issue #615: removing one URL also removed the
// IP range, because ipRanges was sent in full with REMOVE_FROM_LIST.
func TestURLCategories_PartialRemovalKeepsIPRangesAndKeywords(t *testing.T) {
	db := []string{".creditkarma9.com", ".youku9.com"}
	fake, errMsg := runURLCategoryUpdateWithLists(t, "CUSTOM_98",
		[]string{"partial-a.example.com", "partial-b.example.com"}, []string{"partial-a.example.com"}, db,
		map[string][]string{"keywords": {"tfissue615kw"}, "ip_ranges": {"3.217.240.0/25"}}, nil)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "REMOVE_FROM_LIST", urls: []string{"partial-b.example.com"}})
	if sent := fake.updates[0].ipRanges; len(sent) != 0 {
		t.Fatalf("REMOVE_FROM_LIST sent ipRanges %v, want none", sent)
	}
	assertRetainedURLs(t, fake, db)
	if got := fake.list("ipRanges"); !reflect.DeepEqual(got, []string{"3.217.240.0/25"}) {
		t.Fatalf("category ip_ranges = %v, want [3.217.240.0/25]", got)
	}
	if got := fake.list("keywords"); !reflect.DeepEqual(got, []string{"tfissue615kw"}) {
		t.Fatalf("category keywords = %v, want [tfissue615kw]", got)
	}
}

func TestURLCategories_AdditionKeepsIPRanges(t *testing.T) {
	fake, errMsg := runURLCategoryUpdateWithLists(t, "CUSTOM_99",
		[]string{"a.example.com"}, []string{"a.example.com", "b.example.com"}, nil,
		map[string][]string{"ip_ranges": {"3.217.240.0/25"}}, nil)
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "ADD_TO_LIST", urls: []string{"b.example.com"}})
	if sent := fake.updates[0].ipRanges; len(sent) != 0 {
		t.Fatalf("ADD_TO_LIST sent ipRanges %v, want none", sent)
	}
	if got := fake.list("ipRanges"); !reflect.DeepEqual(got, []string{"3.217.240.0/25"}) {
		t.Fatalf("category ip_ranges = %v, want [3.217.240.0/25]", got)
	}
}

// URLs and IP ranges changed in one apply: one full update makes both changes.
func TestURLCategories_URLsAndIPRangesChangedUsesFullUpdate(t *testing.T) {
	fake, errMsg := runURLCategoryUpdateWithLists(t, "CUSTOM_100",
		[]string{"a.example.com", "b.example.com"}, []string{"a.example.com"}, nil,
		map[string][]string{"ip_ranges": {"3.217.240.0/25"}},
		map[string][]string{"ip_ranges": {"3.217.240.0/25", "3.217.241.0/25"}})
	if errMsg != "" {
		t.Fatalf("update failed: %s", errMsg)
	}
	assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "", urls: []string{"a.example.com"}})
	if got := fake.list("ipRanges"); !reflect.DeepEqual(got, []string{"3.217.240.0/25", "3.217.241.0/25"}) {
		t.Fatalf("category ip_ranges = %v", got)
	}
}

// Lists whose behaviour with an action is not verified keep the category on
// full updates, so a partial URL change cannot strip them.
func TestURLCategories_UnverifiedListsUseFullUpdate(t *testing.T) {
	for _, attr := range []string{"keywords_retaining_parent_category", "ip_ranges_retaining_parent_category"} {
		t.Run(attr, func(t *testing.T) {
			value := "retainedkw"
			if attr == "ip_ranges_retaining_parent_category" {
				value = "13.107.10.154/31"
			}
			fake, errMsg := runURLCategoryUpdateWithLists(t, "CUSTOM_"+attr,
				[]string{"a.example.com", "b.example.com"}, []string{"a.example.com"}, nil,
				map[string][]string{attr: {value}}, nil)
			if errMsg != "" {
				t.Fatalf("update failed: %s", errMsg)
			}
			assertURLCategoryUpdates(t, fake, urlCategoryUpdate{action: "", urls: []string{"a.example.com"}})
		})
	}
}

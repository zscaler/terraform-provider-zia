package zia

// Unit test fixture: runs resource CRUD functions against an in-memory mock of
// the ZIA API (the zscaler-sdk-go OneAPI test harness), so provider logic is
// tested offline, without TF_ACC, credentials or a tenant.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

type unitFixture struct {
	t      *testing.T
	api    *common.APITest
	client *Client
}

// newUnitFixture starts a mock ZIA API and returns a provider client wired to
// it. It also resets the rule reorder engine and shortens its tick, so tests
// that create or update ordered rules finish in milliseconds.
func newUnitFixture(t *testing.T) *unitFixture {
	t.Helper()
	api := common.NewZIATest(t)

	tick := reorderTickInterval
	reorderTickInterval = 10 * time.Millisecond
	resetReorderEngine()
	t.Cleanup(func() {
		reorderTickInterval = tick
		resetReorderEngine()
	})
	return &unitFixture{t: t, api: api, client: &Client{Service: api.Service}}
}

func resetReorderEngine() {
	rules.Lock()
	rules.orders = map[string]map[int]orderWithState{}
	rules.orderer = nil
	rules.reorderDone = map[string]chan struct{}{}
	rules.Unlock()
}

// resourceData builds the ResourceData a CRUD function receives.
func (f *unitFixture) resourceData(r *schema.Resource, id string, raw map[string]any) *schema.ResourceData {
	f.t.Helper()
	d := schema.TestResourceDataRaw(f.t, r.Schema, raw)
	d.SetId(id)
	return d
}

// updateData builds the ResourceData an update receives: the prior state plus
// the planned change to config, so HasChange behaves as in a plan.
func (f *unitFixture) updateData(r *schema.Resource, id string, state, config map[string]any) *schema.ResourceData {
	f.t.Helper()
	instance := f.resourceData(r, id, state).State()
	planned, err := r.Diff(context.Background(), instance, terraform.NewResourceConfigRaw(config), f.client)
	if err != nil {
		f.t.Fatalf("plan failed: %v", err)
	}
	d, err := schema.InternalMap(r.Schema).Data(instance, planned)
	if err != nil {
		f.t.Fatalf("building update data: %v", err)
	}
	return d
}

// run calls a CRUD function and returns its first error summary ("" when none).
func (f *unitFixture) run(fn crudFunc, d *schema.ResourceData) string {
	f.t.Helper()
	for _, dg := range fn(context.Background(), d, f.client) {
		if dg.Severity == diag.Error {
			return dg.Summary
		}
	}
	return ""
}

// fakeRuleList is an in-memory ordered rule list served at a list path, with
// each rule at itemPath/<id> (itemPath ends in "/*"). It enforces the API's
// ordering rules: orders are contiguous from 1, a PUT moves the rule to its
// order and shifts the others, and an order past the end of the list is
// rejected with the API's error. Like the API, it processes writes one at a
// time while still answering reads, and a delete can be slow (deleteDelay), so
// writes queued behind it are processed after the list has shrunk. onGetRule
// runs after a read of one rule has built its response, so a test can land a
// delete between a read and the write that follows it.
//
// With parentPath set, the rules are the sub-rules of a parent rule: GET
// parentPath returns the parent with the rules in "subRules".
type fakeRuleList struct {
	mu           sync.Mutex // guards the list
	writeMu      sync.Mutex // serializes writes, as the API does
	deleteDelay  time.Duration
	createDelay  time.Duration
	list         []map[string]any // position = index + 1
	nextID       int
	createOrders []int            // accepted POST positions, in API processing order
	puts         []map[string]any // accepted PUT bodies
	rejected     int              // writes rejected for an invalid order
	onGetRule    func(id int)
	onGetAll     func() // runs after a read of the list has built its response
}

// newFakeRuleList serves n rules with IDs 1001..1000+n, each carrying fields.
func newFakeRuleList(f *unitFixture, n int, listPath, itemPath, parentPath string, fields map[string]any) *fakeRuleList {
	fake := &fakeRuleList{nextID: 1000 + n}
	for i := 1; i <= n; i++ {
		rule := map[string]any{"id": 1000 + i, "name": ruleName(i), "rank": 7}
		for k, v := range fields {
			rule[k] = v
		}
		fake.list = append(fake.list, rule)
	}
	if listPath != "" {
		f.api.OnFunc(http.MethodGet, listPath, fake.getAll)
		f.api.OnFunc(http.MethodPost, listPath, fake.create)
	}
	if parentPath != "" {
		parentID, _ := strconv.Atoi(path.Base(parentPath))
		f.api.OnFunc(http.MethodGet, parentPath, func(*http.Request, []byte) common.MockResponse {
			fake.mu.Lock()
			defer fake.mu.Unlock()
			return common.MockResponse{StatusCode: http.StatusOK, Body: map[string]any{"id": parentID, "subRules": fake.numbered()}}
		})
	}
	f.api.OnFunc(http.MethodGet, itemPath, fake.get)
	f.api.OnFunc(http.MethodPut, itemPath, fake.update)
	f.api.OnFunc(http.MethodDelete, itemPath, fake.remove)
	return fake
}

func ruleName(i int) string { return fmt.Sprintf("Rule %02d", i) }

func asInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	}
	return 0
}

func (f *fakeRuleList) numbered() []map[string]any {
	out := make([]map[string]any, len(f.list))
	for i, r := range f.list {
		c := make(map[string]any, len(r)+1)
		for k, v := range r {
			c[k] = v
		}
		c["order"] = i + 1
		out[i] = c
	}
	return out
}

func (f *fakeRuleList) index(id int) int {
	for i, r := range f.list {
		if asInt(r["id"]) == id {
			return i
		}
	}
	return -1
}

func requestRuleID(req *http.Request) int {
	id, _ := strconv.Atoi(path.Base(req.URL.Path))
	return id
}

func notAllowedAt(order int) common.MockResponse {
	return common.MockResponse{StatusCode: http.StatusBadRequest, Body: map[string]any{
		"code": "INVALID_INPUT_ARGUMENT", "message": fmt.Sprintf("Rule is not allowed at order %d", order)}}
}

func (f *fakeRuleList) getAll(*http.Request, []byte) common.MockResponse {
	f.mu.Lock()
	resp := common.MockResponse{StatusCode: http.StatusOK, Body: f.numbered()}
	hook := f.onGetAll
	f.mu.Unlock()
	if hook != nil {
		hook()
	}
	return resp
}

func (f *fakeRuleList) get(req *http.Request, _ []byte) common.MockResponse {
	id := requestRuleID(req)
	f.mu.Lock()
	resp := common.MockResponse{StatusCode: http.StatusNotFound, Body: common.ZIANotFoundBody()}
	if i := f.index(id); i >= 0 {
		resp = common.MockResponse{StatusCode: http.StatusOK, Body: f.numbered()[i]}
	}
	hook := f.onGetRule
	f.mu.Unlock()
	if hook != nil {
		hook(id)
	}
	return resp
}

func (f *fakeRuleList) update(req *http.Request, body []byte) common.MockResponse {
	in := map[string]any{}
	_ = json.Unmarshal(body, &in)
	f.writeMu.Lock()
	defer f.writeMu.Unlock()
	f.mu.Lock()
	defer f.mu.Unlock()
	i := f.index(requestRuleID(req))
	if i < 0 {
		return common.MockResponse{StatusCode: http.StatusNotFound, Body: common.ZIANotFoundBody()}
	}
	order := asInt(in["order"])
	if order < 1 || order > len(f.list) {
		f.rejected++
		return notAllowedAt(order)
	}
	in["id"] = f.list[i]["id"]
	f.list = append(f.list[:i], f.list[i+1:]...)
	f.list = append(f.list[:order-1], append([]map[string]any{in}, f.list[order-1:]...)...)
	f.puts = append(f.puts, in)
	return common.MockResponse{StatusCode: http.StatusOK, Body: f.numbered()[order-1]}
}

func (f *fakeRuleList) create(_ *http.Request, body []byte) common.MockResponse {
	in := map[string]any{}
	_ = json.Unmarshal(body, &in)
	f.writeMu.Lock()
	defer f.writeMu.Unlock()
	time.Sleep(f.createDelay)
	f.mu.Lock()
	defer f.mu.Unlock()
	order := asInt(in["order"])
	if order < 1 || order > len(f.list)+1 {
		f.rejected++
		return notAllowedAt(order)
	}
	f.createOrders = append(f.createOrders, order)
	f.nextID++
	in["id"] = f.nextID
	f.list = append(f.list[:order-1], append([]map[string]any{in}, f.list[order-1:]...)...)
	return common.MockResponse{StatusCode: http.StatusOK, Body: f.numbered()[order-1]}
}

func (f *fakeRuleList) remove(req *http.Request, _ []byte) common.MockResponse {
	f.writeMu.Lock()
	defer f.writeMu.Unlock()
	time.Sleep(f.deleteDelay)
	f.deleteRule(requestRuleID(req))
	return common.MockResponse{StatusCode: http.StatusNoContent}
}

func (f *fakeRuleList) deleteRule(id int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if i := f.index(id); i >= 0 {
		f.list = append(f.list[:i], f.list[i+1:]...)
	}
}

// orderOf returns the current order of a rule (0 when it does not exist).
func (f *fakeRuleList) orderOf(id int) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.index(id) + 1
}

func (f *fakeRuleList) putCount(id int) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, p := range f.puts {
		if asInt(p["id"]) == id {
			n++
		}
	}
	return n
}

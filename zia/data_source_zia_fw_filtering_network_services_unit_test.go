package zia

import (
	"context"
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

// The network services list endpoint is not paginated. On some clouds it
// rejects a page size outright ("Invalid Page Size 1,000", issue #617), so the
// name lookups used by the data source and by import must send a single request
// without page or pageSize.

const networkServicesPath = "/zia/api/v1/networkServices"

// serveUnpagedNetworkServices answers like that cloud: a request carrying page
// or pageSize is rejected, any other request returns every service. It returns
// a pointer to the number of list requests received.
func serveUnpagedNetworkServices(f *unitFixture) *int {
	calls := 0
	f.api.OnFunc(http.MethodGet, networkServicesPath, func(req *http.Request, _ []byte) common.MockResponse {
		calls++
		q := req.URL.Query()
		if q.Has("pageSize") || q.Has("page") {
			return common.MockResponse{StatusCode: http.StatusBadRequest, Body: map[string]any{
				"code":    "INVALID_INPUT_ARGUMENT",
				"message": "Invalid Page Size " + q.Get("pageSize"),
			}}
		}
		return common.MockResponse{StatusCode: http.StatusOK, Body: []map[string]any{
			{"id": 101, "name": "HTTP", "type": "PREDEFINED", "protocol": "TCP"},
			{"id": 102, "name": "HTTPS", "type": "PREDEFINED", "protocol": "TCP",
				"destTcpPorts": []map[string]any{{"start": 443}}},
		}}
	})
	return &calls
}

func TestNetworkServices_DataSourceLookupByNameSendsNoPageSize(t *testing.T) {
	f := newUnitFixture(t)
	calls := serveUnpagedNetworkServices(f)

	r := dataSourceFWNetworkServices()
	// Same configuration as issue #617: a lower-case name of a predefined service.
	d := f.resourceData(r, "", map[string]any{"name": "https"})
	if errMsg := f.run(r.ReadContext, d); errMsg != "" {
		t.Fatalf("lookup failed: %s", errMsg)
	}
	if d.Id() != "102" {
		t.Fatalf("id = %q, want 102", d.Id())
	}
	if got := d.Get("name"); got != "HTTPS" {
		t.Fatalf("name = %v, want HTTPS", got)
	}
	if *calls != 1 {
		t.Fatalf("list requests = %d, want 1", *calls)
	}
}

func TestNetworkServices_DataSourceLookupByProtocolSendsNoPageSize(t *testing.T) {
	f := newUnitFixture(t)
	calls := serveUnpagedNetworkServices(f)

	r := dataSourceFWNetworkServices()
	d := f.resourceData(r, "", map[string]any{"protocol": "TCP"})
	if errMsg := f.run(r.ReadContext, d); errMsg != "" {
		t.Fatalf("lookup failed: %s", errMsg)
	}
	if d.Id() == "" {
		t.Fatal("no network service returned")
	}
	if *calls != 1 {
		t.Fatalf("list requests = %d, want 1", *calls)
	}
}

func TestNetworkServices_ImportByNameSendsNoPageSize(t *testing.T) {
	f := newUnitFixture(t)
	serveUnpagedNetworkServices(f)

	r := resourceFWNetworkServices()
	d := f.resourceData(r, "HTTPS", map[string]any{})
	got, err := r.Importer.StateContext(context.Background(), d, f.client)
	if err != nil {
		t.Fatalf("import by name failed: %v", err)
	}
	if len(got) != 1 || got[0].Id() != "102" {
		t.Fatalf("imported id = %q, want 102", got[0].Id())
	}
	if v := got[0].Get("network_service_id"); v != 102 {
		t.Fatalf("network_service_id = %v, want 102", v)
	}
}

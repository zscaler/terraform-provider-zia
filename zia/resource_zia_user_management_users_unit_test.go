package zia

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

const usersListPath = "/zia/api/v1/users"

// fakeUsers serves GET /users like the API: the name parameter matches part
// of a user's name or email.
type fakeUsers struct {
	mu    sync.Mutex
	users []map[string]any
}

func newFakeUsers(f *unitFixture, n int) *fakeUsers {
	fake := &fakeUsers{}
	for i := 0; i < n; i++ {
		fake.users = append(fake.users, map[string]any{
			"id":    userID(i),
			"name":  userName(i),
			"email": userName(i) + "@example.com",
		})
	}
	f.api.OnFunc(http.MethodGet, usersListPath, func(req *http.Request, _ []byte) common.MockResponse {
		fake.mu.Lock()
		defer fake.mu.Unlock()
		name := strings.ToLower(req.URL.Query().Get("name"))
		out := []map[string]any{}
		for _, u := range fake.users {
			if name == "" || strings.Contains(strings.ToLower(u["name"].(string)), name) ||
				strings.Contains(strings.ToLower(u["email"].(string)), name) {
				out = append(out, u)
			}
		}
		return common.SuccessResponse(out)
	})
	return fake
}

func (u *fakeUsers) rename(id int, name string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, user := range u.users {
		if user["id"] == id {
			user["name"] = name
		}
	}
}

func (u *fakeUsers) remove(id int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for i, user := range u.users {
		if user["id"] == id {
			u.users = append(u.users[:i], u.users[i+1:]...)
			return
		}
	}
}

func userID(i int) int      { return 1000 + i }
func userName(i int) string { return fmt.Sprintf("user-%03d", i) }
func userState(i int) map[string]any {
	return map[string]any{"user_id": userID(i), "name": userName(i), "email": userName(i) + "@example.com"}
}

// userListReads returns the GET /users requests: filtered by name, and of the
// whole list.
func userListReads(f *unitFixture) (filtered, whole int) {
	for _, r := range f.api.Server.AllRequests() {
		if r.Method != http.MethodGet || r.Path != usersListPath {
			continue
		}
		if strings.Contains(r.Query, "name=") {
			filtered++
		} else {
			whole++
		}
	}
	return filtered, whole
}

func refreshUser(t *testing.T, f *unitFixture, i int) *schema.ResourceData {
	t.Helper()
	r := resourceUserManagement()
	d := f.resourceData(r, strconv.Itoa(userID(i)), userState(i))
	if msg := f.run(r.ReadContext, d); msg != "" {
		t.Fatalf("refresh of user %d: %s", i, msg)
	}
	return d
}

func TestUserManagement_RefreshesShareOneListRead(t *testing.T) {
	f := newUnitFixture(t)
	newFakeUsers(f, 334)

	for i := 0; i < 334; i++ {
		d := refreshUser(t, f, i)
		if d.Id() != strconv.Itoa(userID(i)) || d.Get("name") != userName(i) {
			t.Fatalf("user %d read as id=%q name=%q", i, d.Id(), d.Get("name"))
		}
	}

	filtered, whole := userListReads(f)
	if filtered != userLookupsBeforeSharedList || whole != 1 {
		t.Fatalf("334 refreshes sent %d filtered and %d whole-list requests, want %d and 1", filtered, whole, userLookupsBeforeSharedList)
	}
}

func TestUserManagement_ConcurrentRefreshesShareOneListRead(t *testing.T) {
	f := newUnitFixture(t)
	newFakeUsers(f, 100)

	var wg sync.WaitGroup
	errs := make(chan string, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r := resourceUserManagement()
			d := f.resourceData(r, strconv.Itoa(userID(i)), userState(i))
			if msg := f.run(r.ReadContext, d); msg != "" {
				errs <- msg
			} else if d.Id() != strconv.Itoa(userID(i)) {
				errs <- fmt.Sprintf("user %d removed from state", i)
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for msg := range errs {
		t.Fatal(msg)
	}

	filtered, whole := userListReads(f)
	if filtered != userLookupsBeforeSharedList || whole != 1 {
		t.Fatalf("100 concurrent refreshes sent %d filtered and %d whole-list requests, want %d and 1", filtered, whole, userLookupsBeforeSharedList)
	}
}

func TestUserManagement_DeletedUserRemovedFromState(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeUsers(f, 10)

	// Deleted outside Terraform: the refreshes below read the shared list
	// after the delete.
	fake.remove(userID(7))
	for i := 0; i < 5; i++ {
		refreshUser(t, f, i)
	}

	d := refreshUser(t, f, 7)
	if d.Id() != "" {
		t.Fatalf("deleted user kept in state with id %q", d.Id())
	}
	// The copy predates the lookup, so the user is looked up again in a list
	// read after it started, once, before it is removed from state.
	if _, whole := userListReads(f); whole != 2 {
		t.Fatalf("sent %d whole-list requests, want 2", whole)
	}
}

func TestUserManagement_UserDeletedAfterListReadKeptUntilNextRun(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeUsers(f, 10)

	for i := 0; i < 5; i++ {
		refreshUser(t, f, i)
	}
	// Deleted outside Terraform after the shared list was read: the copy is a
	// snapshot for this command, so the next plan detects the delete.
	fake.remove(userID(7))
	if d := refreshUser(t, f, 7); d.Id() != strconv.Itoa(userID(7)) {
		t.Fatalf("user read as id %q, want it from the shared list", d.Id())
	}
	if _, whole := userListReads(f); whole != 1 {
		t.Fatalf("sent %d whole-list requests, want 1", whole)
	}
}

func TestUserManagement_WriteDiscardsSharedList(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeUsers(f, 10)

	for i := 0; i < 5; i++ {
		refreshUser(t, f, i)
	}
	fake.rename(userID(6), "renamed")

	// Without a write the shared list still serves the old name.
	if d := refreshUser(t, f, 6); d.Get("name") != userName(6) {
		t.Fatalf("name %q, want the shared list's %q", d.Get("name"), userName(6))
	}

	userListFor(f.client.Service).invalidate()
	if d := refreshUser(t, f, 6); d.Get("name") != "renamed" {
		t.Fatalf("name after a write %q, want renamed", d.Get("name"))
	}
	if _, whole := userListReads(f); whole != 2 {
		t.Fatalf("sent %d whole-list requests, want 2", whole)
	}
}

func TestUserManagement_ReadAfterWriteUsesFilteredLookup(t *testing.T) {
	f := newUnitFixture(t)
	newFakeUsers(f, 10)

	for i := 0; i < 5; i++ {
		refreshUser(t, f, i)
	}
	filteredBefore, wholeBefore := userListReads(f)

	r := resourceUserManagement()
	d := f.resourceData(r, strconv.Itoa(userID(8)), userState(8))
	if diags := readUserManagement(t.Context(), d, f.client, false); diags.HasError() {
		t.Fatalf("read after write: %v", diags)
	}
	filtered, whole := userListReads(f)
	if filtered != filteredBefore+1 || whole != wholeBefore {
		t.Fatalf("read after write sent %d filtered and %d whole-list requests, want 1 and 0", filtered-filteredBefore, whole-wholeBefore)
	}
}

func TestUserManagement_DataSourceLookupsShareOneListRead(t *testing.T) {
	f := newUnitFixture(t)
	newFakeUsers(f, 50)

	ds := dataSourceUserManagement()
	for i := 0; i < 50; i++ {
		d := f.resourceData(ds, "", map[string]any{"name": userName(i)})
		if msg := f.run(ds.ReadContext, d); msg != "" {
			t.Fatalf("lookup of user %d: %s", i, msg)
		}
		if d.Id() != strconv.Itoa(userID(i)) {
			t.Fatalf("lookup of %s returned id %q", userName(i), d.Id())
		}
	}
	// By email and by ID from the same list.
	byEmail := f.resourceData(ds, "", map[string]any{"email": userName(3) + "@example.com"})
	byID := f.resourceData(ds, strconv.Itoa(userID(4)), map[string]any{})
	for want, d := range map[int]*schema.ResourceData{userID(3): byEmail, userID(4): byID} {
		if msg := f.run(ds.ReadContext, d); msg != "" {
			t.Fatalf("lookup of %d: %s", want, msg)
		}
		if d.Id() != strconv.Itoa(want) {
			t.Fatalf("lookup of %d returned id %q", want, d.Id())
		}
	}

	filtered, whole := userListReads(f)
	if filtered != userLookupsBeforeSharedList || whole != 1 {
		t.Fatalf("52 lookups sent %d filtered and %d whole-list requests, want %d and 1", filtered, whole, userLookupsBeforeSharedList)
	}
}

func TestUserManagement_DataSourceFindsUserCreatedAfterListRead(t *testing.T) {
	f := newUnitFixture(t)
	fake := newFakeUsers(f, 5)

	for i := 0; i < 5; i++ {
		refreshUser(t, f, i)
	}
	// Created outside Terraform after the shared list was read.
	fake.mu.Lock()
	fake.users = append(fake.users, map[string]any{"id": 2000, "name": "late-user", "email": "late@example.com"})
	fake.mu.Unlock()

	ds := dataSourceUserManagement()
	d := f.resourceData(ds, "", map[string]any{"name": "late-user"})
	if msg := f.run(ds.ReadContext, d); msg != "" {
		t.Fatalf("lookup: %s", msg)
	}
	if d.Id() != "2000" {
		t.Fatalf("id %q, want 2000", d.Id())
	}
}

func TestUserManagement_DataSourceSearchDoesNotUseSharedList(t *testing.T) {
	f := newUnitFixture(t)
	newFakeUsers(f, 10)

	for i := 0; i < 5; i++ {
		refreshUser(t, f, i)
	}
	_, wholeBefore := userListReads(f)

	ds := dataSourceUserManagement()
	d := f.resourceData(ds, "", map[string]any{"name": userName(2), "search": "[?contains(name, 'user-00')]"})
	if msg := f.run(ds.ReadContext, d); msg != "" {
		t.Fatalf("lookup with search: %s", msg)
	}
	if d.Id() != strconv.Itoa(userID(2)) {
		t.Fatalf("id %q", d.Id())
	}
	// A lookup with search reads the list itself, filtered by its expression.
	if _, whole := userListReads(f); whole != wholeBefore+1 {
		t.Fatalf("sent %d whole-list requests, want 1", whole-wholeBefore)
	}
}

func TestUserManagement_ImportByNameUsesSharedList(t *testing.T) {
	f := newUnitFixture(t)
	newFakeUsers(f, 20)

	r := resourceUserManagement()
	for i := 0; i < 20; i++ {
		d := f.resourceData(r, userName(i), map[string]any{})
		out, err := r.Importer.StateContext(t.Context(), d, f.client)
		if err != nil {
			t.Fatalf("import of %s: %v", userName(i), err)
		}
		if out[0].Id() != strconv.Itoa(userID(i)) {
			t.Fatalf("import of %s: id %q", userName(i), out[0].Id())
		}
	}
	filtered, whole := userListReads(f)
	if filtered != userLookupsBeforeSharedList || whole != 1 {
		t.Fatalf("20 imports sent %d filtered and %d whole-list requests, want %d and 1", filtered, whole, userLookupsBeforeSharedList)
	}
}

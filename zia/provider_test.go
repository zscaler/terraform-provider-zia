package zia

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testNamePrefix = "tf-acc-test-"

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"zia": testAccProvider,
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"zia": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = Provider()
}

func testAccPreCheck(t *testing.T) {
	err := accPreCheck()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if v := os.Getenv("ZIA_USERNAME"); v == "" {
		t.Fatal("ZIA_USERNAME must be set for acceptance tests.")
	}
	if v := os.Getenv("ZIA_PASSWORD"); v == "" {
		t.Fatal("ZIA_PASSWORD must be set for acceptance tests.")
	}
	if v := os.Getenv("ZIA_API_KEY"); v == "" {
		t.Fatal("ZIA_API_KEY must be set for acceptance tests.")
	}
	if v := os.Getenv("ZIA_BASE_URL"); v == "" {
		t.Fatal("ZIA_BASE_URL must be set for acceptance tests.")
	}
}

func accPreCheck() error {
	if v := os.Getenv("ZIA_USERNAME"); v == "" {
		return errors.New("ZIA_USERNAME must be set for acceptance tests")
	}
	username := os.Getenv("ZIA_USERNAME")
	password := os.Getenv("ZIA_PASSWORD")
	api_key := os.Getenv("ZIA_API_KEY")
	zia_base_url := os.Getenv("ZIA_BASE_URL")
	if username == "" && (username == "" || password == "" || api_key == "" || zia_base_url == "") {
		return errors.New("either ZIA_USERNAME or ZIA_PASSWORD, ZIA_API_KEY and ZIA_BASE_URL must be set for acceptance tests")
	}
	return nil
}

package zia

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	if err := accPreCheck(); err != nil {
		t.Fatalf("%v", err)
	}
}

// accPreCheck checks if the necessary environment variables for acceptance tests are set.
func accPreCheck() error {
	username := os.Getenv("ZIA_USERNAME")
	password := os.Getenv("ZIA_PASSWORD")
	apiKey := os.Getenv("ZIA_API_KEY")
	ziaCloud := os.Getenv("ZIA_CLOUD")

	// Check for the presence of necessary environment variables.
	if username == "" {
		return errors.New("ZIA_USERNAME must be set for acceptance tests")
	}

	if password == "" {
		return errors.New("ZIA_PASSWORD must be set for acceptance tests")
	}

	if apiKey == "" {
		return errors.New("ZIA_API_KEY must be set for acceptance tests")
	}

	if ziaCloud == "" {
		return errors.New("ZIA_CLOUD must be set for acceptance tests")
	}

	return nil
}

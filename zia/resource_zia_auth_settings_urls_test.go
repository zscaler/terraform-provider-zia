package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/zia/services/user_authentication_settings"
	"github.com/zscaler/terraform-provider-zia/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/zia/common/testing/method"
)

func TestAccResourceAuthSettingsUrlsBasic(t *testing.T) {
	var urls user_authentication_settings.ExemptedUrls
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AuthSettingsURLs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuthSettingsUrlsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAuthSettingsUrlsConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthSettingsUrlsExists(resourceTypeAndName, &urls),
					resource.TestCheckResourceAttr(resourceTypeAndName, "urls.#", "16"),
				),
			},

			// Update test
			{
				Config: testAccCheckAuthSettingsUrlsConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthSettingsUrlsExists(resourceTypeAndName, &urls),
					resource.TestCheckResourceAttr(resourceTypeAndName, "urls.#", "16"),
				),
			},
		},
	})
}

// Need to fix the destroy function : Error running post-test destroy, there may be dangling resources: url exempted_urls already exists
func testAccCheckAuthSettingsUrlsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.AuthSettingsURLs {
			continue
		}

		url, err := apiClient.user_authentication_settings.Get()

		if err == nil {
			return fmt.Errorf("url %s already exists", rs.Primary.ID)
		}

		if url != nil {
			return fmt.Errorf("url %s with id exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAuthSettingsUrlsExists(resource string, url *user_authentication_settings.ExemptedUrls) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedUrls, err := apiClient.user_authentication_settings.Get()

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*url = *receivedUrls

		return nil
	}
}

func testAccCheckAuthSettingsUrlsConfigure(resourceTypeAndName, generatedName string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	urls = [
	  ".okta.com",
	  ".oktacdn.com",
	  ".mtls.oktapreview.com",
	  ".mtls.okta.com",
	  "d3l44rcogcb7iv.cloudfront.net",
	  "pac.zdxcloud.net",
	  ".windowsazure.com",
	  ".fedoraproject.org",
	  "login.windows.net",
	  "d32a6ru7mhaq0c.cloudfront.net",
	  ".kerberos.oktapreview.com",
	  ".oktapreview.com",
	  "login.zdxcloud.net",
	  "login.microsoftonline.com",
	  "smres.zdxcloud.net",
	  ".kerberos.okta.com"
	]
  }

data "%s" "%s" {}
`,
		// resource variables
		resourcetype.AuthSettingsURLs,
		generatedName,

		// data source variables
		resourcetype.AuthSettingsURLs,
		generatedName,
	)
}
*/

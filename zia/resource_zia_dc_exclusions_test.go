package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/dc_exclusions"
)

func TestAccResourceDCExclusions(t *testing.T) {
	var dcExclusions dc_exclusions.DCExclusions
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DCExclusions)

	initialDescription := variable.DCExclusionsDescription + "-" + generatedName
	updatedDescription := variable.DCExclusionsDescription + "-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDCExclusionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCExclusionsConfigure(resourceTypeAndName, initialDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCExclusionsExists(resourceTypeAndName, &dcExclusions),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", initialDescription),
				),
			},

			// Update test
			{
				Config: testAccCheckDCExclusionsConfigure(resourceTypeAndName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCExclusionsExists(resourceTypeAndName, &dcExclusions),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", updatedDescription),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDCExclusionsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DCExclusions {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		all, err := dc_exclusions.GetAll(context.Background(), service)
		if err != nil {
			return err
		}
		for _, ex := range all {
			if ex.DcID == id {
				return fmt.Errorf("dc exclusion with datacenter id %d still exists", id)
			}
		}
	}

	return nil
}

func testAccCheckDCExclusionsExists(resource string, rule *dc_exclusions.DCExclusions) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		all, err := dc_exclusions.GetAll(context.Background(), service)
		if err != nil {
			return fmt.Errorf("failed fetching DC exclusions: %s", err)
		}
		var found *dc_exclusions.DCExclusions
		for i := range all {
			if all[i].DcID == id {
				found = &all[i]
				break
			}
		}
		if found == nil {
			return fmt.Errorf("dc exclusion with datacenter id %d not found", id)
		}
		*rule = *found

		return nil
	}
}

func testAccCheckDCExclusionsConfigure(resourceTypeAndName, description string) string {
	// Generate start_time_utc and end_time_utc using present date/time.
	// Start time must not be in the past; use current time (or very recent).
	// API limit: datacenter cannot be excluded for more than 2 weeks.
	now := time.Now().UTC()
	startTime := now.Format("01/02/2006 03:04:05 pm")
	endTime := now.Add(14 * 24 * time.Hour).Format("01/02/2006 03:04:05 pm") // 2 weeks from now

	resourceName := strings.Split(resourceTypeAndName, ".")[1]
	return fmt.Sprintf(`
data "zia_datacenters" "this" {
  name = "SJC4"
}

resource "%s" "%s" {
  datacenter_id  = data.zia_datacenters.this.datacenter_id
  start_time_utc = "%s"
  end_time_utc   = "%s"
  description    = "%s"
}

data "%s" "%s" {
  id = %s.%s.id
}
`,
		resourcetype.DCExclusions,
		resourceName,
		startTime,
		endTime,
		description,
		resourcetype.DCExclusions,
		resourceName,
		resourcetype.DCExclusions,
		resourceName,
	)
}

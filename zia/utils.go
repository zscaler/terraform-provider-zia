package zia

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
)

func intPtr(n int) *int {
	return &n
}
func SetToStringSlice(d *schema.Set) []string {
	list := d.List()
	return ListToStringSlice(list)
}

func SetToStringList(d *schema.ResourceData, key string) []string {
	setObj, ok := d.GetOk(key)
	if !ok {
		return []string{}
	}
	set, ok := setObj.(*schema.Set)
	if !ok {
		return []string{}
	}
	return SetToStringSlice(set)
}

func ListToStringSlice(v []interface{}) []string {
	if len(v) == 0 {
		return []string{}
	}

	ans := make([]string, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case nil:
			ans[i] = ""
		case string:
			ans[i] = x
		}
	}

	return ans
}

func getIntFromResourceData(d *schema.ResourceData, key string) (int, bool) {
	obj, isSet := d.GetOk(key)
	val, isInt := obj.(int)
	return val, isSet && isInt && val > 0
}

func getStringFromResourceData(d *schema.ResourceData, key string) (string, bool) {
	obj, isSet := d.GetOk(key)
	val, isStr := obj.(string)
	return val, isSet && isStr && val != ""
}

// avoid {"code":"RESOURCE_IN_USE","message":"GROUP is associated with 1 rule(s). Deletion of this group is not allowed."}
func DetachRuleIDNameExtensions(client *Client, id int, resource string, getResources func(*filteringrules.FirewallFilteringRules) []common.IDNameExtensions, setResources func(*filteringrules.FirewallFilteringRules, []common.IDNameExtensions)) error {
	log.Printf("[INFO] Detaching filtering rule from %s: %d\n", resource, id)
	rules, err := client.filteringrules.GetAll()
	if err != nil {
		log.Printf("[error] Error while getting filtering rule")
		return err
	}

	for _, rule := range rules {
		ids := []common.IDNameExtensions{}
		shouldUpdate := false
		for _, destGroup := range getResources(&rule) {
			if destGroup.ID != id {
				ids = append(ids, destGroup)
			} else {
				shouldUpdate = true
			}
		}
		if shouldUpdate {
			setResources(&rule, ids)
			time.Sleep(time.Second * 5)
			_, err = client.filteringrules.Get(rule.ID)
			if err == nil {
				_, err = client.filteringrules.Update(rule.ID, &rule)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func DetachDLPEngineIDNameExtensions(client *Client, id int, resource string, getResources func(*dlp_web_rules.WebDLPRules) []common.IDNameExtensions, setResources func(*dlp_web_rules.WebDLPRules, []common.IDNameExtensions)) error {
	log.Printf("[INFO] Detaching dlp engine from %s: %d\n", resource, id)
	rules, err := client.dlp_web_rules.GetAll()
	if err != nil {
		log.Printf("[error] Error while getting filtering rule")
		return err
	}

	for _, rule := range rules {
		ids := []common.IDNameExtensions{}
		shouldUpdate := false
		for _, dlpEngine := range getResources(&rule) {
			if dlpEngine.ID != id {
				ids = append(ids, dlpEngine)
			} else {
				shouldUpdate = true
			}
		}
		if shouldUpdate {
			setResources(&rule, ids)
			time.Sleep(time.Second * 5)
			_, err = client.dlp_web_rules.Get(rule.ID)
			if err == nil {
				_, err = client.dlp_web_rules.Update(rule.ID, &rule)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func ValidateLatitude(val interface{}, _ string) (warns []string, errs []error) {
	// Directly type assert to float64
	v, ok := val.(float64)
	if !ok {
		errs = append(errs, fmt.Errorf("expected latitude to be a float64"))
		return
	}
	if v < -90 || v > 90 {
		errs = append(errs, fmt.Errorf("latitude must be between -90 and 90"))
	}
	return
}

func ValidateLongitude(val interface{}, _ string) (warns []string, errs []error) {
	// Directly type assert to float64
	v, ok := val.(float64)
	if !ok {
		errs = append(errs, fmt.Errorf("expected longitude to be a float64"))
		return
	}
	if v < -180 || v > 180 {
		errs = append(errs, fmt.Errorf("longitude must be between -180 and 180"))
	}
	return
}

func DiffSuppressFuncCoordinate(_, old, new string, _ *schema.ResourceData) bool {
	o, err := strconv.ParseFloat(old, 64)
	if err != nil {
		return false
	}
	n, err := strconv.ParseFloat(new, 64)
	if err != nil {
		return false
	}
	return math.Round(o*1000000)/1000000 == math.Round(n*1000000)/1000000
}

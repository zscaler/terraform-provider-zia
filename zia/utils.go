package zia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
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

func SetToIntList(d *schema.ResourceData, key string) []int {
	setObj, ok := d.GetOk(key)
	if !ok {
		return []int{}
	}
	set, ok := setObj.(*schema.Set)
	if !ok {
		return []int{}
	}

	intList := make([]int, set.Len())
	for i, v := range set.List() {
		intList[i] = v.(int)
	}
	return intList
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
func DetachRuleIDNameExtensions(ctx context.Context, client *Client, id int, resource string, getResources func(*filteringrules.FirewallFilteringRules) []common.IDNameExtensions, setResources func(*filteringrules.FirewallFilteringRules, []common.IDNameExtensions)) error {
	service := client.Service

	log.Printf("[INFO] Detaching filtering rule from %s: %d\n", resource, id)
	rules, err := filteringrules.GetAll(ctx, service)
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
			_, err = filteringrules.Get(ctx, service, rule.ID)
			if err == nil {
				_, err = filteringrules.Update(ctx, service, rule.ID, &rule)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func DetachDLPEngineIDNameExtensions(ctx context.Context, client *Client, id int, resource string, getResources func(*dlp_web_rules.WebDLPRules) []common.IDNameExtensions, setResources func(*dlp_web_rules.WebDLPRules, []common.IDNameExtensions)) error {
	service := client.Service

	log.Printf("[INFO] Detaching dlp engine from %s: %d\n", resource, id)
	rules, err := dlp_web_rules.GetAll(ctx, service)
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
			_, err = dlp_web_rules.Get(ctx, service, rule.ID)
			if err == nil {
				_, err = dlp_web_rules.Update(ctx, service, rule.ID, &rule)
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

// createValidResourceName converts the given name to a valid Terraform resource name
func createValidResourceName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

// contains checks if a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// Helper function to trigger configuration activation
func triggerActivation(ctx context.Context, zClient *Client) error {
	service := zClient.Service

	// Assuming the activation request doesn't need specific details from the rule labels
	req := activation.Activation{Status: "ACTIVE"}
	log.Printf("[INFO] Triggering configuration activation\n%+v\n", req)

	_, err := activation.CreateActivation(ctx, service, req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Configuration activation triggered successfully.")
	return nil
}

// Helper function to check if we should activate based on the ZIA_ACTIVATION environment variable
func shouldActivate() bool {
	activationEnv, exists := os.LookupEnv("ZIA_ACTIVATION")
	if !exists {
		return false
	}
	activationBool, err := strconv.ParseBool(activationEnv)
	if err != nil {
		log.Printf("[WARN] Error parsing ZIA_ACTIVATION env var as bool: %v", err)
		return false
	}
	return activationBool
}

func validateTimeZone(v interface{}, k string) (ws []string, errors []error) {
	tzStr := v.(string)

	// Try IANA timezone validation
	_, err := time.LoadLocation(tzStr)
	if err != nil {
		// Check if it's a common timezone that should work
		commonTimezones := []string{"UTC", "GMT", "US/Pacific", "US/Eastern", "US/Central", "US/Mountain", "Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Vilnius", "Asia/Tokyo", "Asia/Shanghai"}

		for _, common := range commonTimezones {
			if tzStr == common {
				// This is a known valid timezone, but system might not have it
				errors = append(errors, fmt.Errorf("%q is a valid IANA timezone but your system's timezone database may be incomplete. Please ensure your system has the IANA timezone database installed. Alternatively, you can use Zscaler-specific timezone format like 'LITHUANIA_EUROPE_VILNIUS' for Lithuania", tzStr))
				return
			}
		}

		// For other timezones, provide a more helpful error message
		errors = append(errors, fmt.Errorf("%q is not a valid IANA timezone. Please use IANA format (e.g., 'Europe/Vilnius', 'US/Pacific', 'UTC', 'GMT'). Visit https://nodatime.org/TimeZones for the complete IANA timezone list. If you need to use Zscaler-specific timezone format, please use the 'LITHUANIA_EUROPE_VILNIUS' format instead", tzStr))
	}

	return
}

func ConvertRFC1123ToEpoch(timeStr string) (int, error) {
	t, err := time.Parse(time.RFC1123, timeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid time format: %v. Expected format: RFC1123 (Mon, 02 Jan 2006 15:04:05 MST)", err)
	}
	return int(t.Unix()), nil
}

func convertAndValidateSizeQuota(sizeQuotaMB int) (int, error) {
	const (
		minMB = 10
		maxMB = 100000
	)
	if sizeQuotaMB < minMB || sizeQuotaMB > maxMB {
		return 0, fmt.Errorf("size_quota must be between %d MB and %d MB", minMB, maxMB)
	}
	// Convert MB to KB
	sizeQuotaKB := sizeQuotaMB * 1024
	return sizeQuotaKB, nil
}

func isSingleDigitDay(timeStr string) bool {
	parts := strings.Split(timeStr, " ")
	if len(parts) < 2 {
		return false
	}

	day := parts[1]
	return len(day) == 1
}

// Global semaphore for controlling concurrent API requests. The size of this
// semaphore is configured during provider setup based on the `parallelism`
// configuration option.
var apiSemaphore chan struct{}

// Helper function to process countries
func processCountries(countries []string) []string {
	processedCountries := make([]string, len(countries))
	for i, country := range countries {
		if country != "ANY" && country != "NONE" && len(country) == 2 { // Assuming the 2 letter code is an ISO Alpha-2 Code
			processedCountries[i] = "COUNTRY_" + country
		} else {
			processedCountries[i] = country
		}
	}
	return processedCountries
}

// func handleInvalidInputError(err error) error {
// 	if err == nil {
// 		return nil
// 	}

// 	if strings.Contains(err.Error(), `"code":"INVALID_INPUT_ARGUMENT"`) {
// 		log.Printf("[ERROR] Failing immediately due to INVALID_INPUT_ARGUMENT: %s", err.Error())
// 		return err
// 	}

// 	return nil
// }

var failFastErrorCodes = []string{
	"INVALID_INPUT_ARGUMENT",
	"TRIAL_EXPIRED",
	"EDIT_LOCK_NOT_AVAILABLE",
	"DUPLICATE_ITEM",
	// Add more codes here as needed
}

// failFastOnErrorCodes detects known fatal API error codes and returns the original error to fail immediately.
func failFastOnErrorCodes(err error) error {
	if err == nil {
		return nil
	}

	// Case 1: SDK's structured ErrorResponse (preferred path)
	var apiErr *errorx.ErrorResponse
	if errors.As(err, &apiErr) {
		code := extractErrorCodeFromBody(apiErr.Message)
		for _, c := range failFastErrorCodes {
			if code == c {
				log.Printf("[ERROR] Failing immediately due to API error code '%s': %s", c, apiErr.Message)
				return err
			}
		}
	}

	// Case 2: fallback for unstructured errors
	errMsg := err.Error()
	for _, code := range failFastErrorCodes {
		match := fmt.Sprintf(`"code":"%s"`, code)
		if strings.Contains(errMsg, match) {
			log.Printf("[WARN] Failing due to fallback match for code '%s': %s", code, errMsg)
			return err
		}
	}

	return nil
}

func extractErrorCodeFromBody(body string) string {
	type apiErrorBody struct {
		Code string `json:"code"`
	}
	var parsed apiErrorBody
	if err := json.Unmarshal([]byte(body), &parsed); err == nil {
		return parsed.Code
	}
	return ""
}

func suppressEquivalentIntListOrdering(k, old, new string, d *schema.ResourceData) bool {
	oldList := strings.Split(strings.Trim(old, "[] "), ",")
	newList := strings.Split(strings.Trim(new, "[] "), ",")

	if len(oldList) != len(newList) {
		return false
	}

	sort.Strings(oldList)
	sort.Strings(newList)

	for i := range oldList {
		if strings.TrimSpace(oldList[i]) != strings.TrimSpace(newList[i]) {
			return false
		}
	}
	return true
}

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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewalldnscontrolpolicies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ftp_control_policy"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ips_control_policies/ips_policies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sslinspection"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
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

func ListToStringList(d *schema.ResourceData, key string) []string {
	listObj, ok := d.GetOk(key)
	if !ok {
		return []string{}
	}
	list, ok := listObj.([]interface{})
	if !ok {
		return []string{}
	}
	return ListToStringSlice(list)
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

func getBoolFromResourceData(d *schema.ResourceData, key string) bool {
	obj, isSet := d.GetOk(key)
	if !isSet {
		return false
	}
	val, isBool := obj.(bool)
	return isBool && val
}

// avoid {"code":"RESOURCE_IN_USE","message":"GROUP is associated with 1 rule(s). Deletion of this group is not allowed."}
func DetachRuleIDNameExtensions(ctx context.Context, client *Client, id int, resource string, getResources func(*filteringrules.FirewallFilteringRules) []common.IDNameExtensions, setResources func(*filteringrules.FirewallFilteringRules, []common.IDNameExtensions)) error {
	service := client.Service

	log.Printf("[INFO] Detaching filtering rule from %s: %d\n", resource, id)
	rules, err := filteringrules.GetAll(ctx, service, nil)
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

// DetachURLFilteringRuleRef removes a referenced object (by ID) from a
// []common.IDNameExtensions field on every URL filtering rule that references
// it, so the object can be deleted without the API returning RESOURCE_IN_USE.
// The getResources/setResources closures select which field is swept, allowing
// the same routine to serve any url_filtering association (e.g. HTTP header
// profiles and HTTP header action profiles).
func DetachURLFilteringRuleRef(ctx context.Context, client *Client, id int, resource string, getResources func(*urlfilteringpolicies.URLFilteringRule) []common.IDNameExtensions, setResources func(*urlfilteringpolicies.URLFilteringRule, []common.IDNameExtensions)) error {
	service := client.Service

	log.Printf("[INFO] Detaching %s reference from URL filtering rules: %d\n", resource, id)
	rules, err := urlfilteringpolicies.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] Error while getting URL filtering rules for detach: %v", err)
		return err
	}

	updated := false
	for i := range rules {
		rule := &rules[i]
		ids := []common.IDNameExtensions{}
		shouldUpdate := false
		for _, ref := range getResources(rule) {
			if ref.ID != id {
				ids = append(ids, ref)
			} else {
				shouldUpdate = true
			}
		}
		if !shouldUpdate {
			continue
		}
		setResources(rule, ids)
		if _, err := urlfilteringpolicies.Update(ctx, service, rule.ID, rule); err != nil {
			return fmt.Errorf("failed detaching %s %d from URL filtering rule %d: %w", resource, id, rule.ID, err)
		}
		updated = true
	}

	// Give the API a moment to settle the associations before the caller issues
	// the delete. A single settle is sufficient regardless of how many rules
	// were updated.
	if updated {
		time.Sleep(time.Second * 5)
	}
	return nil
}

// removeStringFromSlice returns a copy of in with every occurrence of target
// removed, plus whether any element was removed.
func removeStringFromSlice(in []string, target string) ([]string, bool) {
	out := make([]string, 0, len(in))
	removed := false
	for _, v := range in {
		if v == target {
			removed = true
			continue
		}
		out = append(out, v)
	}
	return out, removed
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

// detachURLCategoryFromAllResources removes a URL category (matched by its
// string ID, e.g. "CUSTOM_08") from every rule-based resource that can
// reference it, so the category can be deleted without the API returning
// RESOURCE_IN_USE.
//
// Each rule type is checked concurrently. A rule is only updated (a write) when
// the category is actually present on it; rule types where the category is not
// attached incur a single read (GetAll) and no writes. A single settle delay is
// applied at the end if any detach write occurred.
func detachURLCategoryFromAllResources(ctx context.Context, client *Client, categoryID string) error {
	detachers := []func(context.Context, *Client, string) (bool, error){
		detachCategoryFromURLFilteringRules,
		detachCategoryFromSSLInspectionRules,
		detachCategoryFromSandboxRules,
		detachCategoryFromFileTypeControlRules,
		detachCategoryFromBandwidthClasses,
		detachCategoryFromFirewallDNSRules,
		detachCategoryFromFirewallIPSRules,
		detachCategoryFromFTPControlPolicy,
		detachCategoryFromDLPWebRules,
	}

	type result struct {
		updated bool
		err     error
	}
	results := make(chan result, len(detachers))
	for _, fn := range detachers {
		go func(f func(context.Context, *Client, string) (bool, error)) {
			updated, err := f(ctx, client, categoryID)
			results <- result{updated: updated, err: err}
		}(fn)
	}

	anyUpdated := false
	var firstErr error
	for range detachers {
		r := <-results
		if r.err != nil && firstErr == nil {
			firstErr = r.err
		}
		if r.updated {
			anyUpdated = true
		}
	}
	if firstErr != nil {
		return firstErr
	}

	// Give the API a moment to settle the associations before the caller issues
	// the delete. A single settle is sufficient regardless of how many rules
	// across how many resource types were updated.
	if anyUpdated {
		time.Sleep(time.Second * 5)
	}
	return nil
}

func detachCategoryFromURLFilteringRules(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	rules, err := urlfilteringpolicies.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting URL filtering rules for category detach: %w", err)
	}
	updated := false
	for i := range rules {
		rule := &rules[i]
		c1, r1 := removeStringFromSlice(rule.URLCategories, categoryID)
		c2, r2 := removeStringFromSlice(rule.URLCategories2, categoryID)
		if !r1 && !r2 {
			continue
		}
		rule.URLCategories = c1
		rule.URLCategories2 = c2
		log.Printf("[INFO] Detaching URL category %q from URL filtering rule %d\n", categoryID, rule.ID)
		if _, err := urlfilteringpolicies.Update(ctx, service, rule.ID, rule); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from URL filtering rule %d: %w", categoryID, rule.ID, err)
		}
		updated = true
	}
	return updated, nil
}

func detachCategoryFromSSLInspectionRules(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	rules, err := sslinspection.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting SSL inspection rules for category detach: %w", err)
	}
	updated := false
	for i := range rules {
		rule := &rules[i]
		c, removed := removeStringFromSlice(rule.URLCategories, categoryID)
		if !removed {
			continue
		}
		rule.URLCategories = c
		log.Printf("[INFO] Detaching URL category %q from SSL inspection rule %d\n", categoryID, rule.ID)
		if _, err := sslinspection.Update(ctx, service, rule.ID, rule); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from SSL inspection rule %d: %w", categoryID, rule.ID, err)
		}
		updated = true
	}
	return updated, nil
}

func detachCategoryFromSandboxRules(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	rules, err := sandbox_rules.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting sandbox rules for category detach: %w", err)
	}
	updated := false
	for i := range rules {
		rule := &rules[i]
		c, removed := removeStringFromSlice(rule.URLCategories, categoryID)
		if !removed {
			continue
		}
		rule.URLCategories = c
		log.Printf("[INFO] Detaching URL category %q from sandbox rule %d\n", categoryID, rule.ID)
		if _, err := sandbox_rules.Update(ctx, service, rule.ID, rule); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from sandbox rule %d: %w", categoryID, rule.ID, err)
		}
		updated = true
	}
	return updated, nil
}

func detachCategoryFromFileTypeControlRules(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	rules, err := filetypecontrol.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting file type control rules for category detach: %w", err)
	}
	updated := false
	for i := range rules {
		rule := &rules[i]
		c, removed := removeStringFromSlice(rule.URLCategories, categoryID)
		if !removed {
			continue
		}
		rule.URLCategories = c
		log.Printf("[INFO] Detaching URL category %q from file type control rule %d\n", categoryID, rule.ID)
		if _, err := filetypecontrol.Update(ctx, service, rule.ID, rule); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from file type control rule %d: %w", categoryID, rule.ID, err)
		}
		updated = true
	}
	return updated, nil
}

func detachCategoryFromBandwidthClasses(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	classes, err := bandwidth_classes.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting bandwidth classes for category detach: %w", err)
	}
	updated := false
	for i := range classes {
		class := &classes[i]
		c, removed := removeStringFromSlice(class.UrlCategories, categoryID)
		if !removed {
			continue
		}
		class.UrlCategories = c
		log.Printf("[INFO] Detaching URL category %q from bandwidth class %d\n", categoryID, class.ID)
		if _, _, err := bandwidth_classes.Update(ctx, service, class.ID, class); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from bandwidth class %d: %w", categoryID, class.ID, err)
		}
		updated = true
	}
	return updated, nil
}

// detachCategoryFromFirewallDNSRules sweeps both the request (DestIpCategories)
// and response (ResCategories) category fields. The API requires these two to
// stay mirrored, so a detach must strip the ID from both.
func detachCategoryFromFirewallDNSRules(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	rules, err := firewalldnscontrolpolicies.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting firewall DNS rules for category detach: %w", err)
	}
	updated := false
	for i := range rules {
		rule := &rules[i]
		dest, r1 := removeStringFromSlice(rule.DestIpCategories, categoryID)
		res, r2 := removeStringFromSlice(rule.ResCategories, categoryID)
		if !r1 && !r2 {
			continue
		}
		rule.DestIpCategories = dest
		rule.ResCategories = res
		log.Printf("[INFO] Detaching URL category %q from firewall DNS rule %d\n", categoryID, rule.ID)
		if _, err := firewalldnscontrolpolicies.Update(ctx, service, rule.ID, rule); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from firewall DNS rule %d: %w", categoryID, rule.ID, err)
		}
		updated = true
	}
	return updated, nil
}

// detachCategoryFromFirewallIPSRules sweeps both the request (DestIpCategories)
// and response (ResCategories) category fields, which the API requires to stay
// mirrored.
func detachCategoryFromFirewallIPSRules(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	rules, err := ips_policies.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting firewall IPS rules for category detach: %w", err)
	}
	updated := false
	for i := range rules {
		rule := &rules[i]
		dest, r1 := removeStringFromSlice(rule.DestIpCategories, categoryID)
		res, r2 := removeStringFromSlice(rule.ResCategories, categoryID)
		if !r1 && !r2 {
			continue
		}
		rule.DestIpCategories = dest
		rule.ResCategories = res
		log.Printf("[INFO] Detaching URL category %q from firewall IPS rule %d\n", categoryID, rule.ID)
		if _, err := ips_policies.Update(ctx, service, rule.ID, rule); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from firewall IPS rule %d: %w", categoryID, rule.ID, err)
		}
		updated = true
	}
	return updated, nil
}

// detachCategoryFromFTPControlPolicy handles the FTP control policy, which is a
// singleton (no per-rule list); it is read once and updated only if it carries
// the category.
func detachCategoryFromFTPControlPolicy(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	policy, err := ftp_control_policy.GetFTPControlPolicy(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting FTP control policy for category detach: %w", err)
	}
	if policy == nil {
		return false, nil
	}
	c, removed := removeStringFromSlice(policy.UrlCategories, categoryID)
	if !removed {
		return false, nil
	}
	policy.UrlCategories = c
	log.Printf("[INFO] Detaching URL category %q from FTP control policy\n", categoryID)
	if _, _, err := ftp_control_policy.UpdateFTPControlPolicy(ctx, service, policy); err != nil {
		return false, fmt.Errorf("failed detaching URL category %q from FTP control policy: %w", categoryID, err)
	}
	return true, nil
}

// detachCategoryFromDLPWebRules handles DLP web rules, where the URL category is
// stored as []common.IDNameExtensions (nested objects) rather than a string
// list. The category's string ID is compared against each entry's numeric ID
// rendered as a string.
func detachCategoryFromDLPWebRules(ctx context.Context, client *Client, categoryID string) (bool, error) {
	service := client.Service
	rules, err := dlp_web_rules.GetAll(ctx, service)
	if err != nil {
		return false, fmt.Errorf("error getting DLP web rules for category detach: %w", err)
	}
	updated := false
	for i := range rules {
		rule := &rules[i]
		kept := make([]common.IDNameExtensions, 0, len(rule.URLCategories))
		removed := false
		for _, cat := range rule.URLCategories {
			if fmt.Sprintf("%d", cat.ID) == categoryID || cat.Name == categoryID {
				removed = true
				continue
			}
			kept = append(kept, cat)
		}
		if !removed {
			continue
		}
		rule.URLCategories = kept
		log.Printf("[INFO] Detaching URL category %q from DLP web rule %d\n", categoryID, rule.ID)
		if _, err := dlp_web_rules.Update(ctx, service, rule.ID, rule); err != nil {
			return updated, fmt.Errorf("failed detaching URL category %q from DLP web rule %d: %w", categoryID, rule.ID, err)
		}
		updated = true
	}
	return updated, nil
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

// ExclusionTimeUTCLayout is the layout for DC exclusion / subcloud exclusion times (UTC): "MM/DD/YYYY HH:MM:SS am/pm"
const ExclusionTimeUTCLayout = "01/02/2006 03:04:05 pm"

// ParseExclusionTimeUTC parses a string in "MM/DD/YYYY HH:MM:SS am/pm" (UTC) and returns Unix epoch seconds.
// Seconds are required. Used by both zia_dc_exclusions and zia_sub_cloud for exclusion time conversion.
func ParseExclusionTimeUTC(timeStr string) (int, error) {
	s := strings.TrimSpace(timeStr)
	if s == "" {
		return 0, fmt.Errorf("exclusion time string is empty")
	}
	t, err := time.ParseInLocation(ExclusionTimeUTCLayout, s, time.UTC)
	if err != nil {
		return 0, fmt.Errorf("invalid time format: %v. Expected format: MM/DD/YYYY HH:MM:SS am/pm (e.g. 02/19/2026 11:59:00 pm)", err)
	}
	return int(t.Unix()), nil
}

// FormatExclusionTimeUTC formats Unix epoch seconds as "MM/DD/YYYY HH:MM:SS am/pm" in UTC.
// Used by both zia_dc_exclusions and zia_sub_cloud when writing exclusion times to state.
func FormatExclusionTimeUTC(epoch int) string {
	return time.Unix(int64(epoch), 0).UTC().Format(ExclusionTimeUTCLayout)
}

// ResolveExclusionTimeEpoch returns Unix epoch from either the UTC string (MM/DD/YYYY HH:MM:SS am/pm) or the epoch int.
// If utcStr is non-empty it is parsed and overrides epochVal. Otherwise epochVal is used (hasEpoch must be true).
// Used by both zia_dc_exclusions and zia_sub_cloud so exclusion time handling is consistent.
func ResolveExclusionTimeEpoch(hasEpoch bool, epochVal int, utcStr string) (int, error) {
	s := strings.TrimSpace(utcStr)
	if s != "" {
		return ParseExclusionTimeUTC(s)
	}
	if hasEpoch {
		return epochVal, nil
	}
	return 0, fmt.Errorf("either epoch or UTC string must be set")
}

// ValidateExclusionTimeUTC validates that the value is in "MM/DD/YYYY HH:MM:SS am/pm" (UTC) format. Empty string is allowed (optional field).
func ValidateExclusionTimeUTC(v interface{}, k string) (warnings []string, errors []error) {
	s, ok := v.(string)
	if !ok || s == "" {
		return nil, nil
	}
	if _, err := ParseExclusionTimeUTC(s); err != nil {
		errors = append(errors, fmt.Errorf("%s: %w", k, err))
	}
	return nil, errors
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

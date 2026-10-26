package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

// orderedRuleSpec describes one rule resource that uses the shared reorder
// engine, for the tests below that run the same scenario on every one of them.
type orderedRuleSpec struct {
	name       string
	resource   func() *schema.Resource
	update     crudFunc
	listPath   string // GET: the rule list ("" for sub-rules)
	itemPath   string // GET/PUT/DELETE of one rule, ending in "/*"
	parentPath string // sub-rules only: GET returns the parent with its subRules
	idAttr     string
	state      map[string]any // extra state attributes (rule type, parent rule, ...)
	fields     map[string]any // extra fields of every rule served by the API
}

const testCloudAppType = "STREAMING_MEDIA"

func orderedRuleSpecs() []orderedRuleSpec {
	zia := func(name string) (string, string) { return common.ZIAPath(name), common.ZIAPath(name, "*") }
	spec := func(name string, r func() *schema.Resource, update crudFunc, api string) orderedRuleSpec {
		list, item := zia(api)
		return orderedRuleSpec{name: name, resource: r, update: update, listPath: list, itemPath: item, idAttr: "rule_id"}
	}
	withType := func(s orderedRuleSpec, ruleType string) orderedRuleSpec {
		s.state = map[string]any{"type": ruleType}
		s.fields = map[string]any{"type": ruleType}
		return s
	}
	cloudApp := spec("cloud_app_control_rules", resourceCloudAppControlRules, resourceCloudAppControlRulesUpdate, "webApplicationRules")
	cloudApp.listPath = common.ZIAPath("webApplicationRules", testCloudAppType)
	cloudApp.itemPath = common.ZIAPath("webApplicationRules", testCloudAppType, "*")
	forwarding := spec("forwarding_control_rule", resourceForwardingControlRule, resourceForwardingControlRuleUpdate, "forwardingRules")
	forwarding.state = map[string]any{"forward_method": "DIRECT"}
	ssl := spec("ssl_inspection_rules", resourceSSLInspectionRules, resourceSSLInspectionRulesUpdate, "sslInspectionRules")
	ssl.state = map[string]any{"action": []any{map[string]any{"type": "BLOCK"}}}

	return []orderedRuleSpec{
		spec("bandwidth_control_rules", resourceBandwdithControlRules, resourceBandwdithControlRulesUpdate, "bandwidthControlRules"),
		withType(spec("casb_dlp_rules", resourceCasbDlpRules, resourceCasbDlpRulesUpdate, "casbDlpRules"), "OFLCASB_DLP_ITSM"),
		withType(spec("casb_malware_rules", resourceCasbMalwareRules, resourceCasbMalwareRulesUpdate, "casbMalwareRules"), "OFLCASB_AVP_ITSM"),
		withType(cloudApp, testCloudAppType),
		spec("dlp_web_rules", resourceDlpWebRules, resourceDlpWebRulesUpdate, "webDlpRules"),
		spec("endpoint_dlp_rules", resourceEndpointDLPRules, resourceEndpointDLPRulesUpdate, "endPointDlpRules"),
		{
			name: "endpoint_dlp_sub_rules", resource: resourceEndpointDLPSubRules, update: resourceEndpointDLPSubRulesUpdate,
			parentPath: common.ZIAPath("endPointDlpRules", "500"), itemPath: common.ZIAPath("endPointDlpRules", "500", "subRule", "*"),
			idAttr: "sub_rule_id", state: map[string]any{"parent_rule": 500},
		},
		spec("file_type_control_rules", resourceFileTypeControlRules, resourceFileTypeControlRulesUpdate, "fileTypeRules"),
		spec("firewall_dns_rules", resourceFirewallDNSRules, resourceFirewallDNSRulesUpdate, "firewallDnsRules"),
		spec("firewall_filtering_rules", resourceFirewallFilteringRules, resourceFirewallFilteringRulesUpdate, "firewallFilteringRules"),
		spec("firewall_ips_rules", resourceFirewallIPSRules, resourceFirewallIPSRulesUpdate, "firewallIpsRules"),
		forwarding,
		spec("nat_control_rules", resourceNatControlRules, resourceNatControlRulesUpdate, "dnatRules"),
		spec("outbound_email_dlp", resourceOutboundEmailDLP, resourceOutboundEmailDLPUpdate, "emailDlpRules"),
		spec("sandbox_rules", resourceSandboxRules, resourceSandboxRulesUpdate, "sandboxRules"),
		ssl,
		spec("traffic_capture_rules", resourceTrafficCaptureRules, resourceFiresourceTrafficCaptureRulesUpdate, "trafficCaptureRules"),
		spec("url_filtering_rules", resourceURLFilteringRules, resourceURLFilteringRulesUpdate, "urlFilteringRules"),
	}
}

// updateOf builds the ResourceData of an update of rule 1000+i, whose state has
// order i, planned (as Terraform does) from a configuration that change edits.
func (s orderedRuleSpec) updateOf(f *unitFixture, i int, change func(config map[string]any)) *schema.ResourceData {
	f.t.Helper()
	r := s.resource()
	config := map[string]any{"name": ruleName(i), "order": i}
	if _, ok := r.Schema["rank"]; ok {
		config["rank"] = 7
	}
	for k, v := range s.state {
		config[k] = v
	}
	state := map[string]any{s.idAttr: 1000 + i}
	for k, v := range config {
		state[k] = v
	}
	change(config)
	return f.updateData(r, strconv.Itoa(1000+i), state, config)
}

// A rule below a deleted rule is renumbered (order 9 -> 8) after the delete has
// already closed the gap: the update only changes the order, so it sends no
// write and the rule stays where it is.
func TestOrderedRules_OrderOnlyUpdateSendsNoWrite(t *testing.T) {
	for _, s := range orderedRuleSpecs() {
		s := s
		t.Run(s.name, func(t *testing.T) {
			f := newUnitFixture(t)
			fake := newFakeRuleList(f, 20, s.listPath, s.itemPath, s.parentPath, s.fields)
			fake.deleteRule(1008)

			d := s.updateOf(f, 9, func(c map[string]any) { c["order"] = 8 })
			if e := f.run(s.update, d); e != "" {
				t.Fatalf("unexpected error: %s", e)
			}
			if n := fake.putCount(1009); n != 0 {
				t.Errorf("sent %d PUTs for an order-only update, want 0", n)
			}
			if got := fake.orderOf(1009); got != 8 {
				t.Errorf("order %d, want 8", got)
			}
		})
	}
}

// An order past the current end of the list cannot be placed until later
// creates extend the list, so the engine releases the batch and the update
// returns with the rule still at its old position. The state it returns must
// nevertheless carry the declared order: that state is what Terraform
// persists, nothing rewrites it once a later cycle moves the rule, and the
// next plan compares the API against it.
func TestOrderedRules_DeferredPlacementRecordsDeclaredOrder(t *testing.T) {
	for _, s := range orderedRuleSpecs() {
		s := s
		t.Run(s.name, func(t *testing.T) {
			f := newUnitFixture(t)
			fake := newFakeRuleList(f, 20, s.listPath, s.itemPath, s.parentPath, s.fields)

			d := s.updateOf(f, 9, func(c map[string]any) { c["order"] = 21 })
			if e := f.run(s.update, d); e != "" {
				t.Fatalf("unexpected error: %s", e)
			}
			if got := d.Get("order").(int); got != 21 {
				t.Errorf("state order %d after a deferred placement, want declared 21", got)
			}
			if n := fake.putCount(1009); n != 0 {
				t.Errorf("sent %d PUTs for a position that does not exist yet, want 0", n)
			}
			if got := fake.orderOf(1009); got != 9 {
				t.Errorf("upstream order %d, want 9 until the list grows", got)
			}
		})
	}
}

// An update of another attribute writes the rule once, at its current position.
func TestOrderedRules_InPlaceUpdateKeepsPosition(t *testing.T) {
	for _, s := range orderedRuleSpecs() {
		s := s
		t.Run(s.name, func(t *testing.T) {
			f := newUnitFixture(t)
			fake := newFakeRuleList(f, 20, s.listPath, s.itemPath, s.parentPath, s.fields)

			d := s.updateOf(f, 9, func(c map[string]any) { c["name"] = "Renamed" })
			if e := f.run(s.update, d); e != "" {
				t.Fatalf("unexpected error: %s", e)
			}
			if n := fake.putCount(1009); n != 1 {
				t.Fatalf("sent %d PUTs, want 1", n)
			}
			if got := fake.puts[0]; got["name"] != "Renamed" || asInt(got["order"]) != 9 {
				t.Errorf("PUT name %v order %v, want \"Renamed\" at order 9", got["name"], got["order"])
			}
			if fake.rejected != 0 {
				t.Errorf("%d writes rejected", fake.rejected)
			}
		})
	}
}

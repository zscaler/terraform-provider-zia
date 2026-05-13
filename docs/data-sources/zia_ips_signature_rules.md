---
subcategory: "IPS Control Policy"
layout: "zscaler"
page_title: "ZIA: ips_signature_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-custom-ips-signature-rules
  API documentation https://help.zscaler.com/legacy-apis/ips-control-policy#/ipsSignatureRules-get
  Retrieves a custom IPS signature rule.
---

# zia_ips_signature_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-custom-ips-signature-rules)
* [API documentation](https://help.zscaler.com/legacy-apis/ips-control-policy#/ipsSignatureRules-get)

Use the **zia_ips_signature_rules** data source to retrieve a custom IPS signature rule by `id` or `name` from the Zscaler Internet Access cloud. The returned attributes can then be referenced from other resources such as IPS Control rules.

## Example Usage

### Look up by name

```hcl
data "zia_ips_signature_rules" "example" {
  name = "Detect_SSH_Brute_Force"
}
```

### Look up by ID

```hcl
data "zia_ips_signature_rules" "example" {
  id = 1024
}
```

## Argument Reference

The following arguments are supported. Exactly one of `id` or `name` must be provided.

* `id` - (Optional) System-generated identifier for the custom IPS signature rule.
* `name` - (Optional) Custom IPS signature rule name. Matching is case-insensitive.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `rule_text` - (String) The rule text in Suricata/Snort syntax that defines the custom IPS signature.
* `description` - (String) Additional information about the custom signature rule.
* `enabled` - (Bool) Whether the custom signature rule is enabled and ready to be used in IPS Control rules via the assigned threat category.
* `category` - (Block) Threat category assigned to the custom signature rule.
  * `id` - (Int) Unique identifier of the threat category.
  * `name` - (String) Name of the threat category (e.g. `ADVANCED_SECURITY`).
  * `is_name_l10n_tag` - (Bool) Indicates whether the name is a localization tag rather than a literal label.
* `deleted` - (Bool) Whether the custom signature rule is marked as deleted.
* `promote_time` - (Int) Unix timestamp (in seconds) when the rule was promoted. `0` if not yet promoted.
* `rule_text_mod_time` - (Int) Unix timestamp (in seconds) when the rule text was last modified.
* `dynamic_validation_submitted` - (Bool) Whether the rule was submitted for dynamic validation.
* `dynamic_validation_rejected` - (Bool) Whether dynamic validation rejected the rule.
* `dynamic_validation_succeeded` - (Bool) Whether dynamic validation succeeded for the rule.
* `disabled_from_zscm` - (Bool) Whether the rule was disabled from Zscaler Cloud Management.
* `dynamic_val_reject_code` - (Int) Reject code returned by dynamic validation.

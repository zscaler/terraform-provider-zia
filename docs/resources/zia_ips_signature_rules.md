---
subcategory: "IPS Control Policy"
layout: "zscaler"
page_title: "ZIA: ips_signature_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-custom-ips-signature-rules
  API documentation https://help.zscaler.com/legacy-apis/ips-control-policy#/ipsSignatureRules-post
  Creates and manages ZIA custom IPS signature rules.
---

# zia_ips_signature_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/configuring-custom-ips-signature-rules)
* [API documentation](https://help.zscaler.com/legacy-apis/ips-control-policy#/ipsSignatureRules-post)

The **zia_ips_signature_rules** resource creates and manages custom IPS signature rules in the Zscaler Internet Access cloud. Each signature is authored in Suricata/Snort syntax, assigned to a threat category, and — once enabled — picked up by IPS Control rules that reference the same category.

~> **NOTE:** The provider validates `rule_text` against the Zscaler validation endpoint on every create and update. If the rule text is rejected, the create/update aborts before any state is written and the underlying API diagnostic (position, parameter, suggestion) is surfaced in the error.

~> **NOTE:** `rule_text` is whitespace-insensitive across plans. The provider normalizes leading/trailing whitespace and trailing newlines on both sides of the diff, so the same logical rule expressed as a double-quoted string, an indent-stripping heredoc (`<<-EOT`), or a plain heredoc (`<<EOT`) will not produce a perpetual `~ update in-place`.

## Example Usage

`rule_text` typically contains double quotes (`msg:"..."`, `content:"..."`). HCL gives you three equivalent ways to write it; pick the one that's easiest to read in your codebase.

### Double-quoted string — escape inner `"` with `\"`

```hcl
resource "zia_ips_signature_rules" "http_admin_quoted" {
  name        = "Detect_HTTP_Admin_URI_Quoted"
  description = "Alert when /admin is requested over HTTP"
  enabled     = true

  rule_text = "alert http any any -> any any (msg:\"HTTP /admin access\"; content:\"/admin\"; http_uri; nocase; sid:1000010; rev:1;)"

  category {
    id = 64
  }
}
```

### Indent-stripping heredoc (`<<-EOT`) — no escapes, indent must match the closing marker

```hcl
resource "zia_ips_signature_rules" "http_admin_indented_heredoc" {
  name        = "Detect_HTTP_Admin_URI_Indented"
  description = "Alert when /admin is requested over HTTP"
  enabled     = true

  rule_text = <<-EOT
    alert http any any -> any any (msg:"HTTP /admin access"; content:"/admin"; http_uri; nocase; sid:1000010; rev:1;)
  EOT

  category {
    id = 64
  }
}
```

`<<-EOT` strips the minimum common leading whitespace across **all** lines including the closing marker. If the rule line is indented more than `EOT`, the leftover spaces become part of the string and Zscaler's validation endpoint rejects the rule with `Unexpected character(s) at rule header`.

### Plain heredoc (`<<EOT`) — no escapes, no indent stripping

```hcl
resource "zia_ips_signature_rules" "http_admin_plain_heredoc" {
  name        = "Detect_HTTP_Admin_URI_Plain"
  description = "Alert when /admin is requested over HTTP"
  enabled     = true

  rule_text = <<EOT
alert http any any -> any any (msg:"HTTP /admin access"; content:"/admin"; http_uri; nocase; sid:1000010; rev:1;)
EOT

  category {
    id = 64
  }
}
```

The rule body lives at column 0 verbatim. This form has no escape noise and no indent surprises, and is the recommended choice for configurations that contain many signatures.

### Authoring tips

The Zscaler IPS validation endpoint enforces Suricata-style syntax with a few tenant-specific constraints:

* Use literal `any` (or a literal CIDR) for source / destination networks. Suricata variables such as `$EXTERNAL_NET` and `$HOME_NET` are not defined and trigger `Unexpected character(s) at rule header`.
* The custom-signature `sid` range starts at `1000000`. Lower values collide with reserved Talos / ETS ranges and are rejected.
* `msg:"..."` and `sid:N` are required inside the parentheses.
* Direction must be `->`. The bidirectional `<>` operator is not supported.
* No leading whitespace before `alert`. See the heredoc note above for how to avoid the indent trap.

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) Custom IPS signature rule name. Up to 255 characters.
* `rule_text` - (String) The rule text in Suricata/Snort syntax that defines the custom IPS signature. Validated against the Zscaler validation endpoint before every create and update.
* `category` - (Block) Threat category assigned to the custom signature rule. Exactly one block.
  * `id` - (Int, Required) Unique identifier of the threat category.

### Optional

* `description` - (String) Additional information about the custom signature rule.
* `enabled` - (Bool) Whether the signature rule is enabled and ready to be used in IPS Control rules via the assigned threat category. Defaults to `true`.
* `category.name` - (String) Name of the threat category (e.g. `ADVANCED_SECURITY`). The cloud returns the canonical value on read.
* `category.is_name_l10n_tag` - (Bool) Indicates whether the name is a localization tag rather than a literal label. The cloud returns the canonical value on read.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) Terraform resource identifier. Mirrors the API-assigned numeric id of the signature rule.
* `signature_id` - (Int) System-generated identifier for the custom IPS signature rule (same value as `id`, exposed as an integer for use in other resources).
* `deleted` - (Bool) Whether the signature rule is marked as deleted.
* `promote_time` - (Int) Unix timestamp (in seconds) when the rule was promoted. `0` if not yet promoted.
* `rule_text_mod_time` - (Int) Unix timestamp (in seconds) when the rule text was last modified.
* `dynamic_validation_submitted` - (Bool) Whether the rule was submitted for dynamic validation.
* `dynamic_validation_rejected` - (Bool) Whether dynamic validation rejected the rule.
* `dynamic_validation_succeeded` - (Bool) Whether dynamic validation succeeded for the rule.
* `disabled_from_zscm` - (Bool) Whether the rule was disabled from Zscaler Cloud Management.
* `dynamic_val_reject_code` - (Int) Reject code returned by dynamic validation.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_ips_signature_rules** can be imported by using `<SIGNATURE_ID>` or `<SIGNATURE_NAME>` as the import ID.

For example:

```shell
terraform import zia_ips_signature_rules.example <signature_id>
```

or

```shell
terraform import zia_ips_signature_rules.example <signature_name>
```

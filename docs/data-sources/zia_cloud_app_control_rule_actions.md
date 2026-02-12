---
subcategory: "Cloud App Control Policy"
layout: "zscaler"
page_title: "ZIA: cloud_app_control_rule_actions"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-rules-cloud-app-control-policy
  API documentation https://help.zscaler.com/zia/cloud-app-control-policy#/webApplicationRules/{rule_type}-get
  Get information about ZIA Cloud App Control Rules.
---

# zia_cloud_app_control_rule_actions (Data Source)

* [Official documentation](https://help.zscaler.com/zia/adding-rules-cloud-app-control-policy)
* [API documentation](https://help.zscaler.com/zia/cloud-app-control-policy#/webApplicationRules/)

Use the **zia_cloud_app_control_rule_actions** data source to retrieve the available actions for specific cloud applications and rule types. This data source automatically handles action intersections when multiple applications are specified, returning only actions supported by ALL applications.

**NOTE**: Note that some new actions may not be returned in the API response. This is a known issue, and is being investigated via the following issue `ONEAPI-2421`. Please contact Zscaler support for an update if the action you're attempting ton configure isn't supported or returned in the response.

The data source provides multiple output attributes for different use cases:

* **`available_actions_without_isolate`** - Most common use case for standard rules
* **`isolate_actions`** - For Cloud Browser Isolation (CBI) rules
* **`filtered_actions`** - Custom filtering by action type (ALLOW, DENY, etc.)
* **`available_actions`** - Complete list of all actions

## Example Usage - Standard Rule (Most Common)

Use `available_actions_without_isolate` for standard rules that don't require Cloud Browser Isolation:

```hcl
data "zia_cloud_app_control_rule_actions" "chatgpt" {
  type       = "AI_ML"
  cloud_apps = ["CHATGPT_AI"]
}

resource "zia_cloud_app_control_rule" "standard" {
  name         = "ChatGPT Standard Rule"
  type         = "AI_ML"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["CHATGPT_AI"]

  # Use available_actions_without_isolate for standard rules
  actions = data.zia_cloud_app_control_rule_actions.chatgpt.available_actions_without_isolate
}
```

## Example Usage - Isolation Rule (CBI)

Use `isolate_actions` for Cloud Browser Isolation rules:

```hcl
data "zia_cloud_app_control_rule_actions" "chatgpt" {
  type       = "AI_ML"
  cloud_apps = ["CHATGPT_AI"]
}

data "zia_cloud_browser_isolation_profile" "profile" {
  name = "My-CBI-Profile"
}

resource "zia_cloud_app_control_rule" "isolate" {
  name         = "ChatGPT Isolation"
  type         = "AI_ML"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["CHATGPT_AI"]

  # Use isolate_actions for CBI rules
  actions = data.zia_cloud_app_control_rule_actions.chatgpt.isolate_actions

  # Required when using ISOLATE actions
  cbi_profile {
    id   = data.zia_cloud_browser_isolation_profile.profile.id
    name = data.zia_cloud_browser_isolation_profile.profile.name
    url  = data.zia_cloud_browser_isolation_profile.profile.url
  }
}
```

## Example Usage - Multiple Applications (Intersection)

When multiple applications are specified, the API automatically returns only actions supported by ALL applications:

```hcl
data "zia_cloud_app_control_rule_actions" "multi_ai" {
  type       = "AI_ML"
  cloud_apps = ["CHATGPT_AI", "GOOGLE_GEMINI"]
}

resource "zia_cloud_app_control_rule" "multi_ai" {
  name         = "Multiple AI Apps"
  type         = "AI_ML"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["CHATGPT_AI", "GOOGLE_GEMINI"]

  # Returns only actions supported by BOTH applications
  actions = data.zia_cloud_app_control_rule_actions.multi_ai.available_actions_without_isolate
}

# View the intersection
output "common_actions" {
  value = data.zia_cloud_app_control_rule_actions.multi_ai.available_actions_without_isolate
}
```

## Example Usage - Filter by Action Type (ALLOW only)

Use `action_prefixes` to filter actions by type:

```hcl
data "zia_cloud_app_control_rule_actions" "allow_only" {
  type            = "AI_ML"
  cloud_apps      = ["CHATGPT_AI"]
  action_prefixes = ["ALLOW"]  # Filter for ALLOW actions only
}

resource "zia_cloud_app_control_rule" "allow_rule" {
  name         = "ChatGPT Allow Only"
  type         = "AI_ML"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["CHATGPT_AI"]

  # Only ALLOW_ actions
  actions = data.zia_cloud_app_control_rule_actions.allow_only.filtered_actions
}
```

## Example Usage - Filter Multiple Action Types

Filter for multiple action types simultaneously:

```hcl
data "zia_cloud_app_control_rule_actions" "allow_deny" {
  type            = "AI_ML"
  cloud_apps      = ["CHATGPT_AI"]
  action_prefixes = ["ALLOW", "DENY"]  # Get both ALLOW and DENY actions
}

resource "zia_cloud_app_control_rule" "mixed_rule" {
  name         = "ChatGPT Mixed Actions"
  type         = "AI_ML"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["CHATGPT_AI"]

  # ALLOW_ and DENY_ actions only (excludes CAUTION, ISOLATE, ESC)
  actions = data.zia_cloud_app_control_rule_actions.allow_deny.filtered_actions
}
```

## Example Usage - File Sharing Applications

```hcl
data "zia_cloud_app_control_rule_actions" "onedrive" {
  type       = "FILE_SHARE"
  cloud_apps = ["ONEDRIVE"]
}

resource "zia_cloud_app_control_rule" "onedrive_rule" {
  name         = "OneDrive Controls"
  type         = "FILE_SHARE"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["ONEDRIVE"]

  # Get all file sharing actions except ISOLATE
  actions = data.zia_cloud_app_control_rule_actions.onedrive.available_actions_without_isolate
}
```

## Example Usage - Only DENY Actions

```hcl
data "zia_cloud_app_control_rule_actions" "deny_only" {
  type            = "AI_ML"
  cloud_apps      = ["CHATGPT_AI"]
  action_prefixes = ["DENY"]
}

resource "zia_cloud_app_control_rule" "block_chatgpt" {
  name         = "Block ChatGPT Features"
  type         = "AI_ML"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["CHATGPT_AI"]

  # Only DENY_ actions (restrictive)
  actions = data.zia_cloud_app_control_rule_actions.deny_only.filtered_actions
}
```

## Argument Reference

The following arguments are supported:

### Required

* `type` - (String) The rule type for the Cloud App Control policy. Valid values:
  * `AI_ML` - AI/Machine Learning applications
  * `FILE_SHARE` - File sharing applications
  * `ENTERPRISE_COLLABORATION` - Enterprise collaboration tools
  * `SOCIAL_NETWORKING` - Social networking platforms
  * `STREAMING` - Streaming media applications
  * `WEB_MAIL` - Web-based email services
  * And other supported rule types

* `cloud_apps` - (List of Strings) List of cloud application names. When multiple applications are specified, the API automatically returns only actions supported by ALL applications. To retrieve available cloud applications, use the `zia_cloud_applications` data source.

### Optional

* `action_prefixes` - (List of Strings) Optional list of action prefixes to filter results. Valid values:
  * `ALLOW` - Permissive actions (e.g., ALLOW_AI_ML_CHAT)
  * `DENY` - Restrictive actions (e.g., DENY_AI_ML_UPLOAD)
  * `BLOCK` - Block actions (e.g., BLOCK_FILE_SHARE_DOWNLOAD)
  * `CAUTION` - Warning actions (e.g., CAUTION_AI_ML_WEB_USE)
  * `ISOLATE` - Cloud Browser Isolation actions (e.g., ISOLATE_AI_ML_WEB_USE)
  * `ESC` - Conditional access actions

  **Note**: The underscore is automatically added. Multiple prefixes can be specified to include multiple action types.

## Attributes Reference

The following attributes are exported:

* `available_actions` - (List of Strings) Complete list of all available actions for the specified cloud applications and rule type, including ISOLATE actions. Use this when you need the full list or want to apply custom filtering.

* `available_actions_without_isolate` - (List of Strings) **Most common use case**. List of available actions excluding ISOLATE actions. Use this for standard rules. ISOLATE actions cannot be mixed with other actions.

* `isolate_actions` - (List of Strings) List of only ISOLATE actions. Use this for Cloud Browser Isolation (CBI) rules. ISOLATE actions require `cbi_profile` configuration and cannot be mixed with other action types.

* `filtered_actions` - (List of Strings) List of actions filtered by the `action_prefixes` parameter. Only populated when `action_prefixes` is specified. Use this for custom filtering by specific action types (ALLOW only, DENY only, etc.).

## Understanding Action Types

### Action Prefixes

Cloud App Control rules support different action types based on the application and rule type:

| Prefix | Description | Example | Can Mix With |
|--------|-------------|---------|--------------|
| `ALLOW` | Permit specific operations | ALLOW_AI_ML_CHAT | DENY, CAUTION, ESC |
| `DENY` | Block specific operations | DENY_AI_ML_UPLOAD | ALLOW, CAUTION, ESC |
| `BLOCK` | Block operations (some apps) | BLOCK_FILE_SHARE_DOWNLOAD | ALLOW, CAUTION |
| `CAUTION` | Warn before allowing | CAUTION_AI_ML_WEB_USE | ALLOW, DENY, BLOCK |
| `ISOLATE` | Cloud Browser Isolation | ISOLATE_AI_ML_WEB_USE | **Cannot mix** |
| `ESC` | Conditional access | AI_ML_CONDITIONAL_ACCESS | ALLOW, DENY |

### Important Rules

1. **ISOLATE Actions**:
   * Cannot be mixed with any other action type
   * Require `cbi_profile` configuration in the resource
   * Use `isolate_actions` attribute or filter with `action_prefixes = ["ISOLATE"]`

2. **Multiple Applications**:
   * The API automatically returns the intersection of actions
   * Only actions supported by ALL specified applications are returned
   * Always query the data source with the same applications you'll use in the resource

3. **Action Compatibility**:
   * Most actions can be mixed (ALLOW + DENY, ALLOW + CAUTION, etc.)
   * ISOLATE actions are the exception - they must be used alone

## Best Practices

### 1. Use Data Source Instead of Hardcoding

**❌ Avoid hardcoding actions**:

```hcl
resource "zia_cloud_app_control_rule" "example" {
  actions = ["ALLOW_AI_ML_CHAT", "DENY_AI_ML_UPLOAD"]  # May become invalid
}
```

**✅ Use data source**:

```hcl
data "zia_cloud_app_control_rule_actions" "actions" {
  type       = "AI_ML"
  cloud_apps = ["CHATGPT_AI"]
}

resource "zia_cloud_app_control_rule" "example" {
  actions = data.zia_cloud_app_control_rule_actions.actions.available_actions_without_isolate
}
```

### 2. Match Applications Between Data Source and Resource

**❌ Mismatch (will cause validation errors)**:

```hcl
data "zia_cloud_app_control_rule_actions" "actions" {
  cloud_apps = ["CHATGPT_AI"]  # Only one app
}

resource "zia_cloud_app_control_rule" "example" {
  applications = ["CHATGPT_AI", "GOOGLE_GEMINI"]  # Two apps
  actions      = data.zia_cloud_app_control_rule_actions.actions.available_actions_without_isolate
}
```

**✅ Correct match**:

```hcl
data "zia_cloud_app_control_rule_actions" "actions" {
  cloud_apps = ["CHATGPT_AI", "GOOGLE_GEMINI"]  # Same apps
}

resource "zia_cloud_app_control_rule" "example" {
  applications = ["CHATGPT_AI", "GOOGLE_GEMINI"]  # Same apps
  actions      = data.zia_cloud_app_control_rule_actions.actions.available_actions_without_isolate
}
```

### 3. Choose the Right Output Attribute

| Use Case | Attribute to Use | Example |
|----------|------------------|---------|
| Standard rule (no CBI) | `available_actions_without_isolate` | Most common |
| CBI/Isolation rule | `isolate_actions` | Requires cbi_profile |
| Only permissive actions | `filtered_actions` with `action_prefixes = ["ALLOW"]` | Allow-only policy |
| Only restrictive actions | `filtered_actions` with `action_prefixes = ["DENY"]` | Deny-only policy |
| Mixed ALLOW/DENY | `filtered_actions` with `action_prefixes = ["ALLOW", "DENY"]` | Fine-grained control |
| Full list for custom logic | `available_actions` | Manual filtering |

## Complete Examples

### Example 1: Standard Rule with Multiple Action Types

```hcl
data "zia_cloud_app_control_rule_actions" "slack" {
  type       = "ENTERPRISE_COLLABORATION"
  cloud_apps = ["SLACK"]
}

resource "zia_cloud_app_control_rule" "slack_controls" {
  name                    = "Slack Controls"
  description             = "Control Slack usage"
  type                    = "ENTERPRISE_COLLABORATION"
  order                   = 1
  rank                    = 7
  state                   = "ENABLED"
  applications            = ["SLACK"]
  browser_eun_template_id = 5502

  # Returns all actions except ISOLATE
  actions = data.zia_cloud_app_control_rule_actions.slack.available_actions_without_isolate
}
```

### Example 2: Permissive Rule (ALLOW Only)

```hcl
data "zia_cloud_app_control_rule_actions" "dropbox_allow" {
  type            = "FILE_SHARE"
  cloud_apps      = ["DROPBOX"]
  action_prefixes = ["ALLOW"]
}

resource "zia_cloud_app_control_rule" "dropbox_allow" {
  name         = "Dropbox Allow Operations"
  type         = "FILE_SHARE"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["DROPBOX"]

  # Only permissive actions
  actions = data.zia_cloud_app_control_rule_actions.dropbox_allow.filtered_actions
}
```

### Example 3: Restrictive Rule (DENY Only)

```hcl
data "zia_cloud_app_control_rule_actions" "onedrive_deny" {
  type            = "FILE_SHARE"
  cloud_apps      = ["ONEDRIVE"]
  action_prefixes = ["DENY"]
}

resource "zia_cloud_app_control_rule" "onedrive_block_upload" {
  name         = "OneDrive Block Upload"
  type         = "FILE_SHARE"
  order        = 2
  rank         = 7
  state        = "ENABLED"
  applications = ["ONEDRIVE"]

  # Only restrictive DENY actions
  actions = data.zia_cloud_app_control_rule_actions.onedrive_deny.filtered_actions
}
```

### Example 4: Multiple Applications with Intersection

```hcl
# Query actions for two applications
data "zia_cloud_app_control_rule_actions" "multi_file_share" {
  type       = "FILE_SHARE"
  cloud_apps = ["ONEDRIVE", "DROPBOX"]
}

resource "zia_cloud_app_control_rule" "multi_file_share" {
  name         = "File Sharing Controls"
  type         = "FILE_SHARE"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["ONEDRIVE", "DROPBOX"]

  # Returns only actions supported by BOTH OneDrive AND Dropbox
  actions = data.zia_cloud_app_control_rule_actions.multi_file_share.available_actions_without_isolate
}

# Output shows the intersection
output "common_file_share_actions" {
  value = data.zia_cloud_app_control_rule_actions.multi_file_share.available_actions_without_isolate
  # Example output: Actions both apps support
}
```

### Example 5: CAUTION Actions Only

```hcl
data "zia_cloud_app_control_rule_actions" "caution_only" {
  type            = "AI_ML"
  cloud_apps      = ["CHATGPT_AI"]
  action_prefixes = ["CAUTION"]
}

resource "zia_cloud_app_control_rule" "caution_rule" {
  name         = "ChatGPT Caution"
  type         = "AI_ML"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  applications = ["CHATGPT_AI"]

  # Only CAUTION actions (user warnings)
  actions = data.zia_cloud_app_control_rule_actions.caution_only.filtered_actions
}
```

### Example 6: Viewing All Available Attributes

```hcl
data "zia_cloud_app_control_rule_actions" "chatgpt" {
  type            = "AI_ML"
  cloud_apps      = ["CHATGPT_AI"]
  action_prefixes = ["ALLOW", "DENY"]  # Optional filtering
}

# View all output attributes
output "all_actions" {
  value = data.zia_cloud_app_control_rule_actions.chatgpt.available_actions
  # All actions including ISOLATE (17 actions for ChatGPT)
}

output "standard_actions" {
  value = data.zia_cloud_app_control_rule_actions.chatgpt.available_actions_without_isolate
  # All except ISOLATE (16 actions)
}

output "isolate_only" {
  value = data.zia_cloud_app_control_rule_actions.chatgpt.isolate_actions
  # Only ISOLATE actions (1 action)
}

output "custom_filtered" {
  value = data.zia_cloud_app_control_rule_actions.chatgpt.filtered_actions
  # Only ALLOW and DENY actions (based on action_prefixes)
}
```

## Argument Reference

The following arguments are supported:

### Required

* `type` - (String) The rule type for the Cloud App Control policy. Common values include:
  * `AI_ML` - AI/Machine Learning applications
  * `FILE_SHARE` - File sharing and storage applications
  * `ENTERPRISE_COLLABORATION` - Collaboration tools (Slack, Teams, etc.)
  * `SOCIAL_NETWORKING` - Social media platforms
  * `STREAMING` - Streaming media services
  * `WEB_MAIL` - Web-based email services
  * `BUSINESS_PRODUCTIVITY` - Business productivity tools
  * `SALES_AND_MARKETING` - Sales and marketing applications
  * And more...

* `cloud_apps` - (List of Strings) List of cloud application names to retrieve actions for. When multiple applications are specified, the API automatically computes and returns only actions supported by ALL applications (intersection). Use the `zia_cloud_applications` data source to get available application names.

### Optional

* `action_prefixes` - (List of Strings) Optional list of action prefixes to filter results. Valid values: `ALLOW`, `DENY`, `BLOCK`, `CAUTION`, `ISOLATE`, `ESC`. The underscore is automatically added (e.g., `ALLOW` becomes `ALLOW_`). Multiple prefixes can be specified. When specified, results are available in the `filtered_actions` attribute.

## Attributes Reference

The following attributes are exported:

* `available_actions` - (List of Strings) Complete list of all available actions for the specified cloud applications and rule type, including ISOLATE actions. Use when you need the full list or want to apply custom Terraform filtering logic.

* `available_actions_without_isolate` - (List of Strings) **Recommended for most use cases**. List of available actions excluding ISOLATE actions. Use this for standard Cloud App Control rules. ISOLATE actions cannot be mixed with other action types and require separate rules.

* `isolate_actions` - (List of Strings) List of only ISOLATE actions (Cloud Browser Isolation). Use this for CBI rules. When using ISOLATE actions:
  * They **cannot** be mixed with other action types (ALLOW, DENY, etc.)
  * They **require** `cbi_profile` block in the resource
  * They **cannot** have `browser_eun_template_id` set
  * Create separate rules for ISOLATE vs non-ISOLATE actions

* `filtered_actions` - (List of Strings) List of actions filtered by the `action_prefixes` parameter. Only populated when `action_prefixes` is specified. Use this for custom filtering by specific action types (ALLOW only, DENY only, ALLOW+DENY, etc.).

## Notes

### Application Intersection Behavior

When querying multiple applications, the API returns only the intersection of actions:

**Example**:

* `CHATGPT_AI` alone supports 12 actions (including ALLOW_AI_ML_RENAME)
* `GOOGLE_GEMINI` alone supports 11 actions (does NOT support RENAME)
* Query with both: `["CHATGPT_AI", "GOOGLE_GEMINI"]` returns 9 actions (RENAME excluded)

This ensures that rules with multiple applications only use actions that work for all of them.

### ISOLATE Actions Special Requirements

ISOLATE actions have unique requirements:

1. **Cannot be mixed**: ISOLATE actions must be used alone in a rule
2. **Require CBI profile**: Must configure `cbi_profile` block with a valid profile
3. **No EUN template**: Cannot set `browser_eun_template_id` when using ISOLATE
4. **Separate rules**: Create one rule for ISOLATE actions, separate rules for other actions

### Validation

The `zia_cloud_app_control_rule` resource automatically validates actions during `terraform plan`:

* Ensures actions are valid for the specified applications
* Validates ISOLATE action requirements
* Provides helpful error messages with valid action lists
* Suggests using the data source if manual actions are invalid

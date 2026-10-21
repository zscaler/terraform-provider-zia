---
page_title: "Best Practices"
description: "Recommended practices for managing Zscaler Internet Access (ZIA) with the Zscaler Terraform provider"
---

# Best Practices

This guide collects the recommended practices for managing Zscaler Internet Access (ZIA) with Terraform: the general Terraform practices as they apply to this provider, followed by the practices specific to ZIA — activation, the API session timeout, rule ordering and running against one tenant. Each recommendation is consistent with the rest of this provider's documentation, which it links to where relevant.

-> **Disclaimer:** Please refer to our [General Support Statement](support.md) before proceeding with the use of this provider. For specific errors, see the [Troubleshooting Guide](troubleshooting.md).

~> **Before you start** Set the ZIA **API Session Timeout** to its maximum of **20 minutes** before using the provider, and run one Terraform process at a time per tenant. See [Before You Start: Extend the API Session Timeout](https://registry.terraform.io/providers/zscaler/zia/latest/docs#before-you-start-extend-the-api-session-timeout) and [API session timeout and automatic activation](#api-session-timeout-and-automatic-activation).

## General practices

Terraform practice that applies to the ZIA provider as it does to the other Zscaler providers. ZIA-specific behaviour — activation, the session timeout, rule ordering and tenant concurrency — is in [ZIA-specific practices](#zia-specific-practices).

### Authentication model

- **DO** use OneAPI (Zidentity OAuth2) when the tenant has migrated — `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET` (or `ZSCALER_PRIVATE_KEY`) and `ZSCALER_VANITY_DOMAIN`
- **DO** use the legacy framework (`use_legacy_client = true`) for tenants not yet migrated to Zidentity — the docs call it "the recommended authentication method" for those tenants
- **DO** set `zscaler_cloud` only for `beta` or a FedRAMP cloud (`gov` / `govus`, v4.7.25+) — the docs say it is not required for production commercial clouds
- **DON'T** put a legacy cloud name (`zscaler`, `zscalergov`, `zscalerten` …) in `zscaler_cloud` — those belong in `zia_cloud` on the legacy path
- **DON'T** leave both credential sets in one job — the path is chosen **only** by `use_legacy_client` / `ZSCALER_USE_LEGACY_CLIENT`, and the other set is silently ignored
- **DON'T** set both `client_secret` and `private_key` — the argument reference documents them as conflicting
- **DON'T** hard-code credentials in the provider block — the docs warn it "risks secret leakage should this file be committed to public version control"

#### Determining the auth path

The provider decides the path from one input: `use_legacy_client`, falling back to `ZSCALER_USE_LEGACY_CLIENT`. When it is false, `ZIA_USERNAME` / `ZIA_PASSWORD` / `ZIA_API_KEY` are never consulted, however correct they are — which is why a job carrying both namespaces fails with a confusing 401 rather than an explicit error.

API access is granted by role. Create an API Role in ZIA (*Administration > Role Management > Add API Role*), then, if the role does not appear in Zidentity, use *Sync Now* under *Integration > API Resources*.

**Do — empty provider blocks; credentials from the environment**

```hcl
provider "zia" {}   # OneAPI: ZSCALER_CLIENT_ID, ZSCALER_CLIENT_SECRET, ZSCALER_VANITY_DOMAIN
# Legacy: ZIA_USERNAME, ZIA_PASSWORD, ZIA_API_KEY, ZIA_CLOUD, ZSCALER_USE_LEGACY_CLIENT=true
# FedRAMP over OneAPI (v4.7.25+): ZSCALER_CLOUD=gov | govus
```

**Don't — both namespaces in one job**

```sh
export ZSCALER_CLIENT_ID=...       # OneAPI set
export ZIA_USERNAME=...            # legacy set — silently ignored
export ZSCALER_CLOUD=zscalergov    # legacy cloud name — rejected on OneAPI
```

### Credential handling

- **DO** let the provider read credentials from environment variables, injected from a secret store at job start (Vault, AWS Secrets Manager, Azure Key Vault — reached through your CI's OIDC)
- **DO** keep state in an encrypted, access-controlled remote backend — secrets you set on resources are stored in state in plain text
- **DO** protect saved plan files (`tfplan`) like state — they contain the values of the variables used to produce them
- **DO** prefer `private_key` over `client_secret` where the customer can manage keys
- **DON'T** count on write-only (`*_wo`) arguments or ephemeral resources — Terraform 1.11 supports them, but the provider does not expose any today
- **DON'T** commit `terraform.tfvars` with secrets — commit `terraform.tfvars.example`
- **DON'T** pass secrets with `-var "client_secret=…"` — it lands in shell history and CI logs
- **DON'T** output secrets at all, even marked `sensitive`
- **DON'T** leave debug logging on in routine CI. When a debug log is needed for a support case, sanitize it before sharing — see [Troubleshooting and log collection](#troubleshooting-and-log-collection)

#### Two kinds of secret

**Provider credentials** — client ID, client secret, private key, and the legacy credentials. Provider configuration is never written to Terraform state. When the provider reads these from the environment, as the provider documentation recommends, they never appear in state, in plan output, or in code.

**Secrets that are resource attributes** — a VPN `pre_shared_key` on `zia_traffic_forwarding_vpn_credentials`, an admin `password` on `zia_admin_users`. Terraform must store these in state, and stores them in plain text. The provider marks them sensitive, and you can mark variables `sensitive = true`, but both only redact CLI output. `sensitive` is not encryption. The control is the backend: encryption at rest, tight IAM, and state locking.

**Don't**

```hcl
variable "api_key" {
  default = "9f8a7b6c5d4e"     # in git forever
}

output "api_key" {
  value     = var.api_key
  sensitive = true              # redacted on screen, still readable
}
```

**Do**

```hcl
provider "zia" {}   # credentials from ZSCALER_* env vars — never in state

variable "psk" {
  type      = string
  sensitive = true   # from TF_VAR_psk; the value still lands in state
}                    # -> encrypt and restrict the backend

# Output only non-sensitive identifiers
output "rule_id" {
  value = zia_url_filtering_rules.block_gambling.id
}
```

-> **Two operational habits** Add `trivy config --severity HIGH,CRITICAL --exit-code 1 .` to CI to catch secrets in tfvars and insecure backend configuration before merge.

Rotate static credentials every 90 days at most.

### Resource addressing: for_each, not count

- Terraform tracks every resource by an **address**. `count` addresses by position, `for_each` addresses by key
- `count` → `rules[0]`, `rules[1]`, `rules[2]` … **the index is the identity**
- Remove the 2nd of 20 items and every later element shifts down one position, so all 18 addresses now point at different objects — Terraform reads that as **18 destroy/recreates**, not one delete
- `for_each` → `rules["block_gambling"]` — the key is the identity, so removing one key touches **only** that object
- **DO** use `for_each` with a map keyed by **purpose**; use `count` only for on/off (`count = var.enabled ? 1 : 0`)

Terraform does not identify a resource by its name or its content — it identifies it by an address in state. With `count`, that address contains the list index, so **the position in the list becomes the object's identity**. Delete the middle entry from the input list and Terraform sees that every later index now holds a different object, so it destroys and recreates them. One deletion becomes N changes. With `for_each` the address carries a stable key, so removing `"streaming"` removes exactly `rules["streaming"]` and nothing else is re-evaluated.

The cost is higher on a ZIA tenant than in a generic Terraform context. Every needless destroy/recreate is an API write against the documented rate limits, each recreated rule must be placed back into its order — which the docs note "takes longer than for most other resources" — and in the meantime the policy you meant to keep does not exist. Key by purpose (`"block_gambling"`), never by number (`"rule_1"`).

**Don't**

```hcl
resource "zia_url_filtering_rules" "rules" {
  count          = length(var.rule_names)   # address = rules[0], rules[1] …
  name           = var.rule_names[count.index]
  order          = count.index + 1
  url_categories = ["ANY"]
}
# Remove "Streaming" from the list:
#   rules[1] was Streaming -> now Social -> destroy + recreate
```

**Do**

```hcl
resource "zia_url_filtering_rules" "rules" {
  for_each       = var.rules            # map(object({ order = number, … }))
  name           = each.key             # address = rules["block_gambling"]
  order          = each.value.order
  url_categories = each.value.url_categories   # required from v4.8.11
}
```

### Drift and ignore_changes

- **Rare, but it happens:** an attribute you set is missing — or returned empty — in the API response
- Terraform stores what the API returns, so the value never lands in state and every plan shows a change
- **This is an API-side issue, not a provider bug** — the provider can only record what it is given
- Common shapes: booleans omitted when false, API-injected defaults, write-only secrets never returned, list ordering
- **DO** check the provider CHANGELOG and upgrade first — many drift cases are fixed in a later release
- **DON'T** reach for `ignore_changes` as the *first* move, and never `= all` — it hides genuine drift too

#### Two available responses, in order of preference

1. **Raise a bug with the respective engineering team** so it is fixed at the API backend. This is the real fix.
2. **Temporarily** add `lifecycle { ignore_changes = [<field>] }` scoped to that field until the root cause is resolved.

Both the time limit and the scoping matter. A narrowly scoped suppression with a comment and a ticket number is a documented workaround that someone can find and remove later. `ignore_changes = all` is how a configuration silently diverges from the tenant for months, and it will mask the *next*, unrelated problem too.

Write-only fields are a special case and not a defect: passwords, pre-shared keys and API keys are never returned by design. After an import the secret is not recoverable from the API and must be re-supplied through a variable.

**Don't**

```hcl
# Blanket suppression hides genuine drift as well as the bug
lifecycle {
  ignore_changes = all
}
```

**Do**

```hcl
# Narrowly scoped to the affected field, with a comment and a ticket
# reference. Remove once the API fix ships.
lifecycle {
  ignore_changes = [description]   # SUP-XXXX: absent from GET
}
```

### Parallelism and provider tuning

- **DO** run Terraform with its default parallelism — the docs say **do not change `-parallelism`**. Raising it causes rate limiting, edit-lock conflicts and failed rule placement; lowering it slows rule placement and stretches runs across the API session timeout
- `-parallelism` is **global to the whole run**. There is no per-resource, per-endpoint or per-provider setting, so lowering it to help one resource type slows down **everything else** in the same configuration
- **Why ZIA is hurt most:** ZIA orders rules server-side, and the provider reconciles rule placement per **batch** of concurrently created rules. Starve the batch and you pay the reconciliation cost per rule instead of per batch
- Rate limiting needs no help: the docs state that "rate limiting requires no configuration" — on a 429 the provider waits out `Retry-After` and retries automatically
- The provider **argument** `parallelism` is deprecated and ignored from v4.8.1 — remove it
- Seeing lock contention (`EDIT_LOCK_NOT_AVAILABLE`)? The lever is fewer concurrent **runs**, not slower ones — see [Concurrency: one Terraform process per tenant](#concurrency-one-terraform-process-per-tenant)

#### Why lowering parallelism backfires

Rule order is maintained server-side, so after creating rules the provider has to move each rule into its declared position — the rule pages note that "the ZIA API does not provide a way to set the position of several rules in one request". That reconciliation runs on a cycle shared by all rules created concurrently. At default parallelism, ten rules share one cycle; at `-parallelism=1` each rule sits through a cycle by itself. In one support case, ~360 firewall rules took more than six hours at `-parallelism=1`; the same apply at the default finished in well under an hour — and a run that long collides with the API session timeout ([API session timeout and automatic activation](#api-session-timeout-and-automatic-activation)).

**Don't**

```sh
terraform apply -parallelism=1   # ~360 rules: 6+ hours

provider "zia" {
  parallelism = 1                # deprecated, no effect — delete it
}
```

#### Provider tuning arguments: leave them at the defaults

The argument reference states it explicitly: rate limiting requires no configuration; `max_retries` and `request_timeout` are the **only** supported tuning attributes; neither changes request throughput, so lowering them cannot make plan or apply faster. Environment variables belonging to the Zscaler Go SDK, such as `ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES` or `ZSCALER_CLIENT_REQUEST_TIMEOUT`, are **not** part of the provider configuration and are ignored.

| Provider | `max_retries` default | `request_timeout` default (range) | Deprecated and ignored |
|---|---|---|---|
| **ZIA** | 100 (maximum 100) | 1800 s (0–1800) | `parallelism` (v4.8.1+) |

A `request_timeout` of `0` selects the SDK's built-in 60-second timeout; it does **not** disable the timeout.

One documented exception to "leave the defaults": the troubleshooting guide suggests lowering `request_timeout` (for example to `300`) when a single request stalls until Terraform's operation timeout, so that it fails sooner. The default of 1800 seconds exists because reading thousands of firewall rules can take that long.

**Don't — tuning that does nothing, or harm**

```sh
export ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES=2   # ignored by the provider
export ZSCALER_CLIENT_REQUEST_TIMEOUT=15         # ignored by the provider
provider "zia" {
  parallelism = 1                # deprecated, no effect
  max_retries = 2                # only makes runs fail sooner
}
```

#### Conditionally disabled providers: skip_credentials_validation

A shared configuration deployed to several environments may declare the ZIA provider where Zscaler is not used, with every `zia_*` resource disabled through `count = 0`. Terraform still configures the provider, which validates credentials, so a provider block without credentials fails the run. Set `skip_credentials_validation = true` (or `ZSCALER_SKIP_CREDENTIALS_VALIDATION=true`, v4.8.7+) in those environments only; any resource or data source that does attempt an API call fails with an explanatory error.

**Do**

```hcl
provider "zia" {
  skip_credentials_validation = !var.zscaler_enabled
}

resource "zia_rule_labels" "this" {
  count = var.zscaler_enabled ? 1 : 0
  name  = "example"
}
```

If a plan is slow on a very large configuration, the cause is the number of resources refreshed at the documented rate limits, not provider settings. Routine changes can plan with `-refresh=false` (in HCP Terraform / TFE, set `TF_CLI_ARGS_plan="-refresh=false"` on the workspace) as long as a scheduled full-refresh plan still detects out-of-band changes; configurations that keep growing should be split ([State organisation](#state-organisation)).

### State organisation

- How you lay out state and which backend you use is **your call** — any remote backend works
- **DO** split on whatever boundary matches your teams — by product, by rule type, by environment, by approval group
- **DO** use `terraform_remote_state` (or equivalent) to share IDs across states rather than merging them
- **DON'T** import the same object into two states — every object needs exactly one owner

#### ZIA constraints on layout

| Provider | Constraint on state layout and execution | Section |
|---|---|---|
| **ZIA** | All rules of one type in one configuration · one apply at a time per tenant · one activation per tenant, after the last apply | [Rule ordering](#rule-ordering), [Concurrency: one Terraform process per tenant](#concurrency-one-terraform-process-per-tenant) |

This section is deliberately light on prescription. State layout and backend choice are architecture decisions for the operator. What *is* ours to say is the constraint above, and it is a consequence of API behaviour rather than an opinion.

None of it argues for one big state: splitting ZIA firewall rules from URL filtering rules is fine, because their orderings are independent. The execution constraints are the ones people get wrong once they have more than one state — splitting state does not split the tenant.

**One workable layout among many**

```text
infrastructure/zia/
├── prod/policy-network/   # firewall, DNS, IPS rules — each type owned here
├── prod/policy-content/   # URL filtering, cloud app control
├── prod/locations/
└── prod/settings/         # zia_advanced_settings (api_session_timeout = 20)
```

### Modules

- **DO** keep `main.tf`, `variables.tf`, `outputs.tf`, `versions.tf`, `README.md` in every module
- **DO** declare `required_providers` in **every** module, and keep the `version` constraint at the environment (root) level only — the provider troubleshooting guides recommend exactly this to avoid *registry does not have a provider named hashicorp/…*
- **DO** configure providers **only in the root module**
- **DO** give each module a single clear **purpose** — purpose is the boundary, not product
- **DO** declare multiple providers in `required_providers` when the *feature* genuinely spans products
- **DON'T** put a `provider` block in a reusable module
- **DON'T** build a `zscaler_everything` god module bundling unrelated products
- **DON'T** hard-code IDs like `"72058304855457833"` — use data sources

#### Purpose as the module boundary

The blunt formulation — "never mix ZIA and ZPA in one module" — is incorrect. The real boundary is **purpose**: "can I describe what this module does in one sentence without using the word *and*?" A module implementing one feature that spans two products — IP Source Anchoring (a ZIA gateway defined by a ZPA server group and segment), or ZTC forwarding rules that take `src_workload_groups` IDs from the ZIA provider — needs both providers. Even then it declares both only in `required_providers` and configures neither; the root module owns authentication, which is what makes aliasing and multi-tenant composition work.

**Legitimately cross-product — IP Source Anchoring**

```sh
# modules/zia-zpa-source-ip-anchoring/versions.tf — one purpose, two providers
terraform {
  required_providers {
    zia = { source = "zscaler/zia" }
    zpa = { source = "zscaler/zpa" }
  }
}
# no provider block — the root composes and authenticates both
data "zpa_server_group" "this" { name = "Server_Group_IP_Source_Anchoring" }
resource "zia_forwarding_control_zpa_gateway" "this" {
  name = "ZPA_GW01"
  type = "ZPA"
  zpa_server_group {
    external_id = data.zpa_server_group.this.id
    name        = data.zpa_server_group.this.name
  }
}
```

### Version pinning

- **DO** use a pessimistic constraint: `zia = "~> 4.0"` — v4 is the OneAPI floor
- **DO** set `required_version` from the features you use: `>= 1.3` minimum (stable `optional()` defaults); `>= 1.7` if you use `removed {}` or mock providers. Every provider CHANGELOG states support for Terraform **v1.x**
- **DO** commit `.terraform.lock.hcl` — the troubleshooting guides show the *Failed to install provider* checksum error you get without it — and upgrade in a dedicated PR tested in non-prod first
- **DO** read the provider CHANGELOG before filing a bug
- **DON'T** omit version constraints
- **DON'T** pin an exact GA version (`= 4.0.3`) without a specific reason — it blocks bug fixes
- **DON'T** run `init -upgrade` inside a feature PR

| Provider | Latest release (verified 15 Oct 2026) | Recommended constraint |
|---|---|---|
| `zscaler/zia` | v4.8.13 (15 Oct 2026) | `~> 4.0` |

Several fixes described in this guide are version-gated, and "upgrade the provider" is the first step in the troubleshooting guide. Recent examples: v4.8.7 stopped rule-order reconciliation running for hours; v4.8.10 fixed `request_timeout` being capped at 300; v4.8.11 retries automatically when an API session ends mid-run, masks tokens in debug logs and makes `url_categories` required on `zia_url_filtering_rules`; v4.8.12 makes DLP rule updates fail immediately with the API error instead of retrying until the timeout; v4.8.13 skips rule types the tenant is not subscribed to when detaching an object before deleting it.

**Do**

```sh
terraform {
  required_version = ">= 1.7.0, < 2.0.0"   # 1.7+: removed {}, mock providers
  required_providers {
    zia = { source = "zscaler/zia", version = "~> 4.0" }
  }
}
```

### Naming, variables and outputs

- All lowercase; underscores in Terraform identifiers, hyphens in names that appear in the Zscaler portal
- Names communicate purpose: `internal_hr_portal`, not `app1` / `a` / `test`
- Carry environment context in a local rather than repeating it: `prefix = "${var.environment}-${var.project}"`
- Separate required from optional variables, and give every variable a `description` and a `type`
- Use `validation` blocks to fail at plan time rather than at the API
- Output identifiers, never secrets ([Credential handling](#credential-handling))

### CI/CD

1. **Validate** — `terraform fmt -check`, `validate`, `tflint`
2. **Scan** — `trivy config`, `checkov`
3. **Plan** — `terraform plan -out=tfplan`, save the artifact
4. **Apply** — `terraform apply tfplan` (the reviewed artifact)
5. **Activate** — `ziaActivator` as its own stage, once, after the last apply

- **DON'T** re-run plan inside the apply job — you'd apply something nobody reviewed
- **DO** keep host-cloud (backend) auth and Zscaler API auth as separate steps
- **DO** serialise applies with a concurrency group keyed on the **tenant** — the docs: "apply one configuration at a time against a given tenant"

The saved-artifact rule is the one most often bypassed, usually to work around an expired plan. Re-planning in the apply job silently discards human review — it is the CI equivalent of `--auto-approve`.

Authenticating to S3 or Azure Blob for state is a completely different concern from authenticating to the Zscaler API. Keep the two credential paths visibly separate.

With several states against one tenant, activation is one job for the whole pipeline, not one per state. The ZIA activator guide gives the shape: plan in parallel, apply one configuration at a time with a concurrency control keyed on the tenant (`concurrency` in GitHub Actions, `resource_group` in GitLab CI, run queuing in HCP Terraform), and run `ziaActivator` once as the final stage. Avoid cancelling an in-progress apply.

**With several states against one tenant**

```text
plan (parallel — reads are safe) ──▶ apply (serialised on tenant) ──▶ activate (once)
firewall · url · dlp · locations      one at a time                 ziaActivator
```

On HCP Terraform or Terraform Enterprise, routine plans can skip the refresh with the workspace variable `TF_CLI_ARGS_plan="-refresh=false"`. Keep a scheduled full-refresh plan so changes made outside Terraform are still detected.

-> **Two scheduling habits** Run drift detection on a **daily schedule**, not per-commit.

Never run integration tests against production from PR CI.

### Testing

- **DO** run static analysis (`fmt`, `validate`, `tflint`) on every PR
- **DO** run `terraform plan` against a **non-prod tenant** on every PR
- **DO** use `terraform test` (1.6+) for variable validation
- **DO** tag test resources and run a sweeper to clean up
- **DON'T** expect mock providers to catch API behaviour — no rule-order conflicts, no `DUPLICATE_ITEM`, no activation dry-run
- **DON'T** expect plan to catch every ordering conflict — rank and position conflicts surface at apply
- **DON'T** test in production

Mock providers (Terraform 1.7+) validate *module wiring*, not Zscaler semantics. They are valuable and cheap, but a green mock suite says nothing about whether the API accepts the configuration. Pair them with a sandbox tenant.

~> **Activation in test teardown** Test teardown **must include activation** — `ziaActivator`, as in the pipeline. Otherwise the deletes stay as drafts and the next run collides with objects that appear deleted but are not.

Do not use `ZIA_ACTIVATION=true` for tests. The troubleshooting guide documents the result: activations exhaust the 10/min, 40/hr budget and the run hangs until `/status/activate: context deadline exceeded`.

### Import and brownfield adoption

- **DO** use **zscaler-terraformer** for brownfield adoption — `import` writes the HCL *and* imports into state; `generate` writes HCL only
- **DO** adopt **one resource type per PR**, and confirm an empty `terraform plan` before the next
- **DO** import individual objects by ID (many resources also accept the name — check the resource page)
- **DO** keep the default parallelism, even for a large initial import
- **DO** re-supply write-only secrets after import — the API won't return them
- **DO** use `removed {}` (Terraform 1.7+) to stop managing something without destroying it
- **DON'T** hand-write HCL for a tenant that already exists — that is what the tool is for
- **DON'T** use `terraform state rm` on a custom object — it orphans the API object and the next apply hits `DUPLICATE_ITEM`

zscaler-terraformer is a supported Zscaler CLI that authenticates with the same credentials as the providers, reads a live tenant over the API, and emits the Terraform configuration and the import into state, so the first plan after an import should come back clean. Scope it with `--resources`: the product name for a full sweep, or a single resource type when onboarding incrementally — the pattern to recommend. `--prefix` namespaces the generated names, `--validate` runs `terraform init` and `terraform validate` over the output, and `--collect-logs` captures SDK debug logs for a support case.

Generated HCL is correct but not idiomatic: names come from the API, resources are flat rather than grouped with `for_each`, and there are no variables. Refactor it before it becomes the long-term source of truth.

**Do — brownfield: don't hand-write it**

```sh
zscaler-terraformer import --resources "zia"                          # whole product
zscaler-terraformer import --resources "zia_firewall_filtering_rule"  # one type
zscaler-terraformer --prefix "prod" --validate import --resources "zia"
```

-> **Coverage and caveats** List what the installed version covers with `--supported-resources="zia"`, and confirm the resource types you need are in it.

It cannot recover write-only values — passwords, pre-shared keys, API keys — because the API never returns them. Re-supply those from your secret store after the import.

Supply credentials as environment variables rather than CLI flags so they do not land in shell history.

`terraform state rm` is the most frequently misused command in this context. It appears to offer a clean exit from a bad state, but it leaves the object live in the tenant and untracked, and the next apply tries to create it again. `removed {}` expresses the same intent safely and is visible in code review.

**Don't**

```sh
terraform state rm zia_url_filtering_rules.block_gambling
# Object still exists in the tenant. Next apply -> DUPLICATE_ITEM or a duplicate.
```

**Do**

```hcl
removed {
  from = zia_url_filtering_rules.block_gambling
  lifecycle { destroy = false }
}
```

## ZIA-specific practices

ZIA carries the largest configuration surface of the four providers and the two constraints that generate most support volume: changes are drafts until activated, and rule order is a contiguous server-side sequence that the configuration must track exactly. Everything here comes from the ZIA provider documentation — docs/index.md, the zia-activator-overview and troubleshooting guides, the resource pages — the CHANGELOG and the provider source at v4.8.13.

### Scope and resource model

The zscaler/zia provider registers **85 resources and 118 data sources** (zia/provider.go at v4.8.13). The excess of data sources reflects objects owned by another system — identity-provider-driven groups and departments, Zscaler-managed catalogs, and predefined policy objects — which are readable but not creatable.

**Eighteen** resource types are rule-based: their pages state that `order` must be a positive whole number starting at 1 and that orders must be contiguous. They share the ordering contract in [Rule ordering](#rule-ordering) and are listed in [ZIA rule resources sharing ordering semantics](#zia-rule-resources-sharing-ordering-semantics). Resource naming is inconsistent between singular and plural forms (`zia_url_filtering_rules` but `zia_firewall_dns_rule`); check each name against the Terraform Registry.

The documentation states that the provider maintains parity with the **public API**: a setting visible in the Admin Portal but not exposed by the API cannot be managed through Terraform. That is not the provider lagging; it is the scope of the API.

### Authentication

- **DO** use OneAPI (Zidentity OAuth2) when the tenant has migrated
- **DO** use the legacy framework (`use_legacy_client = true`) when the tenant has **not** migrated — the docs call it "the recommended authentication method" for those tenants
- **DO** use OneAPI on the FedRAMP clouds too — supported from **v4.7.25** with `zscaler_cloud = "gov"` or `"govus"`; earlier versions are not supported
- **DON'T** repeat the old "GOV is legacy-only" advice for ZIA — it stopped being true in v4.7.25
- **DON'T** put a legacy cloud name in `zscaler_cloud` — OneAPI takes only `beta`, `gov` or `govus`; `zscalergov` and `zscalerten` are legacy names that belong in `zia_cloud`
- **DON'T** mix `ZSCALER_*` and `ZIA_*` variables in one job — the path is chosen **only** by `use_legacy_client` / `ZSCALER_USE_LEGACY_CLIENT`
- Sandbox authentication is separate: `ZSCALER_SANDBOX_TOKEN` / `ZSCALER_SANDBOX_CLOUD` (the old `ZIA_CLOUD` / `ZIA_SANDBOX_TOKEN` pair is no longer supported)

The two auth paths use *different* cloud values for the same environments: OneAPI takes `gov` or `govus` in `zscaler_cloud`, while `zscalergov` and `zscalerten` are legacy names belonging in `zia_cloud`. Someone who reads "GOV is supported now" and sets `zscalergov` in `zscaler_cloud` gets a rejection and concludes GOV is still not supported.

**Don't**

```hcl
provider "zia" {
  client_id     = "abc123"      # NEVER
  client_secret = "s3cr3t"      # NEVER
  zscaler_cloud = "zscaler"     # legacy name on a OneAPI field
}
```

**Do**

```hcl
# OneAPI — credentials come from the environment
provider "zia" {}
# ZSCALER_CLIENT_ID / ZSCALER_CLIENT_SECRET / ZSCALER_VANITY_DOMAIN

# OneAPI on a FedRAMP cloud — v4.7.25+
#   ZSCALER_CLOUD=gov   <- or govus. NOT zscalergov.

# Legacy — tenants not migrated to Zidentity
#   ZIA_USERNAME / ZIA_PASSWORD / ZIA_API_KEY / ZIA_CLOUD
#   ZSCALER_USE_LEGACY_CLIENT=true
```

| Error | Cause |
|---|---|
| `401 unauthorized` | Wrong credentials or vanity domain — or the wrong auth path selected by `use_legacy_client` |
| `vanity_domain` not found | Typo, or tenant not migrated to Zidentity — use legacy |
| Cloud zscalergov not supported for OneAPI | Legacy cloud name in `zscaler_cloud`. Use `gov` / `govus` — or upgrade if below v4.7.25 |
| `403 Resource Access Blocked` from the first request | The API client's role lacks the permission or functional scope, or a subscription is missing — review the role in the Admin Portal |
| `expected zia_cloud to be one of […]` | Missing or misspelled legacy cloud name |

### Activation

- **DON'T** assume `terraform apply` = enforced — ZIA saves changes in a **pending** state
- Activation is **tenant-wide**: "a single activation call publishes every pending change in the tenant". You never need more than one per run
- **DO** run the `ziaActivator` CLI once after apply — Method 1, "recommended" in the docs
- **DO** use `zia_activation_status` (Method 2) only if a single flat state owns *all* its resources — `depends_on` must list every one of them
- **DON'T** use `ZIA_ACTIVATION=true` (Method 3, "discouraged" and "may be removed in a future major release") — it activates per resource against a 10/min, 40/hr endpoint
- Applies to *everything*: rules, labels, URL categories, locations, admin users — and to deletes

#### Why this happens

Activation is the publish step. Terraform writes the draft; activation applies it. An apply that completes successfully has done its whole job — the objects exist in the tenant and the state file is accurate. They are simply not enforcing traffic yet. Because one call publishes the entire pending state for the organisation, the only question worth asking is *when* you make that call — never *how many* calls you need.

#### Three ways this fails

1. **No activation at all.** Silent. The apply succeeds, the console shows pending changes, and nobody notices until a policy they believed was live is needed.
2. **`zia_activation_status` exists, but `depends_on` is stale.** The docs: "any resource you forget may be applied *after* activation and therefore left pending". In module-based configurations the list is "difficult to maintain and easy to get wrong" — which is why the out-of-band activator is the recommended method.
3. **`ZIA_ACTIVATION=true`.** It fires an activation call after *every* create, update and delete. The activation endpoint allows 10 POST requests per minute and 40 per hour, so a configuration of meaningful size exhausts the budget, and the troubleshooting guide documents the result: the apply hangs, then fails with `/status/activate: context deadline exceeded`.

**Don't — no activation anywhere in the pipeline**

```hcl
resource "zia_url_filtering_rules" "block_gambling" {
  name           = "Block-Gambling-All-Users"
  order          = 1
  url_categories = ["GAMBLING"]
  # …
}
# Apply succeeds. Console shows pending changes. Nothing is enforced.
```

**Do — activation as its own step (build once: make build13 && sudo make ziaActivator)**

```sh
terraform apply && ziaActivator
```

**Also fine — in HCL, but only for a flat state, and depends_on must stay exhaustive**

```hcl
resource "zia_activation_status" "this" {
  status = "ACTIVE"
  depends_on = [
    zia_url_filtering_rules.block_gambling,
    zia_firewall_filtering_rule.allow_engineering,
    # every policy-affecting resource in this state — forget one and it stays pending
  ]
}
```

~> **Two further warnings** Do not mix manual console activation with a managed activation step. Pick one owner, or a human in the portal and a pipeline race each other.

Never `terraform state rm` the activation resource to "skip" activation. That removes Terraform's knowledge of it without changing anything in the tenant.

### Making activation structural

- **DO** make activation an explicit pipeline stage that always runs after the last apply
- **DO** fail or warn in CI when a plan touches `zia_*` but nothing will activate
- **DO** enforce it as policy (Checkov / OPA) rather than relying on reviewers
- **DON'T** rely on humans remembering to update `depends_on`

Activation should be **structural**, not remembered. A dedicated pipeline stage runs on every execution regardless of what changed, and it fails loudly and separately from the apply. Where a customer keeps activation in HCL, a CI gate is the fallback — and the durable version is a policy rule (OPA / Conftest) that denies any plan containing `zia_*` resources without an activation step.

**Do — activation as a stage, so it cannot be forgotten**

```yaml
- run: terraform apply -auto-approve tfplan
- name: Activate ZIA configuration
  run: ziaActivator
```

### API session timeout and automatic activation

- ZIA activates pending changes **when a session ends** — including on timeout, not just logout
- The API session timeout is **5–20 minutes, default 5**
- **DO** set it to the maximum of **20 minutes before using the provider** — the docs' *Before You Start* recommendation, not only for large applies
- A run longer than the timeout gets activated **part-way through**. From **v4.8.11** the provider then obtains a new session and retries automatically; earlier versions failed with `403 Resource Access Blocked`
- **DO** split large configurations across smaller states — the docs: "20 minutes is a ceiling, not a guarantee"
- **DON'T** raise the *UI* session timeout and expect it to help — it is a different field

#### The mechanism

Two reports that present as defects are explained by this: *"changes activated before the apply finished"* and *"we never ran an activation step, but it activated anyway."* Activation on session end is native platform behaviour; the docs state that it "cannot be disabled or overridden by the provider". A 30-minute apply against the default 5-minute timeout crosses that boundary several times, and each crossing publishes whatever was pending at that moment — including a half-written configuration.

#### Where to set it

Admin Portal: **Administration > Advanced Settings**, field **API Session Timeout Duration (In Minutes)** under *Admin Portal Session Timeout*.

![ZIA Admin Portal — Advanced Settings, Admin Portal Session Timeout](https://raw.githubusercontent.com/zscaler/terraform-provider-zia/master/docs/guides/media/advanced_settings.png)

~> **Two fields with similar names** "UI Session Timeout Duration" sits directly above "API Session Timeout Duration". Someone told to "raise the session timeout" will very plausibly change the top one and see no improvement. When triaging, confirm *which* field was edited.

**Or as code**

```hcl
resource "zia_advanced_settings" "this" {
  api_session_timeout = 20   # maximum; default is 5
  # … remaining advanced settings attributes
}
```

This is settable in the portal, not only through Terraform. Customers who do not manage advanced settings with Terraform can change the one field in the portal rather than adopting `zia_advanced_settings`.

### Rule ordering

1. `order` **must be** a positive whole number starting at 1 — zero and negatives are rejected at plan time
2. **Orders must be contiguous** — no gaps, including after a delete ([Deleting a rule](#deleting-a-rule))
3. **Start at 1 and number sequentially** — the provider places your rules around predefined rules automatically

- Applies to all **18** ZIA rule resources ([ZIA rule resources sharing ordering semantics](#zia-rule-resources-sharing-ordering-semantics))
- **DO** keep all rules of the **same type** in one Terraform configuration
- **DO** feel free to split **different** rule types across separate configurations, workspaces or states
- **DON'T** split a **single** rule type across two configurations — each would assert its own ordering for the same list
- **DON'T** leave rules that Terraform does not manage — or duplicate order values — inside the order range your HCL claims
- **DON'T** interrupt a run while rules are being placed — the docs warn that cancelling "can leave the rules partially reordered"

#### Why rule applies take longer

Every rule page carries the same note: "The ZIA API does not provide a way to set the position of several rules in one request, so after a rule is created or updated the provider sends additional requests to move the rules into the order declared in your configuration, and waits until all rules of this type are in place." Terraform shows the resource as still creating or modifying while this happens; the provider is working in the background.

#### Predefined rules and the starting order

The `zia_firewall_filtering_rule` page states that the provider **automatically handles predefined rules** during ordering: "you can still configure your rules starting at `order = 1`", and the provider places new rules at the correct position, reorders around predefined rules, and avoids drift. Do not offset your orders to make room for predefined rules, and avoid the gapped numbering that appears in some examples — a "default" rule at `order = 999`, or deny rules offset by +1000. The API closes those gaps.

#### Orders that cannot be satisfied (v4.8.7+)

When the configured orders cannot be satisfied — typically because rules created in the portal occupy positions inside the range your HCL claims, or because the tenant has duplicate order values — v4.8.7 and later stop after a bounded number of attempts and **complete the apply with a warning** naming the affected rules and their configured versus actual positions. The remaining differences appear in the next plan. The fix is the ownership rule: import the out-of-band rules into the configuration that owns that rule type, or move your range. v4.8.11 also stops waiting, and reports the API error, when the current order cannot be read.

#### Scope of the single-configuration rule

This is **not** "all ZIA rules must live in one state". `order` is scoped per rule type: URL filtering rules are ordered among themselves, firewall rules among themselves. The constraint is only that **one rule type has one owner** — two configurations computing positions for the same server-side list fight over them on every apply.

**Don't**

```hcl
order = 0     # rejected: expected order to be at least (1), got 0
order = 999   # gap after 1, 2 -> the API closes it -> drift
```

**Do**

```hcl
resource "zia_url_filtering_rules" "block_gambling" {
  name           = "Block-Gambling"
  order          = 1
  url_categories = ["GAMBLING"]     # required from v4.8.11 (["ANY"] for all)
}
resource "zia_url_filtering_rules" "block_streaming" {
  name           = "Block-Streaming"
  order          = 2                # contiguous, no gaps
  url_categories = ["STREAMING_MEDIA"]
}
```

### Predefined, default and custom rules

|   | Custom | Predefined | Default |
|---|---|---|---|
| **Create via Terraform** | Yes | No — already exists | No |
| **Update fields** | Yes | Limited — `description`, `state`, `order`, `labels` | No |
| **Change order** | Yes | **Yes** | No |
| **Import / manage** | Yes | Yes — for reordering | **Do not import or manage** |
| **terraform destroy** | Yes | **No** | **No** |

- **DO** keep predefined rules in HCL only if you want to control their position or the few updatable fields
- **DON'T** import or manage **default** rules — the docs: "Management of default rules are not supported via the API"
- **DO** release a predefined rule with a `removed {}` block and `lifecycle { destroy = false }` (Terraform 1.7+)
- **DON'T** `terraform destroy` a predefined rule
- **DON'T** simply delete the resource block — Terraform plans a destroy of the predefined rule

For `zia_firewall_filtering_rule`, the protection lives in the provider: before deleting a rule it checks whether the rule is predefined — by name (*Office 365 One Click Rule*, *UCaaS One Click Rule*, *Block All IPv6*, *Block malicious IPs and domains*, *Default Firewall Filtering Rule*) or by `predefined = true` — and returns `deletion of the predefined rule '…' is not allowed` without calling the API. The other rule resources have no such check in the provider; a destroy of a predefined rule there is sent to the API, which decides whether to reject it. Use `removed {}` for predefined rules of every rule type. Predefined-rule reorder once failed with `Request body is invalid` because read-only fields were echoed back; the provider now strips them. If a customer reports that error, upgrade first.

**Do — stop managing a predefined rule without touching it**

```hcl
removed {
  from = zia_firewall_filtering_rule.office365_one_click
  lifecycle { destroy = false }   # the rule stays in ZIA
}
```

On Terraform older than 1.7, `terraform state rm` is the fallback for a predefined rule — safe there because the rule is meant to stay, and never safe for a custom rule.

### Deleting a rule

- Deleting any rule other than the last leaves a gap that the API closes on its own: every rule below it moves up one
- **DO** renumber the remaining rules so the orders stay contiguous — the rule pages: "When you delete a rule, renumber the remaining rules so the orders stay contiguous"
- **DO** delete the rule and renumber the others **in the same apply** (v4.8.14+)
- **DON'T** delete and leave the gap — the configuration no longer reflects the tenant and every plan shows drift
- **DO** activate once, after the apply

#### What the provider does

Once the API has closed the gap, every rule below the deleted one is already at its renumbered position, so the order updates in the plan change nothing. From v4.8.14, a rule whose only change is `order` is not rewritten when it is already in place, an update of other attributes keeps the rule at its current position, and a write the API rejects because the rule list changed during the apply (`Rule is not allowed at order N`) is retried. This applies to all [rule resources sharing the ordering semantics](#zia-rule-resources-sharing-ordering-semantics), and also covers rules deleted outside Terraform during an apply.

**Do**

```hcl
# rule at order 12 removed from the configuration
# order = 13 -> 12, 14 -> 13, … (contiguous, no gaps)
```

```sh
terraform apply
ziaActivator
```

**Don't**

```hcl
# Deleted the rule at order 12, left the gap
order = 11
order = 13   # gap -> API closes it -> HCL no longer matches -> drift
```

#### Earlier versions

Before v4.8.14, a delete and the renumbering in one apply could fail with `INVALID_INPUT_ARGUMENT` `Rule is not allowed at order N`: each order update first moved its rule to a position read before the delete completed. Upgrade, or destroy the rule on its own with `-target` first and renumber the others in a second apply.

### Country codes

- **The rule: use ISO-3166 Alpha-2** (`"US"`, `"BR"`, `"DE"`) for country attributes — validated at plan time
- The provider adds the API's `COUNTRY_` prefix for you — never write `COUNTRY_US`. `"ANY"` and `"NONE"` are accepted as special values
- Applies to `dest_countries` / `source_countries` on the firewall, DNS, IPS, NAT, forwarding-control, traffic-capture and URL filtering rules; `countries` on `zia_firewall_filtering_destination_groups`; `blocked_countries` on `zia_advanced_threat_settings`
- **The exception is** `zia_location_management`, whose `country` takes the full uppercase name (`"UNITED_STATES"`) — and is **not** validated at plan time
- **DON'T** use three-letter codes (`"USA"`) anywhere

The troubleshooting guide documents both errors: `'AUC' is not a valid ISO-3166 Alpha-2 country code` on the rule resources, and `'UNITED_STATE' is not a valid country name` on `zia_location_management`. The difference is upstream: the policy APIs and the locations API use different country formats.

**Don't**

```hcl
dest_countries = ["USA"]          # not a valid ISO-3166 Alpha-2 country code
dest_countries = ["COUNTRY_US"]   # don't hand-write the internal form
```

**Do**

```hcl
dest_countries = ["US", "BR", "DE"]   # rule-based resources — ISO-3166 Alpha-2

resource "zia_location_management" "branch" {
  country = "UNITED_STATES"         # the one exception
}
```

### Detach before delete

- ZIA refuses to delete an object still referenced by a rule — `RESOURCE_IN_USE` / "associated with N rule(s)"
- **The provider detaches some objects for you** before deleting them — but each object type is detached only from the rule types listed below
- **DO** reference objects by attribute (`zia_url_categories.custom.id`) so Terraform destroys in a sane order — a rule that references the object is then updated or destroyed before the object itself
- **DO** remove the reference from the rule in the same configuration when the rule type is **not** covered below — `RESOURCE_IN_USE` on an uncovered rule type is expected API behaviour, not a provider bug
- **DON'T** report `RESOURCE_IN_USE` as a bug unless it comes from a rule type the table says is covered

| Object being deleted | Detached automatically from |
|---|---|
| `zia_url_categories` | URL filtering, SSL inspection, Sandbox, File Type Control, Firewall DNS and Firewall IPS rules; Bandwidth Classes; FTP control policy; DLP web rules |
| `zia_firewall_filtering_ip_source_groups`, `zia_firewall_filtering_destination_groups`, `zia_firewall_filtering_network_service`, `zia_firewall_filtering_network_service_groups`, `zia_firewall_filtering_network_application_groups`, `zia_dns_application_groups`, `zia_rule_labels`, `zia_location_management`, `zia_user_management` | Firewall filtering rules only |
| `zia_http_header_profile`, `zia_http_header_action_profile` | URL filtering rules only |

On delete the provider first discovers which rules of the covered types reference the object, updates only those rules to remove the reference, then deletes the object. A destroy of a widely referenced object therefore does more API work than its single delete implies — expected, not a fault.

References in any other rule type are not removed automatically — for example a rule label on a URL filtering rule, or an IP source group on a DNS, IPS or NAT rule. If such a reference exists, the delete fails with `RESOURCE_IN_USE`; remove the reference from that rule first.

From **v4.8.13**, a rule type the tenant is not subscribed to, or the admin role cannot access (`403`, such as `NOT_SUBSCRIBED`), is skipped during detach instead of failing the delete — nothing in an unavailable rule type can reference the object. Earlier versions fail the whole delete, for example deleting a URL category on a tenant without DLP.

### DUPLICATE_ITEM and data-source-only objects

- **`DUPLICATE_ITEM` causes:** the troubleshooting guide — "attempting to create a new service that already exists", most commonly one with the name of a predefined object; also a previous failed apply that created it without recording state
- **DO** import the existing object, or choose a unique name
- **DO** check whether a create that failed with an intermittent `401` created the object anyway — the troubleshooting guide: "If it was, import it with terraform import instead of letting Terraform create it again"
- **DON'T** loop `terraform apply` — the conflict repeats identically
- **DON'T** try to *create* IdP-driven or predefined objects — departments and groups are **data sources**
- **DON'T** create a DLP dictionary whose name contains spaces — the guide recommends cloning the predefined dictionary with underscores or dashes

**Don't**

```hcl
resource "zia_department_management" "eng" {   # not a resource
  name = "Engineering"
}
```

**Do**

```hcl
data "zia_department_management" "eng" {
  name = "Engineering"
}

# Adopt an object that already exists
terraform import zia_url_filtering_rules.existing 12345
```

### Concurrency: one Terraform process per tenant

**Splitting state does not split the tenant — many states, still one ZIA org.**

- **DO** split state for blast radius and ownership — the docs call it "fully supported" and "a sound practice"
- **DO** **apply one configuration at a time** per tenant, with a concurrency control keyed on the **tenant**, not the repo or workspace
- **DO** run plan in parallel — "read operations do not take the write lock"
- **DO** **activate once**, after the last apply — not once per configuration
- **DON'T** run several Terraform processes against the same tenant at the same time — the docs name parallel CI jobs, several workspaces, and Terragrunt `run-all` over many modules
- **DON'T** cancel an in-progress apply — the activator guide calls it "itself a cause of partially applied configuration"

#### The two governing behaviours

|   | Behaviour | Consequence |
|---|---|---|
| **Writes** | One tenant-wide write lock. Requests that cannot acquire it are rejected with `EDIT_LOCK_NOT_AVAILABLE`, `Resource Access Blocked` or `Failed during enter Org barrier`; the provider retries them | Contention shows up as a progressively slower run rather than an immediate failure |
| **Activation** | Publishes **all** pending tenant changes; is **queued** while another session has unactivated changes, and a queued activation cannot be cancelled | Your activation waits for the slowest run, then publishes its changes too — including a partially created rule set |

Put those together and activating per configuration does not mean "my changes are live now" — it means "my activation is queued behind whichever run finishes last". When the queue clears, everything pending goes live, including changes from a run that failed part-way through. For rule-based resources the activator guide warns that "a partially created rule set can take effect with incomplete ordering".

Each process also holds its own API session and sends its own requests, so concurrent runs compete for the same rate limits. The troubleshooting guide links running many processes with one API client — for example Terragrunt `run-all` — to intermittent `401` errors with an empty code and message.

**Don't**

```text
09:00:00  workspace: firewall-rules   apply  --+
09:00:00  workspace: url-filtering    apply  --|  all four contend for one write lock,
09:00:00  workspace: dlp              apply  --|  then each activates -> queued,
09:00:00  workspace: locations        apply  --+  uncancellable
```

**Do — serialise on the tenant**

```yaml
# GitHub Actions
concurrency:
  group: zia-apply-${{ inputs.tenant }}
  cancel-in-progress: false   # cancelling mid-apply leaves partial configuration

# GitLab CI
apply:
  resource_group: zia-prod
```

~> **A common incorrect remedy** Faced with lock contention, the instinct is `-parallelism=1`. The activator guide: "Do not lower Terraform's `-parallelism` in response to lock contention. It does not reduce the number of competing runs, and the longer run time increases the chance of crossing an API session boundary and triggering an unplanned activation."

If truly concurrent execution is required, a separate tenant is the only boundary that provides it — each tenant has its own write lock and activation queue.

**Debugging tip:** the Admin Portal's *Queued Activations* list names the administrators with pending activations. Super administrators can **Force Activate**, but that publishes *all* saved changes in the tenant, including any still mid-deployment.

### Troubleshooting and log collection

- **DO** follow the troubleshooting guide's order: check the typical problems, **upgrade the provider**, confirm the API Session Timeout is 20 minutes, confirm default parallelism and one process per tenant, then collect a log
- **DO** provide a debug log (`terraform.log`) of the failing run with every GitHub issue or support case — the docs say it "must be provided"
- **DO** **sanitize** the log before sharing it. From **v4.8.11** access tokens and session cookies are masked; earlier versions log the OAuth access token in clear text in every `Authorization: Bearer` header (valid for up to an hour). The log also contains tenant identifiers and object names
- **DON'T** attach an unsanitized log to a GitHub issue — issues are public. Use a formal Zscaler support case instead
- **DO** include the `x-oneapi-request-id` and `x-oneapi-host` response headers, and the UTC time of the failing request, when reporting an API error

**Enable logging before reproducing the problem**

```sh
export ZSCALER_SDK_LOG=true
export ZSCALER_SDK_VERBOSE=true
export TF_LOG="DEBUG"
export TF_LOG_PATH="terraform.log"
```

[ZIA error reference](#zia-error-reference) lists the errors documented in the troubleshooting guide — `403 Resource Access Blocked`, intermittent `401`, the activation timeout, stalled requests and `Rule is not allowed at order N` — with their causes and fixes.

### ZIA in the general practices

The general practices apply to ZIA without exception. These are the ZIA-specific points that sit on top of them, collected in one place.

#### State layout and execution

1. **All rules of one type belong to one configuration** ([Rule ordering](#rule-ordering)). Different rule types can live wherever you like.
2. **One Terraform process at a time per tenant** — however many states you have ([Concurrency: one Terraform process per tenant](#concurrency-one-terraform-process-per-tenant)).
3. **One activation per tenant, after the last apply** — not one per state ([Concurrency: one Terraform process per tenant](#concurrency-one-terraform-process-per-tenant)).

#### Pipelines and tests

- **DO** add `ziaActivator` as the final stage of the pipeline, and key the concurrency group on the tenant
- **DON'T** set `ZIA_ACTIVATION=true` in a pipeline or a test suite ([Activation](#activation))
- **DO** set `api_session_timeout = 20` before the first run ([API session timeout and automatic activation](#api-session-timeout-and-automatic-activation))
- **DO** include `ziaActivator` in test teardown — otherwise deletes stay as drafts
- **DO** keep Terraform's default `-parallelism` ([Parallelism and provider tuning](#parallelism-and-provider-tuning))

#### Versions and import

- **DO** pin `zia = "~> 4.0"` (latest v4.8.13). Fixes worth an upgrade: v4.8.7 (rule-order reconciliation no longer runs for hours), v4.8.10 (`request_timeout` up to 1800), v4.8.11 (new session on timeout, masked tokens, `url_categories` required, DLP and DNS-rule drift fixes), v4.8.12 (DLP updates fail fast with the API error), v4.8.13 (detach before delete skips rule types the tenant is not subscribed to)
- **DO** use `skip_credentials_validation` (v4.8.7+) only where every `zia_*` resource is disabled
- **DO** adopt brownfield tenants with `zscaler-terraformer import --resources "zia"`, one resource type per PR
- **DO** re-supply `pre_shared_key` and admin `password` values after an import — the API never returns them

**Single-object import**

```sh
terraform import zia_url_filtering_rules.block_gambling 12345
terraform import zia_location_management.hq 67890
```

## ZIA error reference

Rows marked † are quoted from the troubleshooting guide.

| Error | Likely cause | Fix |
|---|---|---|
| Apply succeeds, policy not enforced | No activation, or stale `depends_on` | Run `ziaActivator` after apply; or extend `zia_activation_status` |
| Changes activated part-way through a run | The API session ended mid-apply | Set `api_session_timeout` to 20; split the configuration |
| Activation happened with no activation step | Activation fires when a session ends | Expected platform behaviour, not a bug |
| † `403 Resource Access Blocked` from the first request | API client role lacks permission or functional scope, or a subscription is missing | Review the API client's role in the Admin Portal |
| † `403 Resource Access Blocked` several minutes in, then every request fails | The API session ended (API Session Timeout) | Upgrade to v4.8.11+ (new session, automatic retry); set the timeout to 20 |
| † Intermittent `401` with empty `code` and `message` | Returned by the API mid-run; many processes with one API client | Support case with `x-oneapi-request-id`; avoid concurrent runs; import objects the failed create made |
| † Apply hangs, then `/status/activate: context deadline exceeded` | `ZIA_ACTIVATION=true` exhausting 10/min, 40/hr | Unset it; activate once at the end |
| † A single request takes many minutes before the run fails | The API is not answering that request | Set `request_timeout` (for example 300) to fail sooner |
| † `INVALID_INPUT_ARGUMENT` `Rule is not allowed at order N` | Before v4.8.14: a rule deleted and orders changed in the same apply | Upgrade to v4.8.14 or later; on earlier versions, delete with `-target` first, then re-order in a separate apply |
| `EDIT_LOCK_NOT_AVAILABLE` / `Failed during enter Org barrier` (409) | Concurrent writes to one tenant — not a rate limit | One Terraform process per tenant; don't lower `-parallelism` |
| Activation ran but nothing went live | Another session still editing — activation **queued** | Wait; check *Queued Activations* in the portal |
| Apply completes with a warning naming rules and "configured vs actual" positions | Unmanaged rules in your order range, or duplicate orders (v4.8.7+) | Import the out-of-band rules or move your range; re-plan |
| `expected order to be at least (1), got 0` | `order = 0` or negative | Set `order >= 1` |
| `Request body is invalid` on predefined reorder | Old provider echoing read-only fields | Upgrade the provider |
| † `deletion of the predefined rule '…' is not allowed` | Destroy — or deleting the block — on a predefined `zia_firewall_filtering_rule` | `removed {}` with `destroy = false` |
| † `DUPLICATE_ITEM` | Name collision, or the object exists untracked | Import, or rename |
| † `no dictionary found with name: …` | DLP dictionary name with spaces | Clone the predefined dictionary; use underscores or dashes |
| `RESOURCE_IN_USE` / "associated with N rule(s)" | Object referenced by a rule type the provider does not detach from ([Detach before delete](#detach-before-delete)) | Remove the reference from that rule first; report it only if the rule type is listed as covered |
| Deleting an object fails while reading a rule type the tenant is not subscribed to (`403`, `NOT_SUBSCRIBED`) | Detach sweep reading an unavailable rule type (before v4.8.13) | Upgrade to v4.8.13+ — unavailable rule types are skipped |
| † `'AUC' is not a valid ISO-3166 Alpha-2 country code` | Three-letter or invalid code in a rule | Use Alpha-2 (`"US"`) |
| † `'UNITED_STATE' is not a valid country name` | Wrong name on `zia_location_management` | Use `"UNITED_STATES"` |
| `url_categories` required at plan time | v4.8.11+ rejects URL filtering rules without categories | Set `url_categories` (`["ANY"]` for all) |
| `file_types` and `file_type_categories` both set | v4.8.11+ rejects both on `zia_dlp_web_rules` | Set one; the other is computed |
| Intermittent 429s during bulk apply | Rate limiting | None needed — retries are automatic |
| Bulk rule apply takes many hours | `-parallelism` was lowered | Remove the flag; use the default |
| `request_timeout` above 300 rejected | Provider below v4.8.10 | Upgrade |
| Warning: attribute `parallelism` is deprecated | Old provider block | Delete the argument |
| `Cloud zscalergov not supported for OneAPI` | Legacy cloud name in `zscaler_cloud` | Use `gov` / `govus` (v4.7.25+) |
| One attribute shows a change on every plan | API response omits or empties that field | Upgrade; raise a backend bug; scope `ignore_changes` meanwhile |
| Secret field empty after import | Write-only at the API | Re-supply via variable |

## ZIA rule resources sharing ordering semantics

All 18 resources validate `order >= 1` at plan time, carry the contiguous-order note on their pages, and follow [Rule ordering](#rule-ordering) and [Deleting a rule](#deleting-a-rule).

| Group | Resources |
|---|---|
| **Web & URL** | `zia_url_filtering_rules` · `zia_file_type_control_rules` · `zia_ssl_inspection_rules` · `zia_cloud_app_control_rule` |
| **Firewall** | `zia_firewall_filtering_rule` · `zia_firewall_dns_rule` · `zia_firewall_ips_rule` · `zia_nat_control_rules` |
| **Data protection** | `zia_dlp_web_rules` · `zia_casb_dlp_rules` · `zia_casb_malware_rules` · `zia_endpoint_dlp_rules` · `zia_endpoint_dlp_sub_rules` · `zia_outbound_email_dlp` |
| **Traffic & other** | `zia_forwarding_control_rule` · `zia_sandbox_rules` · `zia_bandwidth_control_rule` · `zia_traffic_capture_rules` |

~> **Plural and singular names are not a typo** `zia_url_filtering_rules` is plural; `zia_firewall_filtering_rule`, `zia_firewall_dns_rule`, `zia_firewall_ips_rule` and `zia_bandwidth_control_rule` are singular. Verify every name against the Terraform Registry before use.

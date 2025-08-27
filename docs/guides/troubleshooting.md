---
page_title: "Troubleshooting Guide"
---

# How to troubleshoot your problem

If you have problems with code that uses ZPA Terraform provider, follow these steps to solve them:

* Check symptoms and solutions in the [Typical problems](#typical-problems) section below.
* Upgrade provider to the latest version. The bug might have already been fixed.
* In case of authentication problems, see the [Authentication Issues](#authentication-issues) below.
* Collect debug information using following command:

```sh
TF_LOG=DEBUG ZSCALER_SDK_VERBOSE=true ZSCALER_SDK_LOG=true terraform apply -no-color 2>&1 |tee tf-debug.log
```

* Open a [new GitHub issue](https://github.com/zscaler/terraform-provider-zia/issues/new/choose) providing all information described in the issue template - debug logs, your Terraform code, Terraform & plugin versions, etc.

## Typical problems

### Authentication Issues

### │ Error: Invalid provider configuration and Error: failed configuring the provider

The most common problem with invalid provider is when the ZIA API credentials are not properly set via one of the supported methods. Please make sure to read the documentation for the supported authentication methods [Authentication Methods](https://registry.terraform.io/providers/zscaler/zia/latest/docs)

```sh
│ Provider "zscaler/zia" requires explicit configuration. Add a provider block to the root module and configure the
│ provider's required arguments as described in the provider documentation.
```

```sh
│ Error: expected zia_cloud to be one of ["zscaler" "zscalerone" "zscalertwo" "zscalerthree" "zscloud" "zscalerbeta" "zscalergov" "zscalerten" "zspreview"], got
│
│   with provider["zscaler.com/zia/zia"],
│   on <input-prompt> line 1:
│   (source code not available)
│
```

## Multiple Provider Configurations

The most common reason for technical difficulties might be related to missing `alias` attribute in `provider "zpa" {}` blocks or `provider` attribute in `resource "zia_..." {}` blocks, when using multiple provider configurations. Please make sure to read [`alias`: Multiple Provider Configurations](https://www.terraform.io/docs/language/providers/configuration.html#alias-multiple-provider-configurations) documentation article.

## Error while installing: registry does not have a provider

```sh
Error while installing hashicorp/zia: provider registry
registry.terraform.io does not have a provider named
registry.terraform.io/hashicorp/zia
```

If you notice below error, it might be due to the fact that [required_providers](https://www.terraform.io/docs/language/providers/requirements.html#requiring-providers) block is not defined in *every module*, that uses ZIA Terraform Provider. Create `versions.tf` file with the following contents:

```hcl
# versions.tf
terraform {
  required_providers {
    zpa = {
      source  = "zscaler/zia"
      version = "2.6.2"
    }
  }
}
```

... and copy the file in every module in your codebase. Our recommendation is to skip the `version` field for `versions.tf` file on module level, and keep it only on the environment level.

```
├── environments
│   ├── sandbox
│   │   ├── README.md
│   │   ├── main.tf
│   │   └── versions.tf
│   └── production
│       ├── README.md
│       ├── main.tf
│       └── versions.tf
└── modules
    ├── first-module
    │   ├── ...
    │   └── versions.tf
    └── second-module
        ├── ...
        └── versions.tf
```

## Error: Failed to install provider

Running the `terraform init` command, you may see `Failed to install provider` error if you didn't check-in [`.terraform.lock.hcl`](https://www.terraform.io/language/files/dependency-lock#lock-file-location) to the source code version control:

```sh
Error: Failed to install provider

Error while installing zscaler/zia: v2.6.2: checksum list has no SHA-256 hash for "https://github.com/zscaler/terraform-provider-zia/releases/download/v2.6.2/terraform-provider-zia_2.6.2_darwin_amd64.zip"
```

You can fix it by following three simple steps:

* Replace `zscaler.com/zia/zia` with `zscaler/zia` in all your `.tf` files with the `python3 -c "$(curl -Ls https://github.com/zscaler/terraform-provider-zia/scripts/upgrade-namespace.py)"` command.
* Run the `terraform state replace-provider zscaler.com/zia/zia zscaler/zia` command and approve the changes. See [Terraform CLI](https://www.terraform.io/cli/commands/state/replace-provider) docs for more information.
* Run `terraform init` to verify everything working.

The terraform apply command should work as expected now.

## Error: Failed to query available provider packages

See the same steps as in [Error: Failed to install provider](#error-failed-to-install-provider).

### Error: Provider registry.terraform.io/zscaler/zia v... does not have a package available for your current platform, windows_386

This kind of errors happens when the 32-bit version of ZIA Terraform provider is used, usually on Microsoft Windows. To fix the issue you need to switch to use of the 64-bit versions of Terraform and ZIA Terraform provider.

### Error: failed configuring the provided

This kind of error happens when the administrator fails to configure the ZIA API credentials via one of the accepted methods such as environment variables, hard-coded method (which is discouraged) or via the `credentials.json` file.

```sh
│   with provider["registry.terraform.io/zscaler/zia"],
│   on zia_location_management.tf line 10, in provider "zia":
│   10: provider "zia" {}
│
│ error:Could not open credentials file, needs to contain one json object with keys: zia_username, zia_password, zia_api_key, and
│ zia_cloud. open /Users/<username>/.zia/credentials.json: no such file or directory
```

### Cloud Firewall `zia_firewall_filtering_rule` Error: 'AUC' is not a valid ISO-3166 Alpha-2 country code

This type of error happens when the administrator fails to provide a valid two letter country code value as part of the attribute `dest_countries`

```sh
╷
│ Error: 'AUC' is not a valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes
```

### Location Management `zia_location_management` Error: 'AUC' is not a valid ISO-3166-1 country name

This type of error happens when the administrator fails to provide a valid uppercase country name value as part of the attribute `country`

```sh
│ Error: 'UNITED_STATE' is not a valid country name. Please refer to ISO 3166-1 for a list of valid country names
```

### │ Error: FAILED: POST, https://zsapi.***.net/api/v1/networkServices, 400, 400 , {"code":"DUPLICATE_ITEM","message":"DUPLICATE_ITEM"}, api responded with code: 400

This type of error happens when the administrator is attempting to create a new service that already exists. This is mostly common when attempting to create a service with the same name of a predefined service or other resource that have been previously created in the tenant.

```sh
│ Error: FAILED: POST, https://zsapi.***.net/api/v1/networkServices, 400, 400 , {"code":"DUPLICATE_ITEM","message":"DUPLICATE_ITEM"}, api responded with code: 400
```

### │ Error: no dictionary found with name: Social Security Numbers (US)

This error is commonly returned when attempting to create a `zia_dlp_dictionaries` where the name contains spaces. To prevent this from happening we recommend to clone the existing predefined DLP dictionary, and provide a name containing underscores or dashes.

```sh
│ Error: no dictionary found with name: Social Security Numbers (US)
```

### │ Error: deletion of the predefined rule i.e 'Office 365 One Click Rule' is not allowed

This error occurs when attempting to delete a predefined firewall filtering rule. Predefined rules such as "Office 365 One Click Rule", "UCaaS One Click Rule", "Block All IPv6", "Block malicious IPs and domains", and "Default Firewall Filtering Rule" cannot be deleted as they are system-managed rules.

**Solution**: Remove the rule from your Terraform configuration and run `terraform apply` instead of `terraform destroy`. The rule will remain in the ZIA system but will no longer be managed by Terraform.

```sh
│ Error: deletion of the predefined rule 'Office 365 One Click Rule' is not allowed
```

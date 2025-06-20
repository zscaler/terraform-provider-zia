---
subcategory: "Subscription Alerts"
layout: "zscaler"
page_title: "ZIA: subscription_alert"
description: |-
  Official documentation https://help.zscaler.com/zia/about-alert-subscriptions
  API documentation https://help.zscaler.com/zia/alerts#/alertSubscriptions-get
  Creates and manages Alert Subscriptions.
---

# zia_subscription_alert (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-alert-subscriptions)
* [API documentation](https://help.zscaler.com/zia/alerts#/alertSubscriptions-get)

Use the **zia_subscription_alert** resource allows the creation and management of Alert Subscriptions in the Zscaler Internet Access.

## Example Usage

```hcl
resource "zia_subscription_alert" "this" {
  email  = "alert@acme.com"
  description = "Terraform Alert"
  pt0_severities = ["CRITICAL"]
  secure_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
  manage_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
  comply_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
  system_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email address of the alert recipient

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `description` - (String) Additional comments or information about the alert subscription
* `pt0_severities` - (List of String) Lists the severity levels of the Patient 0 Severity Alert class information that the recipient receives. Supported Values: `CRITICAL`, `MAJOR`, `MINOR`, `INFO`, `DEBUG`
* `secure_severities` - (List of String) Lists the severity levels of the Secure Severity Alert class information that the recipient receives. Supported Values: `CRITICAL`, `MAJOR`, `MINOR`, `INFO`, `DEBUG`
* `manage_severities` - (List of String) Lists the severity levels of the Manage Severity Alert class information that the recipient receives. Supported Values: `CRITICAL`, `MAJOR`, `MINOR`, `INFO`, `DEBUG`
* `comply_severities` - (List of String) Lists the severity levels of the Comply Severity Alert class information that the recipient receives. Supported Values: `CRITICAL`, `MAJOR`, `MINOR`, `INFO`, `DEBUG`
* `system_severities` - (List of String) Lists the severity levels of the System Severity Alert class information that the recipient receives. Supported Values: `CRITICAL`, `MAJOR`, `MINOR`, `INFO`, `DEBUG`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_subscription_alert** can be imported by using `<ALERT_ID>` or `<ALERT_EMAIL>` as the import ID.

For example:

```shell
terraform import zia_subscription_alert.example <alert_id>
```

or

```shell
terraform import zia_subscription_alert.example <alert_email>
```

---
subcategory: "Subscription Alerts"
layout: "zscaler"
page_title: "ZIA: subscription_alert"
description: |-
  Official documentation https://help.zscaler.com/zia/about-alert-subscriptions
  API documentation https://help.zscaler.com/zia/alerts#/alertSubscriptions-get
  Get information about subscription alerts details.
---

# zia_subscription_alert (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-alert-subscriptions)
* [API documentation](https://help.zscaler.com/zia/alerts#/alertSubscriptions-get)

Use the **zia_subscription_alert** data source to get information about a subscription alert resource in the Zscaler Internet Access cloud or via the API.

## Example Usage - Via Email

```hcl
data "zia_subscription_alert" "this" {
    email = "alert@acme.com"
}
```

## Example Usage - Via ID

```hcl
data "zia_subscription_alert" "this" {
    id = 3271
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Optional) The name of the subscription alert to be exported.
* `id` - (Optional) The unique identifer for the subscription alert.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Additional comments or information about the alert subscription
* `pt0_severities` - (List of String) Lists the severity levels of the Patient 0 Severity Alert class information that the recipient receives
* `secure_severities` - (List of String) Lists the severity levels of the Secure Severity Alert class information that the recipient receives
* `manage_severities` - (List of String) Lists the severity levels of the Manage Severity Alert class information that the recipient receives
* `comply_severities` - (List of String) Lists the severity levels of the Comply Severity Alert class information that the recipient receives
* `system_severities` - (List of String) Lists the severity levels of the System Severity Alert class information that the recipient receives
* `deleted` - (bool) Deletes an existing alert subscription

---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZIA Config Activation"
subcategory: "Activation"
---

# Activation Overview

As of right now, Terraform does not provide native support for commits or post-activation configuration, so configuration and policy activations are handled out-of-band.

In order to handle the activation as part of the provider, a separate source code have been developed to generate a CLI binary.

The activation cli source code can be compiled along with the provider binary and put somewhere in your `$PATH` (such as `$HOME/bin`). By default, the activation cli binary will be installed in the following path: `/usr/local/bin/`. To make it easier, the activation cli binary can be generated by using the `make` command as showed below:

```bash
$ make build13 && sudo make ziaActivator
```

~> You may or may not have to use `sudo` in order to successfully install the `ziaActivator` cli in the path.

Finally, you can invoke this binary after `terraform apply` or `terraform destroy` or invoke all cli commands together.

```bash
$ terraform init && terraform apply && ziaActivator
```

The authentication credentials can be given multiple ways, and if all are present then this is the order, from highest to lowest priority:

!> **WARNING:** Providing authentication credentials via CLI argument is insecure and
is not recommended.

1. CLI arguments
2. Environment variables
3. JSON authentication credential file

Refer to the ZIA provider argument reference documentation for more information on the JSON config file and the environment variables that are used.

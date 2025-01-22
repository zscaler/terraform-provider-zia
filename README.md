[![Release](https://github.com/zscaler/terraform-provider-zia/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/zscaler/terraform-provider-zia/actions/workflows/release.yml)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zscaler/terraform-provider-zia)](https://github.com/zscaler/terraform-provider-zia/v3/blob/master/.go-version)
[![Go Report Card](https://goreportcard.com/badge/github.com/zscaler/terraform-provider-zia)](https://goreportcard.com/report/github.com/zscaler/terraform-provider-zia)
[![codecov](https://codecov.io/gh/zscaler/terraform-provider-zia/graph/badge.svg?token=A9J4AJS7F5)](https://codecov.io/gh/zscaler/terraform-provider-zia)
[![License](https://img.shields.io/github/license/zscaler/terraform-provider-zia?color=blue)](https://github.com/zscaler/terraform-provider-zia/v3/blob/master/LICENSE)
[![Zscaler Community](https://img.shields.io/badge/zscaler-community-blue)](https://community.zscaler.com/)
[![Slack](https://img.shields.io/badge/Join%20Our%20Community-Slack-blue)](https://forms.gle/3iMJvVmJDvmUy36q9)

<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-text.svg" alt="Terraform logo" title="Terraform" height="50" />
</a>

<a href="https://www.zscaler.com/">
    <img src="https://raw.githubusercontent.com/zscaler/zscaler-terraformer/master/images/zscaler_terraformer-logo.svg" alt="Zscaler logo" title="Zscaler" height="50" />
</a>

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](docs/guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](docs/guides/troubleshooting.md) for guidance on typical problems.

## Terraform Provider for ☁️Zscaler Internet Access (ZIA)☁️

The ZIA provider is a Terraform plugin that allows for the full lifecycle management of Zscaler Internet Access resources.

- Website: [https://www.terraform.io](https://registry.terraform.io/providers/zscaler/zia/latest)
- Documentation: <https://help.zscaler.com/zia>
- Zscaler Community: [Zscaler Community](https://community.zscaler.com/)

## Examples

All the resources and data sources have [one or more examples](./examples) to give you an idea of how to use this
provider to build your own Zscaler Internet Access configuration. Provider's official documentation is located in the
[official terraform registry](https://registry.terraform.io/providers/zscaler/zia/latest/docs).

## Requirements

- Install [Terraform](https://www.terraform.io/downloads.html) 0.14.0 or newer (to run acceptance tests)
- [Go](https://golang.org/doc/install) (to build the provider plugin)
- Create a directory, go, follow this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to edit ~/.bash_profile to setup the GOPATH environment variable

## Upgrade

If you have been using version 3.x of the ZIA Terraform Provider, please upgrade to the latest version to take advantage of
all the new features, fixes, and functionality. Please refer to this [Upgrade Guide](https://github.com/zscaler/terraform-provider-zia/issues/1338)
for guidance on how to upgrade to version 4.x. Also, please check our [Releases](https://github.com/zscaler/terraform-provider-zia/releases) page for more details on major, minor, and patch updates.

## Quick Start

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please
check the [requirements](#requirements) before proceeding).

_Note:_ This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside
your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your
home directory outside the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:zscaler/terraform-provider-zia.git
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```sh
$ make tools
```

To compile the provider, run `make build13`. This will build the provider and put the provider binary in the `$GOPATH/bin`
directory.

```sh
$ make build13
...
$ $GOPATH/bin/terraform-provider-zia
...
```

## Testing the Provider

In order to test the provider, you can run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources.

```sh
$ make testacc
```

## Using the Provider

To use a released provider in your Terraform environment,
run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the
provider. To specify a particular provider version when installing released providers, see
the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions)
.

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build
instructions above), follow the instructions
to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). After placing the
custom-built provider into your plugins' directory, run `terraform init` to initialize it.

For either installation method, documentation about the provider specific configuration options can be found on
the [provider's website](https://registry.terraform.io/providers/zscaler/zia/latest/docs).

## Contributing

Terraform is the work of thousands of contributors. We really appreciate your help!

We have these minimum requirements for source code contributions.

Bug fix pull requests must include:

- [Terraform Plugin Acceptance Tests](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests).

Pull requests with new resources and data sources must include:

- Make API calls with the [zscaler-sdk-go v3](https://github.com/zscaler/zscaler-sdk-go) client
- Include [Terraform Plugin Acceptance Tests](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests)

Issues on GitHub are intended to be related to the bugs or feature requests with provider codebase.
See [Plugin SDK Community](https://www.terraform.io/community)
and [Discuss forum](https://discuss.hashicorp.com/c/terraform-providers/31/none) for a list of community resources to
ask questions about Terraform.

License
=========

=======

Copyright (c) 2022 [Zscaler](https://github.com/zscaler)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

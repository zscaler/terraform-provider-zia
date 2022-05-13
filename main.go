package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zscaler/terraform-provider-zia/zia"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zia.Provider,
	})
}

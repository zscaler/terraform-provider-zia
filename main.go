package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/willguibr/terraform-provider-zia/zia"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zia.Provider,
	})
}

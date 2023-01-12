package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zscaler/terraform-provider-zia/v2/zia"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zia.Provider})
}

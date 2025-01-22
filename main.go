package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zscaler/terraform-provider-zia/v4/zia"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(common.Version())
		return
	}
	var debug bool
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		debug = true
	}
	log.Printf(`ZPA Terraform Provider

Version %s

https://registry.terraform.io/providers/zscaler/zia/latest/docs

`, common.Version())
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zia.ZIAProvider,
		ProviderAddr: "registry.terraform.io/zscaler/zia",
		Debug:        debug,
	})
}

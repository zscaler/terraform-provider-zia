package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/zscaler/terraform-provider-zia/v2/zia"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/activation"
)

func getEnvVarOrFail(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	log.Fatalf("[ERROR] Couldn't find environment variable %s\n", k)
	return ""
}

func main() {
	log.Printf("[INFO] Initializing ZIA client\n")
	c := zia.Config{
		Username:   getEnvVarOrFail("ZIA_USERNAME"),
		Password:   getEnvVarOrFail("ZIA_PASSWORD"),
		APIKey:     getEnvVarOrFail("ZIA_API_KEY"),
		ZIABaseURL: getEnvVarOrFail("ZIA_CLOUD"),
		UserAgent:  fmt.Sprintf("(%s %s) cli/ziaActivator", runtime.GOOS, runtime.GOARCH),
	}
	cli, err := client.NewClient(c.Username, c.Password, c.APIKey, c.ZIABaseURL, c.UserAgent)
	if err != nil {
		log.Fatalf("[ERROR] Failed Initializing ZIA client: %v\n", err)
	}
	service := services.New(cli)
	resp, err := activation.CreateActivation(service, activation.Activation{
		Status: "active",
	})
	if err != nil {
		log.Printf("[ERROR] Activation Failed: %v\n", err)
	} else {
		log.Printf("[INFO] Activation succeded: %#v\n", resp)
	}
	log.Printf("[INFO] Destroying session: %#v\n", resp)
	_ = cli.Logout()
	os.Exit(0)
}

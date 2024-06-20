package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

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

	// Here, rather than setting up the client configuration from the external library,
	// we'll simply gather the required details for initializing the client
	username := getEnvVarOrFail("ZIA_USERNAME")
	password := getEnvVarOrFail("ZIA_PASSWORD")
	apiKey := getEnvVarOrFail("ZIA_API_KEY")
	ziaCloud := getEnvVarOrFail("ZIA_CLOUD")
	userAgent := fmt.Sprintf("(%s %s) cli/ziaActivator", runtime.GOOS, runtime.GOARCH)

	// Now, we'll use the local SDK's NewClient method to get the client instance
	cli, err := client.NewClient(username, password, apiKey, ziaCloud, userAgent)
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
		log.Printf("[INFO] Activation succeeded: %#v\n", resp)
	}

	log.Printf("[INFO] Destroying session: %#v\n", resp)
	_ = cli.Logout()
	os.Exit(0)
}

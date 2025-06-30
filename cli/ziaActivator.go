package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
)

func getEnvVarOrFail(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	log.Fatalf("[ERROR] Couldn't find environment variable %s\n", k)
	return ""
}

func main() {
	log.Printf("[INFO] Initializing ZIA activation client")

	useLegacy := strings.ToLower(os.Getenv("ZSCALER_USE_LEGACY_CLIENT")) == "true"

	var (
		service *zscaler.Service
		err     error
	)

	if useLegacy {
		log.Printf("[INFO] Using Legacy Client mode")

		username := getEnvVarOrFail("ZIA_USERNAME")
		password := getEnvVarOrFail("ZIA_PASSWORD")
		apiKey := getEnvVarOrFail("ZIA_API_KEY")
		cloud := getEnvVarOrFail("ZIA_CLOUD")

		ziaCfg, err := zia.NewConfiguration(
			zia.WithZiaUsername(username),
			zia.WithZiaPassword(password),
			zia.WithZiaAPIKey(apiKey),
			zia.WithZiaCloud(cloud),
			zia.WithUserAgent(fmt.Sprintf("(%s %s) cli/ziaActivator", runtime.GOOS, runtime.GOARCH)),
		)
		if err != nil {
			log.Fatalf("Error creating ZIA configuration: %v", err)
		}

		service, err = zscaler.NewLegacyZiaClient(ziaCfg)
		if err != nil {
			log.Fatalf("Error creating ZIA legacy client: %v", err)
		}
	} else {
		log.Printf("[INFO] Using OneAPI Client mode")

		clientID := getEnvVarOrFail("ZSCALER_CLIENT_ID")
		clientSecret := getEnvVarOrFail("ZSCALER_CLIENT_SECRET")
		vanityDomain := getEnvVarOrFail("ZSCALER_VANITY_DOMAIN")
		cloud := getEnvVarOrFail("ZSCALER_CLOUD")

		cfg, err := zscaler.NewConfiguration(
			zscaler.WithClientID(clientID),
			zscaler.WithClientSecret(clientSecret),
			zscaler.WithVanityDomain(vanityDomain),
			zscaler.WithZscalerCloud(cloud),
			zscaler.WithUserAgentExtra(fmt.Sprintf("(%s %s) cli/ziaActivator", runtime.GOOS, runtime.GOARCH)),
		)
		if err != nil {
			log.Fatalf("[ERROR] Failed to build OneAPI configuration: %v", err)
		}

		service, err = zscaler.NewOneAPIClient(cfg)
		if err != nil {
			log.Fatalf("[ERROR] Failed to initialize OneAPI client: %v", err)
		}
	}

	ctx := context.Background()

	resp, err := activation.CreateActivation(ctx, service, activation.Activation{
		Status: "ACTIVE",
	})
	if err != nil {
		log.Fatalf("[ERROR] Activation Failed: %v", err)
	}

	log.Printf("[INFO] Activation succeeded: %#v\n", resp)

	// Perform logout if using Legacy Client
	if useLegacy && service.LegacyClient != nil && service.LegacyClient.ZiaClient != nil {
		log.Printf("[INFO] Destroying session...\n")
		if err := service.LegacyClient.ZiaClient.Logout(ctx); err != nil {
			log.Printf("[WARN] Logout failed: %v\n", err)
		}
	}
}

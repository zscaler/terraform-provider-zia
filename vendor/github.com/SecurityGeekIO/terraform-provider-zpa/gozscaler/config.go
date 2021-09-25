package gozscaler

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	defaultBaseURL = "https://config.private.zscaler.com"
	defaultTimeout = 240 * time.Second
	loggerPrefix   = "zpa-logger: "
)

// Config contains all the configuration data for the API client
type Config struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	// The logger writer interface to write logging messages to. Defaults to standard out.
	Logger *log.Logger
	// Credentials for basic authentication.
	ClientID, ClientSecret, CustomerID string
}

/*
NewConfig returns a default configuration for the client.
By default it will try to read the access and te secret from the environment variable.
*/

// TODO Add healthCheck method to NewConfig
func NewConfig(clientID, clientSecret, customerID, rawUrl string) (*Config, error) {
	if clientID == "" || clientSecret == "" || customerID == "" {
		clientID = os.Getenv("ZPA_CLIENT_ID")
		clientSecret = os.Getenv("ZPA_CLIENT_SECRET")
		customerID = os.Getenv("ZPA_CUSTOMER_ID")
	}
	if rawUrl == "" {
		rawUrl = defaultBaseURL
	}

	var logger *log.Logger
	if loggerEnv := os.Getenv("ZSCALER_SDK_LOG"); loggerEnv == "true" {
		logger = getDefaultLogger()
	}

	baseURL, err := url.Parse(rawUrl)
	return &Config{
		BaseURL:      baseURL,
		HTTPClient:   getDefaultHTTPClient(),
		Logger:       logger,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		CustomerID:   customerID,
	}, err
}

func getDefaultHTTPClient() *http.Client {
	return &http.Client{Timeout: defaultTimeout}
}

func getDefaultLogger() *log.Logger {
	return log.New(os.Stdout, loggerPrefix, log.LstdFlags|log.Lshortfile)
}

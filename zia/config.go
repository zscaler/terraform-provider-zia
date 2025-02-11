package zia

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
)

type (
	// Config contains our provider schema values and Zscaler clients.
	Config struct {
		clientID           string
		clientSecret       string
		vanityDomain       string
		cloud              string
		sandboxToken       string
		sandboxCloud       string
		privateKey         string
		httpProxy          string
		retryCount         int
		parallelism        int
		backoff            bool
		minWait            int
		maxWait            int
		logLevel           int
		requestTimeout     int
		useLegacyClient    bool
		zscalerSDKClientV3 *zscaler.Client
		logger             hclog.Logger
		TerraformVersion   string // New field for Terraform version
		ProviderVersion    string // New field for Provider version

		// Options for Legacy V2 SDK
		Username   string
		Password   string
		APIKey     string
		ZIABaseURL string
		UserAgent  string
	}
)

type Client struct {
	Service *zscaler.Service
}

func NewConfig(d *schema.ResourceData) *Config {
	// defaults
	config := Config{
		backoff:        true,
		minWait:        30,
		maxWait:        300,
		retryCount:     5,
		parallelism:    1,
		logLevel:       int(hclog.Error),
		requestTimeout: 0,
	}
	logLevel := hclog.Level(config.logLevel)
	if os.Getenv("TF_LOG") != "" {
		logLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	config.logger = hclog.New(&hclog.LoggerOptions{
		Level:      logLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})

	if val, ok := d.GetOk("use_legacy_client"); ok {
		config.useLegacyClient = val.(bool)
	} else if os.Getenv("ZSCALER_USE_LEGACY_CLIENT") != "" {
		config.useLegacyClient = strings.ToLower(os.Getenv("ZSCALER_USE_LEGACY_CLIENT")) == "true"
	}

	if val, ok := d.GetOk("client_id"); ok {
		config.clientID = val.(string)
	}
	if config.clientID == "" && os.Getenv("ZSCALER_CLIENT_ID") != "" {
		config.clientID = os.Getenv("ZSCALER_CLIENT_ID")
	}

	if val, ok := d.GetOk("client_secret"); ok {
		config.clientSecret = val.(string)
	}
	if config.clientSecret == "" && os.Getenv("ZSCALER_CLIENT_SECRET") != "" {
		config.clientSecret = os.Getenv("ZSCALER_CLIENT_SECRET")
	}

	if val, ok := d.GetOk("private_key"); ok {
		config.privateKey = val.(string)
	}
	if config.privateKey == "" && os.Getenv("ZSCALER_PRIVATE_KEY") != "" {
		config.privateKey = os.Getenv("ZSCALER_PRIVATE_KEY")
	}

	if val, ok := d.GetOk("vanity_domain"); ok {
		config.vanityDomain = val.(string)
	}
	if config.vanityDomain == "" && os.Getenv("ZSCALER_VANITY_DOMAIN") != "" {
		config.vanityDomain = os.Getenv("ZSCALER_VANITY_DOMAIN")
	}

	if val, ok := d.GetOk("zscaler_cloud"); ok {
		config.cloud = val.(string)
	}
	if config.cloud == "" && os.Getenv("ZSCALER_CLOUD") != "" {
		config.cloud = os.Getenv("ZSCALER_CLOUD")
	}

	if val, ok := d.GetOk("sandbox_token"); ok {
		config.sandboxToken = val.(string)
	}
	if config.sandboxToken == "" && os.Getenv("ZSCALER_SANDBOX_TOKEN") != "" {
		config.sandboxToken = os.Getenv("ZSCALER_SANDBOX_TOKEN")
	}

	if val, ok := d.GetOk("sandbox_cloud"); ok {
		config.sandboxCloud = val.(string)
	}
	if config.sandboxCloud == "" && os.Getenv("ZSCALER_SANDBOX_CLOUD") != "" {
		config.sandboxCloud = os.Getenv("ZSCALER_SANDBOX_CLOUD")
	}

	if val, ok := d.GetOk("username"); ok {
		config.Username = val.(string)
	}
	if config.Username == "" {
		config.Username = os.Getenv("ZIA_USERNAME")
	}

	if val, ok := d.GetOk("password"); ok {
		config.Password = val.(string)
	}
	if config.Password == "" {
		config.Password = os.Getenv("ZIA_PASSWORD")
	}

	if val, ok := d.GetOk("api_key"); ok {
		config.APIKey = val.(string)
	}
	if config.APIKey == "" {
		config.APIKey = os.Getenv("ZIA_API_KEY")
	}

	if val, ok := d.GetOk("zia_cloud"); ok {
		config.ZIABaseURL = val.(string)
	}
	if config.ZIABaseURL == "" {
		config.ZIABaseURL = os.Getenv("ZIA_CLOUD")
	}

	if val, ok := d.GetOk("sandbox_token"); ok {
		config.sandboxToken = val.(string)
	}
	if config.sandboxToken == "" && os.Getenv("ZSCALER_SANDBOX_TOKEN") != "" {
		config.sandboxToken = os.Getenv("ZSCALER_SANDBOX_TOKEN")
	}

	if val, ok := d.GetOk("sandbox_cloud"); ok {
		config.sandboxCloud = val.(string)
	}
	if config.sandboxCloud == "" && os.Getenv("ZSCALER_SANDBOX_CLOUD") != "" {
		config.sandboxCloud = os.Getenv("ZSCALER_SANDBOX_CLOUD")
	}

	if val, ok := d.GetOk("max_retries"); ok {
		config.retryCount = val.(int)
	}

	if val, ok := d.GetOk("parallelism"); ok {
		config.parallelism = val.(int)
	}

	if val, ok := d.GetOk("backoff"); ok {
		config.backoff = val.(bool)
	}

	if val, ok := d.GetOk("min_wait_seconds"); ok {
		config.minWait = val.(int)
	}

	if val, ok := d.GetOk("max_wait_seconds"); ok {
		config.maxWait = val.(int)
	}

	if val, ok := d.GetOk("log_level"); ok {
		config.logLevel = val.(int)
	}

	if val, ok := d.GetOk("request_timeout"); ok {
		config.requestTimeout = val.(int)
	}

	if httpProxy, ok := d.Get("http_proxy").(string); ok {
		config.httpProxy = httpProxy
	}
	if config.httpProxy == "" && os.Getenv("ZSCALER_HTTP_PROXY") != "" {
		config.httpProxy = os.Getenv("ZSCALER_HTTP_PROXY")
	}

	return &config
}

// loadClients initializes SDK clients based on configuration
func (c *Config) loadClients() diag.Diagnostics {
	if c.useLegacyClient {
		log.Println("[INFO] Initializing ZIA V2 (Legacy) client")
		v2Client, err := zscalerSDKV2Client(c)
		if err != nil {
			return diag.Errorf("failed to initialize SDK V2 client: %v", err)
		}
		c.zscalerSDKClientV3 = v2Client.Client
		return nil
	}

	log.Println("[INFO] Initializing ZIA V3 client")
	v3Client, err := zscalerSDKV3Client(c)
	if err != nil {
		return diag.Errorf("failed to initialize SDK V3 client: %v", err)
	}
	c.zscalerSDKClientV3 = v3Client

	return nil
}

// SelectClient returns the appropriate client based on authentication type or other factors.
func (c *Config) SelectClient() (*zscaler.Client, *zia.Client, error) {
	if !c.useLegacyClient && c.zscalerSDKClientV3 != nil {
		return c.zscalerSDKClientV3, nil, nil
	}

	return nil, nil, fmt.Errorf("no valid client configuration provided")
}

// generateUserAgent constructs the user agent string with all required details
func generateUserAgent(terraformVersion string) string {
	// Fetch the provider version dynamically from common.Version()
	providerVersion := common.Version()

	return fmt.Sprintf("(%s %s) Terraform/%s Provider/%s",
		runtime.GOOS,
		runtime.GOARCH,
		terraformVersion,
		providerVersion,
	)
}

func zscalerSDKV2Client(c *Config) (*zscaler.Service, error) {
	customUserAgent := generateUserAgent(c.TerraformVersion)

	// Start with base configuration setters
	setters := []zia.ConfigSetter{
		zia.WithCache(false),
		zia.WithHttpClientPtr(http.DefaultClient),
		zia.WithRateLimitMaxRetries(int32(c.retryCount)),
		zia.WithRequestTimeout(time.Duration(c.requestTimeout) * time.Second),
		zia.WithUserAgent(customUserAgent), // Set the custom user agent
	}

	// Apply credentials and mandatory parameters
	setters = append(
		setters,
		zia.WithZiaUsername(c.Username),
		zia.WithZiaPassword(c.Password),
		zia.WithZiaAPIKey(c.APIKey),
		zia.WithZiaCloud(c.ZIABaseURL),
	)

	// Configure HTTP proxy if provided
	if c.httpProxy != "" {
		_url, err := url.Parse(c.httpProxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}
		setters = append(setters, zia.WithProxyHost(_url.Hostname()))

		// Default to port 80 if not provided
		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		// Parse the port as a 32-bit integer
		port64, err := strconv.ParseInt(sPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %v", err)
		}

		// Optionally, you can also check the port range if needed
		if port64 < 1 || port64 > 65535 {
			return nil, fmt.Errorf("invalid port number: must be between 1 and 65535, got: %d", port64)
		}
		// Safe cast to int32
		port32 := int32(port64)
		setters = append(setters, zia.WithProxyPort(port32))
	}

	// Initialize ZIA configuration
	ziaCfg, err := zia.NewConfiguration(setters...)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZIA configuration: %v", err)
	}
	ziaCfg.UserAgent = customUserAgent
	// Initialize ZIA client
	wrappedV2Client, err := zscaler.NewLegacyZiaClient(ziaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZIA client: %v", err)
	}

	log.Println("[INFO] Successfully initialized ZIA V2 client")
	return wrappedV2Client, nil
}

func zscalerSDKV3Client(c *Config) (*zscaler.Client, error) {
	customUserAgent := generateUserAgent(c.TerraformVersion)

	// Start with base configuration setters
	setters := []zscaler.ConfigSetter{
		zscaler.WithCache(false),
		zscaler.WithHttpClientPtr(http.DefaultClient),
		zscaler.WithRateLimitMaxRetries(int32(c.retryCount)),
		zscaler.WithRequestTimeout(time.Duration(c.requestTimeout) * time.Second),
		zscaler.WithUserAgentExtra(customUserAgent),
	}

	// Configure HTTP proxy if provided
	if c.httpProxy != "" {
		_url, err := url.Parse(c.httpProxy)
		if err != nil {
			return nil, err
		}
		setters = append(setters, zscaler.WithProxyHost(_url.Hostname()))

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		// Parse the port as a 32-bit integer
		port64, err := strconv.ParseInt(sPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %v", err)
		}

		// Optionally, you can also check the port range if needed
		if port64 < 1 || port64 > 65535 {
			return nil, fmt.Errorf("invalid port number: must be between 1 and 65535, got: %d", port64)
		}
		// Safe cast to int32
		port32 := int32(port64)
		setters = append(setters, zscaler.WithProxyPort(port32))
	}

	// Handle Sandbox-only authentication
	if c.sandboxToken != "" && c.sandboxCloud != "" && c.clientID == "" && c.clientSecret == "" && c.privateKey == "" {
		setters = append(setters,
			zscaler.WithSandboxToken(c.sandboxToken),
			zscaler.WithSandboxCloud(c.sandboxCloud),
		)

		// Initialize configuration for Sandbox
		config, err := zscaler.NewConfiguration(setters...)
		if err != nil {
			return nil, fmt.Errorf("failed to create SDK V3 configuration for Sandbox: %v", err)
		}

		config.UserAgent = customUserAgent

		// Create Sandbox-only client
		v3Client, err := zscaler.NewOneAPIClient(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create Zscaler API client for Sandbox: %v", err)
		}

		return v3Client.Client, nil
	}

	// Main switch for OAuth2 authentication
	switch {
	case c.clientID != "" && c.clientSecret != "" && c.vanityDomain != "":
		setters = append(setters,
			zscaler.WithClientID(c.clientID),
			zscaler.WithClientSecret(c.clientSecret),
			zscaler.WithVanityDomain(c.vanityDomain),
			zscaler.WithSandboxToken(c.sandboxToken),
			zscaler.WithSandboxCloud(c.sandboxCloud),
		)

		if c.cloud != "" {
			setters = append(setters, zscaler.WithZscalerCloud(c.cloud))
		}

	case c.clientID != "" && c.privateKey != "" && c.vanityDomain != "":
		setters = append(setters,
			zscaler.WithClientID(c.clientID),
			zscaler.WithPrivateKey(c.privateKey),
			zscaler.WithVanityDomain(c.vanityDomain),
			zscaler.WithSandboxToken(c.sandboxToken),
			zscaler.WithSandboxCloud(c.sandboxCloud),
		)

		if c.cloud != "" {
			setters = append(setters, zscaler.WithZscalerCloud(c.cloud))
		}

	default:
		return nil, fmt.Errorf("invalid authentication configuration: missing required parameters")
	}

	// Create configuration for OAuth2 authentication
	config, err := zscaler.NewConfiguration(setters...)
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK V3 configuration: %v", err)
	}

	config.UserAgent = customUserAgent

	// Initialize the client with the configuration
	v3Client, err := zscaler.NewOneAPIClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zscaler API client: %v", err)
	}

	return v3Client.Client, nil
}

func (c *Config) Client() (*Client, error) {
	// Handle Sandbox-only credentials
	if c.sandboxToken != "" && c.sandboxCloud != "" && c.clientID == "" && c.clientSecret == "" && c.privateKey == "" {
		v3Client, err := zscalerSDKV3Client(c)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Sandbox client: %w", err)
		}
		return &Client{
			Service: zscaler.NewService(v3Client, nil),
		}, nil
	}

	// Legacy client logic
	if c.useLegacyClient {
		wrappedV2Client, err := zscalerSDKV2Client(c)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize legacy v2 client: %w", err)
		}
		return &Client{
			Service: zscaler.NewService(wrappedV2Client.Client, nil),
		}, nil
	}

	// Fallback to V3 client logic
	v3Client, err := zscalerSDKV3Client(c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize v3 client: %w", err)
	}
	return &Client{
		Service: zscaler.NewService(v3Client, nil),
	}, nil
}

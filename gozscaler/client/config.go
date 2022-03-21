package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

const (
	maxIdleConnections    int = 40
	requestTimeout        int = 60
	jSessionIDTimeout         = 30 // minutes
	jSessionTimeoutOffset     = 5 * time.Minute
	contentTypeJSON           = "application/json"
	cookieName                = "JSESSIONID"
	MaxNumOfRetries           = 100
	RetryWaitMaxSeconds       = 20
	RetryWaitMinSeconds       = 5
	// API types
	ziaAPIVersion   = "api/v1"
	defaultProtocol = "https://"
	ziaAPIAuthURL   = "/authenticatedSession"
	loggerPrefix    = "zia-logger: "
)

// Client ...
type Client struct {
	userName         string
	password         string
	apiKey           string
	session          *Session
	sessionRefreshed time.Time     // Also indicates last usage
	sessionTimeout   time.Duration // in minutes
	URL              string
	HTTPClient       *http.Client
	Logger           *log.Logger
	sync.Mutex
}

// Session ...
type Session struct {
	AuthType           string `json:"authType"`
	ObfuscateAPIKey    bool   `json:"obfuscateApiKey"`
	PasswordExpiryTime int    `json:"passwordExpiryTime"`
	PasswordExpiryDays int    `json:"passwordExpiryDays"`
	Source             string `json:"source"`
	JSessionID         string `json:"jSessionID,omitempty"`
}

// Credentials ...
type Credentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	APIKey    string `json:"apiKey"`
	TimeStamp string `json:"timestamp"`
}

func obfuscateAPIKey(apiKey, timeStamp string) (string, error) {
	// check min required size
	if len(timeStamp) < 6 || len(apiKey) < 12 {
		return "", errors.New("time stamp or api key doesn't have required sizes")
	}

	seed := apiKey

	high := timeStamp[len(timeStamp)-6:]
	highInt, _ := strconv.Atoi(high)
	low := fmt.Sprintf("%06d", highInt>>1)
	key := ""

	for i := 0; i < len(high); i++ {
		index, _ := strconv.Atoi((string)(high[i]))
		key += (string)(seed[index])
	}
	for i := 0; i < len(low); i++ {
		index, _ := strconv.Atoi((string)(low[i]))
		key += (string)(seed[index+2])
	}

	return key, nil
}

// NewClientZIA NewClient Returns a Client from credentials passed as parameters
func NewClientZIA(username, password, apiKey, ziaCloud string) (*Client, error) {
	httpClient := getHTTPClient()
	var logger *log.Logger
	if loggerEnv := os.Getenv("ZSCALER_SDK_LOG"); loggerEnv == "true" {
		logger = getDefaultLogger()
	}
	url := fmt.Sprintf("https://zsapi.%s.net/%s", ziaCloud, ziaAPIVersion)
	cli := Client{
		userName:   username,
		password:   password,
		apiKey:     apiKey,
		HTTPClient: httpClient,
		URL:        url,
		Logger:     logger,
	}
	cli.refreshSession()
	return &cli, nil
}

// MakeAuthRequestZIA ...
func MakeAuthRequestZIA(credentials *Credentials, url string, client *http.Client) (*Session, error) {
	if credentials == nil {
		return nil, fmt.Errorf("empty credentials")
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+ziaAPIAuthURL, contentTypeJSON, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("un-successful request with status code: %v", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var session Session
	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, err
	}
	// We get the whole string match as session ID
	session.JSessionID, err = extractJSessionIDFromHeaders(resp.Header)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func extractJSessionIDFromHeaders(header http.Header) (string, error) {
	sessionIdStr := header.Get("Set-Cookie")
	if sessionIdStr == "" {
		return "", fmt.Errorf("no Set-Cookie header received")
	}
	regex := regexp.MustCompile("JSESSIONID=(.*?);")
	// look for the first match we find
	result := regex.FindStringSubmatch(sessionIdStr)
	if len(result) < 2 {
		return "", fmt.Errorf("couldn't find JSESSIONID in header value")
	}
	return result[1], nil
}

func getCurrentTimestampMilisecond() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

// RefreshSession .. the caller should require lock
func (c *Client) refreshSession() error {
	timeStamp := getCurrentTimestampMilisecond()
	obfuscatedKey, err := obfuscateAPIKey(c.apiKey, timeStamp)
	if err != nil {
		return err
	}
	credentialData := Credentials{
		Username:  c.userName,
		Password:  c.password,
		APIKey:    obfuscatedKey,
		TimeStamp: timeStamp,
	}
	session, err := MakeAuthRequestZIA(&credentialData, c.URL, c.HTTPClient)
	if err != nil {
		return err
	}
	c.session = session
	c.sessionRefreshed = time.Now()
	return nil
}

// checkSession synce new session if its over the timeout limit.
func (c *Client) checkSession() error {
	// One call to this function is allowed at a time caller must call lock.
	if c.session == nil {
		err := c.refreshSession()
		if err != nil {
			log.Printf("[ERROR] failed to get session id: %v\n", err)
			return err
		}
	} else {
		now := time.Now()
		// Refresh if session has expire time (diff than -1)  & c.sessionTimeout less than jSessionTimeoutOffset time remaining. You never refresh on exact timeout.
		if c.session.PasswordExpiryTime > 0 && c.sessionRefreshed.Add(c.sessionTimeout-jSessionTimeoutOffset).Before(now) {
			err := c.refreshSession()
			if err != nil {
				log.Printf("[ERROR] failed to refresh session id: %v\n", err)
				return err
			}
		}
	}
	url, err := url.Parse(c.URL)
	if err != nil {
		log.Printf("[ERROR] failed to parse url %s: %v\n", c.URL, err)
		return err
	}
	if c.HTTPClient.Jar == nil {
		c.HTTPClient.Jar, err = cookiejar.New(nil)
		if err != nil {
			log.Printf("[ERROR] failed to create new http cookie jar %v\n", err)
			return err
		}
	}
	c.HTTPClient.Jar.SetCookies(url, []*http.Cookie{
		{
			Name:  cookieName,
			Value: c.session.JSessionID,
		},
	})
	return nil
}

func (c *Client) GetContentType() string {
	return contentTypeJSON
}

func getHTTPClient() *http.Client {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	retryableClient.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	retryableClient.RetryMax = MaxNumOfRetries
	retryableClient.CheckRetry = checkRetry
	retryableClient.HTTPClient.Timeout = time.Duration(requestTimeout) * time.Second
	retryableClient.HTTPClient.Transport = &http.Transport{
		MaxIdleConnsPerHost: maxIdleConnections,
	}
	retryableClient.HTTPClient.Transport = logging.NewTransport("gozscaler-zia", retryableClient.HTTPClient.Transport)

	return retryableClient.StandardClient()
}

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// return empty slice to enable retry on all connection & server errors.
// or return []int{429}  to retry on only TooManyRequests error
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests}
}

// Used to make http client retry on provided list of response status codes
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp != nil && containsInt(getRetryOnStatusCodes(), resp.StatusCode) {
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}
func getDefaultLogger() *log.Logger {
	return log.New(os.Stdout, loggerPrefix, log.LstdFlags|log.Lshortfile)
}

func (c *Client) Logout() error {
	_, err := c.Request(ziaAPIAuthURL, "DELETE", nil, "application/json")
	if err != nil {
		return err
	}
	return nil
}

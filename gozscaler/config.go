package gozscaler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const (
	maxIdleConnections    int    = 10
	requestTimeout        int    = 60
	jSessionIDTimeout            = 30 // minutes
	jSessionTimeoutOffset        = 5 * time.Minute
	configPath            string = ".zscaler/credentials.json"
	contentTypeJSON              = "application/json"

	// API types
	ziaAPI        = "ZIA"
	ziaAPIVersion = "/api/v1"
	ziaAPIAuthURL = "/authenticatedSession"
)

// Needs to refactor this as it is not required.
var (
	api2ContentType = map[string]string{
		ziaAPI: contentTypeJSON,
	}
	mutex = &sync.Mutex{}
)

// Config contains all the configuration data for the API client
type Config struct {
	APIType             string `json:"api_type"` //This field needs to be removed
	UserName            string `json:"username"`
	Password            string `json:"password"`
	APIKey              string `json:"api_key"`
	ZIABaseURL          string `json:"zia_url"`
	ZIASessionIDTimeout string `json:"zia_session_id_timeout"`
}

// Client ...
type Client struct {
	APIType          string //This field needs to be removed
	UserName         string
	Password         string
	APIKey           string
	Session          *Session
	SessionRefreshed time.Time     // Also indicates last usage
	SessionTimeout   time.Duration // in minutes
	URL              string
	HTTPClient       *http.Client
}

// Session ...
type Session struct {
	AuthType           string      `json:"authType"`
	ObfuscateAPIKey    bool        `json:"obfuscateApiKey"`
	PasswordExpiryTime json.Number `json:"passwordExpiryTime"`
	PasswordExpiryDays json.Number `json:"passwordExpiryDays"`
	Source             string      `json:"source"`
	JSessionID         string      `json:"jSessionID,omitempty"`
}

// ZIACredentials ...
type ZIACredentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	APIKey    string `json:"apiKey"`
	TimeStamp string `json:"timestamp"`
}

// Credentials ...
type Credentials struct {
	Type           string         `json:"type"`
	ZIACredentials ZIACredentials `json:"zia_credentials"`
}

func obfuscateAPIKey(apiKey, timeStamp string) (string, error) {
	// check min required size
	if len(timeStamp) < 6 || len(apiKey) < 12 {
		return "", errors.New("time stamp or api key doesn't have required sizes")
	}

	seed := apiKey

	n := timeStamp[len(timeStamp)-6:]
	nInt, _ := strconv.Atoi(n)
	r := fmt.Sprintf("%06d", nInt>>1)
	key := ""

	for i := 0; i < len(n); i++ {
		index, _ := strconv.Atoi((string)(n[i]))
		key += (string)(seed[index])
	}
	for i := 0; i < len(r); i++ {
		index, _ := strconv.Atoi((string)(r[i]))
		key += (string)(seed[index+2])
	}

	return key, nil
}

// NewClientZIA NewClient Returns a Client from credentials passed as parameters
func NewClientZIA(username, password, apiKey, url, customerID string) (*Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}

	timeStamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	newKey, err := obfuscateAPIKey(apiKey, timeStamp)
	if err != nil {
		return nil, err
	}
	credentialData := Credentials{
		Type: ziaAPI,
		ZIACredentials: ZIACredentials{
			Username:  username,
			Password:  password,
			APIKey:    newKey,
			TimeStamp: timeStamp,
		},
	}

	session, err := MakeAuthRequestZIA(&credentialData, url, httpClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		APIType:          ziaAPI,
		UserName:         username,
		Password:         password,
		APIKey:           apiKey,
		Session:          session,
		SessionRefreshed: time.Now(),
		HTTPClient:       httpClient,
		URL:              url,
	}, nil
}

// MakeAuthRequestZIA ...
func MakeAuthRequestZIA(credentials *Credentials, url string, client *http.Client) (*Session, error) {
	if credentials == nil {
		return nil, fmt.Errorf("empty credentials")
	}

	data, err := json.Marshal(credentials.ZIACredentials)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(url+ziaAPIVersion+ziaAPIAuthURL, contentTypeJSON, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("un-successful request with status code: %v", resp.Status)
	}

	var cookieList []string
	var ok bool
	var sessionIdStr string
	if cookieList, ok = resp.Header["Set-Cookie"]; !ok {
		return nil, fmt.Errorf("no Set-Cookie header receieved")
	}
	if len(cookieList) == 0 {
		return nil, fmt.Errorf("empty JSESSIONID receieved")
	}

	var session Session
	sessionIdStr = cookieList[0]
	regex := regexp.MustCompile("JSESSIONID=(.*?);")
	// look for the first match we find
	result := regex.FindStringSubmatch(sessionIdStr)

	if len(result) < 2 {
		return nil, fmt.Errorf("couldn't find JSESSIONID in header value")
	}

	// We get the whole string match as session ID
	session.JSessionID = result[0][:len(result[0])-1]

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetContentType ..
func (c Client) GetContentType() string {
	var contentType string
	var ok bool
	if contentType, ok = api2ContentType[c.APIType]; ok {
		return contentType
	}
	return ""
}

// RefreshSession ..
func (c *Client) RefreshSession() error {
	if c.APIType != ziaAPI {
		return fmt.Errorf("client's API type doesn't include refreshing sessions")
	}

	timeStamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	newKey, err := obfuscateAPIKey(c.APIKey, timeStamp)
	if err != nil {
		return err
	}
	credentialData := Credentials{
		Type: ziaAPI,
		ZIACredentials: ZIACredentials{
			Username:  c.UserName,
			Password:  c.Password,
			APIKey:    newKey,
			TimeStamp: timeStamp,
		},
	}

	session, err := MakeAuthRequestZIA(&credentialData, c.URL, c.HTTPClient)
	if err != nil {
		return err
	}

	c.Session = session
	c.SessionRefreshed = time.Now()
	return nil
}

// GetSession return new session if its over the timeout limit.
func (c *Client) GetSession() *Session {
	// One call to this function is allowed at a time.
	mutex.Lock()
	if c.Session == nil {
		c.RefreshSession()
		return c.Session
	}

	now := time.Now()
	// Refresh if less than jSessionTimeoutOffset time remaining. You never refresh on exact timeout.
	if c.SessionRefreshed.Add(c.SessionTimeout - jSessionTimeoutOffset).Before(now) {
		c.RefreshSession()
	}
	mutex.Unlock()
	return c.Session
}

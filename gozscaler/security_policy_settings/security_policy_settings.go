package security_policy_settings

import (
	"fmt"
	"log"
)

const (
	securityEndpoint                      = "/security"
	securityAdvancedEndpoint              = "/security/advanced"
	securityAdvancedBlackListUrlsEndpoint = "/advanced/blacklistUrls"
)

// TODO: because there isn't an endpoint to get all Urls, we need to have all action types here
var AddRemoveURLFromList []string = []string{
	"ADD_TO_LIST",
	"REMOVE_FROM_LIST",
}

type WhiteListUrls struct {
	WhiteListUrls map[string]interface{} `json:"whitelistUrls"`
}

type BlackListUrls struct {
	BlackListUrls map[string]interface{} `json:"blacklistUrls"`
}

func (service *Service) GetWhiteListUrls() (*WhiteListUrls, error) {
	var whitelist WhiteListUrls
	err := service.Client.Read(fmt.Sprintf(securityEndpoint), &whitelist)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] got whitelisturls:%#v", whitelist)
	return &whitelist, nil
}

func (service *Service) GetBlackListUrls() (*BlackListUrls, error) {
	var blacklist BlackListUrls
	err := service.Client.Read(fmt.Sprintf(securityAdvancedEndpoint), &blacklist)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] got blacklisturls:%#v", blacklist)
	return &blacklist, nil
}

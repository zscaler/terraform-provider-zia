package user_authentication_settings

import (
	"fmt"
	"log"
)

const (
	exemptedUrlsEndpoint = "/authSettings/exemptedUrls"
)

type ExemptedUrls struct {
	URLs []string `json:"urls"`
}

type QueryParameters struct {
	ID string
}

func (service *Service) Get() (*ExemptedUrls, error) {
	var urls ExemptedUrls
	err := service.Client.Read(exemptedUrlsEndpoint, &urls)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning exempted url rules from Get: %v", urls)
	return &urls, nil
}

// return the new items that were added to slice1
func difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}

func (service *Service) Update(urls ExemptedUrls) (*ExemptedUrls, error) {
	currentUrsl, err := service.Get()
	if err != nil {
		return nil, err
	}
	newUrls := difference(urls.URLs, currentUrsl.URLs)
	removedUrls := difference(currentUrsl.URLs, urls.URLs)
	if len(newUrls) > 0 {
		_, err := service.Client.Create(fmt.Sprintf("%s?action=ADD_TO_LIST", exemptedUrlsEndpoint), ExemptedUrls{newUrls})
		if err != nil {
			return nil, err
		}
	}
	if len(removedUrls) > 0 {
		_, err := service.Client.Create(fmt.Sprintf("%s?action=REMOVE_FROM_LIST", exemptedUrlsEndpoint), ExemptedUrls{removedUrls})
		if err != nil {
			return nil, err
		}
	}
	return &urls, nil
}

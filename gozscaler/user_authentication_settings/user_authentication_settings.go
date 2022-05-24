package user_authentication_settings

const (
	exemptedUrlsEndpoint = "/authSettings/exemptedUrls"
)

var AddRemoveURLFromList []string = []string{
	"ADD_TO_LIST",
	"REMOVE_FROM_LIST",
}

type ExemptedUrls struct {
	URLs map[string]interface{} `json:"urls"`
}

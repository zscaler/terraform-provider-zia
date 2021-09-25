package adminrolemgmt

import (
	"fmt"
	"net/http"
)

const (
	adminUsersEndpoint = "/adminUsers"
)

type AdminUsers struct {
	ID                          string                `json:"id,omitempty"`
	LoginName                   string                `json:"loginName,omitempty"`
	UserName                    string                `json:"userName,omitempty"`
	Email                       string                `json:"email,omitempty"`
	Comments                    string                `json:"comments,omitempty"`
	Disabled                    bool                  `json:"disabled"`
	Password                    string                `json:"password,omitempty"`
	IsNonEditable               bool                  `json:"isNonEditable"`
	IsPasswordLoginAllowed      bool                  `json:"isPasswordLoginAllowed"`
	IsPasswordExpired           bool                  `json:"isPasswordExpired"`
	IsAuditor                   bool                  `json:"isAuditor"`
	IsSecurityReportCommEnabled bool                  `json:"isSecurityReportCommEnabled"`
	IsServiceUpdateCommEnabled  bool                  `json:"isServiceUpdateCommEnabled"`
	IsProductUpdateCommEnabled  bool                  `json:"isProductUpdateCommEnabled"`
	IsExecMobileAppEnabled      bool                  `json:"isExecMobileAppEnabled"`
	AdminScope                  AdminScope            `json:"adminScope"`
	ExecMobileAppTokens         []ExecMobileAppTokens `json:"execMobileAppTokens"`
	Role                        Role                  `json:"role,omitempty"`
}

type AdminScope struct {
	ScopeGroupMemberEntities []ScopeGroupMemberEntities `json:"scopeGroupMemberEntities"`
	Type                     string                     `json:"Type,omitempty"`
	ScopeEntities            []ScopeEntities            `json:"ScopeEntities,"`
}

type ScopeEntities struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type ScopeGroupMemberEntities struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type Role struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type ExecMobileAppTokens struct {
	Cloud       string `json:"cloud,omitempty"`
	OrgId       string `json:"orgId,omitempty"`
	Name        string `json:"name,omitempty"`
	TokenId     string `json:"tokenId,omitempty"`
	Token       string `json:"token,omitempty"`
	TokenExpiry string `json:"tokenExpiry,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
	DeviceId    string `json:"deviceId,omitempty"`
	DeviceName  string `json:"deviceName,omitempty"`
}

func (service *Service) Get(userId string) (*AdminUsers, *http.Response, error) {
	v := new(AdminUsers)
	relativeURL := fmt.Sprintf("%s/%s", adminUsersEndpoint, userId)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Create(server AdminUsers) (*AdminUsers, *http.Response, error) {
	v := new(AdminUsers)
	resp, err := service.Client.NewRequestDo("POST", adminUsersEndpoint, nil, server, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(userId string, appServer AdminUsers) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", adminUsersEndpoint, userId)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, appServer, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(userId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", adminUsersEndpoint, userId)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

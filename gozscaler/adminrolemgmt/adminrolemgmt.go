package adminrolemgmt

import (
	"fmt"
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

func (service *Service) Get(userId string) (*AdminUsers, error) {
	v := new(AdminUsers)
	relativeURL := fmt.Sprintf("%s/%s", adminUsersEndpoint, userId)
	err := service.Client.Read(relativeURL, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (service *Service) Create(server AdminUsers) (*AdminUsers, error) {
	resp, err := service.Client.Create(adminUsersEndpoint, server)
	if err != nil {
		return nil, err
	}
	res, ok := resp.(AdminUsers)
	if !ok {
		return nil, fmt.Errorf("could marshal response to a valid object")
	}
	return &res, nil
}

func (service *Service) Update(userId string, appServer AdminUsers) (*AdminUsers, error) {
	path := fmt.Sprintf("%s/%s", adminUsersEndpoint, userId)
	resp, err := service.Client.Update(path, appServer)
	if err != nil {
		return nil, err
	}
	res, _ := resp.(AdminUsers)
	return &res, err
}

func (service *Service) Delete(userId string) error {
	path := fmt.Sprintf("%s/%s", adminUsersEndpoint, userId)
	return service.Client.Delete(path)
}

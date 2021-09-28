package adminuserrolemgmt

import (
	"fmt"
	"strings"
)

const (
	adminUsersEndpoint = "/adminUsers"
)

type AdminUsers struct {
	ID                          int                   `json:"id"`
	LoginName                   string                `json:"loginName,omitempty"`
	UserName                    string                `json:"userName,omitempty"`
	Email                       string                `json:"email,omitempty"`
	Comments                    string                `json:"comments,omitempty"`
	Disabled                    bool                  `json:"disabled"`
	Password                    string                `json:"password,omitempty"`
	PasswordLastModifiedTime    int                   `json:"pwdLastModifiedTime,omitempty"`
	IsNonEditable               bool                  `json:"isNonEditable"`
	IsPasswordLoginAllowed      bool                  `json:"isPasswordLoginAllowed"`
	IsPasswordExpired           bool                  `json:"isPasswordExpired"`
	IsAuditor                   bool                  `json:"isAuditor"`
	IsSecurityReportCommEnabled bool                  `json:"isSecurityReportCommEnabled"`
	IsServiceUpdateCommEnabled  bool                  `json:"isServiceUpdateCommEnabled"`
	IsProductUpdateCommEnabled  bool                  `json:"isProductUpdateCommEnabled"`
	IsExecMobileAppEnabled      bool                  `json:"isExecMobileAppEnabled"`
	AdminScopeType              string                `json:"adminScopeType"`
	AdminScope                  AdminScope            `json:"adminScope"`
	Role                        Role                  `json:"role,omitempty"`
	ExecMobileAppTokens         []ExecMobileAppTokens `json:"execMobileAppTokens"`
}
type Role struct {
	ID           int                    `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	IsNameL10Tag bool                   `json:"isNameL10nTag,omitempty"`
	Extensions   map[string]interface{} `json:"extensions,omitempty"`
}
type AdminScope struct {
	ScopeGroupMemberEntities []ScopeGroupMemberEntities `json:"scopeGroupMemberEntities"`
	Type                     string                     `json:"Type,omitempty"`
	ScopeEntities            []ScopeEntities            `json:"ScopeEntities,"`
}
type ScopeEntities struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type ScopeGroupMemberEntities struct {
	ID         int                    `json:"id,omitempty"`
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

func (service *Service) GetAdminUsers(adminUserId int) (*AdminUsers, error) {
	v := new(AdminUsers)
	relativeURL := fmt.Sprintf("%s/%d", adminUsersEndpoint, adminUserId)
	err := service.Client.Read(relativeURL, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (service *Service) GetAdminUsersByName(adminUsersLoginName string) (*AdminUsers, error) {
	var adminUsers []AdminUsers
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := service.Client.Read(adminUsersEndpoint, &adminUsers)
	if err != nil {
		return nil, err
	}
	for _, adminUser := range adminUsers {
		if strings.EqualFold(adminUser.LoginName, adminUsersLoginName) {
			return &adminUser, nil
		}
	}
	return nil, fmt.Errorf("no location found with name: %s", adminUsersLoginName)
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

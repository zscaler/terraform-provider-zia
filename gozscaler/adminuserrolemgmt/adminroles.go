package adminuserrolemgmt

import (
	"fmt"
	"strings"
)

const (
	adminRolesEndpoint = "/adminRoles/lite"
)

type AdminRoles struct {
	ID              int      `json:"id"`
	Rank            int      `json:"rank,omitempty"`
	Name            string   `json:"name,omitempty"`
	PolicyAccess    string   `json:"policyAccess,omitempty"`
	DashboardAccess string   `json:"dashboardAccess"`
	ReportAccess    string   `json:"reportAccess,omitempty"`
	AnalysisAccess  string   `json:"analysisAccess,omitempty"`
	UsernameAccess  string   `json:"usernameAccess,omitempty"`
	AdminAcctAccess string   `json:"adminAcctAccess,omitempty"`
	IsAuditor       bool     `json:"isAuditor,omitempty"`
	Permissions     []string `json:"permissions,omitempty"`
	IsNonEditable   bool     `json:"isNonEditable,omitempty"`
	LogsLimit       string   `json:"logsLimit,omitempty"`
	RoleType        string   `json:"roleType,omitempty"`
}

func (service *Service) GetAdminRoles(adminRoleId int) (*AdminRoles, error) {
	v := new(AdminRoles)
	relativeURL := fmt.Sprintf("%s/%d", adminRolesEndpoint, adminRoleId)
	err := service.Client.Read(relativeURL, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (service *Service) GetAdminRolesByName(adminRoleName string) (*AdminRoles, error) {
	var adminRoles []AdminRoles
	err := service.Client.Read(adminRolesEndpoint, &adminRoles)
	if err != nil {
		return nil, err
	}
	for _, adminRole := range adminRoles {
		if strings.EqualFold(adminRole.Name, adminRoleName) {
			return &adminRole, nil
		}
	}
	return nil, fmt.Errorf("no admin role found with name: %s", adminRoleName)
}

/*
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
*/

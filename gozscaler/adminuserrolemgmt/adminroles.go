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

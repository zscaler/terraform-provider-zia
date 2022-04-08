package dlp_idmp_profile

import (
	"fmt"
	"log"
	"strings"

	"github.com/willguibr/terraform-provider-zia/gozscaler/common"
)

const (
	dlpIdmProfileEndpoint = "/idmprofile"
)

type DLPIDMProfile struct {
	ProfileID          int                       `json:"id,omitempty"`
	ProfileName        string                    `json:"profileName,omitempty"`
	ProfileDesc        string                    `json:"profileDesc,omitempty"`
	ProfileType        string                    `json:"profileType,omitempty"`
	Host               string                    `json:"host,omitempty"`
	Port               int                       `json:"port,omitempty"`
	ProfileDirPath     string                    `json:"profileDirPath,omitempty"`
	ScheduleType       string                    `json:"scheduleType,omitempty"`
	ScheduleDay        int                       `json:"scheduleDay,omitempty"`
	ScheduleDayOfMonth string                    `json:"scheduleDayOfMonth,omitempty"`
	ScheduleDayOfWeek  string                    `json:"scheduleDayOfWeek,omitempty"`
	ScheduleTime       int                       `json:"scheduleTime,omitempty"`
	ScheduleDisabled   bool                      `json:"scheduleDisabled,omitempty"`
	UploadStatus       string                    `json:"uploadStatus,omitempty"`
	UserName           string                    `json:"userName,omitempty"`
	Version            int                       `json:"version,omitempty"`
	VolumeOfDocuments  int                       `json:"volumeOfDocuments,omitempty"`
	NumDocuments       int                       `json:"numDocuments,omitempty"`
	LastModifiedTime   int                       `json:"lastModifiedTime,omitempty"`
	ModifiedBy         []common.IDNameExtensions `json:"modifiedBy,omitempty"`
	IdmClient          []common.IDNameExtensions `json:"idmClient,omitempty"`
}

func (service *Service) Get(IdmProfileID int) (*DLPIDMProfile, error) {
	var dlpIdmProfile DLPIDMProfile
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpIdmProfileEndpoint, IdmProfileID), &dlpIdmProfile)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning dlp idm profile from Get: %d", dlpIdmProfile.ProfileID)
	return &dlpIdmProfile, nil
}

func (service *Service) GetByName(idmProfileName string) (*DLPIDMProfile, error) {
	var dlpIdmProfile []DLPIDMProfile
	err := service.Client.Read(dlpIdmProfileEndpoint, &dlpIdmProfile)
	if err != nil {
		return nil, err
	}
	for _, idm := range dlpIdmProfile {
		if strings.EqualFold(idm.ProfileName, idmProfileName) {
			return &idm, nil
		}
	}
	return nil, fmt.Errorf("no dlp idm profile found with name: %s", idmProfileName)
}

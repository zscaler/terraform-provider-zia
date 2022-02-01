package devicegroups

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

const (
	deviceGroupEndpoint = "/deviceGroups"
	devicesEndpoint     = "/deviceGroups/devices"
)

type DeviceGroups struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	GroupType   string `json:"groupType,omitempty"`
	Description string `json:"description,omitempty"`
	OSType      string `json:"osType,omitempty"`
	Predefined  bool   `json:"predefined"`
	DeviceNames string `json:"deviceNames,omitempty"`
	DeviceCount int    `json:"deviceCount,omitempty"`
}

type Devices struct {
	ID              int    `json:"id"`
	Name            string `json:"name,omitempty"`
	DeviceGroupType string `json:"deviceGroupType,omitempty"`
	DeviceModel     string `json:"deviceModel,omitempty"`
	OSType          string `json:"osType,omitempty"`
	OSVersion       string `json:"osVersion,omitempty"`
	Description     string `json:"description,omitempty"`
	OwnerUserId     int    `json:"ownerUserId,omitempty"`
	OwnerName       string `json:"ownerName,omitempty"`
}

func (service *Service) GetDeviceGroups(deviceGroupId int) (*DeviceGroups, error) {
	var group DeviceGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", deviceGroupEndpoint, deviceGroupId), &group)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning firewall rule from Get: %d", group.ID)
	return &group, nil
}

func (service *Service) GetDeviceGroupByName(deviceGroupName string) (*DeviceGroups, error) {
	var deviceGroups []DeviceGroups
	err := service.Client.Read(deviceGroupEndpoint, &deviceGroups)
	if err != nil {
		return nil, err
	}
	for _, deviceGroup := range deviceGroups {
		if strings.EqualFold(deviceGroup.Name, deviceGroupName) {
			return &deviceGroup, nil
		}
	}
	return nil, fmt.Errorf("no device group found with name: %s", deviceGroupName)
}

func (service *Service) GetDevicesByID(deviceId int) (*Devices, error) {
	var device Devices
	err := service.Client.Read(fmt.Sprintf("%s/%d", devicesEndpoint, deviceId), &device)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning device from Get: %d", device.ID)
	return &device, nil
}

// Get Devices by Name
func (service *Service) GetDevicesByName(deviceName string) (*Devices, error) {
	var devices []Devices
	// We are assuming this device name will be in the firsy 1000 obejcts
	err := service.Client.Read(fmt.Sprintf("%s?page=1&pageSize=1000", devicesEndpoint), &devices)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if strings.EqualFold(device.Name, deviceName) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found with name: %s", deviceName)
}

func (service *Service) GetDevicesByModel(deviceModel string) (*Devices, error) {
	var models []Devices
	err := service.Client.Read(fmt.Sprintf("%s?model=%s", devicesEndpoint, url.QueryEscape(deviceModel)), &models)
	if err != nil {
		return nil, err
	}
	for _, model := range models {
		if strings.EqualFold(model.DeviceModel, deviceModel) {
			return &model, nil
		}
	}
	return nil, fmt.Errorf("no device found with model: %s", deviceModel)
}

func (service *Service) GetDevicesByOwner(ownerName string) (*Devices, error) {
	var owners []Devices
	err := service.Client.Read(fmt.Sprintf("%s?owner=%s", devicesEndpoint, url.QueryEscape(ownerName)), &owners)
	if err != nil {
		return nil, err
	}
	for _, owner := range owners {
		if strings.EqualFold(owner.OwnerName, ownerName) {
			return &owner, nil
		}
	}
	return nil, fmt.Errorf("no device found for owner: %s", ownerName)
}

func (service *Service) GetDevicesByOSType(osTypeName string) (*Devices, error) {
	var osTypes []Devices
	err := service.Client.Read(fmt.Sprintf("%s?osType=%s", devicesEndpoint, url.QueryEscape(osTypeName)), &osTypes)
	if err != nil {
		return nil, err
	}
	for _, osType := range osTypes {
		if strings.EqualFold(osType.OSType, osTypeName) {
			return &osType, nil
		}
	}
	return nil, fmt.Errorf("no device found for type: %s", osTypeName)
}

func (service *Service) GetDevicesByOSVersion(osVersionName string) (*Devices, error) {
	var osVersions []Devices
	err := service.Client.Read(fmt.Sprintf("%s?osVersion=%s", devicesEndpoint, url.QueryEscape(osVersionName)), &osVersions)
	if err != nil {
		return nil, err
	}
	for _, osVersion := range osVersions {
		if strings.EqualFold(osVersion.OSVersion, osVersionName) {
			return &osVersion, nil
		}
	}
	return nil, fmt.Errorf("no device found for version: %s", osVersionName)
}

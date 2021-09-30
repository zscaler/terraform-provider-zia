package timewindows

import (
	"fmt"
	"log"
	"strings"
)

const (
	timeWindowsEndpoint = "/timeWindows"
)

type TimeWindows struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartTime int    `json:"startTime"`
	EndTime   string `json:"description"`
	DayOfWeek string `json:"dayOfWeek"`
}

func (service *Service) GetNetworkServiceGroups(timeWindowID int) (*TimeWindows, error) {
	var timeWindows TimeWindows
	err := service.Client.Read(fmt.Sprintf("%s/%d", timeWindowsEndpoint, timeWindowID), &timeWindows)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning network application group from Get: %d", timeWindows.ID)
	return &timeWindows, nil
}

func (service *Service) GetNetworkServiceGroupsByName(timeWindowName string) (*TimeWindows, error) {
	var timeWindows []TimeWindows
	err := service.Client.Read(timeWindowsEndpoint, &timeWindows)
	if err != nil {
		return nil, err
	}
	for _, timeWindow := range timeWindows {
		if strings.EqualFold(timeWindow.Name, timeWindowName) {
			return &timeWindow, nil
		}
	}
	return nil, fmt.Errorf("no network service groups found with name: %s", timeWindowName)
}

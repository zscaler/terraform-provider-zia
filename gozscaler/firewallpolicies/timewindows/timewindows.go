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
	Name      string `json:"name,omitempty"`
	StartTime int    `json:"startTime,omitempty"`
	EndTime   string `json:"description,omitempty"`
	DayOfWeek string `json:"dayOfWeek,omitempty"`
}

func (service *Service) Get(timeWindowID int) (*TimeWindows, error) {
	var timeWindows TimeWindows
	err := service.Client.Read(fmt.Sprintf("%s/%d", timeWindowsEndpoint, timeWindowID), &timeWindows)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning time windows from Get: %d", timeWindows.ID)
	return &timeWindows, nil
}

func (service *Service) GetByName(timeWindowName string) (*TimeWindows, error) {
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
	return nil, fmt.Errorf("no time windows found with name: %s", timeWindowName)
}

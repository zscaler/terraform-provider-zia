package timewindow

import (
	"fmt"
	"log"
	"strings"
)

const (
	timeWindowEndpoint = "/timeWindows"
)

type TimeWindow struct {
	ID        int      `json:"id"`
	Name      string   `json:"name,omitempty"`
	StartTime int      `json:"startTime,omitempty"`
	EndTime   string   `json:"description,omitempty"`
	DayOfWeek []string `json:"dayOfWeek,omitempty"`
}

func (service *Service) GetTimeWindow(timeWindowID int) (*TimeWindow, error) {
	var timeWindow TimeWindow
	err := service.Client.Read(fmt.Sprintf("%s/%d", timeWindowEndpoint, timeWindowID), &timeWindow)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning time window from Get: %d", timeWindow.ID)
	return &timeWindow, nil
}

func (service *Service) GetTimeWindowByName(timeWindowName string) (*TimeWindow, error) {
	var timeWindow []TimeWindow
	err := service.Client.Read(timeWindowEndpoint, &timeWindow)
	if err != nil {
		return nil, err
	}
	for _, timeWindow := range timeWindow {
		if strings.EqualFold(timeWindow.Name, timeWindowName) {
			return &timeWindow, nil
		}
	}
	return nil, fmt.Errorf("no time window found with name: %s", timeWindowName)
}

package zia

import (
	"fmt"
	"strings"
	"time"
)

const (
	maxRetries    = 3
	retryInterval = 10 * time.Second
)

// RetryOnError will execute the passed function and check its returned error.
func RetryOnError(fn func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil || !strings.Contains(err.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			break
		}
		time.Sleep(retryInterval)
	}
	return err
}

func condenseError(errorList []error) error {
	if len(errorList) < 1 {
		return nil
	}
	msgList := make([]string, len(errorList))
	for i, err := range errorList {
		if err != nil {
			msgList[i] = err.Error()
		}
	}
	return fmt.Errorf("series of errors occurred: %s", strings.Join(msgList, ", "))
}

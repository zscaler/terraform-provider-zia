package zia

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

type checkUpstream func(string) (bool, error)

func ensureResourceExists(name string, checkUpstream checkUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}
		exist, err := checkUpstream(rs.Primary.ID)
		if err != nil {
			return err
		} else if !exist {
			return missingErr
		}
		return nil
	}
}

func checkResourceDestroy(typeName string, checkUpstream checkUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != typeName {
				continue
			}
			exists, err := checkUpstream(rs.Primary.ID)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("resource still exists, ID: %s", rs.Primary.ID)
			}
		}
		return nil
	}
}

func ensureResourceNotExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return nil
		}
		return fmt.Errorf("Resource found: %s", name)
	}
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

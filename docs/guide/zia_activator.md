---
page_title: "ZIA Config Activation"
subcategory: "Activation"
---

#

```go
package main

import (
 "log"
 "os"

 "github.com/willguibr/terraform-provider-zia/gozscaler/activation"
 "github.com/willguibr/terraform-provider-zia/gozscaler/client"
 "github.com/willguibr/terraform-provider-zia/zia"
)

func getEnvVarOrFail(k string) string {
 if v := os.Getenv(k); v != "" {
  return v
 }
 log.Fatalf("[ERROR] Couldn't find environment variable %s\n", k)
 return ""
}

func main() {
 log.Printf("[INFO] Initializing ZIA client\n")
 c := zia.Config{
  Username:   getEnvVarOrFail("ZIA_USERNAME"),
  Password:   getEnvVarOrFail("ZIA_PASSWORD"),
  APIKey:     getEnvVarOrFail("ZIA_API_KEY"),
  ZIABaseURL: getEnvVarOrFail("ZIA_CLOUD"),
 }
 cli, err := client.NewClientZIA(c.Username, c.Password, c.APIKey, c.ZIABaseURL)
 if err != nil {
  log.Fatalf("[ERROR] Failed Initializing ZIA client: %v\n", err)
 }
 activationService := activation.New(cli)
 resp, err := activationService.CreateActivation(activation.Activation{
  Status: "active",
 })
 if err != nil {
  log.Printf("[ERROR] Activation Failed: %v\n", err)
 } else {
  log.Printf("[INFO] Activation succeded: %#v\n", resp)
 }
 log.Printf("[INFO] Destroying session: %#v\n", resp)
 cli.Logout()
 os.Exit(0)
}

```

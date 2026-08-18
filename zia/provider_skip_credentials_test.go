package zia

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// clearAuthEnv blanks every credential-bearing environment variable so the
// tests behave the same on developer machines and CI.
func clearAuthEnv(t *testing.T) {
	t.Helper()
	for _, v := range []string{
		"ZSCALER_CLIENT_ID", "ZSCALER_CLIENT_SECRET", "ZSCALER_PRIVATE_KEY",
		"ZSCALER_VANITY_DOMAIN", "ZSCALER_CLOUD",
		"ZIA_USERNAME", "ZIA_PASSWORD", "ZIA_API_KEY", "ZIA_CLOUD",
		"ZSCALER_USE_LEGACY_CLIENT", "ZSCALER_SKIP_CREDENTIALS_VALIDATION",
	} {
		t.Setenv(v, "")
	}
}

func TestProviderConfigureSkipCredentialsValidation(t *testing.T) {
	clearAuthEnv(t)

	d := schema.TestResourceDataRaw(t, ZIAProvider().Schema, map[string]interface{}{
		"skip_credentials_validation": true,
	})

	meta, diags := providerConfigure(d, "1.0-test")
	if diags.HasError() {
		t.Fatalf("expected no error in skip mode, got: %v", diags)
	}

	foundWarning := false
	for _, diagnostic := range diags {
		if strings.Contains(diagnostic.Summary, "credentials were not validated") {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Errorf("expected a warning diagnostic about skipped credential validation, got: %v", diags)
	}

	client, ok := meta.(*Client)
	if !ok {
		t.Fatalf("expected *Client meta, got %T", meta)
	}
	if !client.skipCredentialsValidation {
		t.Error("expected client to be marked as inert (skipCredentialsValidation)")
	}
	if client.Service != nil {
		t.Error("expected nil Service on the inert client")
	}
}

func TestProviderConfigureSkipCredentialsValidationEnvVar(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("ZSCALER_SKIP_CREDENTIALS_VALIDATION", "true")

	d := schema.TestResourceDataRaw(t, ZIAProvider().Schema, map[string]interface{}{})

	meta, diags := providerConfigure(d, "1.0-test")
	if diags.HasError() {
		t.Fatalf("expected no error in env-var skip mode, got: %v", diags)
	}
	client, ok := meta.(*Client)
	if !ok || !client.skipCredentialsValidation {
		t.Fatalf("expected inert client via env var, got %#v", meta)
	}
}

func TestProviderConfigureMissingCredentialsStillErrors(t *testing.T) {
	clearAuthEnv(t)

	d := schema.TestResourceDataRaw(t, ZIAProvider().Schema, map[string]interface{}{})

	_, diags := providerConfigure(d, "1.0-test")
	if !diags.HasError() {
		t.Fatal("expected configure to fail without credentials when skip mode is off")
	}
}

func TestInertClientGuardOnResources(t *testing.T) {
	p := ZIAProvider()
	inert := &Client{skipCredentialsValidation: true}
	ctx := context.Background()

	checkGuard := func(kind, name, op string, f func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics, r *schema.Resource) {
		if f == nil {
			return
		}
		diags := f(ctx, r.TestResourceData(), inert)
		if !diags.HasError() {
			t.Errorf("%s %s (%s): expected guard error with inert client, got none", kind, name, op)
		} else if !strings.Contains(diags[0].Summary, "skip_credentials_validation") {
			t.Errorf("%s %s (%s): expected guard error to mention skip_credentials_validation, got: %s", kind, name, op, diags[0].Summary)
		}
	}

	for name, r := range p.ResourcesMap {
		checkGuard("resource", name, "create", r.CreateContext, r)
		checkGuard("resource", name, "read", r.ReadContext, r)
		checkGuard("resource", name, "update", r.UpdateContext, r)
		checkGuard("resource", name, "delete", r.DeleteContext, r)
	}
	for name, ds := range p.DataSourcesMap {
		checkGuard("data source", name, "read", ds.ReadContext, ds)
	}
}

package zia

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestValidateUserAgentTypes(t *testing.T) {
	validate := validateUserAgentTypes()
	path := cty.Path{cty.GetAttrStep{Name: "user_agent_types"}}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"chrome is valid", "CHROME", false},
		{"firefox is valid", "FIREFOX", false},
		{"brave is valid", "BRAVE", false},
		{"unsupported value rejected", "VIVALDI", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validate(tt.value, path)
			if tt.wantErr && !diags.HasError() {
				t.Fatalf("expected error for value %q, got none", tt.value)
			}
			if !tt.wantErr && diags.HasError() {
				t.Fatalf("expected no error for value %q, got %v", tt.value, diags)
			}
		})
	}
}

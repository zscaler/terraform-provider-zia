package zia

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func TestHydrateURLFilteringRuleCBIProfile(t *testing.T) {
	tests := []struct {
		name      string
		rule      urlfilteringpolicies.URLFilteringRule
		listed    []urlfilteringpolicies.URLFilteringRule
		wantCalls int
		wantID    string
	}{
		{
			name: "non-isolate rule does not list",
			rule: urlfilteringpolicies.URLFilteringRule{ID: 11, Action: "ALLOW"},
		},
		{
			name: "existing profile does not list",
			rule: urlfilteringpolicies.URLFilteringRule{
				ID:         12,
				Action:     "ISOLATE",
				CBIProfile: &common.CBIProfile{ID: "existing"},
			},
			wantID: "existing",
		},
		{
			name: "nonzero profile id leaves sdk fallback authoritative",
			rule: urlfilteringpolicies.URLFilteringRule{
				ID:           16,
				Action:       "ISOLATE",
				CBIProfileID: 9001,
			},
		},
		{
			name: "list hydrates profile when detail omits profile and profile id",
			rule: urlfilteringpolicies.URLFilteringRule{
				ID:           13,
				Action:       "ISOLATE",
				CBIProfileID: 0,
			},
			listed: []urlfilteringpolicies.URLFilteringRule{
				{ID: 99, CBIProfile: &common.CBIProfile{ID: "wrong"}},
				{ID: 13, CBIProfile: &common.CBIProfile{ID: "wanted", Name: "Isolation"}},
			},
			wantCalls: 1,
			wantID:    "wanted",
		},
		{
			name: "missing list profile leaves prior-state fallback available",
			rule: urlfilteringpolicies.URLFilteringRule{ID: 14, Action: "ISOLATE"},
			listed: []urlfilteringpolicies.URLFilteringRule{
				{ID: 14},
			},
			wantCalls: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			err := hydrateURLFilteringRuleCBIProfile(
				context.Background(),
				nil,
				&test.rule,
				func(context.Context, *zscaler.Service) ([]urlfilteringpolicies.URLFilteringRule, error) {
					calls++
					return test.listed, nil
				},
			)
			if err != nil {
				t.Fatalf("hydrateURLFilteringRuleCBIProfile() error = %v, want nil", err)
			}
			if calls != test.wantCalls {
				t.Errorf("hydrateURLFilteringRuleCBIProfile() list calls = %d, want %d", calls, test.wantCalls)
			}
			gotID := ""
			if test.rule.CBIProfile != nil {
				gotID = test.rule.CBIProfile.ID
			}
			if gotID != test.wantID {
				t.Errorf("hydrateURLFilteringRuleCBIProfile() profile ID = %q, want %q", gotID, test.wantID)
			}
		})
	}
}

func TestHydrateURLFilteringRuleCBIProfileReturnsListError(t *testing.T) {
	rule := urlfilteringpolicies.URLFilteringRule{ID: 15, Action: "ISOLATE"}
	want := errors.New("list unavailable")
	err := hydrateURLFilteringRuleCBIProfile(
		context.Background(),
		nil,
		&rule,
		func(context.Context, *zscaler.Service) ([]urlfilteringpolicies.URLFilteringRule, error) {
			return nil, want
		},
	)
	if !errors.Is(err, want) {
		t.Errorf("hydrateURLFilteringRuleCBIProfile() error = %v, want wrapped %v", err, want)
	}
	if err == nil || !strings.Contains(err.Error(), "rule 15") {
		t.Errorf("hydrateURLFilteringRuleCBIProfile() error = %v, want rule ID context", err)
	}
}

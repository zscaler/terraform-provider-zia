SWEEP?=global
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep "zia/")
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=zia
GOFMT:=gofumpt
TFPROVIDERLINT=tfproviderlint
STATICCHECK=staticcheck
TF_PLUGIN_DIR=~/.terraform.d/plugins
ZIA_PROVIDER_NAMESPACE=zscaler.com/zia/zia

# Expression to match against tests
# go test -run <filter>
# e.g. Iden will run all TestAccIdentity tests
ifdef TEST_FILTER
	TEST_FILTER := -run $(TEST_FILTER)
endif

TESTARGS?=-test.v

default: build

dep: # Download required dependencies

docs:
	go generate

build: fmtcheck
	go install

clean:
	go clean -cache -testcache ./...

clean-all:
	go clean -cache -testcache -modcache ./...

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -sweep=$(SWEEP) $(SWEEPARGS)

test:
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) $(TEST_FILTER) -timeout=30s -parallel=5

testacc:
	TF_ACC=1 go test $(TEST) $(TESTARGS) $(TEST_FILTER) -timeout 120m

test\:integration\:zia:
	@echo "$(COLOR_ZSCALER)Running zia integration tests...$(COLOR_NONE)"
	@TF_ACC=1 go test -v -race -cover -coverprofile=coverage.out -covermode=atomic ./zia -parallel 5 -timeout 120m
	go tool cover -html=coverage.out -o coverage.html

# Default set of integration tests to run for ZscalerOne
DEFAULT_INTEGRATION_TESTS?=\
  TestAccDataSourceActivationStatus_Basic \
  TestAccDataSourceAdminRoles_Basic \
  TestAccDataSourceAdminUsers_Basic \
  TestAccDataSourceFWApplicationServicesGroupLite_Basic \
  TestAccDataSourceFWApplicationServicesLite_Basic \
  TestAccDataSourceAuthSettingsUrls_Basic \
  TestAccDataSourceCBIProfile_Basic \
  TestAccDataSourceDeviceGroups_Basic \
  TestAccDataSourceDLPDictionaries_Basic \
  TestAccDataSourceDLPEngines_Basic \
  TestAccDataSourceDLPICAPServers_Basic \
  TestAccDataSourceDLPIncidentReceiverServers_Basic \
  TestAccDataSourceDLPNotificationTemplates_Basic \
  TestAccDataSourceDlpWebRules_Basic \
  TestAccDataSourceFWIPDestinationGroups_Basic \
  TestAccDataSourceFWIPSourceGroups_Basic \
  TestAccDataSourceFWNetworkApplicationGroups_Basic \
  TestAccDataSourceFWNetworkServiceGroups_Basic \
  TestAccResourceFWNetworkServicesBasic \
  TestAccDataSourceFWTimeWindow_Basic \
  TestAccDataSourceLocationGroup_Basic \
  TestAccDataSourceLocationLite_Basic \
  TestAccDataSourceLocationManagement_Basic \
  TestAccDataSourceRuleLabels_Basic \
  TestResourceSandboxSettings_basic \
  TestAccDataSourceTrafficGreInternalIPRangeList_Basic \
  TestAccDataSourceTrafficForwardingGreTunnels_Basic \
  TestAccDataSourceTrafficForwardingStaticIP_Basic \
  TestAccDataSourceTrafficForwardingVPNCredentials_Basic \
  TestAccDataSourceURLCategories_Basic \
  TestAccDataSourceURLFilteringRules_Basic \
  TestAccDataSourceDepartmentManagement_Basic \
  TestAccDataSourceGroupManagement_Basic \
  TestAccDataSourceUserManagement_Basic \
  TestAccResourceAdminUsersBasic \
  TestAccResourceAuthSettingsUrls_basic \
  TestAccResourceDLPDictionariesBasic \
  TestAccResourceDLPEnginesBasic \
  TestAccResourceDLPNotificationTemplatesBasic \
  TestAccResourceDlpWebRules_Basic \
  TestAccResourceFWIPDestinationGroupsBasic \
  TestAccResourceFWIPSourceGroupsBasic \
  TestAccResourceFWNetworkApplicationGroupsBasic \
  TestAccResourceFWNetworkServiceGroupsBasic \
  TestAccResourceFWNetworkServicesBasic \
  TestAccResourceLocationManagementBasic \
  TestAccResourceRuleLabelsBasic \
  TestResourceSandboxSettings_basic \
  TestAccResourceSecurityPolicySettings_basic \
  TestAccResourceTrafficForwardingGRETunnelBasic \
  TestAccResourceTrafficForwardingStaticIPBasic \
  TestAccResourceTrafficForwardingVPNCredentialsBasic \
  TestAccResourceURLCategoriesBasic \
  TestAccResourceURLFilteringRulesBasic \
  TestAccResourceUserManagementBasic

ifeq ($(strip $(INTEGRATION_TESTS)),)
  INTEGRATION_TESTS = $(DEFAULT_INTEGRATION_TESTS)
endif

space := $(subst ,, )
integration_tests := $(subst $(space),\|,$(INTEGRATION_TESTS))

# Target to run integration tests for ZscalerOne
test\:integration\:zscalerone:
	@echo "Running integration tests for ZscalerOne..."
	@TF_ACC=1 go test -v -race -cover -coverprofile=coverage.out -covermode=atomic ./zia -parallel 5 -timeout 120m -run ^$(integration_tests)$$
	go tool cover -html=coverage.out -o coverage.html

# Default set of integration tests to run for ZscalerOne
ZS2_INTEGRATION_TESTS?=\
  TestAccDataSourceActivationStatus_Basic \
  TestAccDataSourceAdminRoles_Basic \
  TestAccDataSourceAdminUsers_Basic \
  TestAccDataSourceFWApplicationServicesGroupLite_Basic \
  TestAccDataSourceFWApplicationServicesLite_Basic \
  TestAccDataSourceAuthSettingsUrls_Basic \
  TestAccDataSourceCBIProfile_Basic \
  TestAccDataSourceDeviceGroups_Basic \
  TestAccDataSourceDLPDictionaries_Basic \
  TestAccDataSourceDLPEngines_Basic \
  TestAccDataSourceDLPICAPServers_Basic \
  TestAccDataSourceDLPIncidentReceiverServers_Basic \
  TestAccDataSourceDLPNotificationTemplates_Basic \
  TestAccDataSourceDlpWebRules_Basic \
  TestAccDataSourceFWIPDestinationGroups_Basic \
  TestAccDataSourceFWIPSourceGroups_Basic \
  TestAccDataSourceFWNetworkApplicationGroups_Basic \
  TestAccDataSourceFWNetworkServiceGroups_Basic \
  TestAccResourceFWNetworkServicesBasic \
  TestAccDataSourceFWTimeWindow_Basic \
  TestAccDataSourceLocationGroup_Basic \
  TestAccDataSourceLocationLite_Basic \
  TestAccDataSourceLocationManagement_Basic \
  TestAccDataSourceRuleLabels_Basic \
  TestResourceSandboxSettings_basic \
  TestAccDataSourceTrafficGreInternalIPRangeList_Basic \
  TestAccDataSourceTrafficForwardingStaticIP_Basic \
  TestAccDataSourceTrafficForwardingVPNCredentials_Basic \
  TestAccDataSourceTrafficForwardingGreTunnels_Basic \
  TestAccDataSourceURLCategories_Basic \
  TestAccDataSourceURLFilteringRules_Basic \
  TestAccDataSourceDepartmentManagement_Basic \
  TestAccDataSourceGroupManagement_Basic \
  TestAccDataSourceUserManagement_Basic \
  TestAccResourceAdminUsersBasic \
  TestAccResourceAuthSettingsUrls_basic \
  TestAccResourceDLPDictionariesBasic \
  TestAccResourceDLPEnginesBasic \
  TestAccResourceDLPNotificationTemplatesBasic \
  TestAccResourceDlpWebRules_Basic \
  TestAccResourceFWIPDestinationGroupsBasic \
  TestAccResourceFWIPSourceGroupsBasic \
  TestAccResourceFWNetworkApplicationGroupsBasic \
  TestAccResourceFWNetworkServiceGroupsBasic \
  TestAccResourceFWNetworkServicesBasic \
  TestAccResourceLocationManagementBasic \
  TestAccResourceRuleLabelsBasic \
  TestAccZiaSandboxFileSubmission_basic \
  TestResourceSandboxSettings_basic \
  TestAccResourceSecurityPolicySettings_basic \
  TestAccResourceTrafficForwardingStaticIPBasic \
  TestAccResourceTrafficForwardingVPNCredentialsBasic \
  TestAccResourceTrafficForwardingGRETunnelBasic \
  TestAccResourceURLCategoriesBasic \
  TestAccResourceURLFilteringRulesBasic \
  TestAccResourceUserManagementBasic

ifeq ($(strip $(ZS_INTEGRATION_TESTS)),)
  ZS_INTEGRATION_TESTS = $(ZS2_INTEGRATION_TESTS)
endif

space := $(subst ,, )
integration_zs2_tests := $(subst $(space),\|,$(ZS_INTEGRATION_TESTS))

# Target to run integration tests for ZscalerTwo
test\:integration\:zscalertwo:
	@echo "Running integration tests for ZscalerTwo..."
	@TF_ACC=1 go test -v -race -cover -coverprofile=coverage.out -covermode=atomic ./zia -parallel 5 -timeout 120m -run ^$(integration_zs2_tests)$$
	go tool cover -html=coverage.out -o coverage.html

build13: GOOS=$(shell go env GOOS)
build13: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
build13: DESTINATION=$(APPDATA)/terraform.d/plugins/$(ZIA_PROVIDER_NAMESPACE)/4.5.3/$(GOOS)_$(GOARCH)
else
build13: DESTINATION=$(HOME)/.terraform.d/plugins/$(ZIA_PROVIDER_NAMESPACE)/4.5.3/$(GOOS)_$(GOARCH)
endif
build13: fmtcheck
	@echo "==> Installing plugin to $(DESTINATION)"
	@mkdir -p $(DESTINATION)
	go build -o $(DESTINATION)/terraform-provider-zia_v4.5.3

coverage: test
	@echo "✓ Opening coverage for unit tests ..."
	@go tool cover -html=coverage.txt

vet:
	@echo "==> Checking source code against go vet and staticcheck"
	@go vet ./...
	@staticcheck ./...

imports:
	goimports -w $(GOFMT_FILES)

fmt: tools # Format the code
	@echo "formatting the code with $(GOFMT)..."
	@$(GOFMT) -l -w .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

fmt-docs:
	@echo "✓ Formatting code samples in documentation"
	@terrafmt fmt -p '*.md' .

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

lint:
	@echo "==> Checking source code against linters..."
	@$(TFPROVIDERLINT) \
		-c 1 \
		-AT001 \
    -R004 \
		-S001 \
		-S002 \
		-S003 \
		-S004 \
		-S005 \
		-S007 \
		-S008 \
		-S009 \
		-S010 \
		-S011 \
		-S012 \
		-S013 \
		-S014 \
		-S015 \
		-S016 \
		-S017 \
		-S019 \
		./$(PKG_NAME)

tools:
	@which $(GOFMT) || go install mvdan.cc/gofumpt@v0.5.0
	@which $(TFPROVIDERLINT) || go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	@which $(STATICCHECK) || go install honnef.co/go/tools/cmd/staticcheck@v0.4.6

tools-update:
	@go install mvdan.cc/gofumpt@v0.5.0
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	@go install honnef.co/go/tools/cmd/staticcheck@v0.4.6

ziaActivator: GOOS=$(shell go env GOOS)
ziaActivator: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
ziaActivator: DESTINATION=C:\Windows\System32
else
ziaActivator: DESTINATION=/usr/local/bin
endif
ziaActivator:
	@echo "==> Installing ziaActivator cli $(DESTINATION)"
	@mkdir -p $(DESTINATION)
	@rm -f $(DESTINATION)/ziaActivator
	@go build -o $(DESTINATION)/ziaActivator  ./cli/ziaActivator.go

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck tools vendor-status test-compile website-lint website website-test
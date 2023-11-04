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
	go mod tidy

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
		xargs -t -n4 go test $(TESTARGS) $(TEST_FILTER) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) $(TESTARGS) $(TEST_FILTER) -timeout 120m

build13: GOOS=$(shell go env GOOS)
build13: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
build13: DESTINATION=$(APPDATA)/terraform.d/plugins/$(ZIA_PROVIDER_NAMESPACE)/2.6.5/$(GOOS)_$(GOARCH)
else
build13: DESTINATION=$(HOME)/.terraform.d/plugins/$(ZIA_PROVIDER_NAMESPACE)/2.6.5/$(GOOS)_$(GOARCH)
endif
build13: fmtcheck
	go mod tidy && go mod vendor
	@echo "==> Installing plugin to $(DESTINATION)"
	@mkdir -p $(DESTINATION)
	go build -o $(DESTINATION)/terraform-provider-zia_v2.6.5

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

tools:
	@which $(GOFMT) || go install mvdan.cc/gofumpt@v0.4.0
	@which $(TFPROVIDERLINT) || go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.28.1
	@which $(STATICCHECK) || go install honnef.co/go/tools/cmd/staticcheck@v0.4.2

tools-update:
	@go install mvdan.cc/gofumpt@v0.4.0
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.28.1
	@go install honnef.co/go/tools/cmd/staticcheck@v0.4.2

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

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
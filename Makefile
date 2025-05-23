PLUGIN_NAME := steampipe-plugin-yandexcloud
PLUGIN_VERSION := 0.0.1
PLUGIN_BIN := $(PLUGIN_NAME).plugin
STEAMPIPE_HOME ?= $(HOME)/.steampipe/plugins/local/$(PLUGIN_NAME)/$(PLUGIN_VERSION)

GO      ?= go
LDFLAGS := -s -w -X "main.version=$(PLUGIN_VERSION)"
GO_MODFILE ?= go.mod

.PHONY: all build install clean deps

all: build

# Ensure go.sum is up to date before building
deps:
	$(GO) mod tidy

build: deps
	$(GO) build -modfile=$(GO_MODFILE) -ldflags="$(LDFLAGS)" -o $(PLUGIN_BIN) .

install: build
	@mkdir -p $(STEAMPIPE_HOME)
	@mv $(PLUGIN_BIN) $(STEAMPIPE_HOME)/
	@echo "Installed $(PLUGIN_BIN) to $(STEAMPIPE_HOME)"

clean:
	rm -f $(PLUGIN_BIN)

# ------------------------------------------------------------
# LOCAL BUILD AND INSTALLATION FOR STEAMPIPE
# ------------------------------------------------------------

# Plugin version (can be passed as make VERSION=0.0.2)
VERSION ?= 0.0.1

# Target platform (default is current)
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Where to copy .plugin
PLUGIN_NAME := yandexcloud
PLUGIN_DIR  := $(HOME)/.steampipe/plugins/local/$(PLUGIN_NAME)
PLUGIN_FILE := $(PLUGIN_DIR)/$(PLUGIN_NAME).plugin

# Build static binary for local development and run tests
local-build:
	@echo "=> Building $(PLUGIN_FILE) (GOOS=$(GOOS) GOARCH=$(GOARCH))"
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		$(GO) build -modfile=$(GO_MODFILE) -ldflags="-s -w -X 'main.version=$(VERSION)'" \
		-o $(PLUGIN_BIN) .
	$(GO) test -modfile=$(GO_MODFILE) ./yandexcloud

# Copy to Steampipe directory
local-install: local-build
	@echo "=> Installing into $(PLUGIN_DIR)"
	mkdir -p $(PLUGIN_DIR)
	cp -f $(PLUGIN_BIN) $(PLUGIN_FILE)
	chmod +x $(PLUGIN_FILE)
	@echo "✓ Installed. Restart Steampipe service:  steampipe service restart"
	steampipe service restart
	make tests

tests:	
	steampipe query yandexcloud-test/tests/yandexcloud_compute_instance/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_snapshot/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_image/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_disk/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_filesystem/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_placement_group/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_host_group/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_gpu_cluster/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_disk_placement_group/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_snapshot_schedule/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_reserved_instance_pool/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_zone/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_disk_type/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_host_type/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_compute_operation/test-get-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_vpc_network/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_vpc_subnet/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_vpc_route_table/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_vpc_security_group/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_vpc_address/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_vpc_gateway/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_vpc_operation/test-list-query.sql
	steampipe query yandexcloud-test/tests/yandexcloud_billing_resource_usage/test-list-query.sql	
	steampipe query yandexcloud-test/tests/yandexcloud_billing_account/test-list-query.sql

# Clean local plugin directory
local-clean:
	rm -f $(PLUGIN_BIN)
	rm -f $(PLUGIN_FILE)
	rmdir -p --ignore-fail-on-non-empty $(PLUGIN_DIR)

.PHONY: local-build local-install local-clean
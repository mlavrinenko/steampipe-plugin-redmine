PLUGIN_NAME = steampipe-plugin-redmine
INSTALL_DIR ?= $(HOME)/.steampipe
PLUGIN_DIR = $(INSTALL_DIR)/plugins/ghcr.io/mlavrinenko/redmine

.PHONY: build test test-cover fmt vet lint tidy install clean uninstall setup-refs

build:
	go build -o $(PLUGIN_NAME).plugin -tags netgo *.go

test:
	go test ./... -v

test-cover:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(PLUGIN_NAME).plugin $(PLUGIN_DIR)/
	@# Also install to versioned path if it exists (Steampipe may use either)
	@if [ -d "$(INSTALL_DIR)/plugins/ghcr.io/mlavrinenko/$(PLUGIN_NAME)@"* ]; then \
		for d in $(INSTALL_DIR)/plugins/ghcr.io/mlavrinenko/$(PLUGIN_NAME)@*/; do \
			cp $(PLUGIN_NAME).plugin "$$d"; \
			echo "Also installed to $$d"; \
		done; \
	fi
	cp -rn config/* $(INSTALL_DIR)/config/ 2>/dev/null || true
	@echo "Plugin installed to $(PLUGIN_DIR)"
	@echo "Ensure $(INSTALL_DIR)/config/redmine.spc exists with your credentials"

clean:
	rm -f $(PLUGIN_NAME).plugin
	rm -f coverage.out

uninstall:
	rm -rf $(PLUGIN_DIR)
	@echo "Plugin uninstalled from $(PLUGIN_DIR)"

setup-refs:
	@mkdir -p .res
	@echo "Cloning reference repositories to .res/..."
	@for repo_spec in \
		"https://github.com/mattn/go-redmine.git:.res/go-redmine" \
		"https://github.com/nixys/nxs-go-redmine.git:.res/nxs-go-redmine" \
		"https://github.com/turbot/steampipe-docs.git:.res/steampipe-docs" \
		"https://github.com/turbot/steampipe-plugin-github.git:.res/steampipe-plugin-github" \
		"https://github.com/turbot/steampipe-plugin-sdk.git:.res/steampipe-plugin-sdk"; \
	do \
		url=$${repo_spec%:*}; \
		dest=$${repo_spec##*:}; \
		repo_name=$$(basename "$$dest"); \
		if [ -d "$$dest/.git" ]; then \
			echo "  $$repo_name already exists, skipping"; \
		else \
			echo "  Cloning $$repo_name..."; \
			git clone --depth 1 "$$url" "$$dest"; \
		fi; \
	done
	@echo "Reference repositories ready in .res/"

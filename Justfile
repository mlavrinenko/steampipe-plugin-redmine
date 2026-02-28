# steampipe-plugin-redmine

plugin_name := "steampipe-plugin-redmine"
install_dir := env("STEAMPIPE_INSTALL_DIR", env("HOME") / ".steampipe")
plugin_dir := install_dir / "plugins" / "local" / "redmine"

# List available commands
default:
    @just --list

# Clone reference repositories to .res/ (for development reference)
setup-refs:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p .res
    echo "Cloning reference repositories to .res/..."

    repos=(
        "https://github.com/mattn/go-redmine.git:.res/go-redmine"
        "https://github.com/nixys/nxs-go-redmine.git:.res/nxs-go-redmine"
        "https://github.com/turbot/steampipe-docs.git:.res/steampipe-docs"
        "https://github.com/turbot/steampipe-plugin-github.git:.res/steampipe-plugin-github"
        "https://github.com/turbot/steampipe-plugin-sdk.git:.res/steampipe-plugin-sdk"
    )

    for repo_spec in "${repos[@]}"; do
        url="${repo_spec%%:*}"
        dest="${repo_spec##*:}"
        repo_name="$(basename "$dest")"

        if [ -d "$dest/.git" ]; then
            echo "✓ $repo_name already exists, skipping"
        else
            echo "→ Cloning $repo_name..."
            git clone --depth 1 "$url" "$dest"
        fi
    done

    echo "✓ Reference repositories ready in .res/"

# Build the plugin binary
build:
    go build -o {{ plugin_name }}.plugin -tags netgo *.go

# Run all tests
test:
    go test ./... -v

# Run tests with coverage
test-cover:
    go test ./... -v -coverprofile=coverage.out
    go tool cover -func=coverage.out

# Format code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Run golangci-lint
lint:
    golangci-lint run ./...

# Tidy module dependencies
tidy:
    go mod tidy

# Build and install plugin locally for steampipe
install: build
    mkdir -p {{ plugin_dir }}
    cp {{ plugin_name }}.plugin {{ plugin_dir }}/
    cp -rn config/* {{ install_dir }}/config/ 2>/dev/null || true
    @echo "Plugin installed to {{ plugin_dir }}"
    @echo "Ensure {{ install_dir }}/config/redmine.spc exists with your credentials"

# Remove built artifacts
clean:
    rm -f {{ plugin_name }}.plugin
    rm -f coverage.out

# Uninstall plugin from steampipe
uninstall:
    rm -rf {{ plugin_dir }}
    @echo "Plugin uninstalled from {{ plugin_dir }}"

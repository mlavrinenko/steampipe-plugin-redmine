# steampipe-plugin-redmine

plugin_name := "steampipe-plugin-redmine"
install_dir := env("STEAMPIPE_INSTALL_DIR", env("HOME") / ".steampipe")
plugin_dir := install_dir / "plugins" / "local" / "redmine"

# List available commands
default:
    @just --list

# Build the plugin binary
build:
    go build -o {{plugin_name}}.plugin -tags netgo *.go

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

# Run linter
vet:
    go vet ./...

# Tidy module dependencies
tidy:
    go mod tidy

# Build and install plugin locally for steampipe
install: build
    mkdir -p {{plugin_dir}}
    cp {{plugin_name}}.plugin {{plugin_dir}}/
    cp -rn config/* {{install_dir}}/config/ 2>/dev/null || true
    @echo "Plugin installed to {{plugin_dir}}"
    @echo "Ensure {{install_dir}}/config/redmine.spc exists with your credentials"

# Remove built artifacts
clean:
    rm -f {{plugin_name}}.plugin
    rm -f coverage.out

# Uninstall plugin from steampipe
uninstall:
    rm -rf {{plugin_dir}}
    @echo "Plugin uninstalled from {{plugin_dir}}"

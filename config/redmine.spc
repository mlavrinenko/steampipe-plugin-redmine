connection "redmine" {
  plugin = "ghcr.io/mlavrinenko/steampipe-plugin-redmine@latest"

  # Redmine instance URL (required).
  # Can also be set with the REDMINE_ENDPOINT environment variable.
  # endpoint = "https://www.redmine.org"

  # API key from My Account -> API access key (required).
  # Can also be set with the REDMINE_API_KEY environment variable.
  # api_key = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
}

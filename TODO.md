# Pre-Publication TODO

Manual tasks required before publishing.

## GHCR OCI Registry Setup

One-time setup to enable `steampipe plugin install ghcr.io/mlavrinenko/steampipe-plugin-redmine@latest`.

- [x] **Make the GHCR package public.** After the first release pushes the OCI image,
  the package is created as **private** by default. Go to:
  `https://github.com/users/mlavrinenko/packages/container/steampipe-plugin-redmine/settings`
  and change visibility to **Public**.

- [x] **Verify installation** after the first release:
  ```bash
  steampipe plugin install ghcr.io/mlavrinenko/steampipe-plugin-redmine@latest
  ```
  Then confirm the connection config was created at `~/.steampipe/config/redmine.spc`
  and the plugin binary landed in `~/.steampipe/plugins/`.

### How it works

The release workflow (`.github/workflows/release.yml`) triggers on `v*` tags:

1. GoReleaser builds binaries for linux/darwin x amd64/arm64.
2. ORAS pushes a single OCI artifact to `ghcr.io/mlavrinenko/steampipe-plugin-redmine:{version}`
   containing platform-specific binary layers, docs, and config with Steampipe's custom media types.
3. Non-RC versions are also tagged `latest`.

To create a release: `git tag v0.1.0 && git push origin v0.1.0`

## Hub Assets (request from Turbot)

- [ ] **Request plugin icon (SVG)** via [Turbot Community Slack](https://turbot.com/community/join).
  Once received, add to `docs/index.md` front matter:
  ```yaml
  icon_url: "/images/plugins/mlavrinenko/redmine.svg"
  ```

- [ ] **Request social graphic (PNG)** via [Turbot Community Slack](https://turbot.com/community/join).
  Once received:
  1. Add to `docs/index.md` front matter:
     ```yaml
     og_image: "/images/plugins/mlavrinenko/redmine-social-graphic.png"
     ```
  2. Add to the top of `README.md` as an image.
  3. Upload as the GitHub repository Social Preview under Settings -> General.

## GitHub Repository Settings

- [ ] **Set repository topics**: `postgresql`, `postgresql-fdw`, `sql`, `steampipe`, `steampipe-plugin`.

- [ ] **Set repository website** to the Hub URL once published:
  `https://hub.steampipe.io/plugins/mlavrinenko/redmine`

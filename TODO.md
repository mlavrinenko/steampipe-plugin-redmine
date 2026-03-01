# Pre-Publication TODO

Manual tasks required before publishing to the Steampipe Hub.

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

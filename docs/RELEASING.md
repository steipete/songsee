---
summary: "Release checklist for songsee (GitHub release + Homebrew tap + Pages)"
---

# Releasing songsee

Follow these steps for each release. Title GitHub releases as `songsee <version>`.

## Checklist
- Ensure `CHANGELOG.md` has the new version section.
- Tag the release: `git tag -a v<version> -m "Release <version>"` and push tags after commits.
- Verify tests + lint: `go test ./... -cover` and `golangci-lint run`.
- Build source archive for Homebrew: `git archive --format=tar.gz --output /tmp/songsee-<version>.tar.gz v<version>`.
- Compute checksums: `shasum -a 256 /tmp/songsee-<version>.tar.gz`.
- Update `../homebrew-tap/Formula/songsee.rb` with the new tarball URL + sha256 and version.
- Update tap README to list songsee if needed.
- Commit + push changes in songsee and the tap.
- Create GitHub release for `v<version>`:
  - Title: `songsee <version>`
  - Body: bullets from `CHANGELOG.md` for that version
- Verify Homebrew install: `brew update && brew reinstall steipete/tap/songsee && songsee --version`.
- Verify Pages build is green (workflow `pages-build-deployment`).

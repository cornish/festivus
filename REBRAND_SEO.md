# Textivus Rebrand + GitHub Discovery (SEO) Guide

This doc is a **practical checklist** for rebranding the project to **Textivus** and optimizing the repo for **GitHub + Google discovery**.

Goal: people searching for a **terminal/TUI text editor** should find this project quickly.

---

## 1) Naming decisions

### Recommended product identity
- **Project name:** Textivus
- **Tagline:** “A text editor for the rest of us!”

Reasoning:
- “Textivus” communicates *text editor* immediately, while still being a Festivus homage.
- “Festivus” is culturally crowded and has package-name collisions.
- “Textivus” is more unique + searchable.

---

## 2) Repo naming + URL

### Best repo names (choose one)
1. `textivus` (clean, memorable)
2. `textivus-editor` (more keyword searchable)

If you already have a repo named “festivus”, rename it to `textivus` or create a new `textivus` repo and archive/redirect the old one.

---

## 3) GitHub repo metadata (TOP RIGHT of repo page)

These fields strongly influence GitHub search and Google snippets.

### Description (copy/paste option)
> A friendly terminal (TUI) text editor for Linux — fast, familiar, and easy to use.

### Website
- Use GitHub Pages if you can: `https://<user>.github.io/textivus/`
- Otherwise link to your release page or README anchor.

### Topics (do this aggressively)
Add these GitHub Topics:
- `text-editor`
- `terminal`
- `tui`
- `cli`
- `console`
- `linux`
- `editor`

If applicable:
- `bubbletea`
- `golang`
- `vim` (ONLY if you support a vim keymap)
- `nano`
- `micro`

Tip: topics are the closest thing GitHub has to “SEO tags”.

---

## 4) README = your landing page

The README should function as a landing page (pitch + demo + install) rather than a manual.

### The first 10 lines matter most
People search for terms like:
- “terminal text editor”
- “TUI editor”
- “linux console editor”
- “nano-like editor”
- “micro-like editor”

Make sure those exact phrases appear near the top.

### Suggested README opening block
```md
# Textivus

**Textivus** is a fast, friendly **terminal (TUI) text editor for Linux** inspired by the simplicity of **nano/micro**, with modern comforts like multi-file buffers, incremental search/replace, and syntax highlighting.

> A text editor for the rest of us!
```

### Include ONE demo GIF
Add `docs/demo.gif` and show it near the top:
```md
![Textivus terminal text editor demo](docs/demo.gif)
```

Use ALT text with keywords (“terminal text editor”, “TUI editor”) — it is indexed.

---

## 5) Install UX (high impact for visibility)

### Provide GitHub releases with binaries
At minimum:
- Linux x86_64
- Linux arm64

Nice-to-have:
- macOS arm64 (low effort, large audience)
- macOS x86_64 (optional)

### Install script
Host:
- `install.sh` at repo root

Document:
```sh
curl -fsSL https://raw.githubusercontent.com/<you>/textivus/main/install.sh | sh
```

Also provide the safer two-step:
```sh
curl -fsSL https://raw.githubusercontent.com/<you>/textivus/main/install.sh -o install_textivus.sh
sh install_textivus.sh
```

### Installer behavior checklist
- installs to `~/.local/bin` by default
- supports `--version` and `--bin-dir`
- verifies download via `checksums.txt` (sha256)
- prints PATH instructions if needed

---

## 6) Canonical repo files people expect

Add these at the repo root:
- `LICENSE` (MIT)
- `README.md`
- `CHANGELOG.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md` (optional but good)
- `SECURITY.md` (optional)

These increase trust and help discovery/contribution.

---

## 7) Releases strategy (visibility + trust)

- Tag releases: `v0.1.0`, `v0.1.1`, ...
- Attach binaries + `checksums.txt`
- Maintain `CHANGELOG.md`

GitHub promotes repos that have releases and activity.

---

## 8) GitHub Pages mini-site (optional but strong SEO)

Even a single-page site helps dominate Google for your project name.

Minimal content:
- Textivus logo/banner
- demo gif
- install snippet
- feature bullets
- link to GitHub + releases

---

## 9) Packaging roadmap (discovery + adoption)

Packaging pages rank on Google and function as “distribution channels”.

Recommended order:
1. **AUR** (`textivus-bin`)
2. **Homebrew** formula
3. `.deb` + `.rpm` assets
4. Fedora **COPR**
5. Snap/Flatpak (optional)

Important: don’t overcommit early. Ask for community maintainers.

---

## 10) Search collision mitigation

### Make “Textivus” unambiguous
Use consistent naming everywhere:
- binary: `textivus`
- project name in help output: `Textivus`
- repo name: `textivus`
- tagline: “text editor” words near top

This ensures Google results point to your project.

---

## 11) Short “non-goals” section (prevents flamewars)

Add to README:
```md
## Non-goals
Textivus is not an IDE.
- No always-on language servers
- No background project indexing
- No plugin marketplace
```

This reduces hostile “yet another editor” takes and keeps expectations aligned.

---

## 12) Action checklist (copy into issue)

- [ ] Rename repo to `textivus` or `textivus-editor`
- [ ] Update GitHub description
- [ ] Add topics (`text-editor`, `tui`, `terminal`, `linux`, `cli`, ...)
- [ ] Update README first paragraph with keywords
- [ ] Add demo GIF
- [ ] Add `install.sh`
- [ ] Ensure releases include binaries + checksums
- [ ] Add `CHANGELOG.md`
- [ ] Add `CONTRIBUTING.md`
- [ ] (Optional) GitHub Pages

---

## 13) Claude Code instructions (for the rename)

When doing the rename, ensure:
- replace user-facing name strings: “Festivus” → “Textivus”
- update binary name if applicable
- update config file naming (e.g., `.textivus.toml`)
- update docs and examples
- update module path (if Go): `module github.com/<you>/textivus`
- add redirect notes for old name (if necessary)

---

*End of doc.*

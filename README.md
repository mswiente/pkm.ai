# pkm

A lightweight CLI for a markdown-first personal knowledge management system built around [Obsidian](https://obsidian.md), Claude Code, and GitHub Copilot CLI.

The core idea: **capture fast, curate later**. Every note lands in `00-inbox/` first. You decide where it goes.

## Installation

```bash
git clone https://github.com/mswiente/pkm.ai
cd pkm.ai
make install
```

Requires Go 1.25+.

## Setup

### 1. Create your vault

```bash
bash setup/setup_pkm_vault.sh ~/Documents/pkm-vault
```

This creates the folder structure and copies the default templates into the vault.

### 2. Point pkm at your vault

```bash
pkm config --set-vault-path ~/Documents/pkm-vault
```

### 3. (Optional) Install Claude Code skills

```bash
pkm skill install
```

Installs available slash commands (e.g. `/capture-plan`) into `~/.claude/commands/`.

## Commands

### `pkm capture`

Capture a new note to `00-inbox/`.

```bash
pkm capture "Quick thought about X"
pkm capture --title "Meeting summary" --source claude-code --tags "ai,work"
echo "paste from clipboard" | pkm capture
pkm capture --clipboard
pkm capture --editor          # opens in $VISUAL/$EDITOR before saving
```

**Flags:** `--title`, `--source`, `--tags`, `--type-hint`, `--editor`, `--clipboard`

Valid sources: `manual`, `chatgpt`, `claude-code`, `copilot-cli`, `readwise`, `other`

---

### `pkm daily create`

Create today's daily note in `01-daily/`. Idempotent — safe to run multiple times.

```bash
pkm daily create
pkm daily create --date 2026-03-15
pkm daily create --open          # open in editor after creation
```

---

### `pkm meeting create`

Create a meeting note from the meeting template, placed in `00-inbox/`.

```bash
pkm meeting create --title "Q2 Planning" --participants "Alice, Bob" --project "roadmap"
pkm meeting create --date 2026-04-10 --title "Retro"
```

**Flags:** `--title`, `--date`, `--participants`, `--project`

---

### `pkm decision create`

Create a decision note from the decision template, placed in `00-inbox/`.

```bash
pkm decision create --title "Switch to PostgreSQL" --status accepted
pkm decision create --title "API versioning strategy" --project "platform" --from-stdin < notes.md
```

**Flags:** `--title`, `--project`, `--status` (`draft`|`accepted`|`superseded`), `--from-stdin`

---

### `pkm process inbox`

Analyze inbox notes and print a structured report — useful as context for a Claude Code session.

```bash
pkm process inbox
pkm process inbox --file 2026-04-01-1200-manual-meeting-notes.md
pkm process inbox --interactive   # walk through each note with a single-key action menu
pkm process inbox --dry-run
```

No files are modified unless `--apply` is passed.

---

### `pkm note move`

Move a note to a target vault folder and update its frontmatter automatically.

```bash
pkm note move 2026-04-01-1200-manual-ai-notes.md knowledge
pkm note move some-note.md 4                  # shorthand: 4 = 04-knowledge
pkm note move some-note.md archive --dry-run
```

Target folders: `knowledge` (4), `projects` (2), `areas` (3), `resources` (5), `decisions` (6), `archive` (9)

Frontmatter updated automatically:
- `type` — inferred from destination folder
- `status` — `inbox → draft`; any `→ archived` when moving to archive
- `updated` — set to today

**Flags:** `--dry-run`, `--type`, `--status`

---

### `pkm sync readwise`

Sync saved articles from [Readwise Reader](https://readwise.io/read) to your vault inbox. Incremental: only new and updated articles are fetched on each run.

#### First-time setup

```bash
pkm sync readwise auth
```

Opens a prompt for your API token (get it at `https://readwise.io/access_token`). Validates the token and saves it to `~/.config/pkm/config.yaml`. Alternatively, set the `READWISE_TOKEN` environment variable.

#### Syncing

```bash
pkm sync readwise                         # incremental sync (uses state file)
pkm sync readwise --dry-run               # preview without writing files
pkm sync readwise --since 2026-01-01      # override state file; sync from a specific date
pkm sync readwise --limit 20              # cap at 20 articles
```

Each synced article becomes a note in `00-inbox/` with:
- Frontmatter: `type: resource`, `status: inbox`, `source: readwise`, plus any Readwise tags
- Metadata block: source URL, author, saved date, reading time
- Article summary
- `## Highlights` section with all your Reader highlights

Sync state is stored at `~/.config/pkm/readwise_sync_state.json`. Subsequent runs only fetch articles updated since the last successful sync.

---

### `pkm config`

Show or update pkm configuration.

```bash
pkm config --show
pkm config --set-vault-path ~/Documents/pkm-vault
```

Config file: `~/.config/pkm/config.yaml`

| Key | Default | Description |
|-----|---------|-------------|
| `vault_path` | — | Path to your Obsidian vault (required) |
| `inbox_path` | `00-inbox` | Inbox subfolder |
| `daily_path` | `01-daily` | Daily notes subfolder |
| `templates_path` | `07-templates` | Templates subfolder |
| `editor` | `$VISUAL`/`$EDITOR`/`vi` | Editor for `--editor` flag |
| `filename_timezone` | `UTC` | Timezone used in filenames |
| `default_source` | `manual` | Default source for captures |
| `readwise_token` | — | Readwise API token (or use `READWISE_TOKEN` env var) |

---

### `pkm skill`

Manage Claude Code slash commands bundled with pkm.

```bash
pkm skill list                   # list available skills
pkm skill install                # install all skills
pkm skill install capture-plan   # install a specific skill
```

After installation, invoke in a Claude Code session with e.g. `/capture-plan`.

## Vault Structure

```
00-inbox/       # all new captures land here first
01-daily/       # daily notes (YYYY-MM-DD.md)
02-projects/    # project tracking
03-areas/       # areas of responsibility
04-knowledge/   # knowledge base (evergreen notes)
05-resources/   # external resources and references
06-decisions/   # decision log
07-templates/   # note templates
08-attachments/ # images and files
09-archive/     # archived notes
```

## Filename Convention

| Type | Pattern | Example |
|------|---------|---------|
| Inbox | `YYYY-MM-DD-HHMM-source-slug.md` | `2026-04-05-0900-manual-api-design-notes.md` |
| Daily | `YYYY-MM-DD.md` | `2026-04-05.md` |

## Frontmatter Schema

```yaml
title: Note title
type: resource           # inbox | daily | project | knowledge | resource | decision | template | meeting | troubleshooting
status: inbox            # inbox | draft | active | evergreen | archived
source: manual           # manual | chatgpt | claude-code | copilot-cli | readwise | other
created: 2026-04-05
tags: [tag1, tag2]
```

## Development

```bash
make build          # build to bin/pkm
make install        # install to $GOPATH/bin
make test           # run tests
make fmt            # format code
make vet            # run go vet
```

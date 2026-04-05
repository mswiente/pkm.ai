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

Installs available slash commands (`/capture-plan`, `/distill-inbox`) into `~/.claude/commands/`.

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

**Flags:** `--title`, `--source`, `--tags`, `--type-hint`, `--editor`, `--clipboard`, `--update`

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

Prints a structured report of inbox notes — the primary context-gathering tool for a Claude Code session.

```bash
pkm process inbox                  # summary report of all inbox notes
pkm process inbox --full           # full note bodies + current 04-knowledge/index.md
pkm process inbox --file note.md   # single note
pkm process inbox --interactive    # walk through each note with a single-key action menu
pkm process inbox --dry-run
```

**`--full` mode** is designed for distillation sessions: it outputs the complete body of every inbox note followed by the current `04-knowledge/index.md`, giving Claude Code everything it needs in a single pass to propose how each note should be processed and which topic pages to update.

After reviewing the report, Claude Code calls `pkm knowledge` and `pkm note move` to apply changes.

---

### `pkm knowledge`

Primitives for maintaining the `04-knowledge/` wiki. Called by Claude Code after it has analyzed inbox notes.

#### `pkm knowledge append-topic`

Appends markdown content (from stdin) to a `04-knowledge/<slug>.md` topic page. Creates the page with frontmatter if it does not exist.

```bash
echo "## From [[2026-04-03-readwise-article]]

Key insight extracted as evergreen prose." \
  | pkm knowledge append-topic ai-agents --title "AI Agents"

pkm knowledge append-topic resilience-patterns --title "Resilience Patterns" \
  --dry-run < content.md
```

#### `pkm knowledge update-index`

Adds a `[[slug]] — description` entry to `04-knowledge/index.md`. No-op if the slug is already present.

```bash
pkm knowledge update-index ai-agents \
  --description "agentic systems, LLM autonomy, multi-agent patterns"
```

#### `pkm knowledge append-log`

Appends a processing entry to `04-knowledge/log.md` (the audit trail).

```bash
pkm knowledge append-log \
  --note 2026-04-03-readwise-it-has-never-been-about-code.md \
  --action "distill + file" \
  --filed-to resources \
  --updated software-engineering-philosophy \
  --created llm-patterns
```

| Flag | Description |
|------|-------------|
| `--note` | Basename of the processed note (required) |
| `--action` | Action taken, e.g. `"distill + file"` (required) |
| `--filed-to` | Target folder the source note was moved to |
| `--updated` | Comma-separated slugs of updated topic pages |
| `--created` | Comma-separated slugs of newly created topic pages |

---

### `pkm project`

Manage project context notes in `02-projects/`. Each project note has four sections:

- **Intent** — what the project is trying to achieve (stable)
- **Current Status** — what was last worked on (updated each session)
- **Next Steps** — what to do next (updated each session)
- **Plan History** — dated log of captured plans (append-only)

#### `pkm project update`

Create or update a project note. Only sections with provided flags are updated; others are left unchanged. If stdin contains content, it is appended to the Plan History section.

```bash
# Create a new project
pkm project update pkm-ai \
  --title "pkm.ai CLI" \
  --intent "A lightweight CLI for a markdown-first PKM system." \
  --current-status "Readwise sync and distill workflow implemented." \
  --next-steps "- [ ] Add project management commands"

# Update after a work session (with plan content from stdin)
cat plan.md | pkm project update pkm-ai \
  --current-status "Implemented pkm project commands." \
  --next-steps "- [ ] Update distill skill\n- [ ] Write tests" \
  --plan-heading "2026-04-06 — Project management"

# Preview without writing
pkm project update pkm-ai --current-status "..." --dry-run
```

**Flags:** `--title`, `--intent`, `--current-status`, `--next-steps`, `--plan-heading`, `--status` (`active`|`on-hold`|`archived`), `--dry-run`

#### `pkm project list`

List all project notes with their status.

```bash
pkm project list
# ● pkm-ai                         pkm.ai CLI
# ○ raycast-pkm                    Raycast PKM Extension
```

Markers: `●` active, `○` on-hold, `–` archived

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
02-projects/    # project tracking (managed by pkm project)
03-areas/       # areas of responsibility
04-knowledge/   # knowledge base (evergreen topic pages)
  index.md      #   topic discovery map (managed by pkm knowledge update-index)
  log.md        #   processing audit trail (managed by pkm knowledge append-log)
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

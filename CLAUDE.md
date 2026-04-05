# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Repo Is

`pkm.ai` is the source repository for the `pkm` CLI — a lightweight command-line tool for a markdown-first personal knowledge system built around Obsidian, Claude Code, and GitHub Copilot CLI.

The Obsidian vault lives separately at the user-configured path (default: `~/Library/Mobile Documents/iCloud~md~obsidian/Documents/pkm-vault`). This repo contains only the CLI source and templates.

## Actual CLI Structure (Go)

```
cmd/pkm/              # main entry point
internal/capture/     # pkm capture logic
internal/process/     # pkm process inbox + knowledge management
internal/daily/       # pkm daily create
internal/meeting/     # pkm meeting create
internal/decision/    # pkm decision create
internal/note/        # pkm note move
internal/readwise/    # Readwise Reader API client and sync
internal/cli/         # cobra command wiring
internal/config/      # config loading (~/.config/pkm/config.yaml)
internal/vault/       # vault filesystem operations
internal/frontmatter/ # frontmatter parsing and marshalling
internal/slug/        # filename slug generation
internal/templates/   # markdown template rendering
templates/            # markdown templates (inbox, daily, meeting, decision, knowledge, troubleshooting, project, resource)
skills/               # Claude Code slash command definitions
```

## Implemented CLI Commands

```bash
pkm capture [text] [--title] [--source] [--tags] [--type-hint] [--editor] [--clipboard] [--update]
pkm process inbox [--file] [--all] [--full] [--dry-run] [--apply] [--interactive]
pkm daily create [--date] [--open]
pkm meeting create [--title] [--date] [--project] [--participants]
pkm decision create [--title] [--project] [--status] [--from-stdin]
pkm note move <filename> <folder> [--type] [--status] [--dry-run]
pkm knowledge append-topic <slug> --title <title> [--dry-run]   # reads content from stdin
pkm knowledge update-index <slug> --description <desc> [--dry-run]
pkm knowledge append-log --note <file> --action <action> [--filed-to] [--updated] [--created]
pkm project update <slug> [--title] [--intent] [--current-status] [--next-steps] [--plan-heading] [--status] [--dry-run]
pkm project list
pkm sync readwise [--dry-run] [--since] [--limit]
pkm sync readwise auth
pkm config --show | --set-vault-path <path>
pkm skill list | install [name]
```

## Core Behavioral Rules for the CLI

- `pkm capture` always writes to `00-inbox/` — never to other folders directly
- `pkm process inbox` must never move or modify files without explicit user confirmation
- No file is silently created outside the configured vault path
- `--dry-run` must be available for all mutating operations
- AI suggestions must remain inspectable before any changes are applied

## Vault Folder Structure

```
00-inbox/       # all new captures land here first
01-daily/       # daily notes (YYYY-MM-DD.md)
02-projects/    # project tracking
03-areas/       # areas of responsibility
04-knowledge/   # evergreen knowledge base
  index.md      #   topic discovery map (managed by pkm knowledge update-index)
  log.md        #   processing audit trail (managed by pkm knowledge append-log)
05-resources/   # external articles, references, saved reading
06-decisions/   # decision records
07-templates/   # vault-local copies of templates/
08-attachments/ # images and files
09-archive/     # archived notes
```

## Filename Convention

Inbox notes: `YYYY-MM-DD-HHMM-source-slug.md`
Daily notes: `YYYY-MM-DD.md`

Valid `source` values: `manual`, `chatgpt`, `claude-code`, `copilot-cli`, `readwise`, `other`

## Frontmatter Schema

Mandatory fields: `title`, `type`, `status`, `source`, `created`
Optional: `tags`, `updated`, `type_hint`

Valid `type` values: `inbox`, `daily`, `project`, `knowledge`, `resource`, `decision`, `template`, `meeting`, `troubleshooting`
Valid `status` values: `inbox`, `draft`, `active`, `evergreen`, `archived`

## Config File

`~/.config/pkm/config.yaml` — keys: `vault_path`, `inbox_path`, `daily_path`, `templates_path`, `editor`, `filename_timezone`, `default_source`, `readwise_token`

## Inbox Distillation Workflow (Karpathy pattern)

The primary way to process inbox notes in a Claude Code session. Use `/distill-inbox` skill or follow this pattern manually:

1. `pkm process inbox --full` — outputs every inbox note's full content + `04-knowledge/index.md`
2. For each note, classify as **plan note** or **regular note** (see below), then decide action
3. Apply using CLI primitives — Claude Code calls these, user reviews before each step

### Regular notes
Decide: **distill** (extract insights into knowledge base) | **file-only** (move, no extraction) | **skip** (archive)

```bash
printf '%s\n' '<content>' | pkm knowledge append-topic <slug> --title "<Title>"
pkm knowledge update-index <slug> --description "<desc>"   # new topics only
pkm note move <filename> <folder>
pkm knowledge append-log --note <file> --action <action> --filed-to <folder> ...
```

### Plan notes (captured plans from Claude Code sessions)

A note is a **plan note** if: `source: claude-code` AND body contains a `# Plan:` heading OR `type_hint: knowledge`.

Do NOT distill plan notes into `04-knowledge/`. Instead, route to `02-projects/`:

```bash
printf '%s\n' '<plan-body>' | pkm project update <slug> \
  --title "<Project Title>" \
  --current-status "<what was built/decided>" \
  --next-steps "- [ ] Next action" \
  --plan-heading "<YYYY-MM-DD — session description>"
pkm note move <filename> archive
```

Each project note has: **Intent** (stable goal), **Current Status** (last session), **Next Steps** (checklist), **Plan History** (append-only dated entries).

`04-knowledge/index.md` is the topic discovery map — always read it before creating new topic slugs to avoid duplicates.
`04-knowledge/log.md` is append-only — never modify, only append via `pkm knowledge append-log`.

## Skills (Claude Code slash commands)

Install with: `pkm skill install`

| Skill | Description |
|-------|-------------|
| `/capture-plan` | Capture the current plan as a PKM inbox note |
| `/distill-inbox` | Process inbox notes using the Karpathy wiki pattern |

## Design Principles

- Markdown is the source of truth — no proprietary formats
- Capture and curation are intentionally separated (capture fast → process later)
- Human-in-the-loop: AI assists, never autonomously restructures the vault
- The CLI must work non-interactively (for pipes/scripts) and interactively
- Output must be tool-independent: works with Claude Code, Copilot CLI, shell scripts, Raycast

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Repo Is

`pkm.ai` is the source repository for the `pkm` CLI — a lightweight command-line tool for a markdown-first personal knowledge system built around Obsidian, Claude Code, and GitHub Copilot CLI. The CLI has been scaffolded (see `setup/setup_pkm_repo.sh`) but the Go implementation is not yet written.

The Obsidian vault itself lives separately at `~/Documents/Obsidian/pkm-vault` (or a user-configured path). This repo contains the CLI source and templates only.

## Planned CLI Structure (Go)

```
cmd/pkm/          # main entry point
internal/capture/ # pkm capture logic
internal/process/ # pkm process inbox logic
internal/templates/ # template rendering
internal/config/  # config loading (~/.config/pkm/config.yaml)
internal/fs/      # vault file system operations
templates/        # markdown templates (inbox, daily, meeting, decision, knowledge, troubleshooting, project, resource)
```

## CLI Commands to Implement

```bash
pkm capture [text] [--title] [--source] [--tags] [--type-hint] [--editor] [--clipboard]
pkm process inbox [--file] [--all] [--dry-run] [--apply] [--interactive]
pkm daily create [--date] [--open]
pkm meeting create [--title] [--date] [--project] [--participants] [--inbox]
pkm decision create [--title] [--project] [--status] [--from-stdin] [--inbox]
```

## Core Behavioral Rules for the CLI

- `pkm capture` always writes to `00-inbox/` — never to other folders directly
- `pkm process inbox` must never move or modify files without explicit user confirmation (unless `--apply` is passed)
- No file is silently created outside the configured vault path
- `--dry-run` must be available for all AI-assisted mutating operations
- AI suggestions must remain inspectable before any changes are applied

## Vault Folder Structure

```
00-inbox/       # all new captures land here first
01-daily/       # daily notes (YYYY-MM-DD.md)
02-projects/
03-areas/
04-knowledge/
05-resources/
06-decisions/
07-templates/   # copies of templates/ in this repo
08-attachments/
09-archive/
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

## Setup Scripts

```bash
./setup/setup_pkm_vault.sh [vault_path]   # creates vault folder structure and templates
./setup/setup_pkm_repo.sh [repo_path]     # scaffolds this repo's directory layout
```

## Design Principles

- Markdown is the source of truth — no proprietary formats
- Capture and curation are intentionally separated (capture fast → process later)
- Human-in-the-loop: AI assists, never autonomously restructures the vault
- The CLI must work non-interactively (for pipes/scripts) and interactively
- Output must be tool-independent: works with Claude Code, Copilot CLI, shell scripts, Raycast

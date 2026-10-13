---
description: Process inbox notes using the Karpathy wiki pattern — distill insights into 04-knowledge/ topic pages, file sources to the right vault folder, and merge captured plans into project notes
allowed-tools: [Bash, WebFetch]
---

Process inbox notes using the Karpathy-inspired wiki pattern.

Works with all note types: Readwise articles, meeting notes, daily notes, decisions, raw captures, and **captured plan notes** (which get special treatment — see below).

## Step 1 — Load context

Run:

```bash
pkm process inbox --full
```

Read the output carefully. It contains:
- Every inbox note's filename, metadata, and **full body**
- The current `04-knowledge/index.md` at the end

Also run to see existing projects:

```bash
pkm project list
```

If the inbox is empty, report that and stop.

## Step 2 — Classify every note

For each inbox note, determine its **type**:

### A) Plan notes (special treatment)

A note is a **plan note** if it meets ANY of these criteria:
- `source: claude-code` AND body contains a heading like `# Plan: ...`
- `source: claude-code` AND `type_hint: knowledge`

**Do not distill plan notes into 04-knowledge/.** Instead, route them to `02-projects/` — see Step 3A.

### B) Regular notes (standard distillation)

Everything else. Determine:

**Action** (pick one):
- `distill` — contains generalizable insights worth adding to the knowledge base
- `file-only` — worth keeping as a reference but nothing extractable
- `skip` — low value, noise, outdated, or duplicate → archive

**Target folder**: `resources`, `knowledge`, `projects`, `areas`, `decisions`, `archive`

**Topic updates** (for `distill` only, 1–3 per note):
- **Slug**: use an existing slug from the index if the insight fits; otherwise propose a new `kebab-case` slug
- **Content to merge**: write **evergreen prose** — extracted insight useful in 6 months. Attribute with `> Source: [[note-slug-without-extension]]`

**Fetch the source article for Readwise/resource notes**: before distilling, use `WebFetch`
on the `> **Source:** <url>` link in the note body to read the full article. Base the
distilled insight on the full article, not just Readwise's auto-summary/highlights —
they're often too thin to distill well. Do NOT save the fetched article text anywhere
(not in the inbox note, not in the knowledge page) — only the distilled evergreen prose
gets written to `04-knowledge/`. If the fetch fails (paywall, dead link, non-article
content), fall back to the existing summary/highlights.

**Topic relationships** (when creating a NEW topic):
- Scan `04-knowledge/index.md` for existing topics that are a **parent** (broader topic this is a special case of), **child** (narrower topic this generalizes), or close **sibling** of the new topic.
- If found, propose a `## Related` entry for the new page, e.g. `- [[parent-slug]] — parent topic: <why>`, and a reciprocal `- [[new-slug]] — <why>` entry to add to the existing page.
- Only propose this for genuinely useful navigational links (shared subject matter a reader would want to jump between) — not every tangential overlap.

**Per-type guidance:**
| Type | Target | What to distill |
|------|--------|-----------------|
| Readwise / resource | `resources` | Key claims, patterns, mental models |
| Meeting note | `projects` or `areas` | Decisions made, open questions, action items |
| Daily note | `01-daily` (`file-only`) | Non-trivial insights only |
| Decision | `decisions` | Rationale to knowledge if reusable |
| Raw capture | depends | Classify by substance |

## Step 3 — Present the full plan before acting

**Do not run any commands yet.** Present the complete plan for all notes, separated into plan notes and regular notes:

```
=== Plan Notes → 02-projects/ ===

[1/N] 2026-04-06-claude-code-pkm-feature.md  (source: claude-code)
  Project: pkm-ai  (EXISTING — last updated 2026-04-05)
  Heading: 2026-04-06 — Project management feature  ← also becomes Timeline entry
  Intent:  (unchanged)
  Status → "pkm project commands implemented. Distill skill updated."
  Next Steps → "- [ ] Write tests\n- [ ] Update README"

=== Regular Notes → knowledge / resources / archive ===

[2/N] 2026-04-03-readwise-article.md  (resource)
  Action:  distill → 05-resources/
  Topics:  [update] software-engineering-philosophy
           ## From [[...]]: Key insight here...

[3/N] 2026-02-01-readwise-spam.md
  Action:  skip → archive
```

If any **new topics** are proposed, also list relationship suggestions:

```
=== Topic Relationships ===

[new] ai-agent-evaluation
  Related → [[ai-agents]] (parent: general agent patterns)
  Reciprocal: add "- [[ai-agent-evaluation]] — evaluating agent behavior and outputs" to ai-agents.md
```

Ask: **"Apply all at once, or step through note by note? (all / step / cancel)"**

## Step 3A — Plan note routing details

For each plan note:

1. **Extract project name** from the `# Plan: <Name>` heading in the body.
   Derive a `kebab-case` slug (e.g. "pkm.ai CLI" → `pkm-ai`).

2. **Check if project exists**: `pkm project list` (already done in Step 1).

3. **If project exists** — read it with:
   ```bash
   cat "$(pkm config --show | grep vault_path | sed 's/^vault_path: *//')/02-projects/<slug>.md"
   ```
   Understand the current Intent, Status, and Next Steps.

4. **Propose**:
   - **Intent**: keep unchanged unless this plan significantly reframes the project goal
   - **Current Status**: synthesise from the plan — what was just built/decided/designed?
   - **Next Steps**: extract from the plan's "Next Steps", "Follow-ups", or "Open Questions" sections; rewrite as a short `- [ ]` checklist
   - **Plan heading**: `YYYY-MM-DD — <brief session description>` — this is also used as the Timeline entry automatically; no separate flag needed

5. **If project is new**:
   - Propose an **Intent** (what problem is this solving? what is the end goal?)
   - Propose initial Status and Next Steps from the plan content

## Step 4 — Apply each approved note

As you process each **regular note**, collect its log entry (note filename, action,
filed-to folder, updated/created topic slugs) instead of logging it immediately.
At the end of Step 4, write all collected entries to `04-knowledge/log.md` in a
**single batched call** so they share one `## <date>` heading:

```bash
printf '%s\n' '[
  {"note": "<filename>", "action": "distill + file", "filed_to": "resources", "updated": ["<slug>"], "created": []},
  {"note": "<filename>", "action": "file-only", "filed_to": "resources"},
  {"note": "<filename>", "action": "skip", "filed_to": "archive"}
]' | pkm knowledge append-log --from-stdin
```

Plan notes (Step 3A) are NOT logged to `04-knowledge/log.md` — they go to
`02-projects/` Plan History instead.

### Plan notes

```bash
# Add a wikilink to the plan note in Plan History (not the full body inline)
# Use $'...' for --next-steps so \n expands to real newlines
printf '%s\n' '[[<note-slug-without-extension>]]' | pkm project update <slug> \
  --title "<Project Title>" \
  --current-status "<synthesised status>" \
  --next-steps $'- [ ] First action\n- [ ] Second action' \
  --plan-heading "<YYYY-MM-DD — session description>"

# Then move the original inbox note to archive
pkm note move <filename> archive
```

### Regular notes — distill

```bash
printf '%s\n' '<content-to-merge>' | pkm knowledge append-topic <slug> --title "<Title>"
pkm knowledge update-index <slug> --description "<desc>"  # new topics only
pkm note move <filename> <target-folder>
```
Record a log entry: `{"note": "<filename>", "action": "distill + file", "filed_to": "<folder>", "updated": [<slugs>], "created": [<slugs>]}`

### New topics — add Related links

For each approved topic relationship (new topic ↔ existing topic), append a `## Related`
section to **both** pages. Before appending to an **existing** page, `cat` it first —
if it already has a `## Related` section, fold the new line into that section manually
(via a full rewrite with `printf | pkm knowledge append-topic` is append-only, so avoid
creating a second `## Related` heading; if one exists, skip the automated append and
note it in the report for manual addition).

```bash
# New page — safe to append, no existing ## Related section yet
printf '%s\n' '## Related

- [[<existing-slug>]] — <relationship description>' | pkm knowledge append-topic <new-slug> --title "<Title>"

# Existing page — only if it has no ## Related section yet
printf '%s\n' '## Related

- [[<new-slug>]] — <relationship description>' | pkm knowledge append-topic <existing-slug> --title "<Title>"
```

### Regular notes — file-only

```bash
pkm note move <filename> <target-folder>
```
Record a log entry: `{"note": "<filename>", "action": "file-only", "filed_to": "<folder>"}`

### Regular notes — skip

```bash
pkm note move <filename> archive
```
Record a log entry: `{"note": "<filename>", "action": "skip", "filed_to": "archive"}`

### After all regular notes are processed

Append all collected entries in one batched call (see Step 4):

```bash
printf '%s\n' '[...]' | pkm knowledge append-log --from-stdin
```

## Step 5 — Report

```
Processed N notes:
  Plan notes merged: X  (into Y projects)
  Distilled:         X
  Filed:             Y
  Archived:          Z
```

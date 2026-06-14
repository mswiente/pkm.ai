---
description: Health-check 04-knowledge/ — find broken links, orphan pages, index drift, stale resources, and concepts that deserve their own page
allowed-tools: [Bash]
---

Run a periodic health check over the knowledge base (Karpathy's "lint" step). This
is read-only analysis followed by an optional, confirmed cleanup — it never edits
files without asking first.

## Step 1 — Gather data

```bash
VAULT="$(pkm config --show | grep vault_path | sed 's/^vault_path: *//')"

echo "=== 04-knowledge/index.md ==="
cat "$VAULT/04-knowledge/index.md"

echo "=== 04-knowledge/*.md (all topic pages, full content) ==="
for f in "$VAULT"/04-knowledge/*.md; do
  case "$(basename "$f")" in index.md|log.md) continue ;; esac
  echo "--- $(basename "$f" .md) ---"
  cat "$f"
done

echo "=== All vault note slugs (for resolving [[wikilinks]]) ==="
find "$VAULT" -name '*.md' | sed "s|$VAULT/||;s|\.md$||;s|.*/||" | sort -u

echo "=== 05-resources/ notes never cited by any knowledge page ==="
for f in "$VAULT"/05-resources/*.md; do
  slug=$(basename "$f" .md)
  if ! grep -rq "\[\[$slug\]\]" "$VAULT/04-knowledge/"*.md 2>/dev/null; then
    echo "$slug"
  fi
done
```

## Step 2 — Run checks

### A) Broken wikilinks

For every `[[slug]]` referenced anywhere in `04-knowledge/*.md` (body text, `> Source:`
attributions, `## Related` sections), check whether `slug` appears in the "all vault
note slugs" list from Step 1. Report any that don't resolve to any file in the vault.

### B) Index drift

- `04-knowledge/*.md` files (excluding `index.md`, `log.md`) not listed in `index.md`
  → "missing from index"
- `index.md` entries with no corresponding `04-knowledge/<slug>.md` file
  → "dangling index entry"

### C) Orphan pages

A topic page is an orphan if it's listed in `index.md` but **no other** `04-knowledge/*.md`
page links to it via `[[slug]]` (in body or `## Related`) — i.e. nothing in the wiki
itself points to it; it's only reachable via the flat index.

### D) Undistilled resources

List `05-resources/*.md` notes not cited by any `04-knowledge/*.md` page (the
`/redistill-resources` orphan check) — these were filed but never distilled into the
wiki. Flag as candidates for `/redistill-resources`.

### E) Concepts mentioned but lacking their own page

Read through all topic page bodies (not just headings). Look for named concepts,
patterns, or tools that:
- appear as their own `##`/`###` subsection in **two or more different** topic pages
  (suggesting the concept has outgrown being a subsection of multiple unrelated pages), or
- get a full paragraph of substantive discussion in a page but are tangential to
  that page's main subject and would stand alone as their own topic

For each, propose a new `kebab-case` slug, a one-line description for `index.md`,
and list which existing pages/sections would seed its initial content (as
`## Related` links and/or content to move).

### F) Contradictions (best-effort)

While reading, note any claims across different topic pages that directly conflict
(not just differing emphasis). Flag both pages and the conflicting statements —
this is for human judgment, not auto-resolution.

## Step 3 — Present findings

```
=== Knowledge Base Lint Report ===

Broken wikilinks (N):
  [[slug]] referenced in <page> — no matching file found

Index drift (N):
  <slug> — page exists but missing from index.md
  <slug> — listed in index.md but file does not exist

Orphan pages (N):
  <slug> — in index.md but not linked from any other topic page
    Suggested fix: add [[<slug>]] to ## Related on <other-page> (because ...)

Undistilled resources (N):
  <slug> — not cited by any knowledge page

Concepts lacking a page (N):
  <new-slug> — "<description>"
    Currently discussed in: <page>#<section>, <page>#<section>

Possible contradictions (N):
  <page A>: "<claim>"
  <page B>: "<claim>"
```

Ask: **"Which of these should I fix now? (all / select / none)"**

Concepts-lacking-a-page and contradictions are **suggestions for discussion** —
present them but don't bundle them into the same yes/no as the mechanical fixes
unless the user explicitly asks to create the new pages.

## Step 4 — Apply approved fixes

- **Index drift (missing from index)**: `pkm knowledge update-index <slug> --description "<desc>"`
- **Index drift (dangling entry)**: remove the line from `index.md` directly (no CLI primitive for removal)
- **Orphan pages**: append a `## Related` entry on the linking page(s), following the
  same append/fold pattern as `/distill-inbox` (skip automated append if `## Related`
  already exists and fold manually instead)
- **Broken wikilinks**: present options per link — fix the slug (typo), remove the
  link, or leave it if the target is intentionally not-yet-created (e.g. a concept
  flagged in part E that the user wants to create)
- **New concept pages**: only create if the user explicitly approves — use
  `pkm knowledge append-topic <new-slug> --title "<Title>"` seeded with content
  drawn from the source pages, then `pkm knowledge update-index`, then add
  `## Related` links back from the source pages

## Step 5 — Report

```
Lint complete:
  Broken links fixed:   X / N
  Index drift fixed:    X / N
  Orphans linked:       X / N
  New pages created:    X (from M concept suggestions)
  Undistilled resources: N flagged for /redistill-resources
  Contradictions:       N flagged for manual review (not auto-resolved)
```

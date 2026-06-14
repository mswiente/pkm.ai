---
description: Re-distill already-filed 05-resources/ notes by fetching their source articles and extracting insights into 04-knowledge/ topic pages
allowed-tools: [Bash, WebFetch]
---

Re-distill resource notes that were filed before the WebFetch-based distillation
existed (or were originally filed as `file-only` with nothing extracted).

## Step 1 — Find candidates

```bash
VAULT="$(pkm config --show | grep vault_path | sed 's/^vault_path: *//')"
ls "$VAULT/05-resources/"
```

A resource is an **orphan candidate** if no `04-knowledge/*.md` page cites it as a
source. Check with:

```bash
for f in "$VAULT"/05-resources/*.md; do
  slug=$(basename "$f" .md)
  if ! grep -rq "\[\[$slug\]\]" "$VAULT/04-knowledge/"*.md 2>/dev/null; then
    echo "$slug"
  fi
done
```

Present the list of orphan candidates to the user. Ask which to process:
**"Re-distill all N orphans, a specific list, or cancel?"**

If the user names specific resource files/slugs directly, skip the orphan scan and
use those instead.

## Step 2 — For each selected resource

1. `cat "$VAULT/05-resources/<slug>.md"` to read frontmatter + existing summary/highlights.
2. Extract the `> **Source:** <url>` line and `WebFetch` that URL to read the full article.
   Do NOT save the fetched text anywhere — it's only used as context for this step.
   If the fetch fails (paywall, dead link, non-article content), fall back to the
   existing summary/highlights, or propose `skip` if there's nothing usable.
3. Read `04-knowledge/index.md` to find existing topics this article relates to.
4. Propose **1–3 topic updates** (existing or new slugs), following the same
   evergreen-prose style as `/distill-inbox`:
   `> Source: [[<slug>]]`
5. Apply the same **topic relationship** check as `/distill-inbox` (Step 2 of that
   skill) for any newly created topics — propose `## Related` cross-links.

## Step 3 — Present the full plan before acting

```
=== Re-distill Plan ===

[1/N] 2026-04-03-readwise-article (resources/)
  Topics:  [update] software-engineering-philosophy
           ## From [[...]]: Key insight here...
  (or)
  Action:  no new insight found — leave as-is
```

Ask: **"Apply all at once, or step through note by note? (all / step / cancel)"**

## Step 4 — Apply each approved resource

```bash
printf '%s\n' '<content-to-merge>' | pkm knowledge append-topic <slug> --title "<Title>"
pkm knowledge update-index <slug> --description "<desc>"  # new topics only
```

Apply any approved `## Related` cross-links the same way as `/distill-inbox` Step 4
(append-only, skip existing-page append if a `## Related` section already exists).

Collect a log entry per processed resource and append them in **one batched call**
at the end:

```bash
printf '%s\n' '[
  {"note": "<slug>.md", "action": "redistill", "updated": ["<slug>"], "created": []}
]' | pkm knowledge append-log --from-stdin
```

Resources with no new insight are not logged.

## Step 5 — Report

```
Re-distilled N of M candidates:
  Updated topics: ...
  Created topics: ...
  No new insight:  X
```

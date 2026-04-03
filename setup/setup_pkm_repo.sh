#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   ./setup_pkm_repo.sh [repo_path]
#
# Example:
#   ./setup_pkm_repo.sh "$HOME/projects/pkm.ai"
#
# If no repo_path is provided, the default is:
#   ~/projects/pkm.ai

REPO_PATH="${1:-$HOME/projects/pkm.ai}"

mkdir -p \
  "$REPO_PATH/cmd/pkm" \
  "$REPO_PATH/internal/capture" \
  "$REPO_PATH/internal/process" \
  "$REPO_PATH/internal/templates" \
  "$REPO_PATH/internal/config" \
  "$REPO_PATH/internal/fs" \
  "$REPO_PATH/templates" \
  "$REPO_PATH/examples" \
  "$REPO_PATH/scripts"

cat > "$REPO_PATH/README.md" <<'EOF'
# pkm.ai

CLI for a markdown-first personal knowledge system built around:
- Obsidian
- Claude Code
- GitHub Copilot CLI

## Planned commands

- `pkm capture`
- `pkm process inbox`
- `pkm daily create`
- `pkm meeting create`
- `pkm decision create`

## Goals

- fast capture into `00-inbox/`
- AI-assisted inbox processing
- markdown as source of truth
- tool-independent workflows
EOF

cat > "$REPO_PATH/.gitignore" <<'EOF'
.DS_Store
bin/
dist/
build/
.coverage
coverage.out
.env
.idea/
.vscode/
EOF

cat > "$REPO_PATH/LICENSE" <<'EOF'
TODO: add license
EOF

cat > "$REPO_PATH/scripts/dev.sh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Add local dev commands here."
EOF

chmod +x "$REPO_PATH/scripts/dev.sh"

cat > "$REPO_PATH/templates/inbox.md" <<'EOF'
---
title: 
type: inbox
status: inbox
source: manual
created: 
updated: 
tags: []
type_hint: 
---

## Context

- Herkunft:
- Anlass:
- Erwarteter Zieltyp:

## Content


## Open Questions

- 

## Follow-ups

- 
EOF

cat > "$REPO_PATH/templates/daily.md" <<'EOF'
---
title: 
type: daily
status: active
created: 
updated: 
tags: [daily]
---

# Daily Note

## Focus Today

- 

## Schedule / Meetings

- 

## Notes

- 

## Decisions / Learnings

- 

## Open Loops

- 

## Tasks

- [ ] 
EOF

cat > "$REPO_PATH/templates/meeting.md" <<'EOF'
---
title: 
type: meeting
status: draft
created: 
updated: 
tags: [meeting]
date: 
participants: []
project: 
source: manual
---

# Meeting Note

## Context

- Date:
- Participants:
- Project / Area:
- Purpose:

## Agenda / Topics

- 

## Notes

- 

## Decisions

- 

## Risks / Issues

- 

## Next Steps

- [ ] 

## Related Notes

- 
EOF

cat > "$REPO_PATH/templates/decision.md" <<'EOF'
---
title: 
type: decision
status: draft
created: 
updated: 
tags: [decision]
decision_date: 
project: 
related_notes: []
---

# Decision

## Status

Draft / Accepted / Superseded

## Context


## Decision


## Options Considered

- Option A:
- Option B:

## Rationale


## Consequences

### Positive

- 

### Negative / Trade-offs

- 

## Follow-ups

- [ ] 

## References

- 
EOF

cat > "$REPO_PATH/templates/knowledge.md" <<'EOF'
---
title: 
type: knowledge
status: evergreen
created: 
updated: 
tags: []
related_notes: []
---

# Summary


## Key Points

- 

## Explanation


## Examples / Applications

- 

## Related Concepts

- 

## Open Questions

- 

## Sources

- 
EOF

cat > "$REPO_PATH/templates/troubleshooting.md" <<'EOF'
---
title: 
type: troubleshooting
status: draft
created: 
updated: 
tags: [troubleshooting]
systems: []
project: 
source: claude-code
related_notes: []
---

# Troubleshooting

## Problem


## Context

- Environment:
- System / Component:
- Trigger:

## Symptoms

- 

## Root Cause


## Resolution

1. 
2. 
3. 

## Validation

- 

## Learnings

- 

## Follow-ups

- [ ] 

## References

- 
EOF

cat > "$REPO_PATH/templates/project.md" <<'EOF'
---
title: 
type: project
status: active
created: 
updated: 
tags: [project]
area: 
owner: 
related_notes: []
---

# Project

## Goal


## Context


## Current Status


## Key Notes

- 

## Decisions

- 

## Risks / Issues

- 

## Next Steps

- [ ] 

## Links

- 
EOF

cat > "$REPO_PATH/templates/resource.md" <<'EOF'
---
title: 
type: resource
status: draft
created: 
updated: 
tags: [resource]
url: 
author: 
published: 
related_notes: []
---

# Resource

## Source

- URL:
- Author:
- Published:

## Summary


## Why It Matters


## Key Extracts

- 

## Related Notes

- 
EOF

cat <<EOF
pkm.ai repo scaffold created successfully.

Location:
  $REPO_PATH

Created directories:
  cmd/pkm
  internal/capture
  internal/process
  internal/templates
  internal/config
  internal/fs
  templates
  examples
  scripts

Created files:
  README.md
  .gitignore
  LICENSE
  scripts/dev.sh

Template files:
  templates/inbox.md
  templates/daily.md
  templates/meeting.md
  templates/decision.md
  templates/knowledge.md
  templates/troubleshooting.md
  templates/project.md
  templates/resource.md

Next steps:
  1. cd "$REPO_PATH"
  2. git init
  3. choose implementation language for pkm CLI
EOF

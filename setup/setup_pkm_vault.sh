#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   ./setup_pkm_vault.sh [vault_path]
#
# Example:
#   ./setup_pkm_vault.sh "$HOME/Documents/Obsidian/pkm-vault"
#
# If no vault_path is provided, the default is:
#   ~/Documents/Obsidian/pkm-vault

VAULT_PATH="${1:-$HOME/Documents/Obsidian/pkm-vault}"
TEMPLATES_DIR="$VAULT_PATH/07-templates"
GITIGNORE_FILE="$VAULT_PATH/.gitignore"

mkdir -p \
  "$VAULT_PATH/00-inbox" \
  "$VAULT_PATH/01-daily" \
  "$VAULT_PATH/02-projects" \
  "$VAULT_PATH/03-areas" \
  "$VAULT_PATH/04-knowledge" \
  "$VAULT_PATH/05-resources" \
  "$VAULT_PATH/06-decisions" \
  "$VAULT_PATH/07-templates" \
  "$VAULT_PATH/08-attachments" \
  "$VAULT_PATH/09-archive"

cat > "$TEMPLATES_DIR/inbox.md" <<'EOF'
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

cat > "$TEMPLATES_DIR/daily.md" <<'EOF'
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

cat > "$TEMPLATES_DIR/meeting.md" <<'EOF'
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

cat > "$TEMPLATES_DIR/decision.md" <<'EOF'
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

cat > "$TEMPLATES_DIR/knowledge.md" <<'EOF'
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

cat > "$TEMPLATES_DIR/troubleshooting.md" <<'EOF'
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

cat > "$TEMPLATES_DIR/project.md" <<'EOF'
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

cat > "$TEMPLATES_DIR/resource.md" <<'EOF'
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

if [[ ! -f "$GITIGNORE_FILE" ]]; then
  cat > "$GITIGNORE_FILE" <<'EOF'
.DS_Store
.obsidian/workspace.json
EOF
fi

cat <<EOF
PKM vault created successfully.

Location:
  $VAULT_PATH

Created folders:
  00-inbox
  01-daily
  02-projects
  03-areas
  04-knowledge
  05-resources
  06-decisions
  07-templates
  08-attachments
  09-archive

Created templates:
  inbox.md
  daily.md
  meeting.md
  decision.md
  knowledge.md
  troubleshooting.md
  project.md
  resource.md

Next steps:
  1. Open the folder in Obsidian.
  2. Optionally run: cd "$VAULT_PATH" && git init
EOF


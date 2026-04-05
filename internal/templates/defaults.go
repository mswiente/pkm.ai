package templates

// Default Go template strings for each note type.
// These use BODY_PLACEHOLDER as a sentinel that is replaced after template
// execution, so user-provided content is never interpreted as template syntax.

const defaultInboxTemplate = `---
title: {{.Title}}
type: inbox
status: inbox
source: {{.Source}}
created: {{.Date}}
updated: {{.Date}}
tags: {{.TagsYAML}}
type_hint: {{.TypeHint}}
---

## Context

- Herkunft:
- Anlass:
- Erwarteter Zieltyp:

## Content

BODY_PLACEHOLDER

## Open Questions

-

## Follow-ups

-
`

const defaultDailyTemplate = `---
title: {{.Title}}
type: daily
status: active
created: {{.Date}}
updated: {{.Date}}
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
`

const defaultMeetingTemplate = `---
title: {{.Title}}
type: meeting
status: draft
source: {{.Source}}
created: {{.Date}}
updated: {{.Date}}
tags: [meeting]
date: {{.Date}}
participants: {{.ParticipantsYAML}}
project: {{.Project}}
---

# Meeting Note

## Context

- Date: {{.Date}}
- Participants: {{.ParticipantsYAML}}
- Project / Area: {{.Project}}
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
`

const defaultDecisionTemplate = `---
title: {{.Title}}
type: decision
status: {{.Status}}
created: {{.Date}}
updated: {{.Date}}
tags: [decision]
decision_date: {{.Date}}
project: {{.Project}}
related_notes: []
---

# Decision

## Status

{{.StatusLabel}}

## Context


## Decision

BODY_PLACEHOLDER

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
`

const defaultKnowledgeTemplate = `---
title: {{.Title}}
type: knowledge
status: evergreen
created: {{.Date}}
updated: {{.Date}}
tags: {{.TagsYAML}}
related_notes: []
---

# Summary

BODY_PLACEHOLDER

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
`

var defaultTemplates = map[string]string{
	"inbox":    defaultInboxTemplate,
	"daily":    defaultDailyTemplate,
	"meeting":  defaultMeetingTemplate,
	"decision": defaultDecisionTemplate,
	"knowledge": defaultKnowledgeTemplate,
}

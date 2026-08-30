package project

import (
	"strings"
	"testing"
)

// ---- effectiveTimelineEntry ------------------------------------------------

func TestEffectiveTimelineEntry(t *testing.T) {
	tests := []struct {
		name        string
		explicit    string
		planHeading string
		want        string
	}{
		{"explicit overrides heading", "my summary", "2026-04-06 — session", "my summary"},
		{"falls back to planHeading", "", "2026-04-06 — session", "2026-04-06 — session"},
		{"both empty returns empty", "", "", ""},
		{"whitespace explicit is ignored", "   ", "2026-04-06 — session", "2026-04-06 — session"},
		{"trims whitespace from explicit", "  summary  ", "", "summary"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := effectiveTimelineEntry(tc.explicit, tc.planHeading)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// ---- appendToTimeline -------------------------------------------------------

func TestAppendToTimeline_NoSections(t *testing.T) {
	// Neither Timeline nor Plan History present — append at end.
	content := "---\ntitle: Foo\n---\n\n## Next Steps\n\n- [ ] do something\n"
	got := appendToTimeline(content, "2026-04-06 — first session")
	if !strings.Contains(got, "\n## Timeline\n") {
		t.Error("expected ## Timeline section to be created")
	}
	if !strings.Contains(got, "- 2026-04-06 — first session\n") {
		t.Error("expected bullet entry in Timeline")
	}
}

func TestAppendToTimeline_InsertedBeforePlanHistory(t *testing.T) {
	content := "---\ntitle: Foo\n---\n\n## Next Steps\n\n- [ ] do something\n\n## Plan History\n\n### 2026-04-05\n\n[[old-plan]]\n\n"
	got := appendToTimeline(content, "2026-04-06 — new session")

	timelineIdx := strings.Index(got, "\n## Timeline\n")
	planHistoryIdx := strings.Index(got, "\n## Plan History\n")

	if timelineIdx < 0 {
		t.Fatal("## Timeline section not found")
	}
	if planHistoryIdx < 0 {
		t.Fatal("## Plan History section not found")
	}
	if timelineIdx > planHistoryIdx {
		t.Error("## Timeline should appear before ## Plan History")
	}
	if !strings.Contains(got, "- 2026-04-06 — new session\n") {
		t.Error("expected bullet entry in Timeline")
	}
}

func TestAppendToTimeline_AppendsToExistingSection(t *testing.T) {
	content := "---\ntitle: Foo\n---\n\n## Timeline\n\n- 2026-04-05 — first session\n\n## Plan History\n\n### 2026-04-05\n\n[[old-plan]]\n\n"
	got := appendToTimeline(content, "2026-04-06 — second session")

	firstIdx := strings.Index(got, "- 2026-04-05 — first session")
	secondIdx := strings.Index(got, "- 2026-04-06 — second session")
	planIdx := strings.Index(got, "\n## Plan History\n")

	if firstIdx < 0 || secondIdx < 0 {
		t.Fatal("one or both bullet entries missing")
	}
	if firstIdx > secondIdx {
		t.Error("first entry should appear before second entry")
	}
	if secondIdx > planIdx {
		t.Error("second Timeline entry should appear before ## Plan History")
	}
}

func TestAppendToTimeline_TwoSequentialCalls(t *testing.T) {
	content := "---\ntitle: Foo\n---\n\n## Next Steps\n\n(none)\n\n## Plan History\n\n"
	content = appendToTimeline(content, "2026-04-05 — first")
	content = appendToTimeline(content, "2026-04-06 — second")

	firstIdx := strings.Index(content, "- 2026-04-05 — first")
	secondIdx := strings.Index(content, "- 2026-04-06 — second")
	planIdx := strings.Index(content, "\n## Plan History\n")

	if firstIdx < 0 || secondIdx < 0 {
		t.Fatal("one or both bullet entries missing")
	}
	if firstIdx > secondIdx {
		t.Error("first entry should appear before second")
	}
	if secondIdx > planIdx {
		t.Error("both Timeline entries should appear before ## Plan History")
	}
}

// ---- buildNote -------------------------------------------------------------

func TestBuildNote_NoPlanHeading(t *testing.T) {
	opts := UpdateOptions{
		Slug:   "my-proj",
		Title:  "My Project",
		Intent: "Do something useful.",
	}
	got := buildNote(opts, "2026-04-06")

	if strings.Contains(got, "## Timeline") {
		t.Error("## Timeline should not be created when no PlanHeading or TimelineEntry")
	}
	if strings.Contains(got, "## Plan History") {
		t.Error("## Plan History should not be created when no PlanContent")
	}
}

func TestBuildNote_WithPlanHeadingAndContent(t *testing.T) {
	opts := UpdateOptions{
		Slug:        "my-proj",
		Title:       "My Project",
		PlanHeading: "2026-04-06 — Initial setup",
		PlanContent: "Set up the repo.",
	}
	got := buildNote(opts, "2026-04-06")

	if !strings.Contains(got, "## Timeline") {
		t.Error("expected ## Timeline section")
	}
	if !strings.Contains(got, "- 2026-04-06 — Initial setup\n") {
		t.Error("expected Timeline bullet entry derived from PlanHeading")
	}
	if !strings.Contains(got, "## Plan History") {
		t.Error("expected ## Plan History section")
	}

	// Timeline must appear before Plan History
	tIdx := strings.Index(got, "## Timeline")
	pIdx := strings.Index(got, "## Plan History")
	if tIdx > pIdx {
		t.Error("## Timeline should appear before ## Plan History")
	}
}

func TestBuildNote_ExplicitTimelineEntryOverridesHeading(t *testing.T) {
	opts := UpdateOptions{
		Slug:          "my-proj",
		PlanHeading:   "2026-04-06 — Raw heading",
		TimelineEntry: "Custom summary for the timeline",
	}
	got := buildNote(opts, "2026-04-06")

	if !strings.Contains(got, "- Custom summary for the timeline\n") {
		t.Error("expected explicit TimelineEntry to appear in Timeline")
	}
	if strings.Contains(got, "- 2026-04-06 — Raw heading\n") {
		t.Error("PlanHeading should not appear as Timeline entry when TimelineEntry is explicit")
	}
}

// ---- patchNote -------------------------------------------------------------

func TestPatchNote_Timeline_CreatedOnFirstEntry(t *testing.T) {
	existing := "---\ntitle: P\ntype: project\nstatus: active\nsource: claude-code\ncreated: 2026-04-01\nupdated: 2026-04-01\ntags: [project]\n---\n\n## Intent\n\nDo things.\n\n## Current Status\n\n(to be described)\n\n## Next Steps\n\n(to be defined)\n\n"
	opts := UpdateOptions{
		PlanHeading: "2026-04-06 — First session",
	}
	got, patched := patchNote(existing, opts, "2026-04-06")

	if !strings.Contains(got, "## Timeline") {
		t.Error("## Timeline section should be created")
	}
	if !strings.Contains(got, "- 2026-04-06 — First session\n") {
		t.Error("expected Timeline bullet")
	}
	if !containsString(patched, "Timeline") {
		t.Errorf("patched list should include 'Timeline', got %v", patched)
	}
}

func TestPatchNote_Timeline_AppendedOnSecondEntry(t *testing.T) {
	existing := "---\ntitle: P\ntype: project\nstatus: active\nsource: claude-code\ncreated: 2026-04-01\nupdated: 2026-04-05\ntags: [project]\n---\n\n## Next Steps\n\n(none)\n\n## Timeline\n\n- 2026-04-05 — First session\n\n## Plan History\n\n### 2026-04-05 — First session\n\n[[note]]\n\n"
	opts := UpdateOptions{
		PlanHeading: "2026-04-06 — Second session",
	}
	got, _ := patchNote(existing, opts, "2026-04-06")

	if !strings.Contains(got, "- 2026-04-05 — First session\n") {
		t.Error("first Timeline entry should be preserved")
	}
	if !strings.Contains(got, "- 2026-04-06 — Second session\n") {
		t.Error("second Timeline entry should be appended")
	}
}

func TestPatchNote_TimelineBeforePlanHistory_InPatchedSlice(t *testing.T) {
	existing := "---\ntitle: P\ntype: project\nstatus: active\nsource: claude-code\ncreated: 2026-04-01\nupdated: 2026-04-01\ntags: [project]\n---\n\n## Next Steps\n\n(none)\n\n"
	opts := UpdateOptions{
		PlanHeading: "2026-04-06 — Session",
		PlanContent: "Did something.",
	}
	_, patched := patchNote(existing, opts, "2026-04-06")

	timelinePos := indexOfString(patched, "Timeline")
	planHistoryPos := indexOfString(patched, "Plan History")

	if timelinePos < 0 {
		t.Fatal("Timeline missing from patched list")
	}
	if planHistoryPos < 0 {
		t.Fatal("Plan History missing from patched list")
	}
	if timelinePos > planHistoryPos {
		t.Error("Timeline should appear before Plan History in patched list")
	}
}

// ---- helpers ----------------------------------------------------------------

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func indexOfString(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1
}

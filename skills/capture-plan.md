---
description: Capture the current plan as a PKM inbox note
allowed-tools: [Bash]
---

Capture the current plan as a PKM inbox note.

Steps:
1. Find the most recently modified `.md` file in `~/.claude/plans/` by running:
   `ls -t ~/.claude/plans/*.md 2>/dev/null | head -1`
   If no files are found, tell the user there is no current plan to capture.

2. Extract the plan title: read the file and find the first line starting with `# `.
   Strip the leading `# ` to get the title text.
   If no H1 heading exists, use the filename without extension as the title.

3. Check if this session already captured a plan note by running:
   `cat /tmp/pkm-last-captured-plan.txt 2>/dev/null`
   If the file exists and the path it contains points to an existing file, this is an update — go to step 4a.
   Otherwise, this is a new capture — go to step 4b.

4a. **Update** (plan was already captured this session):
   Run: `cat "<plan-file-path>" | pkm capture --update "<existing-note-path>"`
   Report to the user that the existing note was updated.
   Skip to step 5.

4b. **New capture**:
   Run: `cat "<plan-file-path>" | pkm capture --title "<extracted-title>" --source claude-code --type-hint knowledge`
   The command will print the created note path. Save that path for this session:
   `echo "<created-note-path>" > /tmp/pkm-last-captured-plan.txt`
   Report the created note path to the user.

5. Done.

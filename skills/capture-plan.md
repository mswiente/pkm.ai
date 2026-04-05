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

3. Capture the note by running:
   `cat "<plan-file-path>" | pkm capture --title "<extracted-title>" --source claude-code --type-hint knowledge`

4. Report the created note path to the user.

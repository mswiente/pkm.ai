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

5. **Auto-create project if inside a git repo** (new captures only — skip on updates):

   a. Check if we're inside a git repo: `git rev-parse --show-toplevel 2>/dev/null`
      If that fails, skip to step 6.

   b. Derive a slug from the repo:
      - Get the remote URL: `git remote get-url origin 2>/dev/null`
      - Strip `.git` suffix, then take the last path segment (e.g. `pkm.ai` → `pkm-ai`, `my_project` → `my-project`).
      - Convert underscores and dots to hyphens and lowercase everything.
      - If no remote, fall back to the basename of the repo root directory.

   c. Check if the project already exists: `pkm project list`
      Parse the output. If a line contains the derived slug, the project exists — skip to step 6.

   d. Gather repo info to seed the project:
      - **Title**: use the repo name (last segment of remote URL or directory name, human-readable form).
      - **Intent**: try these sources in order, stopping at the first non-empty result:
        1. First non-empty paragraph of `README.md` (run: `head -50 README.md 2>/dev/null`)
        2. `description` field from `package.json` (`jq -r '.description // empty' package.json 2>/dev/null`)
        3. `description` field from `pyproject.toml` or `Cargo.toml` if present.
        4. Fall back to: "Work on the <repo-name> project."

   e. Create the project:
      ```
      pkm project update <slug> \
        --title "<Title>" \
        --intent "<Intent>" \
        --status active
      ```
      Report the created project to the user.

6. Done.

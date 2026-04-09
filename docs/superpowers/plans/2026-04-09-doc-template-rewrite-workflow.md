# Document Template Rewrite Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Execute the coordination workflow to rewrite all 5 documentation files using the template system with verified runnable examples.

**Architecture:** Coordinator launches 5 parallel subagents (one per document), aggregates error reports by category, handles template bugs (C-errors) immediately, presents implementation mismatches (B-errors) to user for decision, iterates until all documents complete.

**Tech Stack:** Go text/template, bash scripts, testkit.Launch, YAML config for normalization, make docs-generate

---

## File Structure

**Coordinator outputs (this plan):**
- Task tracking via TaskCreate/TaskUpdate
- Aggregated error reports (A, B, C lists)
- Final commit of all generated docs

**Subagent outputs per document:**
- `docs/generate/<doc>.md.tmpl` - Template file with {{output}} placeholders
- `docs/generate/<doc>_examples/*.sh` - Runnable scripts (minimal, executable)
- `docs/generate/<doc>_examples/config.yaml` - Normalization rules
- `docs/generate/<doc>_examples/*_out.txt` - Pre-staged outputs
- Error report: categorized findings

**Documents to rewrite:**
- `autoui.md` (938 lines) - EbitenUI widget automation
- `commands.md` (452 lines) - CLI commands reference
- `tutorial.md` (517 lines) - Getting started tutorial
- `testkit.md` (596 lines) - Go testing framework
- `integration.md` (391 lines) - Setup and integration guide

---

## Task 1: Preparation and Task List Setup

**Files:**
- Read: `docs/generate/original_docs/*.md` (all 5 originals)
- Verify: `docs/generate/autoui_examples/` (existing reference)

- [ ] **Step 1: Create task tracking for 5 documents**

Use TaskCreate to create 5 tasks (one per document):

| ID | Subject | Description |
|----|---------|-------------|
| 1 | Rewrite autoui.md template | Convert autoui.md to template system with scripts, outputs |
| 2 | Rewrite commands.md template | Convert commands.md to template system |
| 3 | Rewrite tutorial.md template | Convert tutorial.md to template system |
| 4 | Rewrite testkit.md template | Convert testkit.md to template system |
| 5 | Rewrite integration.md template | Convert integration.md to template system |

- [ ] **Step 2: Verify backup originals exist**

Run: `ls -la docs/generate/original_docs/`
Expected: 5 .md files present (autoui.md, commands.md, tutorial.md, testkit.md, integration.md)

- [ ] **Step 3: Verify template system is working**

Run: `make docs-generate`
Expected: Successfully processes existing autoui.md.tmpl, generates docs/autoui.md

---

## Task 2: Launch Subagents in Parallel

**Files:**
- Dispatch: 5 Agent tool calls in parallel

- [ ] **Step 1: Launch all 5 subagents simultaneously**

Use Agent tool to dispatch 5 subagents in a single message. Each subagent receives identical prompt structure with document name varying.

**Subagent assignments:**

| Agent Name | Document | Model |
|------------|----------|-------|
| autoui-writer | autoui.md | sonnet |
| commands-writer | commands.md | sonnet |
| tutorial-writer | tutorial.md | sonnet |
| testkit-writer | testkit.md | sonnet |
| integration-writer | integration.md | sonnet |

**Prompt template for each subagent:**

```
You are rewriting the documentation file [DOC_NAME] to use the template system.

## Your Assignment

**Document:** docs/generate/original_docs/[DOC_NAME]
**Lines:** [DOC_LINES]

## Template System Overview

The template system uses Go text/template with a custom `{{output}}` function:

- `{{output "path" "lang"}}` - Reads file from examples directory, wraps in code fence
- Example: `{{output "autoui_examples/autoui_tree_out.txt" "xml"}}`
- Output: ```xml\n<content>\n```

**config.yaml format:**
```yaml
game_dir: examples/autoui  # Relative path from project root
normalize:
  - pattern: '_addr="[^"]*"'  # Regex pattern
    replace: '_addr="<ADDR>"'  # Replacement
```

**Script requirements:**
- Minimal: only essential CLI commands (no setup/teardown - handled by docgen)
- Executable: created with chmod +x or via Write tool
- One script per example/output
- Named: `<doc>_<example>.sh`

**docgen workflow:**
1. `./docgen --generate docs/generate/<doc>_examples` - Runs scripts, captures outputs
2. `./docgen --process` - Processes all templates, generates final .md files
3. Makefile: `make docs-generate` runs both steps

## Your Deliverables

1. **Template file:** `docs/generate/[DOC_NAME].tmpl`
   - Convert original document sections
   - Replace hardcoded examples with `{{output}}` placeholders
   - Preserve all narrative text, headings, explanations

2. **Scripts directory:** `docs/generate/[DOC_EXAMPLES]/*.sh`
   - Create minimal scripts for each example
   - Each script runs one CLI command sequence
   - No setup code - docgen handles launch/wait/shutdown

3. **Config file:** `docs/generate/[DOC_EXAMPLES]/config.yaml`
   - Set game_dir if examples need running game
   - Add normalization rules for dynamic values (addresses, timestamps, PIDs)

4. **Pre-staged outputs:** `docs/generate/[DOC_EXAMPLES]/*_out.txt`
   - Run scripts manually or via docgen to populate
   - Verify outputs match expected format

## Process

1. Read original document from `docs/generate/original_docs/[DOC_NAME]`
2. Identify all code examples that should be runnable
3. For each example:
   - Create minimal shell script
   - Determine if normalization needed
   - Add `{{output}}` placeholder in template
4. Create config.yaml with game_dir and normalization rules
5. Run `./docgen --generate docs/generate/[DOC_EXAMPLES]` to capture outputs
6. Create template file preserving all original content
7. Test: `./docgen --process` and verify docs/[DOC_NAME] generated

## Error Handling Protocol

Track errors in three categories:

**Category A: Minor Format/Naming Issues**
- Cosmetic differences (whitespace, attribute order, slight naming variations)
- Handling: Accept and proceed. Note in error report: "Accepted: <issue>"
- Example: Extra newline, `id` before `x` vs `x` before `id`

**Category B: Major Autoebiten Errors**
- Documented behavior doesn't match actual implementation
- Handling: Document with full context, continue working on rest
- Report format:
  ```
  B-ERROR: <brief description>
  Expected: <what doc says>
  Actual: <what actually happened>
  Command: <CLI command that revealed it>
  Section: <doc section location>
  ```
- Example: Command documented but doesn't exist, flag behavior differs

**Category C: Template System Bugs**
- The docgen system itself has issues
- Handling: Signal immediately, describe the bug, STOP work on this document
- Report format:
  ```
  C-ERROR: <bug description>
  Template: <what you tried>
  Error: <actual error/behavior>
  ```
- Example: `{{output}}` fails, normalization regex doesn't work

## Self-Check Before Finishing

1. All scripts executable: `ls -la docs/generate/[DOC_EXAMPLES]/*.sh` shows -rwxr-xr-x
2. All outputs captured: `ls docs/generate/[DOC_EXAMPLES]/*_out.txt` shows files
3. Template syntax correct: `./docgen --process` succeeds
4. Final .md generated: `docs/[DOC_NAME]` exists and contains output placeholders filled

## Report Your Status

At the end, report one of:
- **DONE**: All deliverables complete, no B or C errors found
- **DONE_WITH_CONCERNS**: Complete but has B-errors to report (list them)
- **BLOCKED**: Hit C-error, stopped work (describe the C-error)

Include your error report:
```
## Error Report for [DOC_NAME]

### Category A (Accepted):
- Accepted: <issue 1>
- Accepted: <issue 2>

### Category B (Major Errors):
- B-ERROR: <description> (if any)

### Category C (Template Bugs):
- C-ERROR: <description> (if any)
```
```

---

## Task 3: Aggregate Results

**Files:**
- Collect: Error reports from all 5 subagents

- [ ] **Step 1: Wait for all 5 subagents to return**

Each subagent returns with status (DONE, DONE_WITH_CONCERNS, BLOCKED) and error report.

- [ ] **Step 2: Collect and categorize all error reports**

Create consolidated lists:

**List A (Minor - informational):**
- Aggregate all "Accepted" entries from all subagents
- No action required

**List B (Major - must discuss with user):**
- Aggregate all B-ERROR entries
- Format: `<doc>: <description> | Expected: ... | Actual: ...`

**List C (Template bugs - fix immediately):**
- Aggregate all C-ERROR entries
- Format: `<doc>: <description>`

- [ ] **Step 3: Assess completion status**

Count documents by status:
- DONE count
- DONE_WITH_CONCERNS count
- BLOCKED count

---

## Task 4: Handle C-Errors (Template System Bugs)

**Condition:** Only execute if List C has entries

**Files:**
- Modify: Template system files (internal/docgen/*.go, cmd/docgen/main.go)

- [ ] **Step 1: Review each C-error**

For each C-error:
- Understand the bug description
- Determine if it's a genuine template system bug or usage error
- If usage error: provide guidance to subagent, mark as not a bug

- [ ] **Step 2: Fix genuine template bugs**

If confirmed bug:
- Implement fix in template system code
- Run tests: `go test ./internal/docgen/... ./cmd/docgen/...`
- Verify fix works

- [ ] **Step 3: Re-launch blocked subagents**

Use Agent tool to re-launch only subagents that were BLOCKED by C-errors.
Use same prompt structure but add note about fixed bug.

---

## Task 5: Handle B-Errors (User Decision Required)

**Condition:** Only execute if List B has entries AND no remaining C-errors

**Files:**
- None (decision phase)

- [ ] **Step 1: Present B-errors to user**

Display consolidated B-error list with full context:

```
## Major Errors Found (Requires Your Decision)

| Doc | Issue | Expected | Actual |
|-----|-------|----------|--------|
| autoui.md | Command 'xyz' doesn't exist | ... | ... |
| commands.md | Flag '--timeout' behaves differently | ... | ... |

For each error, you can choose:
1. Fix in autoebiten (update implementation to match documentation)
2. Fix in documentation (update doc to match actual behavior)
3. Defer (leave as-is, note discrepancy)

Please provide your decisions.
```

- [ ] **Step 2: Implement user decisions**

For each B-error, based on user's choice:
- **Fix in autoebiten:** Create new task to fix implementation
- **Fix in documentation:** Subagent will adjust template in next iteration
- **Defer:** Note in documentation, no code changes

- [ ] **Step 3: Re-launch affected subagents**

Use Agent tool to re-launch subagents for documents with B-errors.
Add context about user's decision for relevant issues.

---

## Task 6: Iterate Loop

**Condition:** Execute if any documents not DONE

- [ ] **Step 1: Check completion status**

If all 5 documents are DONE (no B-errors, no C-errors): proceed to Task 7

Otherwise: return to Task 2 with modified prompt for affected documents

- [ ] **Step 2: Track iteration count**

Log iteration number. No limit per spec - continue until all complete.

---

## Task 7: Finalization

**Condition:** All 5 documents DONE

**Files:**
- Generate: `docs/*.md` (all 5 final documents)
- Commit: All template files, scripts, configs, outputs

- [ ] **Step 1: Run full generation**

Run: `make docs-generate`
Expected: All 5 templates processed, docs/*.md generated

- [ ] **Step 2: Verify all outputs**

Run: `make docs-verify`
Expected: All outputs match normalized expected

- [ ] **Step 3: Commit all changes**

```bash
git add docs/generate/*.tmpl docs/generate/*_examples/ docs/*.md
git commit -m "docs: rewrite all documentation with template system"
```

- [ ] **Step 4: Mark all tasks complete**

Use TaskUpdate to mark all 5 document tasks as completed.

---

## Spec Coverage Check

| Spec Requirement | Task Coverage |
|------------------|---------------|
| 5 documents rewritten | Task 2 (launch), Task 7 (finalize) |
| Template + scripts + outputs per doc | Task 2 subagent prompts |
| Category A errors handled | Task 2 subagent protocol, Task 3 aggregation |
| Category B errors to user | Task 5 (user decision gate) |
| Category C errors fixed immediately | Task 4 (template fix) |
| No iteration limit | Task 6 (loop without limit) |
| Signal-and-wait for C-errors | Task 2 subagent protocol |
| User permission for B-errors | Task 5 Step 1 (present + wait) |

---

## Self-Review

1. **Placeholder scan:** No TBD, TODO, or vague instructions. All subagent prompts contain complete template system documentation.

2. **Internal consistency:** Task 4 C-error handling loops back correctly. Task 5 B-error handling loops back correctly. Task 6 iteration logic connects all phases.

3. **Type consistency:** All file paths use consistent naming (`docs/generate/<doc>_examples/`). Agent names consistent across tasks.

No gaps found.
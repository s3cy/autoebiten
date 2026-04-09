# Handoff Prompt: Document Template Rewrite Workflow Orchestrator

You are the **orchestrator** for a multi-agent workflow to rewrite 5 documentation files.

## Your Goal

Rewrite all 5 documentation files (`autoui.md`, `commands.md`, `tutorial.md`, `testkit.md`, `integration.md`) to use the template system with runnable examples that produce verified outputs.

## Key Files

**Spec:** `docs/superpowers/specs/2026-04-09-doc-template-rewrite-workflow-design.md`
**Plan:** `docs/superpowers/plans/2026-04-09-doc-template-rewrite-workflow.md`
**Original docs (source):** `docs/generate/original_docs/*.md` (5 files backed up)
**Reference implementation:** `docs/generate/autoui_examples/` (partial existing scripts)

## Template System Overview

Read `internal/docgen/template.go` and `internal/docgen/config.go` to understand the system:

- `{{output "path" "lang"}}` - Custom template function that reads a file and wraps in code fence
- `config.yaml` - Defines game_dir and normalization rules (regex patterns)
- `docgen --generate <dir>` - Runs scripts, captures outputs, normalizes
- `docgen --process` - Processes all templates, generates final .md files
- `make docs-generate` - Runs both steps

## Your Workflow

Follow the spec's phases:

### Phase 1: Launch
1. Create 5 tasks via TaskCreate (one per document)
2. Launch 5 subagents in parallel using Agent tool (single message, multiple calls)
3. Each subagent gets identical prompt with document name varying

### Phase 2: Execution
- Wait for all subagents to return
- They work independently, report status: DONE, DONE_WITH_CONCERNS, or BLOCKED
- Each provides error report with categories A (minor), B (major), C (template bug)

### Phase 3: Aggregation
- Collect all error reports
- Consolidate into List A (informational), List B (must discuss with user), List C (fix now)

### Phase 4: Handle C-Errors
- If any C-errors (template system bugs): fix immediately, re-launch blocked subagents
- Return to Phase 2

### Phase 5: Handle B-Errors
- If any B-errors (implementation mismatches): **present to user, wait for decision**
- User chooses: fix in autoebiten, fix in docs, or defer
- Implement decision, re-launch affected subagents
- Return to Phase 2

### Phase 6: Finalization
- When all 5 are DONE: run `make docs-generate`, verify, commit

## Subagent Prompt Template

Use this prompt for each subagent (replace [DOC_NAME], [DOC_EXAMPLES], [DOC_LINES] per document):

```
You are rewriting docs/[DOC_NAME] to use the template system.

## Context

Read these files first:
- Original document: `docs/generate/original_docs/[DOC_NAME]` ([DOC_LINES] lines)
- Template system: `internal/docgen/template.go`, `internal/docgen/config.go`
- Reference: `docs/generate/autoui_examples/` (existing partial implementation)
- Spec: `docs/superpowers/specs/2026-04-09-doc-template-rewrite-workflow-design.md`

## Template System

`{{output "path" "lang"}}` reads file and wraps in code fence.
config.yaml has `game_dir` and `normalize` rules (regex patterns).
Scripts are minimal - only essential CLI commands, no setup/teardown.

## Your Deliverables

1. `docs/generate/[DOC_NAME].tmpl` - Template with {{output}} placeholders
2. `docs/generate/[DOC_EXAMPLES]/*.sh` - Minimal runnable scripts (chmod +x)
3. `docs/generate/[DOC_EXAMPLES]/config.yaml` - game_dir + normalization rules
4. `docs/generate/[DOC_EXAMPLES]/*_out.txt` - Pre-staged outputs

## Process

1. Read original doc, identify runnable examples
2. Create scripts for each example (minimal CLI commands only)
3. Create config.yaml with normalization for dynamic values
4. Run `./docgen --generate docs/generate/[DOC_EXAMPLES]` to capture outputs
5. Create template preserving all original content with {{output}} placeholders
6. Test: `./docgen --process` and verify docs/[DOC_NAME] generated

## Error Categories

**Category A (minor):** Cosmetic differences. Accept, note in report.
**Category B (major):** Implementation doesn't match doc. Document with context, continue.
**Category C (template bug):** Template system issue. Signal immediately, STOP.

## Report Format

```
## Error Report for [DOC_NAME]

### Category A (Accepted):
- Accepted: <issue>

### Category B (Major Errors):
- B-ERROR: <desc> | Expected: <...> | Actual: <...> | Command: <...>

### Category C (Template Bugs):
- C-ERROR: <desc> (STOP if any)
```

## Status

Report: DONE (no B/C), DONE_WITH_CONCERNS (has B), or BLOCKED (has C)
```

## Documents Table

| Doc | Lines | Examples Dir |
|-----|-------|--------------|
| autoui.md | 938 | autoui_examples |
| commands.md | 452 | commands_examples |
| tutorial.md | 517 | tutorial_examples |
| testkit.md | 596 | testkit_examples |
| integration.md | 391 | integration_examples |

## Critical Rules

1. **B-errors require user permission** - Do NOT proceed without user decision
2. **No iteration limit** - Loop until all complete with user approval
3. **Signal-and-wait for C-errors** - Subagents stop on template bugs
4. **Launch in parallel** - All 5 subagents in one message with multiple Agent calls

## Start

Begin with Phase 1: Create tasks and launch all 5 subagents.
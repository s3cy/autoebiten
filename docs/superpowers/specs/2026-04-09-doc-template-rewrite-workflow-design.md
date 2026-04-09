# Document Template Rewrite Workflow Design

> **Purpose:** Define the coordination workflow for rewriting 5 documentation files to use the template system with verified runnable examples

---

## Goal

Rewrite all 5 documentation files (`autoui.md`, `commands.md`, `tutorial.md`, `testkit.md`, `integration.md`) to use the template system with runnable examples that produce verified outputs.

## Inputs

- Original documents: `docs/generate/original_docs/*.md` (backed up originals)
- Template system: `{{output "path" "lang"}}` syntax, config.yaml normalization
- Existing reference: `docs/generate/autoui_examples/` (partial implementation)

## Outputs per Document

- Template file: `docs/generate/<doc>.md.tmpl`
- Scripts directory: `docs/generate/<doc>_examples/*.sh`
- Config file: `docs/generate/<doc>_examples/config.yaml`
- Pre-staged outputs: `docs/generate/<doc>_examples/*_out.txt`
- Final document: `docs/<doc>.md` (generated via `make docs-generate`)

---

## Error Categorization

### Category A: Minor Format/Naming Issues

Cosmetic differences that don't affect functionality or correctness.

**Handling:** Accept and proceed. Document in error report with note.

**Examples:**
- Extra newline in output
- Widget attribute order differs
- Whitespace variations

### Category B: Major Autoebiten Errors

Documented behavior doesn't match actual implementation - the code works differently than documented.

**Handling:** Document with full context (expected, actual, command). Continue working. **Must discuss with user before proceeding to next iteration.**

**Examples:**
- Command documented but doesn't exist
- Flag behavior differs from documentation
- Output format mismatch
- Feature documented but not implemented

### Category C: Template System Bugs

The docgen system itself has issues preventing valid use.

**Handling:** Signal immediately, wait for coordinator decision. If confirmed, stop work on that document until fixed.

**Examples:**
- `{{output}}` function fails on valid path
- Normalization regex doesn't work
- Missing needed feature

---

## Coordination Workflow

### Phase 1: Launch

1. Coordinator creates task tracking list (5 tasks, one per document)
2. Launch 5 subagents in parallel with identical prompt structure
3. Each subagent assigned one document

### Phase 2: Execution

- Subagents work independently, no inter-agent communication
- Each tracks its own error list while working
- Completion states:
  - **DONE**: All deliverables complete, no B or C errors
  - **DONE_WITH_CONCERNS**: Complete but has B-errors to report
  - **BLOCKED**: Hit C-error, waiting for coordinator

### Phase 3: Aggregation

1. Wait for all 5 subagents to return
2. Collect error reports into consolidated list:
   - List A: Minor issues (informational)
   - List B: Major errors (must discuss with user)
   - List C: Template bugs (fix immediately)
3. **If C-errors**: fix template system, re-launch blocked subagents, return to Phase 2
4. **If B-errors**: present to user, do not proceed without permission

### Phase 4: User Decision on B-errors

1. Present consolidated B-error list to user
2. User decides: fix in autoebiten, fix in documentation, or defer
3. User gives permission to proceed
4. Implement user's decision
5. Re-launch subagents that had B-errors
6. Return to Phase 2

### Phase 5: Finalization

1. Run `make docs-generate` for all templates
2. Verify all 5 final `.md` files generated correctly
3. Commit changes

**No iteration limit** - loop continues until all documents complete with user approval.

---

## Subagent Prompt Structure

Each subagent receives identical instructions with document name varying:

**Part 1: Context**
- Template system overview
- Project structure and existing reference
- Available CLI commands

**Part 2: Task**
- Document assignment
- Deliverables list
- Process steps

**Part 3: Error Handling Protocol**
- Three categories with definitions
- Signal-and-wait for Category C
- Compile error report at end

**Part 4: Self-Check**
- Verify scripts executable, outputs captured
- Test local generation
- Report completion status

---

## Success Criteria

- All 5 documents have complete templates
- All scripts run successfully
- Outputs normalized correctly
- Final `.md` files match intent of originals
- No unresolved B or C errors
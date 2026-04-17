## Writing docs

- Never touch docs/*.md directly. Update the template files in docs/generate instead. Run `make docs` to re-generate. 
- When documenting an example that includes outputs, don't hard code. Always use template brackets `{{}}`. The template system will expand the brackets to real inputs/outputs. After `make docs`, read the re-generated files to make sure they match your intention.

## Running example games

- You are not in a headless environment. Games should run without problems. 
- Ignore `[CAMetalLayer nextDrawable] returning nil because allocation failed.` errors. They don't harm.

## graphify

This project has a graphify knowledge graph at graphify-out/.

Rules:
- Before answering architecture or codebase questions, read graphify-out/GRAPH_REPORT.md for god nodes and community structure
- If graphify-out/wiki/index.md exists, navigate it instead of reading raw files
- After modifying code files in this session, run `graphify update .` to keep the graph current (AST-only, no API cost)

# CLAUDE.md

Orientation for Claude Code sessions. Minimal by design — the README and
CONTRIBUTING carry what humans and agents share; this file only adds
agent-specific framing.

## Read first, in order

1. [README.md](./README.md) — what the library does + public surface
2. [CONTRIBUTING.md](./CONTRIBUTING.md) — layout, dev loop, conventions, PR flow
3. [`v1/api.go`](./v1/api.go) — public `Builder[T]` interface + `Result[T]`
4. [`examples/hello/main.go`](./examples/hello/main.go) — minimal call site

## Before you touch anything

- File map, module layout, and the conventions that bite (generics + godoc,
  e2e `-count=1`, skip release anchoring, annotated cosign tags) are all in
  [CONTRIBUTING.md](./CONTRIBUTING.md). Don't re-derive them from the diff.
- **Every change starts with an issue**, then a branch + PR with `Closes #<n>`.
  Don't push to `main`. Full flow in
  [CONTRIBUTING.md → Branch / PR flow](./CONTRIBUTING.md#branch--pr-flow).
- Don't commit secrets. [`.gitignore`](./.gitignore) covers `.env*`, `.claude/`.

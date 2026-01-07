# Contributing to gpx

Thanks for considering contributing to **gpx**!

This document describes how to work with the project,
run tests, and submit changes.

---

## Development setup

Requirements:
- Go 1.22 or newer
- Unix-like OS (macOS/Linux recommended for development)

Clone the repository and run tests:

```bash
go test ./...
```

---

## Project structure

```
cmd/gpx            # CLI entrypoint
internal/app       # use-cases and orchestration
internal/config    # config load/save/validate
internal/envx      # env parsing, quoting, export/unset
internal/shell     # apply to rc files (atomic replace)
internal/state     # active profile state
```

---

## Coding guidelines

- Prefer small, focused functions.
- Wrap errors with context:
  ```go
  return fmt.Errorf("read config: %w", err)
  ```
- Use sentinel errors + `errors.Is` / `errors.As` where appropriate.
- Avoid introducing external dependencies unless clearly justified.
- Keep CLI contracts stable (flags order, output format).

---

## Tests

Before submitting changes, ensure:

```bash
go test ./...
```

All tests must pass.

---

## Submitting changes

1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/my-change
   ```
3. Make your changes.
4. Run tests.
5. Commit with a clear message.
6. Open a Pull Request describing:
    - what changed
    - why it changed
    - how to test it

---

## Reporting issues

When opening an issue, please include:
- OS and shell (zsh/bash)
- Go version
- `gpx version` output
- exact commands executed
- actual vs expected behavior

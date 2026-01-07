# gpx

[Русская версия](README.ru.md)

`gpx` is a CLI tool to manage **environment variable presets** for Go workflows
(`GOPROXY`, `GOPRIVATE`, `GONOSUMDB`, `GOTOOLCHAIN`, and any other env vars).

It supports:
- switching **profiles** (multiple variables at once)
- temporary application via `eval "$(gpx use <profile>)"`
- persistent application via `gpx apply` (managed block in shell rc file)
- editing config via CLI (`gpx profile ...`)
- highlighting the **active profile** in `gpx list` (`*` = last used/applied)

---

## Important limitation

A CLI program **cannot permanently change the parent shell environment**.

Therefore:
- the safe default is printing `export ...` and using `eval`
- the persistent mode is explicit: `gpx apply`

---

## Install

### Build from source

```bash
go test ./...
go build -o gpx ./cmd/gpx
install -m 0755 gpx /usr/local/bin/gpx
```

---

## Quick start

### Create config

```bash
gpx init
```

### List profiles (active profile is marked with `*`)

```bash
gpx list
```

### Use a profile temporarily (recommended)

```bash
eval "$(gpx use public)"
```

### Apply a profile persistently

```bash
gpx apply public
source ~/.zshrc
```

---

## Commands

### gpx list

Lists profiles.  
The active profile (last used or applied) is marked with `*`.

### gpx status

Shows current values for environment variables
present in any profile.

### gpx use <profile>

Prints `export ...` lines for the profile
(using shell-safe quoting).

### gpx set KEY=VALUE [KEY=VALUE ...]

One-off exports without a profile.

### gpx unset KEY [KEY ...]

Prints `unset KEY` commands.

### gpx diff <profile>

Shows what would change compared to the current environment.

Example:

```
* GOPROXY: "https://proxy.corp.local,direct" -> "https://proxy.golang.org,direct"
  GOPRIVATE: "" -> ""
```

`*` means the value would change.

### gpx apply [flags] <profile>

Writes a managed block to shell rc file (`~/.zshrc` or `~/.bashrc`):

```
# GPX_BEGIN
export GOPROXY='...'
export GOPRIVATE='...'
# GPX_END
```

Flags:
- `--rc PATH` – explicit rc file path (overrides `--shell`)
- `--shell zsh|bash` – choose default rc file if `--rc` not provided
- `--dry-run` – show resulting content without writing files
- `--backup` – create timestamped backup before modification (disabled by default)

**CLI contract:** flags must come before positional arguments.

Correct:
```bash
gpx apply --rc /tmp/test.rc public
```

Incorrect:
```bash
gpx apply public --rc /tmp/test.rc
```

---

## Config editing

### Manage profiles

```bash
gpx profile add corp
gpx profile rm corp
gpx profile rename corp avito
```

### Show profile content

```bash
gpx profile show public
```

### Set / unset variables in a profile

```bash
gpx profile set corp GOPROXY=https://proxy.corp.local,direct GOPRIVATE=github.com/mycorp/*
gpx profile unset corp GOPRIVATE GONOSUMDB
```

---

## Configuration file

Default path:

```
~/.config/gpx/config.json
```

Example:

```json
{
  "profiles": {
    "public": {
      "GOPROXY": "https://proxy.golang.org,direct",
      "GOPRIVATE": "",
      "GONOSUMDB": "",
      "GOTOOLCHAIN": "auto"
    }
  }
}
```

---

## Version

```bash
gpx version
```

Example output:

```
gpx 0.1.0 (commit=abc123, date=2026-01-07)
```

---

## Project layout

```
cmd/gpx            # CLI entrypoint
internal/app       # use-cases and orchestration
internal/config    # config load/save/validate
internal/envx      # env parsing, quoting, export/unset
internal/shell     # apply to rc files (atomic replace)
internal/state     # active profile state
```

---

## Development

```bash
go test ./...
```

---

## License

MIT

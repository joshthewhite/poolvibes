# Taskfile Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a Taskfile.yml and .air.toml to standardize the build/test/dev workflow, with a `bin/` output directory for built artifacts.

**Architecture:** Single `Taskfile.yml` at project root defines all tasks. `.air.toml` configures live reload for `task dev`. Build output goes to `bin/poolvibes`. Docs and config files updated to reference `task` commands.

**Tech Stack:** [Task](https://taskfile.dev) v3, [air](https://github.com/air-verse/air)

---

### Task 1: Create Taskfile.yml

**Files:**
- Create: `Taskfile.yml`

**Step 1: Create the Taskfile**

```yaml
version: '3'

vars:
  BINARY: poolvibes
  BUILD_DIR: bin

tasks:
  default:
    desc: Show available tasks
    cmds:
      - task --list

  build:
    desc: Build the binary
    cmds:
      - go build -o {{.BUILD_DIR}}/{{.BINARY}} .
    sources:
      - '**/*.go'
      - go.mod
      - go.sum
    generates:
      - '{{.BUILD_DIR}}/{{.BINARY}}'

  test:
    desc: Run tests
    cmds:
      - go test ./...

  test:verbose:
    desc: Run tests with verbose output
    cmds:
      - go test -v ./...

  lint:
    desc: Run linter
    cmds:
      - go vet ./...

  fmt:
    desc: Format code
    cmds:
      - gofmt -w .

  templ:
    desc: Generate templ templates
    cmds:
      - templ generate

  dev:
    desc: Start dev server with live reload
    deps: [templ]
    cmds:
      - air

  run:
    desc: Build and run the server
    deps: [build]
    cmds:
      - ./{{.BUILD_DIR}}/{{.BINARY}} serve

  clean:
    desc: Remove build artifacts
    cmds:
      - rm -rf {{.BUILD_DIR}}

  tidy:
    desc: Tidy Go module dependencies
    cmds:
      - go mod tidy

  docker:build:
    desc: Build Docker image
    cmds:
      - docker build -t {{.BINARY}} .

  docker:up:
    desc: Start all Docker services
    cmds:
      - docker-compose up

  docker:down:
    desc: Stop all Docker services
    cmds:
      - docker-compose down
```

**Step 2: Verify it works**

Run: `task --list`
Expected: All tasks listed with descriptions.

Run: `task build`
Expected: Binary created at `bin/poolvibes`.

Run: `task test`
Expected: All tests pass.

**Step 3: Commit**

```bash
git add Taskfile.yml
git commit -m "Add Taskfile for standardized build/test/dev workflow"
```

---

### Task 2: Create .air.toml

**Files:**
- Create: `.air.toml`

**Step 1: Create air config**

```toml
root = "."
tmp_dir = "bin"

[build]
  pre_cmd = ["templ generate"]
  cmd = "go build -o ./bin/poolvibes ."
  bin = "./bin/poolvibes"
  args_bin = ["serve"]
  include_ext = ["go", "templ"]
  exclude_dir = ["bin", "vendor", "e2e", "docs", "node_modules", "migrations"]
  exclude_regex = ["_templ\\.go$", "_test\\.go$"]
  delay = 1000

[log]
  time = false

[misc]
  clean_on_exit = true
```

**Step 2: Verify it works**

Run: `task dev`
Expected: templ generates, air starts watching, server starts on :8080. Ctrl+C to stop.

**Step 3: Commit**

```bash
git add .air.toml
git commit -m "Add air config for live reload dev server"
```

---

### Task 3: Update .gitignore

**Files:**
- Modify: `.gitignore:1-8`

**Step 1: Replace binary entries with bin/**

Replace lines 7-8 (`poolio` and `poolvibes`) with `bin/`. Also add `.air.toml` tmp directory entry isn't needed since we use `bin/` which is already covered. Add the air tmp file:

In `.gitignore`, replace:
```
poolio
poolvibes
```
With:
```
bin/
```

Also add under the OS section:
```
# Air
tmp/
```

**Step 2: Verify**

Run: `git status`
Expected: `.gitignore` modified, `poolvibes` binary (if present) now untracked or gone.

**Step 3: Commit**

```bash
git add .gitignore
git commit -m "Update gitignore: use bin/ directory for build output"
```

---

### Task 4: Update CLAUDE.md

**Files:**
- Modify: `CLAUDE.md:35-49`

**Step 1: Update the Development/Commands section**

Replace the Commands section (lines 40-49) with:

```markdown
### Commands

Uses [Task](https://taskfile.dev) for build automation. Run `task --list` to see all available tasks.

- **Generate templates**: `task templ` (required after editing `.templ` files; generated `*_templ.go` files are committed)
- **Build**: `task build` (outputs to `bin/poolvibes`)
- **Test**: `task test`
- **Test (verbose)**: `task test:verbose`
- **Test (single)**: `go test -v -run TestName ./path/to/package`
- **Lint**: `task lint`
- **Format**: `task fmt`
- **Tidy deps**: `task tidy`
- **Dev server (live reload)**: `task dev` (runs templ generate + air)
- **Build and run**: `task run`
- **Clean**: `task clean`
```

**Step 2: Commit**

```bash
git add CLAUDE.md
git commit -m "Update CLAUDE.md with task commands"
```

---

### Task 5: Update README.md

**Files:**
- Modify: `README.md:30-33` (Getting Started)
- Modify: `README.md:116-125` (Development)

**Step 1: Update Getting Started**

Replace lines 30-33:
```markdown
```sh
go build -o poolvibes .
./poolvibes serve
```
```

With:
```markdown
```sh
task build
task run
```
```

**Step 2: Update Development section**

Replace lines 116-125:
```markdown
## Development

```sh
templ generate          # regenerate templates (after editing .templ files)
go build ./...          # build
go vet ./...            # lint
go test ./...           # test
gofmt -w .              # format
go mod tidy             # tidy deps
```
```

With:
```markdown
## Development

Uses [Task](https://taskfile.dev) for build automation. Install: `go install github.com/go-task/task/v3/cmd/task@latest`

```sh
task --list             # show all tasks
task build              # build to bin/poolvibes
task test               # run tests
task lint               # lint
task fmt                # format
task templ              # regenerate templates
task dev                # dev server with live reload (air)
task run                # build and run server
task tidy               # tidy deps
task clean              # remove build artifacts
task docker:build       # build Docker image
task docker:up          # start Docker services
task docker:down        # stop Docker services
```
```

**Step 3: Commit**

```bash
git add README.md
git commit -m "Update README with task commands"
```

---

### Task 6: Update docs/development.md

**Files:**
- Modify: `docs/development.md:1-19`

**Step 1: Update prerequisites and commands table**

Replace lines 1-19 with:

```markdown
# Development

## Prerequisites

- Go 1.25 or later
- [Task](https://taskfile.dev) — `go install github.com/go-task/task/v3/cmd/task@latest`
- [templ](https://templ.guide) CLI — `go install github.com/a-h/templ/cmd/templ@latest` (only needed when editing `.templ` files)
- [air](https://github.com/air-verse/air) — `go install github.com/air-verse/air@latest` (only needed for `task dev`)

## Commands

| Command | Description |
|---------|-------------|
| `task --list` | Show all available tasks |
| `task build` | Build binary to `bin/poolvibes` |
| `task test` | Run tests |
| `task test:verbose` | Run tests with verbose output |
| `go test -v -run TestName ./path/to/package` | Run a single test |
| `task lint` | Lint |
| `task fmt` | Format code |
| `task templ` | Regenerate Go code from `.templ` files |
| `task dev` | Start dev server with live reload (templ + air) |
| `task run` | Build and run the server |
| `task clean` | Remove build artifacts |
| `task tidy` | Tidy dependencies |
| `task docker:build` | Build Docker image |
| `task docker:up` | Start Docker services |
| `task docker:down` | Stop Docker services |
```

**Step 2: Commit**

```bash
git add docs/development.md
git commit -m "Update development docs with task commands"
```

---

### Task 7: Final verification

**Step 1: Run all key tasks**

```bash
task build
task test
task lint
task clean
```

Expected: All pass, `bin/` created and cleaned.

**Step 2: Verify git is clean**

```bash
git status
```

Expected: No untracked or modified files (except possibly `bin/` which is gitignored).

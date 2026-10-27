# Contributing to the ModelHelper CLI

Thanks for wanting to help. Contributions are welcome: bug reports, fixes, new features and documentation. Please read this page before you open a pull request.

- [Code of Conduct](./CODE_OF_CONDUCT.md)
- [Reporting issues](#reporting-issues)
- [Pull requests](#pull-requests)
- [Development setup](#development-setup)
- [Committing changes](#committing-changes)
- [Project structure](#project-structure)

## Reporting issues

Use [GitHub issues](https://github.com/Modelhelper/cli/issues/new/choose). For a bug, include the `mh` command you ran, what you expected, what happened, and the output of `mh version`.

Templates and agent skills live in their own repositories: [Modelhelper/templates](https://github.com/Modelhelper/templates) and [Modelhelper/skills](https://github.com/Modelhelper/skills). Report template problems there.

## Pull requests

- Branch from `main` and open the pull request against `main`.
- For a new feature, open an issue first and describe what you want to add and why. That way you know it fits before you spend time on it.
- For a bug fix, describe the bug in the pull request, and say how you reproduced it. Reference the issue if there is one.
- Keep a pull request to one change. Several small commits are fine, they are squashed when merged.
- Add or update tests when you change behaviour, and update the docs when a command or option changes.
- Run `gofmt` and `go vet ./...` on what you change.

## Development setup

You need Go (see `go.mod` for the version).

```bash
git clone https://github.com/Modelhelper/cli.git
cd cli
go build -o mh ./cmd/mh
go test ./...
```

You can try the CLI without a database or a config for your own sources by using the demo source:

```bash
go run ./cmd/mh code generate -t go-struct --demo -s order --screen
```

Code generation needs templates, see [Modelhelper/templates](https://github.com/Modelhelper/templates) and run `mh` once to set up your config.

A few tests fail on a clean checkout today because they depend on the machine they run on (`ports/code`, `ports/connection` and `utils/path`). If a test you did not touch fails, check that it also fails on `main` before you spend time on it.

## Committing changes

Commit messages follow the [commit message convention](./COMMIT_CONVENTION.md), for example `fix(project): keep empty folders` or `feat(code): add --db option`. The convention is not enforced by a hook, so check your messages before you push.

## Project structure

| Folder | Contents |
|--------|----------|
| `cmd/mh` | The `main` package, builds the `mh` binary |
| `modelhelper` | Interfaces and models shared by the rest of the code |
| `ports` | The implementations: `command` (cobra commands), `code`, `template`, `project`, `source` (mssql, postgres, demo), `config`, `connection` and so on |
| `ui` | Terminal tables, prompts and colours |
| `utils` | Small helpers: text casing, template functions, paths |
| `schemas` | JSON schemas for config files |

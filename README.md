# ModelHelper CLI

ModelHelper (`mh`) reads the structure of a database and turns it into code. Point it at a table, pick a template, and get a struct, a repository, a controller, or whatever the template describes. The templates are plain YAML files, so you can change them or write your own.

> **Status:** beta. The CLI is being rewritten from C# to Go, so some commands are unfinished and details may change.

## What it does

- **Explore a database.** List tables, schemas and columns from a connection.
- **Generate code from templates.** Combine a source table with one or more templates, and print the result, copy it to the clipboard, or write it to the right folders in your project.
- **Create projects from templates.** Scaffold a whole project, folders and files, from a project template.
- **Copy data between sources.** Copy tables from one SQL Server connection to another.

Supported sources: SQL Server (`mssql`) and PostgreSQL (`postgres`). There is also a built-in demo source, so you can try the CLI without a database.

## Install

There are no packaged releases yet, so build it from source. You need Go (see `go.mod` for the version).

```bash
git clone https://github.com/Modelhelper/cli.git
cd cli
go build -o mh ./main.go
```

Put the `mh` binary somewhere on your `PATH`.

You also need a set of templates. They live in a separate repository, [Modelhelper/templates](https://github.com/Modelhelper/templates). Clone it next to the CLI.

## Quick start

The first time you run `mh` without a config, it starts a wizard and creates `~/.modelhelper/config.yaml`. The wizard asks for a connection, the folders with your templates, and your programming languages.

Try it with the demo source first. No database needed:

```bash
# list the templates for Go
mh code template list --lang go

# preview a struct for the demo table "order"
mh code generate -t go-struct --demo -s order --screen
```

Then use a real database:

```bash
# add a connection (see `mh connection create --help` for the types)
mh connection create --help

# list the tables
mh source list

# generate a repository for dbo.Order and copy it to the clipboard
mh code generate -t go-db-repo -s dbo.Order --copy
```

Useful flags for `mh code generate`:

| Flag | Meaning |
|------|---------|
| `-t`, `--template` | Template to use. Repeat it for several templates |
| `-s`, `--source` | Table to use as the model, as `Schema.Table`. Repeat it for several tables |
| `-c`, `--connection` | Connection to use. The default connection is used if you leave it out |
| `--screen` | Print the generated code |
| `--copy` | Copy the generated code to the clipboard |
| `--export-bykey` | Write the files into the folders defined in the project config (see below) |
| `--overwrite` | Replace files that already exist |
| `--code-only` | Print only the code, with no statistics |
| `--demo` | Use the built-in sample tables |

## Using it in a project

Add a `.modelhelper` folder at the root of your project, next to `.git`. A `config.yaml` in it tells `mh` where each kind of file belongs:

```yaml
version: "3.0"
name: flight-api
language: go

locations:
  model: ./src/core/models
  core-command: ./src/core/commands
  api-controller: ./src/adapters/api/controllers
```

Every template has a `key` and a `fileName`. With `--export-bykey`, `mh` looks up the key in `locations` and writes the file there, using the template's file name:

```bash
mh code generate --export-bykey \
  -t go-hex-feat-db-get-model \
  -t go-hex-gin-controller \
  -s dbo.Order
```

Run `mh` from the folder that contains `.modelhelper`. Templates in `.modelhelper/templates/` replace global templates with the same name, which lets one project tweak a template without touching the shared set.

## Create a project from a template

```bash
mh project template list
mh project create -t go-hex-api -n flight-api --dry-run   # show what would be created
mh project create -t go-hex-api -n flight-api             # create it
```

A project template is a folder of files with placeholders such as the project name. The folder structure is kept, including empty folders. Add `-d <dir>` to choose where it is created.

## Commands

| Command | What it does |
|---------|--------------|
| `mh code generate` | Generate code from templates and a source |
| `mh code template list` | List templates, filter by language, key, type or model |
| `mh code template detail <name>` | Show one template |
| `mh source list` | List the tables in a connection |
| `mh source copy` | Copy data from one source to another |
| `mh connection create/list/default/delete` | Manage connections |
| `mh project create` | Create a project from a template |
| `mh project` | Show the project in the current folder |
| `mh language list` | List the configured languages |
| `mh config open` | Open the config file in an editor |
| `mh serve` | Start the ModelHelper web server |

Run `mh <command> --help` for all options.

## Templates

Templates are YAML files with a Go [`text/template`](https://pkg.go.dev/text/template) body. They have a language, a type (`file`, `block` or `snippet`), a key, an optional file name, and the code to generate. Blocks can be included in other templates, which is how SQL for SQL Server and PostgreSQL is kept apart from the Go code that uses it.

The template format and the full set of templates are documented in [Modelhelper/templates](https://github.com/Modelhelper/templates).

## Agent skills

[Modelhelper/skills](https://github.com/Modelhelper/skills) has skills for Claude Code that help you turn existing code into a template and set up the `.modelhelper` folder in a project.

## Contributing

You are welcome to contribute, but for now I will not approve any pull requests. The project is in the middle of a rewrite from C# to Go.

## Community

Join the community on Slack: [Join ModelHelper](https://join.slack.com/t/model-helper/shared_invite/zt-ruud19ok-RFe871ZyFjLHvNGilcjMlg)

## License

MIT, see [LICENSE](LICENSE).

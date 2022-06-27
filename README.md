# GMY

> Work in progress!

A tiny command line tool to setup the needed files to quickly create Telegram
bots powered by the [grammY bot framework](https://grammy.dev). This tool allows
you to create projects from several [templates](templates.json),
maintained by third-party users.

This is almost a complete Go re-write of the
[official grammY CLI](https://github.com/grammyjs/create-grammy).

## Install

Install using [Go](https://go.dev).

```shell
go install github.com/dcdunkan/gmy@latest
```

After installation, run **gmy** command to use the tool. You can provide a
project name as the first argument.

## Templates

Open a pull request by adding your own templates to the
[templates.json](templates.json) file. There are currently three platforms that
you can add templates to: Deno, Node.js, and other templates.

Each template should contain the following fields:

- `name` — Name to be shown in the templates list in CLI. Recommended to use
  "owner/repository" as the name if it's a repository.
- `type` — **repository** or **subfolder**. If your template is an entire
  repository, use "repository" as type, or if it is a subfolder in a repository
  use "subfolder" as the type.
- `owner` — GitHub repository owner.
- `repository` — GitHub repository name.
- `docker_prompt` — Should the CLI prompt the user to add
  [default docker files](internal/files/dockerfiles.go).
- `tsconfig_prompt` — Should the CLI prompt the user to add the default
  [tsconfig.json](configs/tsconfig.json) file.

#### "repository" type

- `branch` — Primary repository branch name.

#### "subfolder" type

- `path` — Path to the subfolder the template is located at.

#### Deno templates

- `cache_file` — When the user chose to Cache dependencies, `deno cache` command
  will get executed for the specified file.

## GMY?

Short form of grammY, for a command name.

## TODO

- [ ] Documentation
- [ ] New templates
- [ ] Setup Releases

## License

Licensed under MIT © 2022 Dunkan.

[See License →](LICENSE)

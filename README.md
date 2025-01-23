# Modo🧯

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/modo/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/modo/actions/workflows/tests.yml)
[![stable](https://img.shields.io/github/actions/workflow/status/mlange-42/modo/test-stable.yml?branch=main&label=stable&logo=github)](https://github.com/mlange-42/modo/actions/workflows/test-stable.yml)
[![nightly](https://img.shields.io/github/actions/workflow/status/mlange-42/modo/test-nightly.yml?branch=main&label=nightly&logo=github)](https://github.com/mlange-42/modo/actions/workflows/test-nightly.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/modo)](https://goreportcard.com/report/github.com/mlange-42/modo)
[![User Guide](https://img.shields.io/badge/user_guide-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://mlange-42.github.io/modo/)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/modo)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/modo)
[![MIT license](https://img.shields.io/badge/MIT-brightgreen?label=license)](https://github.com/mlange-42/modo/blob/main/LICENSE)

Modo🧯 is a documentation generator (DocGen) for the [Mojo](https://www.modular.com/mojo)🔥 programming language.
It generates Markdown for static site generators (SSGs) from `mojo doc` JSON output.

[This example](https://mlange-42.github.io/modo/mypkg/) in the [User guide](https://mlange-42.github.io/modo/) shows a Mojo🔥 package processed with Modo🧯 and rendered with [Hugo](https://gohugo.io), to demonstrate Modo🧯's features.

## Features

* Generates [Mojo](https://www.modular.com/mojo)🔥 API docs for [Hugo](https://mlange-42.github.io/modo/guide/formats#hugo), [mdBook](https://mlange-42.github.io/modo/guide/formats#mdbook) or just [plain](https://mlange-42.github.io/modo/guide/formats#plain-markdown) Markdown.
* Provides a simple syntax for code [cross-references](https://mlange-42.github.io/modo/guide/cross-refs).
* Optionally structures API docs according to [package re-exports](https://mlange-42.github.io/modo/guide/re-exports).
* Optionally extracts [doc-tests](https://mlange-42.github.io/modo/guide/doctests) for `mojo test` from code blocks.
* Customizable output through [user templates](https://mlange-42.github.io/modo/guide/templates).

See the [User guide](https://mlange-42.github.io/modo/) for more information.

## Installation

### Using Python

Modo🧯 is available on PyPI as [`pymodo`](https://pypi.org/project/pymodo/).
Install it with pip:

```
pip install pymodo
```

> This installs the `modo` command. If the command is not found, try:  
> `python -m pymodo`

### Using Go

With [Go](https://go.dev) installed, you can install Modo🧯 like this:
```shell
go install github.com/mlange-42/modo@latest
```

### Precompiled binaries

Pre-compiled binaries for manual installation are available in the
[Releases](https://github.com/mlange-42/modo/releases)
for Linux, Windows and MacOS.

## Usage

In your Mojo🔥 project, set up Modo🧯:

```shell
modo init
```

This sets up the project with default settings and paths.
See the generated `modo.yaml` file to modify them.

Next, run `mojo doc` to extract the API docs in JSON format:

```shell
mojo doc src/ -o api.json
```

Finally, build the Markdown documentation:

```shell
modo build
```

See the [User guide](https://mlange-42.github.io/modo/) for more information.

## Packages using Modo🧯

- [Larecs](https://github.com/samufi/larecs) -- a performance-centred archetype-based ECS ([docs](https://samufi.github.io/larecs/)).

## License

This project is distributed under the [MIT license](./LICENSE).

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
* Super easy to [set up](https://mlange-42.github.io/modo/guide/setup) for an existing Mojo🔥 project.
* Provides a simple syntax for code [cross-references](https://mlange-42.github.io/modo/guide/features/crossrefs).
* Optionally structures API docs according to [package re-exports](https://mlange-42.github.io/modo/guide/features/reexports).
* Optionally extracts [doc-tests](https://mlange-42.github.io/modo/guide/features/doctests) for `mojo test` from code blocks.
* Customizable output through [user templates](https://mlange-42.github.io/modo/guide/features/templates).

See the [User guide](https://mlange-42.github.io/modo/) for more information.
See the Modo🧯 [slides](https://mlange-42.github.io/modo/slides/) for a feature overview.

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

```
go install github.com/mlange-42/modo@latest
```

> With Go, you can also install the latest **development version**:  
> `go install github.com/mlange-42/modo@main`

### Precompiled binaries

Pre-compiled binaries for manual installation are available in the
[Releases](https://github.com/mlange-42/modo/releases)
for Linux, Windows and MacOS.

## Usage

To initialize an existing Mojo🔥 project for Modo🧯 and an SSG like [Hugo](https://gohugo.io), run command `init` once:

```
modo init hugo
```

This sets up everything to be able to build Markdown files for the target SSG with command `build`:

```
modo build
```

Finally, serve or build the site with the target SSG (here Hugo):

```
hugo serve -s docs/site/
```

See [Project setup](https://mlange-42.github.io/modo/guide/setup/) for details and other supported SSGs.

## Packages using Modo🧯

- [Larecs🌲](https://github.com/samufi/larecs) -- a performance-centred archetype-based ECS ([docs](https://samufi.github.io/larecs/)).
- [ExtraMojo](https://github.com/ExtraMojo/ExtraMojo) -- a collection of useful things that aren't (yet) in the standard library ([docs](https://extramojo.github.io/ExtraMojo/)).

## License

This project is distributed under the [MIT license](./LICENSE).

# Modo

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/modo/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/modo/actions/workflows/tests.yml)
[![stable](https://img.shields.io/github/actions/workflow/status/mlange-42/modo/test-stable.yml?branch=main&label=stable&logo=github)](https://github.com/mlange-42/modo/actions/workflows/test-stable.yml)
[![nightly](https://img.shields.io/github/actions/workflow/status/mlange-42/modo/test-nightly.yml?branch=main&label=nightly&logo=github)](https://github.com/mlange-42/modo/actions/workflows/test-nightly.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/modo)](https://goreportcard.com/report/github.com/mlange-42/modo)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/modo)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/modo)
[![MIT license](https://img.shields.io/badge/MIT-brightgreen?label=license)](https://github.com/mlange-42/modo/blob/main/LICENSE)

Modo is a documentation generator (DocGen) for the [Mojo](https://www.modular.com/mojo) programming language.
It generates Markdown for static site generators (SSGs) from `mojo doc` JSON output.

[This example](https://mlange-42.github.io/modo/) shows the Mojo [stdlib](https://github.com/modularml/mojo) processed with Modo and rendered with [mdBook](https://github.com/rust-lang/mdBook).

**! Early work in progress !**

## Installation

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/modo/releases).

> Alternatively, install using [Go](https://go.dev):
> ```shell
> go install github.com/mlange-42/modo/cmd/modo@latest
> ```

## Usage

Pipe `mojo doc` to Modo:

```
mojo doc <src-path> | modo <out-dir>
```

Alternatively, use a file:

```
mojo doc <src-path> -o docs.json
modo <out-dir> -i docs.json
```

Command line help:

```
modo -h
```

## Output formats

Modo emits Markdown files.
These files need to be processed further to generate an HTML site that can be served on GitHub pages (or elsewhere).
Modo supports different formats to make this step easier:

### Plain Markdown

Just plain markdown files.
This is Modo's default output format.

### mdBook

Markdown files as well as auxiliary files for [mdBook](https://github.com/rust-lang/mdBook),
with flag `--mdbook`.
Modo's output folder can be used by mdBook without any further steps:

```
modo docs-out -i docs.json --mdbook
mdbook serve docs-out --open
```

### Hugo

Not yet implemented.

## Cross-referencing

Modo supports cross-refs within the documentation of a project.
Absolute as well as relative references are supported.
Relative references follow Mojo's import syntax, with a leading dot denoting the current module, and further dots navigating upwards.

Some examples:

| Ref | Explanation |
|-----|-------------|
| `[pkg.mod.A]` | Absolute reference. |
| `[.A]` | Struct `A` in the current module. |
| `[.A.method]` | Method `method` of struct `A` in the current module. |
| `[..mod.A]` | Struct `A` in sibling module `mod`. |
| `[.A.method method]` | Method `method` of struct `A`, with custom text. |

Leading dots are stripped from the link text if no custom text is given.

Besides that, normal Markdown links can be used.

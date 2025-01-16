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

## Features

* Generates API docs websites for [Hugo](#hugo), [mdBook](#mdbook) or just [plain](#plain-markdown) Markdown.
* Resolves and renders [cross-references](#cross-referencing).
* Optionally structures docs according to [package re-exports](#package-re-exports).

## Installation

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/modo/releases).

> Alternatively, install using [Go](https://go.dev):
> ```shell
> go install github.com/mlange-42/modo@latest
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
Modo supports different formats to make this step easier, via the flag `--format`:

### Plain Markdown

Just plain markdown files.
This is Modo's default output format.

### mdBook

Markdown files as well as auxiliary files for [mdBook](https://github.com/rust-lang/mdBook),
with `--format=mdbook`.
Modo's output folder can be used by mdBook without any further steps:

```
modo docs-out -i docs.json --format=mdbook
mdbook serve docs-out --open
```

### Hugo

Markdown files with front matter and cross-references for [Hugo](https://gohugo.io/),
with flag `--format=hugo`.

You should first set up a Hugo project in a sub-folder of your repository.
Then, run Modo with the Hugo `content` folder as output path:

```
modo <hugo-project>/content -i docs.json --format=hugo
```

Further, in your `hugo.toml`, add `disablePathToLower = true` to the main section
to prevent lower case members (like functions) and upper case members (like structs)
overwrite each other.
Alternatively, run Modo with switch `--case-insensitive`.

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

Leading dots are stripped from the link text if no custom text is given, so `.mod.Type` becomes `mod.Type`.
With flag `--short-links`, modules are also stripped, so `.mod.Type` becomes just `Type`.

Besides cross-references, normal Markdown links can be used in doc-strings.

## Package re-exports

In mojo, package-level re-exports (or rather, imports) can be used
to flatten the structure of a package and shorten import paths for users.

Modo can structure documentation output according to re-exports using the flag `--exports`.
However, as we don't look at the actual code but just `mojo doc` JSON,
these re-exports must be documented in an `Exports:` section in the package docstring.

In a package's `__init__.mojo`, document re-exports like this:

```python
"""
Package creatures demonstrates Modo re-exports.

Exports:
 - animals.vertebrates.Cat
 - animals.vertebrates.Dog
 - plants.vascular
 - fungi
"""
from .animals.vertebrates import Cat, Dog
from .plants import vascular
```

> Note that `Exports:` should not be the first line of the docstring, as it is considered the summary and is not processed.

When processed with `--exports`, only exported members are included in the documentation.
Re-exports are processed recursively.
This means that sub-packages need an `Exports:` section too if they are re-exported directly,
like `fungi` in the example.
For exporting members from a sub-package (like `Cat` and `Doc`), the sub-package `Exports:` are ignored.
Re-exported modules (like `plants.vascular`) are included completely.

[Cross-references](#cross-referencing) should still use the original structure of the package.
They are automatically transformed to match the altered structure.

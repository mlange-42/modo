---
title: Getting started
type: docs
summary: Installation and basic usage of Modo🧯.
prev: guide
weight: 1
---

## Installation

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/modo/releases).

> Alternatively, install using [Go](https://go.dev):
> ```shell {class="no-wrap"}
> go install github.com/mlange-42/modo@latest
> ```

## Usage

In your Mojo🔥 project, set up Modo🧯:

```shell {class="no-wrap"}
modo init
```

This sets up the project with default settings and paths.
See the generated `modo.yaml` file to modify them.

Next, run `mojo doc` to extract the API docs in JSON format:

```shell {class="no-wrap"}
mojo doc src/ -o api.json
```

Finally, build the Markdown documentation:

```shell {class="no-wrap"}
modo build
```

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

Pipe `mojo doc` to Modo🧯:

``` {class="no-wrap"}
mojo doc src/ | modo -o docs/
```

Alternatively, use a file:

``` {class="no-wrap"}
mojo doc src/ -o api.json
modo -i api.json -o docs/
```

Get CLI help with `modo --help`.

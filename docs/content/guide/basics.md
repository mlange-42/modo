---
title: Getting started
type: docs
summary: Installation and basic usage of Modo🧯.
weight: 1
---

## Installation

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/modo/releases).

> Alternatively, install using [Go](https://go.dev):
> ```shell
> go install github.com/mlange-42/modo@latest
> ```

## Usage

Pipe `mojo doc` to Modo🧯:

```
mojo doc <src-path> | modo <out-dir>
```

Alternatively, use a file:

```
mojo doc <src-path> -o docs.json
modo <out-dir> -i docs.json
```

Get CLI help with `modo --help`.

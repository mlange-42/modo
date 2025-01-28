---
title: Getting started
type: docs
summary: Installation and basic usage of Modo🧯.
prev: guide
weight: 1
---

## Installation

### Using Python

Modo🧯 is available on PyPI as [`pymodo`](https://pypi.org/project/pymodo/).
Install it with pip:

```{class="no-wrap"}
pip install pymodo
```

> This installs the `modo` command. If the command is not found, try:  
> `python -m pymodo`

### Using Go

With [Go](https://go.dev) installed, you can install Modo🧯 like this:

```{class="no-wrap"}
go install github.com/mlange-42/modo@latest
```

### Precompiled binaries

Pre-compiled binaries for manual installation are available in the
[Releases](https://github.com/mlange-42/modo/releases)
for Linux, Windows and MacOS.

## Usage

In your Mojo🔥 project, set up Modo🧯:

```{class="no-wrap"}
modo init
```

This sets up the project with default settings and paths.
See the generated `modo.yaml` file to modify them.

Next, run `mojo doc` to extract the API docs in JSON format:

```{class="no-wrap"}
mojo doc src/ -o api.json
```

Finally, build the Markdown documentation:

```{class="no-wrap"}
modo build
```

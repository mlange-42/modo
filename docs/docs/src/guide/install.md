---
title: Installation
type: docs
summary: Installation options for Modo🧯.
prev: guide
weight: 10
---
Modo🧯 offers different installation options.

## Using Python

Modo🧯 is available on PyPI as [`pymodo`](https://pypi.org/project/pymodo/).
Install it with pip:

```shell {class="no-wrap"}
pip install pymodo
```

> This installs the `modo` command. If the command is not found, try:  
> `python -m pymodo`

## Using Go

With [Go](https://go.dev) installed, you can install Modo🧯 like this:

```shell {class="no-wrap"}
go install github.com/mlange-42/modo@latest
```

> With Go, you can also install the latest **development version**:
> ```
> go install github.com/mlange-42/modo@main
> ```

## Precompiled binaries

Pre-compiled binaries for manual installation are available in the
[Releases](https://github.com/mlange-42/modo/releases)
for Linux, Windows and MacOS.

## Using Pixi

Not supported yet. Help is appreciated, see issue [#133](https://github.com/mlange-42/modo/issues/133).

---
title: Modo
type: docs
summary: Modo is a documentation generator (DocGen) for the Mojo🔥 programming language.
---
Modo🧯 is a documentation generator (DocGen) for the [Mojo](https://www.modular.com/mojo)🔥 programming language.

It generates Markdown for static site generators (SSGs) from `mojo doc` JSON output.

{{< html >}}
<div style="width 100%; text-align: center; margin-top: 1rem; margin-bottom: 1rem;">

<a href="https://github.com/mlange-42/modo/actions/workflows/tests.yml" style="display:inline-block">
<img alt="Test status" src="https://img.shields.io/github/actions/workflow/status/mlange-42/modo/tests.yml?branch=main&label=Tests&logo=github" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>

<a href="https://github.com/mlange-42/modo/actions/workflows/test-stable.yml" style="display:inline-block">
<img alt="stable" src="https://img.shields.io/github/actions/workflow/status/mlange-42/modo/test-stable.yml?branch=main&label=stable&logo=github" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>

<a href="https://github.com/mlange-42/modo/actions/workflows/test-nightly.yml" style="display:inline-block">
<img alt="nightly" src="https://img.shields.io/github/actions/workflow/status/mlange-42/modo/test-nightly.yml?branch=main&label=nightly&logo=github" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>

<a href="https://goreportcard.com/report/github.com/mlange-42/modo" style="display:inline-block">
<img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mlange-42/modo" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>
<br />
<a href="https://mlange-42.github.io/modo/" style="display:inline-block">
<img alt="User Guide" src="https://img.shields.io/badge/user_guide-%23007D9C?logo=go&logoColor=white&labelColor=gray" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>

<a href="https://pkg.go.dev/github.com/mlange-42/modo" style="display:inline-block">
<img alt="Go Reference" src="https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>

<a href="https://github.com/mlange-42/modo" style="display:inline-block">
<img alt="GitHub" src="https://img.shields.io/badge/github-repo-blue?logo=github" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>

<a href="https://github.com/mlange-42/modo/blob/main/LICENSE" style="display:inline-block">
<img alt="MIT license" src="https://img.shields.io/badge/MIT-brightgreen?label=license" style="margin-top: 0.2rem; margin-bottom: 0.2rem;"></img></a>

</div>
{{< /html >}}

## Features

* Generates [Mojo](https://www.modular.com/mojo)🔥 API docs for [Hugo](guide/formats#hugo), [mdBook](guide/formats#mdbook) or just [plain](guide/formats#plain-markdown) Markdown.
* Super easy to [set up](guide/setup) for an existing Mojo🔥 project.
* Provides a simple syntax for code [cross-references](guide/features/crossrefs).
* Optionally structures API docs according to [package re-exports](guide/features/reexports).
* Optionally extracts [doc-tests](guide/features/doctests) for `mojo test` from code blocks.
* Customizable output through [user templates](guide/features/templates).

See also the Modo🧯 [slides](https://mlange-42.github.io/modo/slides/) for a feature overview.

## Usage

For usage examples and details, see the [User guide](guide).

## Demo

See the [Example API docs](mypkg) for a demo of Modo🧯's features.

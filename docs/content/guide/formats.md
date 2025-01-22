---
title: Output formats
type: docs
summary: Modo🧯's output formats.
weight: 3
---

Modo🧯 emits Markdown files.
These files need to be processed further to generate an HTML site that can be served on GitHub Pages (or elsewhere).
Modo🧯 supports different formats to make this step easier, via the config field `format` or flag `--format`.

## Plain Markdown

Just plain markdown files.
This is Modo🧯's default output format.
The generated files are suitable for GitHub's Markdown rendering.

## mdBook

Markdown files as well as auxiliary files for [mdBook](https://github.com/rust-lang/mdBook),
with `format: mdbook` in the `modo.yaml` or with flag `--format=mdbook`.
The generated files can be used by mdBook without any further steps:

``` {class="no-wrap"}
modo build
mdbook serve docs/ --open
```

[Templates](../templates) can be used to customize the mdBook configuration file `book.toml`.

## Hugo

Markdown files with front matter and cross-references for [Hugo](https://gohugo.io/),
with `format: hugo` in the `modo.yaml` or with flag `--format=hugo`.

You should first set up a Hugo project in a sub-folder of your repository.
Then, use the Hugo `content` folder as output path,
either by editing `modo.yaml` or via flag `--output`.

``` {class="no-wrap"}
modo build --output=docs/content
```

Further, in your `hugo.toml`, add `disablePathToLower = true` to the main section
to prevent lower case members (like functions) and upper case members (like structs)
overwrite each other.
Alternatively, run Modo🧯 with `case-insensitive: true` or flag `--case-insensitive`.

[Templates](../templates) can be used to customize the Hugo front matter of each page.

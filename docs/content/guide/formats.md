---
title: Output formats
type: docs
summary: Modo🧯's output formats.
weight: 2
---

Modo🧯 emits Markdown files.
These files need to be processed further to generate an HTML site that can be served on GitHub pages (or elsewhere).
Modo🧯 supports different formats to make this step easier, via the flag `--format`:

## Plain Markdown

Just plain markdown files.
This is Modo🧯's default output format.
The generated files are suitable for GitHub's Markdown rendering.

## mdBook

Markdown files as well as auxiliary files for [mdBook](https://github.com/rust-lang/mdBook),
with `--format=mdbook`.
The generated files can be used by mdBook without any further steps:

``` {class="no-wrap"}
modo build -i api.json -o docs/ --format=mdbook
mdbook serve docs-out --open
```

[Templates](../templates) can be used to customize the mdBook configuration file `book.toml`.

## Hugo

Markdown files with front matter and cross-references for [Hugo](https://gohugo.io/),
with flag `--format=hugo`.

You should first set up a Hugo project in a sub-folder of your repository.
Then, run Modo🧯 with the Hugo `content` folder as output path:

``` {class="no-wrap"}
modo build -i api.json -o <hugo-project>/content --format=hugo
```

Further, in your `hugo.toml`, add `disablePathToLower = true` to the main section
to prevent lower case members (like functions) and upper case members (like structs)
overwrite each other.
Alternatively, run Modo🧯 with switch `--case-insensitive`.

[Templates](../templates) can be used to customize the Hugo front matter of each page.

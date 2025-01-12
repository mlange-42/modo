# Modo

[Mojo](https://www.modular.com/mojo) documentation generator.

Generates Markdown for static site generators (SSGs) from `mojo doc` JSON output.

As an example, [here](https://mlange-42.github.io/modo/) is the Mojo [stdlib](https://github.com/modularml/mojo) processed with Modo and [mdBook](https://github.com/rust-lang/mdBook).

**! Early work in progress !**

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
Modo's output folder can be used to render a book instantly:

```
modo docs-out -i docs.json --mdbook
mdbook serve docs-out --open
```

---
title: Commands
type: docs
summary: Available Modo🧯 commands.
weight: 30
---
Modo🧯 uses commands to structure the command line interface.
Most commands also accept flags. These are covered in chapter [Configuration](../config).

Commands are generally used like this:

```
modo <command> [ARGS] [FLAGS]
```

## `init`

Command `init` serves to prepare an existing Mojo🔥 project for instant usage
with Modo🧯 and a static site generator (SSG).
It takes a mandatory argument for the intended [output format](../formats).
Details are covered in chapter [Project setup](../setup).

## `build`

Command `build` builds documentation from `mojo doc` JSON files
and additional Markdown files.
Optionally, it also extracts [doc-tests](..doctests).
Takes an optional path argument for the project to build.
See chapter [Configuration](../config) for details.

## `test`

Command `test` extracts [doc-tests](..doctests) `mojo doc` JSON files
and additional Markdown files.
Takes an optional path argument for the project to extract tests from.
See chapter [Configuration](../config) for details.

## `clean`

Command `clean` removes Markdown and test files created by Modo🧯.
Takes an optional path argument for the project to clean.
This is particularly useful to get rid of old artifacts
after moving, removing or renaming documentation files or API members.

---
title: Cross-references
type: docs
summary: Cross-referencing in the API docs.
prev: features
weight: 10
---

Modo🧯 supports cross-refs within the documentation of a project.
Absolute as well as relative references are supported.
Relative references follow Mojo🔥's import syntax, with a leading dot denoting the current module, and further dots navigating upwards.

Some examples:

| Ref | Explanation |
|-----|-------------|
| `[pkg.mod.A]` | Absolute reference. |
| `[.A]` | Struct `A` in the current module. |
| `[.A.method]` | Method `method` of struct `A` in the current module. |
| `[..mod.A]` | Struct `A` in sibling module `mod`. |
| `[.A.method link text]` | Method `method` of struct `A`, with custom text. |

Leading dots are stripped from the link text if no custom text is given, so `.mod.Type` becomes `mod.Type`.
With flag `--short-links`, packages and modules are also stripped, so `.mod.Type` becomes just `Type`.

Besides cross-references, normal Markdown links can be used in doc-strings.
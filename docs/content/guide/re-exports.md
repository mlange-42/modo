---
title: Re-exports
type: docs
summary: Restructure package according to re-exports.
weight: 5
---

In Mojo🔥, package-level re-exports (or rather, imports) can be used
to flatten the structure of a package and shorten import paths for users.

Modo🧯 can structure documentation output according to re-exports
by using `exports: true` in the config, or flag `--exports`.
However, as we don't look at the actual code but just `mojo doc` JSON,
these re-exports must be documented in an `Exports:` section in the package docstring.

In a package's `__init__.mojo`, document re-exports like this:

```python {class="no-wrap"}
"""
Package creatures demonstrates Modo re-exports.

Exports:
 - animals.mammals.Cat
 - animals.mammals.Dog
 - plants.vascular
 - fungi
"""
from .animals.mammals import Cat, Dog
from .plants import vascular
```

> Note that `Exports:` should not be the first line of the docstring, as it is considered the summary and is not processed.

When processed with `--exports`, only exported members are included in the documentation.
Re-exports are processed recursively.
This means that sub-packages need an `Exports:` section too if they are re-exported directly,
like `fungi` in the example.
For exporting members from a sub-package (like `Cat` and `Doc`), the sub-package `Exports:` are ignored.

Re-exported modules (like `plants.vascular`) are fully included with all members.

[Cross-references](../cross-refs) should still use the original structure of the package.
They are automatically transformed to match the altered structure.

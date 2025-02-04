---
title: Re-exports
type: docs
summary: Restructure package according to re-exports.
weight: 20
---

In Mojo🔥, package-level re-exports (or rather, imports) can be used
to flatten the structure of a package and to shorten import paths for users.

Modo🧯 can structure documentation output according to re-exports
by using `exports: true` in the config, or flag `--exports`.
However, as we don't look at the actual code but just `mojo doc` JSON,
these re-exports must be documented in an `Exports:` section in the package docstring.

In a package's `__init__.mojo`, document re-exports like this:

```mojo {class="no-wrap"}
"""
Package creatures demonstrates Modo re-exports.

Exports:
 - animals.Dog
 - animals.Cat as Kitten
 - fungi
 - plants.vascular
 - plants.bryophytes.Moss
"""
from .animals import Dog, Cat as Kitten
from .plants import vascular
from .plants.bryophytes import Moss
```

> Note that `Exports:` should not be the first line of the docstring, as it is considered the summary and is not processed.

When processed with `--exports`, only exported members are included in the documentation.

Re-exports are processed recursively.
Re-exported modules (like `plants.vascular`) are fully included with all members.
Sub-packages (like `fungi`) need an `Exports:` section too if they are re-exported as a whole.
When directly exporting members from a sub-package (like `plants.bryophytes.Moss`), the sub-package exports are ignored.

Renaming (`X as Y`) can be used like in Mojo🔥 imports.
Member names, link paths and cross-reference labels are changed accordingly.
Signatures, however, are not changed.

[Cross-references](../crossrefs) should still use the original structure of the package.
They are automatically transformed to match the altered structure.
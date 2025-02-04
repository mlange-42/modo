"""
Package doctests tests doctests

```mojo {doctest="test" global=true hide=true}
from testing import *

fn add(a: Int, b: Int) -> Int:
    return a + b
```

```mojo {doctest="test" hide=true}
var a = 1
var b = 2
```

```mojo {doctest="test"}
var c = add(a, b)
```

```mojo {doctest="test" hide=true}
assert_equal(c, 3)
```
"""
alias ModuleAlias2 = Int


struct Struct:
    alias StructAlias = Int

    fn func(self):
        """
        Doctests in a struct member.

        ```mojo {doctest="func"}
        var a = 1
        ```
        """
        pass

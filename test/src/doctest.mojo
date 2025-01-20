"""
Package doctests tests doctests

```mojo {doctest="test-global"}
from testing import *

struct Test:
    fn __init__(out self):
        pass
```

```mojo {doctest="test-setup"}
alias T = Test
```

```mojo {doctest="test-setup"}
var b = 1
```

```mojo {doctest="test"}
var a = 1
```

```mojo {doctest="test-teardown"}
assert_equal(a, b)
```
"""


struct Struct:
    fn func(self):
        """
        Doctests in a struct member.

        ```mojo {doctest="func"}
        var a = 1
        ```
        """
        pass

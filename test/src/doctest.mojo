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
b = T()
```

```mojo {doctest="test"}
a = T()
```

```mojo {doctest="test-teardown"}
assert_equal(a, b)
```
"""

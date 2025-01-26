Additional Markdown file. 

```mojo {doctest="markdown" global=true hide=true}
from testing import *

fn add(a: Int, b: Int) -> Int:
    return a + b
```

```mojo {doctest="markdown" hide=true}
var a = 1
var b = 2
```

```mojo {doctest="markdown"}
var c = add(a, b)
```

```mojo {doctest="markdown" hide=true}
assert_equal(c, 3)
```
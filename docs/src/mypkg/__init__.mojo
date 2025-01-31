"""
Package `mypkg` provides an example for Modo🧯 API docs.

The package has some re-exported dummy members (optional) as well as a [.subpkg].
The docstring for the package looks like this:

```mojo
\"\"\"
Package `example` provides an example for Modo🧯 API docs.

The package has some re-exported dummy members (optional) as well as a [.subpkg].
The source code of the package looks like this:

<code block excluded>

Exports:
 - mymodule.MyPair
 - mymodule.MyIntAlias
 - subpkg
\"\"\"
from .mymodule import MyPair, MyIntAlias
```

Exports:
 - mymodule.MyPair
 - mymodule.MyIntAlias
 - subpkg
"""
from .mymodule import MyPair, MyIntAlias

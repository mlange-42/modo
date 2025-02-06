---
title: Inheritance
type: docs
summary: Inherit method docs from parent traits.
weight: 40
---

When implementing trait methods, documentation of the method
in the implementing struct is often redundent.
For these cases, Modo🧯 offers doc inheritance.

In the method docstring, reference the parent trait in double square brackets,
somewhere in the summary (i.e. in the first sentence).
The syntax is the same as for [cross-references](../crossrefs).
See line 14 below:

```mojo {doctest="inherit" global=true class="no-wrap" linenos=true}
trait MyTrait:
    fn do_something(self) -> Int:
        """
        A method that returns an integer.

        Returns:
            An arbitrary integer.
        """
        ...

@value
struct MyStruct(MyTrait):
    fn do_something(self) -> Int:
        """See [[.MyTrait]]."""
        return 100
```

```mojo {doctest="inherit" hide=true}
var s = MyStruct()
var result = s.do_something()
if result != 100:
    raise Error("failed")
```

This copies the entire documentation from the trait's method with the same signature into the struct's method. This includes docs for arguments, parameters, return and raises.
Any original docs of the struct's method are replaced completely.

Inheriting documentation is only possible from traits in the same root package.

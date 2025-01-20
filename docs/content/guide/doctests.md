---
title: Doc testing
type: docs
summary: Extract doc tests from code examples in the API docs.
weight: 5
---

To keep code examples in docstrings up to date, Modo🧯 can generate test files for `mojo test` from them.
Doctests are enabled by flag `--doctest`, which takes an output directory for test files as an argument:

Code block attributes are used to identify code blocks to be tested.
Any block that should be included in the tests needs a name:

````markdown
```mojo {doctest="mytest"}
var a = 0
```
````

Multiple code blocks with the same name are concatenated.
Individual blocks can be hidden with an attribute `hide=true`:

````markdown
```mojo {doctest="mytest" hide=true}
# hidden code block
```
````

Further, for code examples that can't be put into a test function, attribute `global=true` can be used:

````markdown
```mojo {doctest="mytest" global=true}
struct MyStruct:
    pass
```
````

Combining multiple code blocks using these attributes allows for flexible tests with hidden setup, teardown and assertions.

CLI usage example:

```
mojo doc src/ -o docs.json                 # generate doc JSON
modo docs -i docs.json --doctest=doctest   # render to Markdown and extract doctests
mojo test -I src doctest                   # run the doctests
```

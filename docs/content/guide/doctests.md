---
title: Doc testing
type: docs
summary: Extract doc tests from code examples in the API docs.
weight: 6
---

To keep code examples in docstrings up to date, Modo🧯 can generate test files for `mojo test` from them.
Doctests are enabled by `tests` in the `modo.yaml` or flag `--tests`, which take an output directory for test files as an argument:

```{class="no-wrap"}
modo build --tests doctest/       # render to Markdown and extract doctests
mojo test -I . doctest/           # run the doctests
```

Alternatively, Modo🧯's `test` command can be used to extract tests without building the Markdown docs:

```{class="no-wrap"}
modo test --tests doctest/        # only extract doctests
```

In both cases, flag `--tests` can be omitted if `tests: doctest/` is set in the `modo.yaml` file.

## Tested blocks

Code block attributes are used to identify code blocks to be tested.
Any block that should be included in the tests needs a name:

````{class="no-wrap"}
```mojo {doctest="mytest"}
var a = 0
```
````

Multiple code blocks with the same name are concatenated.

## Hidden blocks

Individual blocks can be hidden with an attribute `hide=true`:

````{class="no-wrap"}
```mojo {doctest="mytest" hide=true}
# hidden code block
```
````

## Global blocks

Further, for code examples that can't be put into a test function, attribute `global=true` can be used:

````{class="no-wrap"}
```mojo {doctest="mytest" global=true}
struct MyStruct:
    pass
```
````

## Full example

Combining multiple code blocks using these attributes allows for flexible tests with hidden setup, teardown and assertions.
Here is a full example:

````mojo {class="no-wrap"}
fn add(a: Int, b: Int) -> Int:
    """
    Function `add` sums up its arguments.

    ```mojo {doctest="add" global=true hide=true}
    from testing import assert_equal
    from mypkg import add
    ```

    ```mojo {doctest="add"}
    var result = add(1, 2)
    ```
    
    ```mojo {doctest="add" hide=true}
    assert_equal(result, 3)
    ```
    """
    return a + b
````

This generates the following docs content:

----
Function `add` sums up its arguments.

```mojo {doctest="add"}
var result = add(1, 2)
```
----

Further, Modo🧯 creates a test file with this content:

```mojo
from testing import assert_equal
from mypkg import add

fn test_add() raises:
    result = add(1, 2)
    assert_equal(result, 3)
```

## Markdown files

A completely valid Modo🧯 use case is a site with not just API docs, but also other documentation.
Thus, code examples in Markdown files that are not produced by Modo🧯 can also be processed for doctests.

For that sake, Modo🧯 can use an entire directory as input, instead of one or more JSON files.
The input directory should be structured like the intended output, with API docs folders replaced by `mojo doc` JSON files.
Here is an example for a Hugo site with a user guide and API docs for `mypkg`:

{{< filetree/container >}}
  {{< filetree/folder name="docs-in" >}}
    {{< filetree/folder name="guide" >}}
      {{< filetree/file name="_index.md" >}}
      {{< filetree/file name="installation.md" >}}
      {{< filetree/file name="usage.md" >}}
    {{< /filetree/folder >}}
    {{< filetree/file name="_index.md" >}}
    {{< filetree/file name="mypkg.json" >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

With a directory as input, Modo🧯 does the following:

- For each JSON file (`.json`), generate API docs, extract doctests, and write Markdown to the output folder and tests to the tests folder.
- For each Markdown (`.md`) file, extract doctests, and write processed Markdown to the output folder and tests to the tests folder.
- For any other files, copy them to the output folder.

Note that this feature is not available with the [mdBook](../formats#mdbook) format.

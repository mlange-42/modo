---
title: Doc-tests
type: docs
summary: Extract doc tests from code examples in the API docs.
weight: 30
---

To keep code examples in docstrings up to date, Modo🧯 can generate test files for `mojo test` from them.
Doctests are enabled by `tests` in the `modo.yaml` or flag `--tests`. Doctests are enabled by default.
Further, the default setup also contains a post-processing [script](../scripts) that runs `mojo test`.

Alternatively to `modo build`, Modo🧯's `test` command can be used to extract tests without building the Markdown docs:

```shell {class="no-wrap"}
modo test           # only extract doctests
```

## Tested blocks

Code block attributes are used to identify code blocks to be tested.
Any block that should be included in the tests needs a `doctest` identifier:

````{class="no-wrap"}
```mojo {doctest="mytest"}
var a = 0
```
````

Multiple code blocks with the same identifier are concatenated.

## Hidden blocks

Individual blocks can be hidden with `hide=true`:

````{class="no-wrap"}
```mojo {doctest="mytest" hide=true}
# hidden code block
```
````

## Global blocks

Further, for code examples that can't be put into a test function, `global=true` can be used:

````{class="no-wrap"}
```mojo {doctest="mytest" global=true}
struct MyStruct:
    pass
```
````

## Full example

Combining multiple code blocks using these attributes allows for flexible tests with imports, hidden setup, teardown and assertions.
Here is a full example:

````mojo {doctest="add" global=true class="no-wrap"}
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

{{<html>}}<div style="border: 2px solid grey; padding: 1rem; margin: 1rem 0;">{{</html>}}

Function `add` sums up its arguments.

```mojo
var result = add(1, 2)
```

{{<html>}}</div>{{</html>}}

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
  {{< filetree/folder name="docs" >}}
    {{< filetree/folder name="src" >}}
      {{< filetree/folder name="guide" state="closed" >}}
        {{< filetree/file name="_index.md" >}}
        {{< filetree/file name="installation.md" >}}
        {{< filetree/file name="usage.md" >}}
      {{< /filetree/folder >}}
      {{< filetree/file name="_index.md" >}}
      {{< filetree/file name="mypkg.json" >}}
    {{< /filetree/folder >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

With a directory as input, Modo🧯 does the following:

- For each JSON file (`.json`), generate API docs, extract doctests, and write Markdown to the output folder and tests to the tests folder.
- For each Markdown (`.md`) file, extract doctests, and write processed Markdown to the output folder and tests to the tests folder.
- For any other files, copy them to the output folder.

Note that this feature is not available with the [mdBook](../../formats#mdbook) format.
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

```mojo {doctest="add" global=true}
from testing import assert_equal

fn test_add() raises:
    var result = add(1, 2)
    assert_equal(result, 3)
```

Note that the two code blocks around the box form a doc-test themselves,
to ensure that this guide is correct and up to date.
Both blocks have the attributes `{doctest="add" global=true}`,
which concatenates them into one test file.

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

## Modo🧯 vs. `mojo test`

Mojo🔥 can also test code examples directly,
see the manual's section on [testing](https://docs.modular.com/mojo/tools/testing/#writing-api-documentation-tests).
So why would you want to use Modo🧯's doc-testing feature?

First and foremost, `mojo test` will only test code blocks in docstrings,
while Modo🧯 also allows for doc-testing of general [Markdown files](#markdown-files).
Further, Modo🧯 is probably a bit more flexible and offers more control than `mojo test`.
E.g., multiple tests suites are possible per docstring, and code blocks can be hidden completely
where `mojo test` produces empty code blocks.

Note that `mojo test` does not recognize docstrings with attributes (```` ```mojo {...} ````).
This means that code blocks prepared for Modo🧯's doc-testing are not tested by `mojo test`.
But it also means that you can't use `mojo test` if you want to add attributes for your SSG,
e.g. for line numbers, line highlights, or other features.

The doc-testing features of Modo🧯 and `mojo test` can be used in the same project without conflicts.

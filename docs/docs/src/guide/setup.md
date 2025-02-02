---
title: Setup
type: docs
summary: Setting up a Mojo🔥 project for Modo🧯.
weight: 20
---

The [command](../commands) `init` serves to prepare an existing Mojo🔥 project for instant usage
with Modo🧯 and a static site generator (SSG).

## Hugo example

As an example, we use [Hugo](https://gohugo.io) as SSG.
For all supported options, see chapter [formats](../formats).
Navigate into your Mojo🔥 project's root folder and run:

```shell {class="no-wrap"}
modo init hugo
```

Modo🧯 analyzes the structure of your project and tries to find Mojo🔥 packages.
It then sets up a `modo.yaml` [config file](../config) and a directory `docs`, containing a minimal Hugo project as well as sub-directories for auxiliary documentation files and extracted [doc-tests](../features/doctests).

After that, you should be able to instantly generate your API docs with Modo🧯
and render them with Hugo:

```shell {class="no-wrap"}
modo build
hugo serve -s docs/site/
```

If your project has a GitHub repository, Modo🧯 will set up the project so
that it can be deployed to GitHub Pages instantly.

For more details, see the [Hugo format](../formats#hugo).

## mdBook example

Similarly, with [mdBook](https://github.com/rust-lang/mdBook) as SSG, these three commands should be sufficient to view your API docs in a web browser:

```shell {class="no-wrap"}
modo init mdbook
modo build
mdbook serve docs/
```

For more details on the generated directory structure and files, see chapter [formats](../formats).

## Detected packages

These are the possible project layouts the `init` command can work with:

{{<html>}}<div style="display: flex;"><div style="flex: 25%;">{{</html>}}

{{< filetree/container >}}
  {{< filetree/folder name="root" >}}
    {{< filetree/folder name="src" state="closed" >}}
      {{< filetree/file name="`__init__.mojo`" >}}
    {{< /filetree/folder >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

{{<html>}}</div><div style="flex: 25%;">{{</html>}}

{{< filetree/container >}}
  {{< filetree/folder name="root" >}}
    {{< filetree/folder name="pkg_a" >}}
      {{< filetree/folder name="src" state="closed" >}}
        {{< filetree/file name="`__init__.mojo`" >}}
      {{< /filetree/folder >}}
    {{< /filetree/folder >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

{{<html>}}</div><div style="flex: 25%;">{{</html>}}

{{< filetree/container >}}
  {{< filetree/folder name="root" >}}
    {{< filetree/folder name="pkg_a" state="closed" >}}
      {{< filetree/file name="`__init__.mojo`" >}}
    {{< /filetree/folder >}}
    {{< filetree/folder name="pkg_b" state="closed" >}}
      {{< filetree/file name="`__init__.mojo`" >}}
    {{< /filetree/folder >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

{{<html>}}</div><div style="flex: 25%;">{{</html>}}

{{< filetree/container >}}
  {{< filetree/folder name="root" >}}
    {{< filetree/folder name="src" >}}
      {{< filetree/folder name="pkg_a" state="closed" >}}
        {{< filetree/file name="`__init__.mojo`" >}}
      {{< /filetree/folder >}}
      {{< filetree/folder name="pkg_b" state="closed" >}}
        {{< filetree/file name="`__init__.mojo`" >}}
      {{< /filetree/folder >}}
    {{< /filetree/folder >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

{{<html>}}</div></div>{{</html>}}

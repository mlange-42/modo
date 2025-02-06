---
title: Templates
type: docs
summary: Use templates to customize Modo🧯's output.
next: mypkg
weight: 60
---

Modo🧯 relies heavily on templating.
With the option `templates` in the `modo.yaml` or flag `--templates`, custom template folders can be specified to (partially) overwrite the embedded templates.
Simply use the same files names, and alter the content.
Embedded templates that can be overwritten can be found in folder [assets/templates](https://github.com/mlange-42/modo/tree/main/assets/templates).

Modo🧯 uses the templating syntax of [Go](https://go.dev) and [Hugo](https://gohugo.io).
See Go's [`text/template`](https://pkg.go.dev/text/template) docs and
Hugo's [Introduction to templating](https://gohugo.io/templates/introduction/)
for details.

As an example, here is the builtin template for rendering a trait.
It uses nested templates for clarity and modularity.

```md
Mojo trait

# `{{.Name}}`

{{template "summary" . -}}
{{template "description" . -}}
{{template "fields" . -}}
{{template "parent_traits" . -}}
{{template "methods" . -}}
```

Besides changing the page layout and content, templates can also be used to alter the [Hugo](../../formats#hugo) front matter of individual pages, e.g. to change the document type or to add more information for Hugo.

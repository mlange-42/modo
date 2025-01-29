---
title: Templates
type: docs
summary: Use templates to customize Modo🧯's output.
next: mypkg
weight: 50
---

Modo🧯 relies heavily on templating.
With the option `templates` in the `modo.yaml` or flag `--templates`, custom template folders can be specified to (partially) overwrite the embedded templates.
Simply use the same files names, and alter the content.
Embedded templates that can be overwritten can be found in folder [assets/templates](https://github.com/mlange-42/modo/tree/main/assets/templates).

Besides changing the page layout and content, this feature can also be used to alter the [Hugo](../../formats#hugo) front matter of individual pages.

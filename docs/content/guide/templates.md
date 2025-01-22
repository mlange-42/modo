---
title: Templates
type: docs
summary: Use templates to customize Modo🧯's output.
weight: 7
---

Modo🧯 relies heavily on templating.
With the option `templates` in the `modo.yaml` or flag `--templates`, custom template folders can be specified to (partially) overwrite the embedded templates.
Simply use the same files names, and alter the content.
Embedded templates that can be overwritten can be found in folder [assets/templates](https://github.com/mlange-42/modo/tree/main/assets/templates).

Besides changing the page layout, this feature can be used to alter the [Hugo](../formats#hugo) front matter, or to adapt the [mdBook](../formats#mdbook) configuration file.

---
title: Bash commands
type: docs
summary: Configure bash commands to before and after processing.
weight: 7
---

Modo🧯 can be configured to automatically run bash commands before and/or after processing.

This feature can be used to run all necessary steps with a single `modo build` or `modo test` command.
Particularly, `mojo doc` can be executed before processing, and `mojo test` after extracting [doc-tests](../doctests).

## Configuration

The `modo.yaml` [config file](../config) provides the following fields for bash commands:

- `pre-run`: runs before `build` as well as `test`.
- `pre-build`: runs before `build`.
- `pre-test`: runs before `test`. Also runs before build if `tests` is given.
- `post-test`: runs after `test`. Also runs after build if `tests` is given.
- `post-build`: runs after `build`.
- `post-run`: runs after `build` as well as `test`.

Each of those takes an array of bash commands.
Each bash command can be comprised of multiple lines.

Here is an example that runs `mojo doc` before builds and tests:

```yaml
pre-run:
  - mojo doc -o api.json src/
```

And here is how to run `mojo test` after doc-tests extraction:

```yaml
post-test:
  - |
    echo Running 'mojo test'...
    mojo test -I . doctest/
    echo Done.
```

Here, we use a single command that consists of 3 lines.

## Error trap

Each (potentially multi-line) command starts a new bash process.
Each process is initialized with an error trap via `set -e`.
This means that any failing line (or sub-command) causes the command to fail with that error.

To let sub-command errors pass, use `set +e` as the first line of your command.

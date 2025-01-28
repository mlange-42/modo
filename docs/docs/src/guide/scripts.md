---
title: Bash scripts
type: docs
summary: Configure bash scripts to run before and after processing.
weight: 7
---

Modo🧯 can be configured to automatically run bash scripts before and/or after processing.

This feature can be used to run all necessary steps with a single `modo build` or `modo test` command.
Particularly, `mojo doc` can be executed before processing, and `mojo test` after extracting [doc-tests](../doctests).

## Configuration

The `modo.yaml` [config file](../config) provides the following fields for bash scripts:

- `pre-run`: runs before `build` as well as `test`.
- `pre-build`: runs before `build`.
- `pre-test`: runs before `test`. Also runs before build if `tests` is given.
- `post-test`: runs after `test`. Also runs after build if `tests` is given.
- `post-build`: runs after `build`.
- `post-run`: runs after `build` as well as `test`.

Each of those takes an array of bash scripts.
Each bash script can be comprised of multiple commands.

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

Here, we use a single script that consists of 3 commands.

## Skipping scripts

Using the flag `--bare` (`-B`), shell commands can be skipped
so that only the Modo🧯 command is executed.
This can be useful to skip scripts that are intended for the CI
when working locally.

## Error trap

Each script starts a new bash process.
Each process is initialized with an error trap via `set -e`.
This means that any failing command causes the script to fail with that error.

To let errors of individual commands pass, use `set +e` as the first line of your script.

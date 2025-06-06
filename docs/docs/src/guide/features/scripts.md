---
title: Bash scripts
type: docs
summary: Configure bash scripts to run before and after processing.
weight: 50
---

Modo🧯 can be configured to automatically run bash scripts before and/or after processing.

This feature can be used to run all necessary steps with a single `modo build` or `modo test` command.
Particularly, `mojo doc` can be executed before processing, and `mojo test` after extracting [doc-tests](../doctests).

## Configuration

The `modo.yaml` [config file](../../config) provides the following fields for bash scripts:

- `pre-run`: runs before `build` as well as `test`.
- `pre-build`: runs before `build`.
- `pre-test`: runs before `test`. Also runs before build if `tests` is given.
- `post-test`: runs after `test`. Also runs after build if `tests` is given.
- `post-build`: runs after `build`.
- `post-run`: runs after `build` as well as `test`.

Each of those takes an array of bash scripts.
Each bash script can be comprised of multiple commands.

Here is an example that shows how `mojo doc` is executed in the default setup before build and test:

```yaml
# Bash scripts to run before build as well as test.
pre-run:
  - |
    echo Running 'mojo doc'...
    pixi run mojo doc -o docs/src/mypkg.json src/mypkg
    echo Done.
```

And here is how `mojo test` is executed in the default setup after doc-tests extraction:

```yaml
# Bash scripts to run after test.
# Also runs after build if 'tests' is given.
post-test:
  - |
    echo Running 'mojo test'...
    pixi run mojo test -I src docs/test
    echo Done.
```

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
## [[v0.8.1]](https://github.com/mlange-42/modo/compare/v0.8.0...v0.8.1)

Modo🧯 is now installable with `pip`, as `pymodo` (#102, #104)

## [[v0.8.0]](https://github.com/mlange-42/modo/compare/v0.7.0...v0.8.0)

### Breaking changes

* Rework of the command line interface to use sub-commands `init`, `build` and `test` (#91)

### Features

* Optionally extracts doc-tests to run with `mojo test` (#81, #88)
* CLI flags can be stored in a `modo.yaml` config file (#90, #94)
* Shell commands to run before and after docs processing can be given in `modo.yaml` (#90)

### Documentation

* Moves most of the documentation to the user guide on GitHub Pages (#82, #83, #84)

### Performance

* Speedup by avoiding repeated RegExp compilation (#89)

## [[v0.7.0]](https://github.com/mlange-42/modo/compare/v0.6.0...v0.7.0)

### Features

* Adds CLI flag `--strict` to break with an error instead of warnings (#68)
* Adds CLI flag `--dry-run` to run without file output (#71)
* Allows to re-export aliases on package level (#73)
* Allows for cross-refs to aliases in packages, modules and structs (#69, #73)

### Documentation

* Replaces the mdBook stdlib example by a custom one using Hugo (#64)

### Bugfixes

* Don't trim docstrings during extraction and removal of exports section (#60)
* Fixes broken links when referencing parents or the root package (#62)

### Other

* Checks package re-exports for correct paths (#68)
* Checks package re-exports for resulting name conflicts (#72)
* More consistent line breaks in Markdown output (#75)

## [[v0.6.0]](https://github.com/mlange-42/modo/compare/v0.5.0...v0.6.0)

### Features

* Implements output customization by (partial) template overwrites (#48)
* Adds support for YAML input, particularly for shorter test files (#49)

### Bugfixes

* Fixes overwritten cross-refs in case of multiple re-exports of the same member (#46)
* Performs link replacement as the very last rendering step, allowing for refs in summaries (#47)

## [[v0.5.0]](https://github.com/mlange-42/modo/compare/v0.4.0...v0.5.0)

### Features

* Adds CLI flag `--short-links` to strip packages and modules from link labels (#41)
* Supports re-structuring doc output according to package re-exports (#42)

## [[v0.4.0]](https://github.com/mlange-42/modo/compare/v0.3.0...v0.4.0)

### Breaking changes

* Replaces format flags `--mdbook` and `--hugo` by `--format=(plain|mdbook|hugo)` (#37)

### Features

* Adds support to export for [Hugo](https://gohugo.io/) (#36)

### Other

* Moves the binary package to the module root, simplifying the install path (#39)

## [[v0.3.0]](https://github.com/mlange-42/modo/compare/v0.2.0...v0.3.0)

### Features

* Adds support for cross-references in docstrings (#28, #30)

### Formats

* Adds CSS to mdBook output to enable text wrapping in code blocks (#33)

## [[v0.2.0]](https://github.com/mlange-42/modo/compare/v0.1.1...v0.2.0)

### Features

* Adds CLI flag `--case-insensitive` to append hyphen `-` at the end of capitalized file names, as fix for case-insensitive systems (#20, #21)
* Uses templates to generate package, module and member paths (#22, #23)

### Formats

* Removes numbering from navigation entries (#16)
* Navigation, top-level headings and method headings use inline code style (#18, #19)

### Bugfixes

* Generates struct signatures if not present due to seemingly `modo doc` bug (#20)

### Other

* Simplifies templates to use `.Name` instead of `.GetName` (#24)

## [[v0.1.1]](https://github.com/mlange-42/modo/compare/v0.1.0...v0.1.1)

### Documentation

* Adds a CHANGELOG.md file (#14)

### Other

* Re-release due to pkg.go.dev error (#14)

## [[v0.1.0]](https://github.com/mlange-42/modo/tree/v0.1.0)

First minimal usable release of Modo, a Mojo documentation generator.

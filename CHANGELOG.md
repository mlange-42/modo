## [[unpublished]](https://github.com/mlange-42/modo/compare/v0.2.0...main)

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

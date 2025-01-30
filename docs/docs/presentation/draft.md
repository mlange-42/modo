# Modo presentation

Thanks for the opportunity...

This is not a Mojo project, it is a project for Mojo projects.
Modo is a DocGen...

## Why I built/build Modo

- No standard tool for API docs so far, mostly in-house/project solutions (partially published)
- Need API docs for (probably) first Mojo ECS by colleague: Larecs
- Want simple, low-tech, generic solution without "stack dependency" (however this is called professionally)

## What it does

- From mojo doc JSON, creates a hierarchy of Markdown files suitable for SSGs

Flow chart, animated:

```
                                     actual Modo
                                   _______/\_______
                                  /                \

sources --> mojo doc --> JSON --> Modo ---> Markdown --> SSG --> HTML
                                  ^    \
                     other docs __|     '-> Tests   --> mojo test

            \__________________________  __________________________/
                                       \/
                       (potentially) one command: `modo build`
```

- Several additional features, compatible with known Mojo (docstring) conventions.

## Demo

1. Clone a 3rd party a Mojo package.
    - ~~stdlib~~ -- takes too long to clone.
    - ~~lightbug_http~~ -- has package `testutils` that does not work out of the box.
1. Delete folder `docs` if present.
1. Run `init`, `build` and `hugo serve`.
1. Show docs in browser.
1. Show structure of `docs` folder?

## Features beyond that

At this point, we get what the Mojo stdlib API docs show.
But there is more...

### Cross-refs

- Very simple syntax, equivalent to Python/Mojo package imports
  - Show examples
- Transformed to Markdown links or SSG-specific refs
- No need for changes when re-exports change

### Re-exports

- Structure the docs as the user sees your package
- Must be represented in docs, by `Exports:` section
- Again, very simple syntax, equivalent to Python/Mojo package imports
  - Show example

### Doc-tests

- Extract unit test files from doc comments
- Highly flexible with hidden and global blocks
- Uses standard block attributes, alias "info string"
  - Show example

### Scripts

- Pre- and post-processing can be specified in config file
- All (incl. external) steps in a single command
  - Show default scripts

### Templates

- Heavily based on Go's templating
- Highly customizable via template overwrites
- Same syntax (and Go library) as Hugo
  - Show example (original) template? Just to show how simple they are...

### 2nd demo?

- Maybe show Modo API example for features?
- Or put screenshots into feature sections? Probably better.

## How to get Modo

- `pip install pymodo`
- `go install github.com/mlange-42/modo`
- `magic instal ...` (not yet)
- Download binaries from releases

## Performance

"fast enough", no optimizations done so far

- stdlib: 1.5sec
- small package: <50ms

## Outlook

- Feature-compete regarding initial plan
- A few potential near-term features: commands `clean` and `watch`
- Set up for Community Packages (no patience for this ATM)
- Open for contributions!

## @modular

Messages/request @Modular:

- Specify cross-ref syntax
- Include package re-exports into JSON
- Inconsistency: `out` arg vs. return docs?
- Bug: currently no signature for structs in JSON
- Allow lists in `Raises:`

## Contribute

- Feedback on tool and docs
- "Playtest"
- Make issues & PRs

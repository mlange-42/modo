---
title: Modo🧯 -- DocGen for Mojo🔥
theme: night
css:
  - https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/font-awesome.min.css
scripts:
  - https://cdn.jsdelivr.net/npm/mermaid@11.4.1/dist/mermaid.min.js
  - lib/reveal-mermaid.js
  - https://kit.fontawesome.com/f4816f3363.js
revealOptions:
  transition: 'convex'
  controls: true
  progress: false
  history: true
  center: true
  slide-number: false
mermaid:
  look: handDrawn
  theme: dark
  themeVariables:
    fontSize: 24px
---
<style>
.reveal {
  font-size: 36px;
}
p code, li code {
  padding-left: 0.5rem;
  padding-right: 0.5rem;
  background: #303030;
  border-radius: 0.2em;
}
.reveal .code-wrapper code:not(.mermaid) {
	white-space: preserve;
  font-size: 120%;
  background: #303030;
  border-radius: 0.33em;
}
.reveal .code-wrapper code .nowrap {
  text-wrap: nowrap;
}

.reveal .slides section .fragment.step-fade-in-then-out {
	opacity: 0;
	display: none;
}
.reveal .slides section .fragment.step-fade-in-then-out.current-fragment {
	opacity: 1;
	display: inline;
}
.columns {
  display: flex;
}
.col {
  flex: 1;
  text-align: left;
  font-size: 90%;
}
</style>

# Modo🧯

<br />

DocGen for Mojo🔥

[<i class="fa fa-github"></i>](https://github.com/mlange-42/modo)

---

## What is Modo🧯?

### <big>&darr;</big>

----

This is not a Mojo🔥 project!

It is a project for Mojo🔥 projects.
<!-- .element: class="fragment" data-fragment-index="1" -->

<br />

A DocGen for Mojo🔥, written in <i class="fa-brands fa-golang" style="font-size: 200%; position: relative; top: 12px; color: #00ADD8;"></i>
<!-- .element: class="fragment" data-fragment-index="2" -->

---

<!-- .slide: data-visibility="hidden" -->

## Why I built Modo🧯

### <big>&darr;</big>

----

<!-- .slide: data-visibility="hidden" -->

No standard tool for API docs so far

Need API docs for first(?) Mojo🔥 ECS: [Larecs](https://github.com/samufi/larecs)
<!-- .element: class="fragment" data-fragment-index="1" -->

Want simple, low-tech, generic solution
<!-- .element: class="fragment" data-fragment-index="2" -->

---

## What it does

### <big>&darr;</big>

----

From  `mojo doc`  JSON...
- creates Markdown files suitable for SSGs
<!-- .element: class="fragment" data-fragment-index="1" -->
- converts code examples to unit tests
<!-- .element: class="fragment" data-fragment-index="2" -->

<br />
<br />

```mermaid
graph LR
  sources[(Sources)]
  mojo_doc[mojo doc]
  JSON[(JSON)]
  Modo[Modo🧯]
  Markdown[(Markdown)]
  Tests[(Tests)]
  mojo_test[mojo test]
  HTML[(HTML)]
  SSG["`SSG
(e.g. Hugo)`"]

  sources-->mojo_doc
  subgraph cmd [modo build]
    mojo_doc-->JSON

    JSON-->Modo
    Modo-->Markdown
    Modo-->Tests

    Tests-->mojo_test
    Markdown-->SSG
  end
  SSG-->HTML
```

---

## Demo

---

## Features

### <big>&darr;</big>

----

### Cross-references

Very simple syntax, resembling Mojo imports
<!-- .element: class="fragment" data-fragment-index="1" -->

<div><div class="columns"><div class="col">

```md
Relative ref to [.Struct.method] in the current module.
```

</div><div class="col">

Relative ref to [Struct.method]() in the current module.

</div></div></div>
<!-- .element: class="fragment" data-fragment-index="2" -->
<div><div class="columns"><div class="col">

```md
Absolute ref to module [pkg.mod].
```

</div><div class="col">

Absolute ref to module [mod]().

</div></div></div>
<!-- .element: class="fragment" data-fragment-index="3" -->
<div><div class="columns"><div class="col">

```md
Ref with [pkg.mod custom text].
```

</div><div class="col">

Ref with [custom text]().

</div></div></div>
<!-- .element: class="fragment" data-fragment-index="4" -->

----

### Re-exports

<div class="columns" style="align-items: center; justify-content: center;"><div class="col">

<pre style="width:100%; font-size: 0.65em;">
- pkg
  - mod
    - Struct
  - subpkg
    - submod
      - Trait
</pre>

</div><!-- .element: class="fragment" data-fragment-index="1" -->
<div class="col" style="flex:0.2;">

#### <i class="fa-solid fa-arrow-right"></i>

</div><!-- .element: class="fragment" data-fragment-index="2" -->
<div class="col" style="flex:2.0">

```python
"""
Package mypkg...

Exports:
 - mod.Struct
 - subpkg.submod.Trait
"""
from .mod import Struct
from .subpkg.submod import Trait
```

</div><!-- .element: class="fragment" data-fragment-index="2" -->
<div class="col" style="flex:0.2">

#### <i class="fa-solid fa-arrow-right"></i>

</div><!-- .element: class="fragment" data-fragment-index="3" -->
<div class="col">

Modo🧯

<pre style="width:100%; font-size: 0.65em;">
- pkg
  - Struct
  - Trait
</pre>

</div><!-- .element: class="fragment" data-fragment-index="3" -->
</div>

----

### Doc-tests

----

### Scripts

----

### Templates

---

## How to get Modo🧯

### <big>&darr;</big>

----

#### Python/pip

`pip install pymodo`

<br/>

#### Go
<!-- .element: class="fragment" data-fragment-index="2" -->
`go install github.com/mlange-42/modo`
<!-- .element: class="fragment" data-fragment-index="2" -->
<br/>

#### Pre-compiled binaries
<!-- .element: class="fragment" data-fragment-index="3" -->
GitHub Releases
<!-- .element: class="fragment" data-fragment-index="3" -->

---

## @Modular

### <big>&darr;</big>

----

- Specify cross-ref syntax <!-- .element: class="fragment" data-fragment-index="1" -->
- Include package re-exports into JSON <!-- .element: class="fragment" data-fragment-index="2" -->
- Bug: currently no signature for structs in JSON <!-- .element: class="fragment" data-fragment-index="3" -->
- Allow lists in <!-- .element: class="fragment" data-fragment-index="4" -->`Raises:` <!-- .element: class="fragment" data-fragment-index="4" -->

---

## Contributing

### <big>&darr;</big>

---

## Thank you!

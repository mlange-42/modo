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
  history: true,
  center: true
  slide-number: false
mermaid:
  theme: dark
  look: handDrawn
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
.reveal .slides section .fragment.step-fade-in-then-out {
	opacity: 0;
	display: none;
}
.reveal .slides section .fragment.step-fade-in-then-out.current-fragment {
	opacity: 1;
	display: inline;
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

From  `mojo doc`  JSON:
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
<!-- .element: class="fragment" data-fragment-index="3" -->

---

## Demo

---

## Features

### <big>&darr;</big>

----

#### Cross-references

----

#### Re-exports

----

#### Doc-tests

----

#### Scripts

----

#### Templates

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

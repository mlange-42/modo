---
title: Modo🧯 -- DocGen for Mojo🔥
theme: night
scripts:
  - https://kit.fontawesome.com/f4816f3363.js
  - https://cdn.jsdelivr.net/npm/mermaid@11.4.1/dist/mermaid.min.js
  - https://mlange-42.github.io/reveal-svg-smil/_assets/reveal-svg-smil.js
  #- lib/reveal-mermaid.js
revealOptions:
  transition: 'convex'
  controls: true
  progress: false
  history: true
  center: true
  slide-number: false
  width: 1024
  height: 700
mermaid:
  look: handDrawn
  theme: dark
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
.reveal .code-wrapper {
  width: 100%;
  margin-left: 0;
  margin-right: 0;
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
.reveal .fragment {
	transition: all 1s ease;
}
code.mermaid {
  text-align: center;
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

### 🔥

---

## What is Modo🧯?

### <big>&darr;</big>

----

Modo🧯 is not a Mojo🔥 project!

It is a project for Mojo🔥 projects.
<!-- .element: class="fragment" data-fragment-index="1" -->

<br />

Modo🧯 is a DocGen for Mojo🔥, written in <i class="fa-brands fa-golang" style="font-size: 200%; position: relative; top: 12px; color: #00ADD8;"></i>
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

Using `mojo doc` JSON...

- creates Markdown files suitable for SSGs
<!-- .element: class="fragment" data-fragment-index="4" -->
- converts code examples to unit tests
<!-- .element: class="fragment" data-fragment-index="7" -->
- can run the whole chain with a single command
<!-- .element: class="fragment" data-fragment-index="9" -->

<br />
<br />

<object type="image/svg+xml" data="img/flowchart.svg">
    <img src="img/flowchart.svg" />
</object>

<span class="fragment" data-fragment-index="3" data-svg-classes="show-modo"></span>
<span class="fragment" data-fragment-index="4" data-svg-classes="show-markdown"></span>
<span class="fragment" data-fragment-index="5" data-svg-classes="show-ssg"></span>
<span class="fragment" data-fragment-index="6" data-svg-classes="show-html"></span>
<span class="fragment" data-fragment-index="7" data-svg-classes="show-tests"></span>
<span class="fragment" data-fragment-index="8" data-svg-classes="show-mojo-test"></span>
<span class="fragment" data-fragment-index="9" data-svg-classes="show-cluster"></span>

---

## Demo

### <big>&darr;</big>

----

Clone a random Mojo🔥 project

```bash
# use a fresh directory
mkdir demo; cd demo
# clone EmberJson
git clone https://github.com/bgreni/EmberJson.git
# cd into it
cd EmberJson
```

----

Set up a Modo🧯 project and build API docs

```bash
# initialize Modo project in Hugo format
modo init hugo
# build API docs for the project
modo build
# serve API docs with Hugo
hugo serve -s docs/site/
```

Open the displayed URL in a web browser

---

## Features

### <big>&darr;</big>

----

### Output formats

Several fully supported rendering target

<br />

Plain Markdown<!-- .element: class="fragment" data-fragment-index="1" -->

<br />

<div class="columns" style="align-items: center; justify-content: center;">
<div class="col" style="margin-right: 2rem; flex: 1;">

[Hugo](https://gohugo.io)  
popular open-source SSG

</div><!-- .element: class="fragment" data-fragment-index="2" -->
<div class="col" style="margin-left: 2rem flex: 1;">

[mdBook](https://github.com/rust-lang/mdBook)  
known from the <i class="fa-brands fa-rust"></i>Rust book

</div><!-- .element: class="fragment" data-fragment-index="3" -->
</div>

<br />

Your wishes?<!-- .element: class="fragment" data-fragment-index="4" -->

----

### Cross-references

Very simple syntax, resembling Mojo imports

<div><div class="columns"><div class="col">

```python
"""
Relative ref to [.Struct.method] in the current module.
"""
```

</div><div class="col" style="flex:0.1;"></div><div class="col">

Relative ref to [Struct.method]() in the current module.

</div></div></div>
<!-- .element: class="fragment" data-fragment-index="2" -->
<div><div class="columns"><div class="col">

```python
"""
Absolute ref to module [pkg.mod].
"""
```

</div><div class="col" style="flex:0.1;"></div><div class="col">

Absolute ref to module [mod]().

</div></div></div>
<!-- .element: class="fragment" data-fragment-index="3" -->
<div><div class="columns"><div class="col">

```python
"""
Ref with [pkg.mod custom text].
"""
```

</div><div class="col" style="flex:0.1;"></div><div class="col">

Ref with [custom text]().

</div></div></div>
<!-- .element: class="fragment" data-fragment-index="4" -->

----

### Re-exports

Re-structure docs the way users see your package

<br />

<div class="columns" style="align-items: center; justify-content: center;"><div class="col">

<pre style="width:100%; font-size: 0.65em;">
- mypkg
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
<div class="col" style="flex:2.0;">

<div style="display:flex; flex-direction: column; align-items: center; justify-content: center; position: relative;">

```python
"""
Package mypkg...
"""
from .mod import Struct
from .subpkg.submod import Trait
```

<div style="position:absolute; width: 100%; z-index: 10;">

```python [4-6]
"""
Package mypkg...

Exports:
 - mod.Struct
 - subpkg.submod.Trait
"""
from .mod import Struct
from .subpkg.submod import Trait
```

</div><!-- .element: class="fragment" data-fragment-index="4" -->

</div><!-- .element: class="fragment" data-fragment-index="2" -->

</div>
<div class="col" style="flex:0.2">

#### <i class="fa-solid fa-arrow-right"></i>

</div><!-- .element: class="fragment" data-fragment-index="3" -->
<div class="col">

Modo🧯

<pre style="width:100%; font-size: 0.65em;">
- mypkg
  - Struct
  - Trait
</pre>

</div><!-- .element: class="fragment" data-fragment-index="3" -->
</div>

----

### Doc-tests

Convert code examples into unit tests

<div class="columns" style="align-items: center; justify-content: center;"><div class="col">

````python
"""
Doc-test example.

```mojo {doctest="sum"}
var a = 1 + 2
```

```mojo {doctest="sum" hide=true}
if a != 3:
    raise Error("failed")
```
"""
````

</div><!-- .element: class="fragment" data-fragment-index="1" -->
<div class="col" style="flex:0.4;">

#### <i class="fa-solid fa-arrow-right"></i>

</div><!-- .element: class="fragment" data-fragment-index="2" -->
<div class="col">

Doc-test example.

```python
var a = 1 + 2
```

<div>
<hr />

`..._test.mojo`

```python
fn test_sum() raises:
    var a = 1 + 2
    if a != 3:
        raise Error("failed")
```

</div><!-- .element: class="fragment" data-fragment-index="3" -->
</div><!-- .element: class="fragment" data-fragment-index="2" -->
</div>

----

### Scripts

Configure pre- and post-processing bash scripts

```yaml
# Bash scripts to run before build as well as test.
pre-run:
  - |
    echo Running 'mojo doc'...
    pixi run mojo doc -o docs/src/mypkg.json src/mypkg
    echo Done.
```
<!-- .element: class="fragment" data-fragment-index="1" -->

```yaml
# Bash scripts to run after test.
# Also runs after build if 'tests' is given.
post-test:
  - |
    echo Running 'mojo test'...
    pixi run mojo test -I src docs/test
    echo Done.
```
<!-- .element: class="fragment" data-fragment-index="2" -->

----

### Templates

Customize Markdown output via templates

```template
Mojo struct

# `{{.Name}}`

{{template "signature_struct" .}}

{{template "summary" . -}}
{{template "description" . -}}
{{template "aliases" . -}}
{{template "parameters" . -}}
{{template "fields" . -}}
{{template "parent_traits" . -}}
{{template "methods" . -}}
```
<!-- .element: class="fragment" data-fragment-index="1" -->

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

#### Please...

<br />

Specify cross-ref syntax <!-- .element: class="fragment" data-fragment-index="1" -->

Include package re-exports in JSON <!-- .element: class="fragment" data-fragment-index="2" -->

Support Markdown lists in <!-- .element: class="fragment" data-fragment-index="3" --> `Raises`<!-- .element: class="fragment" data-fragment-index="3" --> section

---

## Contributing

### <big>&darr;</big>

----

"Playtest"

Give feedback on tool and docs
<!-- .element: class="fragment" -->

Create issues & PRs
<!-- .element: class="fragment" -->

---

## Thank you!

[<i class="fa fa-github"></i>/mlange-42/modo](https://github.com/mlange-42/modo)

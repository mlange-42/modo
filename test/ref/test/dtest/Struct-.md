Mojo struct [🡭](https://github.com/mlange-42/modo/blob/main/test/src/doctest.mojo)

# `Struct`

```mojo
@memory_only
struct Struct
```

## Aliases

- `StructAlias = Int`

## Implemented traits

`AnyType`, `UnknownDestructibility`

## Methods

### `func`

```mojo
fn func(self)
```

Doctests in a struct member.

```mojo {doctest="func"}
var a = 1
```

**Args:**

- **self** (`Self`)



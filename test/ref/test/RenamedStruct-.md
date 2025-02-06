Mojo struct [🡭](https://github.com/mlange-42/modo/blob/main/test/src/mod.mojo)

# `RenamedStruct`

```mojo
@memory_only
struct Struct[StructParameter: Intable]
```

[`ModAlias`](_index.md#aliases).

[`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method)

## Aliases

- `StructAlias = StructParameter`: [`ModAlias`](_index.md#aliases). [`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method)

## Parameters

- **StructParameter** (`Intable`)

## Fields

- **struct_field** (`Int`): [`ModAlias`](_index.md#aliases). [`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method)

## Implemented traits

`AnyType`, `UnknownDestructibility`

## Methods

### `struct_method`

```mojo
fn struct_method[T: Intable](self, arg: StructParameter) -> Int
```

[`ModAlias`](_index.md#aliases).

[`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method)

**Parameters:**

- **T** (`Intable`): [`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method).

**Args:**

- **self** (`Self`)
- **arg** (`StructParameter`): [`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method).

**Returns:**

`Int`: Bla [`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method).

**Raises:**

Error [`RenamedStruct.struct_method`](RenamedStruct-.md#struct_method).

### `impl_method`

```mojo
fn impl_method(self, x: Int)
```

Test method for transclusions.

More details...

**Args:**

- **self** (`Self`)
- **x** (`Int`): Itself.



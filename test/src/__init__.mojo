"""
Package test.

Self ref [test].

Exports (rel) [.mod.ModuleAlias], [.mod.Struct], [.mod.Trait], [.mod.module_function], [.pkg].
Exports (abs) [test.mod.ModuleAlias], [test.mod.Struct], [test.mod.Trait], [test.mod.module_function], [test.pkg].

 - [.mod.Struct.StructParameter]
 - [.mod.Struct.struct_field]
 - [.mod.Struct.struct_method]

 - [.mod.Trait.trait_method]

 - [.pkg.submod]
 - [.pkg.submod.Struct]
 - [.pkg.submod.Struct.struct_method]

Exports:
 - mod.ModuleAlias
 - mod.Struct
 - mod.Trait
 - mod.module_function
 - pkg
 - doctest
"""
from .mod import ModuleAlias, Struct, Trait, module_function
from .pkg import submod

"""
Package test.

Self ref [test].

Exports (rel) [.mod.ModuleAlias], [.mod.Struct], [.mod.Trait], [.mod.module_function], [.pkg].
Exports (abs) [test.mod.ModuleAlias], [test.mod.Struct], [test.mod.Trait], [test.mod.module_function], [test.pkg].

 - [.mod.Struct.StructParameter]
 - [.mod.Struct.struct_field]
 - [.mod.Struct.struct_method]
 - [.mod.Struct.StructAlias]

 - [.mod.Trait.trait_method]

 - [.pkg.submod]
 - [.pkg.submod.Struct]
 - [.pkg.submod.Struct.struct_method]
 - [.pkg.submod.ModuleAlias]

Exports:
 - mod.ModuleAlias as ModAlias
 - mod.Struct as RenamedStruct
 - mod.Trait
 - mod.module_function
 - pkg
 - doctest as dtest
 - doctest.ModuleAlias2 as ModAlias2
"""
from .mod import (
    ModuleAlias as ModAlias,
    Struct as RenamedStruct,
    Trait,
    module_function,
)
from .doctest import ModuleAlias2 as ModAlias2
from .pkg import submod

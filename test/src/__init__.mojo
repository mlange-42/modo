"""
Package test.

Self ref [test].

Exports (rel) [.mod.ModuleAlias], [.mod.Struct], [.mod.Trait], [.mod.module_function], [.pkg].
Exports (abs) [test.mod.ModuleAlias], [test.mod.Struct], [test.mod.Trait], [test.mod.module_function], [test.pkg].

Exports:
 - mod.ModuleAlias
 - mod.Struct
 - mod.Trait
 - mod.module_function
 - pkg
"""
from .mod import ModuleAlias, Struct, Trait, module_function
from .pkg import submod

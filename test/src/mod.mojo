"""
Module mod.

Contains (rel) [.ModuleAlias], [.Struct], [.Trait], [.module_function].
Contains (abs) [test.mod.ModuleAlias], [test.mod.Struct], [test.mod.Trait], [test.mod.module_function].
"""
alias ModuleAlias = Int


struct Struct[StructParameter: Intable]:
    """[.Struct], [.ModuleAlias].

    [..pkg.submod.Struct].
    """

    alias StructAlias = StructParameter

    var struct_field: Int

    fn struct_method(self, arg: StructParameter) -> Int:
        return self.struct_field


trait Trait:
    # TODO: fields in traits are not supported yet
    # var trait_field: Int

    fn trait_method(self, arg: Int) -> Int:
        ...


fn module_function[FunctionParameter: Intable](arg: Int) -> Int:
    return arg

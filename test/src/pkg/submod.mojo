alias ModuleAlias = Int


struct Struct[StructParameter: Intable]:
    """[..submod.Struct], [..submod.ModuleAlias].

    [.Struct], [.ModuleAlias]
    """

    alias StructAlias = StructParameter

    var struct_field: Int

    fn struct_method(self, arg: StructParameter) -> Int:
        return self.struct_field

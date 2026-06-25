comptime ModuleAlias = Int


struct Struct[StructParameter: Intable]:
    """[..submod.Struct], [..submod.ModuleAlias].

    [.Struct], [.ModuleAlias]
    """

    comptime StructAlias = Self.StructParameter

    var struct_field: Int

    fn struct_method(self, arg: Self.StructParameter) -> Int:
        return self.struct_field

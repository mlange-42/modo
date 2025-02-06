"""
Module mod.

Contains (rel) [.ModuleAlias], [.Struct], [.Trait], [.module_function].
Contains (abs) [test.mod.ModuleAlias], [test.mod.Struct], [test.mod.Trait], [test.mod.module_function].

[.Struct.StructAlias]
"""
alias ModuleAlias = Int


struct Struct[StructParameter: Intable]:
    """[.ModuleAlias].

    [.Struct.struct_method]"""

    alias StructAlias = StructParameter
    """[.ModuleAlias].

    [.Struct.struct_method]"""

    var struct_field: Int
    """[.ModuleAlias].

    [.Struct.struct_method]"""

    fn struct_method[T: Intable](self, arg: StructParameter) raises -> Int:
        """[.ModuleAlias].

        [.Struct.struct_method]

        Parameters:
            T: [.Struct.struct_method].

        Args:
            arg: [.Struct.struct_method].

        Returns:
            Bla [.Struct.struct_method].

        Raises:
            Error [.Struct.struct_method].
        """
        return self.struct_field

    fn impl_method(self, x: Int):
        """[[.Trait]]."""
        pass


trait Trait:
    """[.ModuleAlias].

    [.Struct.struct_method]"""

    # TODO: fields in traits are not supported yet
    # var trait_field: Int

    fn trait_method[T: Intable](self, arg: T) raises -> Int:
        """[.ModuleAlias].

        [.Struct.struct_method]

        Parameters:
            T: [.Struct.struct_method].

        Args:
            arg: [.Struct.struct_method].

        Returns:
            Bla [.Struct.struct_method].

        Raises:
            Error [.Struct.struct_method].
        """
        ...

    fn impl_method(self, x: Int):
        """Test method for transclusions.

        More details...

        Args:
          x: Itself.
        """
        ...


fn module_function[FunctionParameter: Intable](arg: Int) raises -> Int:
    """[.ModuleAlias].

    [.Struct.struct_method]

    Parameters:
        FunctionParameter: [.Struct.struct_method].

    Args:
        arg: [.Struct.struct_method].

    Returns:
        Bla [.Struct.struct_method].

    Raises:
        Error [.Struct.struct_method].
    """
    return arg

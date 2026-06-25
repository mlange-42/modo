comptime MyIntAlias = Int
"""An example module alias."""


struct MyPair[T: Intable & ImplicitlyCopyable]:
    """
    A simple example struct.

    This struct is re-exported by its [mypkg parent package], so shows up where the user expects it due to the import path.
    It has [aliases](#aliases), [parameters](#parameters), [fields](#fields) and [methods](#methods).

    Linking to individual members is as easy as this:

    ```
    Method: [.MyPair.dump], field: [.MyPair.first].
    ```

    which gives:

    Method: [.MyPair.dump], field: [.MyPair.first].

    Parameters:
      T: The [.MyPair]'s element type.
    """

    comptime MyInt = MyIntAlias
    """An example struct alias. Alias for [.MyIntAlias]."""

    var first: Self.T
    """First struct field."""
    var second: Self.T
    """Second struct field."""

    fn __init__(out self, first: Self.T, second: Self.T):
        """
        Creates a new [.MyPair].

        Args:
            first: The value for [.MyPair.first].
            second: The value for [.MyPair.second].
        """
        self.first = first
        self.second = second

    fn dump(self):
        """Prints fields [.MyPair.first `first`] and [.MyPair.second `second`].
        """
        print(Int(self.first), Int(self.second))

    fn format(self, fmt: String) raises -> String:
        """
        Formats this [.MyPair].

        A longer description of the method,
        explaining how to use it.

        Args:
            fmt: Desired output format.

        Returns:
            A formatted string.

        Raises:
            When called with an unknown `fmt` argument.
        """
        return ""

    @staticmethod
    fn static_method(self):
        pass

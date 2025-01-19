struct MyPair[T: Intable]:
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

    alias MyInt = Intable
    """An example alias."""

    var first: T
    """First struct field."""
    var second: T
    """Second struct field."""

    fn __init__(out self, first: T, second: T):
        """
        Creates a new [.MyPair].

        Args:
            first: The value for [.MyPair.first].
            second: The value for [.MyPair.second].
        """
        self.first = first
        self.second = second

    fn dump(self):
        """Creates a new [.MyPair]'s fields [.MyPair.first `first`] and [.MyPair.second `second`].
        """
        print(Int(self.first), Int(self.second))

struct MyPair[T: Intable]:
    """
    A simple example struct.

    This struct is re-exported by its [mypkg parent package], so it where the user expects it due to the import path.
    It has [aliases](#aliases), [parameters](#parameters), [fields](#fields) and [methods](#methods).

    Linking to individual members is as easy as this:

    ```
    Method: [.MyPair.dump], field: [.MyPair.first].
    ```

    which gives:

    Method: [.MyPair.dump], field: [.MyPair.first].
    """

    alias MyInt = Intable

    var first: T
    var second: T

    fn __init__(out self, first: T, second: T):
        self.first = first
        self.second = second

    fn dump(self):
        print(Int(self.first), Int(self.second))

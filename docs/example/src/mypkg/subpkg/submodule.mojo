"""
Module in subpkg.

Link to [.dummy].

Link to [...subpkg].
"""


fn dummy():
    """
    Dummy function.

    Link to [...subpkg].

    Abs link to [mypkg.subpkg]

    Link to [mypkg].

    Link to containing module: [..submodule]

    Link to cstruct in mypkg: [...mymodule.MyPair]
    """
    pass


struct SubStruct:
    """
    Dummy struct.

    Link to [.dummy].
    """

    fn test(self):
        """
        Link in struct member: [.dummy].
        """
        pass

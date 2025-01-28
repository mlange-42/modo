"""
Module in subpkg.

This module is re-exported by [...subpkg],
so all its members are included.
"""


fn simple_func():
    """A simple function."""
    pass


fn overloaded_func():
    """An overloaded function, variant 1."""
    pass


fn overloaded_func(v: Int):
    """
    An overloaded function, variant 2.

    Args:
        v: Some arguments...
    """
    pass


fn overloaded_func(v: String):
    """
    An overloaded function, variant 3.

    Args:
        v: Some arguments...
    """
    pass

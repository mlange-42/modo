from setuptools import setup

setup(
    name="pymodo",
    version="0.8.1",
    packages=["pymodo"],
    entry_points = {
        'console_scripts': ['pymodo=pymodo.run:main'],
    },
)

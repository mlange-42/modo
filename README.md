# Modo

[Mojo](https://www.modular.com/mojo) documentation generator.

Generates markdown for static site generators (SSGs) from `mojo doc` JSON output.

**! Early work in progress !**

## Usage

Piping `mojo doc` to Modo:

```
mojo doc <src-path> | modo <out-dir>
```

Reading from a file:

```
modo <out-dir> -I file.json
```

Command line help:

```
modo -h
```

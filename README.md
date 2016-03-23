# simplebp
Minimal example of using Blueprint

[Blueprint](https://github.com/google/blueprint) is a meta-build system that uses [Ninja](https://ninja-build.org/) 
as its backend. The project itself doesn't build anything (other than itself) -- build system developers use it to 
define build rules, and project developers create Blueprint files to describe the inputs and outputs of the build
system.

**Simplebp** is a minimalistic builder that uses Blueprint to compile C/C++ binaries and shared libraries and
to run scripts. I don't intend for this to become a production-quality build system, but rather to offer
a relatively simple implementation of Blueprint. I'll probably add a few features as I go along, to explore
what Blueprint can do, but for more advanced usage, be sure to look at
[Android's soong](https://android.googlesource.com/platform/build/soong/).

## Getting started

First, initialize the git repos.
```bash
git clone http://github.com/TKilbourn/simplebp
cd simplebp
git submodule update --init
```

Now create an `out` directory and run the bootstrap script.
```bash
mkdir out && cd out
../bootstrap.bash
```

Now at any time, just execute `out/simplebp` to execute a build.

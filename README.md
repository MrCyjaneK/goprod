# GoProd - Let's go live on production!

GoProd is a simple tool that let you package your Golang app in seconds!

## Usage

```bash
$ goprod \
     -combo="linux/arm;linux/i386;linux/arm64;linux/amd64;windows/amd64;windows/i386" \
     -builddir="build" \
     -tags="gui"
```

# Installation

To install from source: `go run main.go -binname="goprod" -combo="$(go env GOOS)/$(go env GOARCH) && apt install ./build/deb/goprod*.deb"`

To get a prebuild binary:

```bash
# wget 'https://static.mrcyjanek.net/abstruse/apt-repository/mrcyjanek-repo/mrcyjanek-repo_2.0-1_all.deb' && \
  apt install ./mrcyjanek-repo_2.0-1_all.deb && \
  rm ./mrcyjanek-repo_2.0-1_all.deb && \
  apt update \
  apt install goprod
```

### Advanced Usage

Goprod exist because I want to package my programs with ease, and I expect thhem to 'just work', I don't care about proper debian packaging (well I do - but not here.). Here are all possible options

You can build and package for:

 - [x] Android (.apk)
 - [x] Debian (.deb)
 - [x] MacOS (partial, .zip)
 - [x] Windows (aka portable exe)
 - [ ] Ubuntu Touch (.click)
 - [ ] Flatpak
 - [ ] Sailfish OS

I generally recommend you using the docker image if you have some CI/CD solution installed (btw you should use [Abstruse](https://github.com/bleenco/abstruse)). If you want to run something on your machine just for development - then grab a copy of `gcc` and `g++` and `checkinstall`, and make sure that `/bin/x86_64-linux-gnu-g++` is a symlink to your gcc, on debian it's automatic (afaik) but on arch I had to do that manually. So keep this in mind when something goes wrong with GCC.

## Reference app

For an reference app check out [magnetgraph](https://github.com/MrCyjaneK/magnetgraph) I use goprod to package it

## Preparing

You need to install depedencies, for debian check `docker/` for an up-to-date list of **all** depedencies (you usually don't need all of them, unless you deploy from the machine on which you work.)

So list of required things:

 - To compile native pure-go apps: `golang`
 - To compile (native) apps that require c/c++ stuff: `gcc` and `g++`
   - To compile apps for the adware os: `gcc-mingw-w64-x86-64-win32` `g++-mingw-w64-x86-64-win32`
   - To compile apps for the bloarware os:
     - This is actually the hard path, but chill - we can handle that for you:
     - `goprod update-ndk` - this command will `update` (if required) the ndk and place it by default in ~/Android/Sdk/ndk
     This will also download commandline tools and place it by default in ~/Android/Sdk/cmdline-tools/
     - `goprod accept-license` will accept the sdk licenses, so you will be able to compile things for android.
     - To package for android `default-jre` is required.

`goprod` is accessible there as a command, together with everything configured, except for your email - you may want to do `git config user.email your-email@selfhostbtw.onion` and `git config user.name "Your Name"`, because this name and email is used to package debian package.

### A word about docker

There are a couple of docker images available, and if you feel bad about using the public docker registry - `make docker`:

 - `mrcyjanek/goprod:core`
   This is the smallest image - it only contain `goprod` executable, and (sometimes) not it's up to date version (it's a cached version from my apt, in the other versions it is getting updated by compiling what's in this repo). It's used to build other docker images.
   It's fine to use it for `darwin` builds afaik, but I didn't really test it.
 - `mrcyjanek/goprod:nodejs`
   This is like `core` but with nodejs installed. Why? Sometimes you may want to use some nodejs thing to build frontend, that's why.
 - `mrcyjanek/goprod:core-android`
   This image contain things to build stuff for android, that include android ndk, sdk, java and an emulator.
 - `mrcyjanek/goprod:nodejs-android`
   Same as above but with nodejs
 - `mrcyjanek/goprod:core-linux`
   This image contain all the things required for linux (cross) compiling and packaging.
 - `mrcyjanek/goprod:nodejs-linux`
   Same as above but with nodejs
 - `mrcyjanek/goprod:core-windows`
   This image contain 32 and 64bit mingw for cross compiling.
 - `mrcyjanek/goprod:nodejs-windows`
   same as above but with nodejs
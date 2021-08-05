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

### Docker

You can also use it directly in docker, `mrcyjanek/goprod:latest` this image is updated on each commit to this repo, and contain gcc/g++ for amd64,aaarch64,arm and i386 linux and windows, i686 and x64


### Advanced Usage

Goprod exist because I want to package my programs with ease, and I expect thhem to 'just work', I don't care about proper debian packaging (well I do - but not here.). Here are all possible options

You can build and package for:

 - [x] Android 
 - [x] Debian (.deb)
 - [x] MacOS (partial)
 - [x] Windows (aka portable exe)
 - [ ] Ubuntu Touch (.click)
 - [ ] Flatpak
 - [ ] Sailfish OS

I generally recommend you using the docker image if you have some CI/CD solution installed (btw you should use [Abstruse](https://github.com/bleenco/abstruse)). If you want to run something on your machine just for development - then grab a copy of `gcc` and `g++` and `checkinstall`, and make sure that `/bin/x86_64-linux-gnu-g++` is a symlink to your gcc, on debian it's automatic (afaik) but on arch I had to do that manually. So keep this in mind when something goes wrong with GCC.

## Preparing

You need to install depedencies, for debian check `Dockerfile` for an up-to-date list of **all** depedencies (you usually don't need all of them, unless you deploy from the machine on which you work.)

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
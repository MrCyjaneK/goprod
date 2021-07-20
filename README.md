# GoProd - Let's go live on production!

GoProd is a simple tool that let you package your Golang app in seconds!

## Usage

```bash
$ goprod
$ goprod golang build \
     -combo="linux/arm;linux/i386;linux/arm64;linux/amd64;windows/amd64;windows/i386" \
     -builddir="build" \
     -tags="gui"
```

# Installation

To install from source: `go run main.go -binname="goprod" -combo="$(go env GOOS)/$(go env GOARCH) && apt install ./build/deb/goprod*.deb"`

To get a prebuild binary:

```bash
# wget 'https://static.mrcyjanek.net/laminarci/apt-repository/cyjan_repo/mrcyjanek-repo-latest.deb' && \
  apt install ./mrcyjanek-repo-latest.deb && \
  rm ./mrcyjanek-repo-latest.deb && \
  apt update
```

```bash
# apt install goprod
```

You can also use it directly in docker, `mrcyjanek/goprod:latest` this image is updated on each commit to this repo, and contain gcc/g++ for amd64,aaarch64,arm and i386 linux and windows, i686 and x64

`goprod` is accessible there as a command, together with everything configured, except for your email - you may want to do `git config user.email your-email@selfhostbtw.onion` and `git config user.name "Your Name"`, because this name and email is used to package debian package.
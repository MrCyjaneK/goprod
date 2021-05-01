# GoProd - Let's go live on production!

GoProd is a simple tool that let you package your Golang app in seconds!

## Usage

```bash
$ goprod golang build # This fill build for your host target ONLY.
$ goprod golang build \
     -combo="linux/arm;linux/i386;linux/arm64;linux/amd64;windows/amd64;windows/i386" \
     -builddir="build" \
     -tags="gui"
```

# Installation

You can either run `go build` inside of this project, `go get` it, or install it [from my repo](https://git.mrcyjanek.net/mrcyjanek/mrcyjanekrepo):

```bash
# wget 'https://static.mrcyjanek.net/laminarci/apt-repository/cyjan_repo/mrcyjanek-repo-latest.deb' && \
  apt install ./mrcyjanek-repo-latest.deb && \
  rm ./mrcyjanek-repo-latest.deb && \
  apt update
```

```bash
# apt install goprod
```
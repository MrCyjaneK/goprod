FROM golang:latest

RUN echo deb http://deb.debian.org/debian buster-backports main > /etc/apt/sources.list.d/backports.list
# Mingw is only in bullseye
RUN echo deb http://deb.debian.org/debian bullseye main > /etc/apt/sources.list.d/bullseye.list
RUN apt update
RUN apt install checkinstall \
        gcc-aarch64-linux-gnu g++-aarch64-linux-gnu \
        gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf \
        gcc-i686-linux-gnu g++-i686-linux-gnu \ 
        gcc-mingw-w64-x86-64-win32 g++-mingw-w64-x86-64-win32 \
        gcc-mingw-w64-i686-win32 g++-mingw-w64-i686-win32 \
        gcc g++ -y


COPY . /go/src/git.mrcyjanek.net/goprod/
WORKDIR /go/src/git.mrcyjanek.net/goprod/
RUN ls
RUN go run main.go -combo="linux/amd64" -builddir="build" -binname="goprod"
RUN apt install ./build/deb/goprod*.deb -y
WORKDIR /go
RUN goprod ndk-update
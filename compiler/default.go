// +build !linux

package compiler

import (
	"log"
	"runtime"
)

func Build(combo string, tags string, binname string, builddir string, ndk string, ldflags string, shoulddel bool) {
	log.Fatal(runtime.GOOS, runtime.GOARCH, "")
}

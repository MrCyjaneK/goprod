// +build !linux

package debpackage

import (
	"log"
	"runtime"
)

func Build(combo string, tags string, binname string, builddir string, version string) {
	log.Fatal(runtime.GOOS, runtime.GOARCH, "")
}

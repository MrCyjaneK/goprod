// +build !linux

package debpackage

import (
	"log"
	"runtime"
)

func Build(combo string, binname string, bindir string, debdir string, version string) {
	log.Fatal(runtime.GOOS, runtime.GOARCH, "")
}

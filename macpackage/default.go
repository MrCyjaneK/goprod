package macpackage

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"strings"
)

func Package(combo string, binname string, bindir string, zipdir string, version string) {
	spl := strings.Split(combo, "/")
	if len(spl) != 2 {
		log.Fatal("Invalid " + combo + " provided.")
	}
	GOOS := spl[0]
	GOARCH := spl[1]
	if GOOS != "darwin" {
		log.Fatal("no")
	}
	filepath := bindir + "/" + binname + "_" + GOOS + "_" + GOARCH
	target := zipdir + "/" + binname + "_" + GOOS + "_" + GOARCH + ".zip"
	os.MkdirAll(zipdir, 0750)
	zipfile, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}

	zipwriter := zip.NewWriter(zipfile)
	defer zipwriter.Close()

	fileToZip, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		log.Fatal(err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		log.Fatal(err)
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = header.FileInfo().Name()

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipwriter.CreateHeader(header)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		log.Fatal(err)
	}
	// readme
	header, err = zip.FileInfoHeader(info)
	if err != nil {
		log.Fatal(err)
	}
	header.Name = "readme-macos.txt"
	header.Method = zip.Deflate
	writer, err = zipwriter.CreateHeader(header)
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.Write([]byte(`Welcome to simple guide to running MacOS binaries!
	
1. Extract the executable from this archive
2. Double-click on it
3. You will be warned, click cancel
4. Double-click again
5. Click Open

Enjoy.

Autogenerated by git.mrcyjanek.net/mrcyjanek/goprod`))
	if err != nil {
		log.Fatal(err)
	}
}

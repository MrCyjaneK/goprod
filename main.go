package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"git.mrcyjanek.net/mrcyjanek/goprod/apkpackage"
	"git.mrcyjanek.net/mrcyjanek/goprod/appimage"
	"git.mrcyjanek.net/mrcyjanek/goprod/compiler"
	"git.mrcyjanek.net/mrcyjanek/goprod/debpackage"
	"git.mrcyjanek.net/mrcyjanek/goprod/macpackage"
	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"
)

const (
	default_combo       = runtime.GOOS + "/" + runtime.GOARCH
	default_builddir    = "build"
	default_ldflags     = ""
	default_binname     = "binname"
	default_tags        = "goprod"
	default_version     = "1.0.0"
	default_shouldpkg   = true
	default_apkit       = true
	default_appurl      = "http://127.0.0.1:8081"
	default_deltmp      = true
	default_buildcmd    = ""
	default_apktemplate = "console"
	default_appimageit  = false
)

var (
	combo     = flag.String("combo", default_combo, "Which combo should I serve?")
	builddir  = flag.String("builddir", default_builddir, "Where should the files get saved.")
	ldflags   = flag.String("ldflags", default_ldflags, "Things to get passwd by with --ldflags")
	binname   = flag.String("binname", default_binname, "What is the program name?")
	tags      = flag.String("tags", default_tags, "Tags that are passed to go build command")
	version_a = flag.String("version", default_version, "Version of your program.")
	//TODO: Do **NOT** hardcode the path here, bruh.
	ndk_a       = flag.String("ndk", "~/Android/Sdk/ndk/android-ndk-r22b/toolchains/llvm/prebuilt/linux-x86_64/bin/", "Path to android toolchain")
	sdk_a       = flag.String("sdkpath", "~/Android/Sdk/", "Path to android Sdk")
	shouldpkg   = flag.Bool("shouldpkg", default_shouldpkg, "Should we create a package out of the binary?")
	apkit       = flag.Bool("apkit", default_apkit, "Should I create android app?")
	appurl      = flag.String("appurl", default_appurl, "What url should I open in native app?")
	deltmp      = flag.Bool("deltmp", default_deltmp, "Should I delete tmp files?")
	buildcmd    = flag.String("buildcmd", default_buildcmd, "What command should be used to build the program? Defaults to 'go build`")
	apktemplate = flag.String("apktemplate", default_apktemplate, "Which template should I use for building the .apk?")
	appimageit  = flag.Bool("appimageit", default_appimageit, "Should I create tha appimage?")
)
var ndk string
var sdk string
var version string

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.Parse()
	verout, err := exec.Command("git", "show", "-s", "--date=format:%Y%m%d%H%M", "--format=%cd").Output()
	if err != nil {
		verout = []byte("99999notagitrepo")
	}
	commit := strings.Split(string(verout), "\n")[0]
	version = *version_a + "+git" + commit

	usr, err := user.Current()
	if err != nil {
		log.Fatal("user.Current():", err)
	}
	ndk = strings.ReplaceAll(*ndk_a, "~", usr.HomeDir)
	sdk = strings.ReplaceAll(*sdk_a, "~", usr.HomeDir)
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "ndk-update":
			{
				updateNdk()
				return
			}
		case "accept-license":
			{
				acceptLicense()
				return
			}
		}
	}

	//os.RemoveAll(*builddir)
	os.MkdirAll(*builddir, 0750)
	log.Println(*combo)
	for _, i := range strings.Split(*combo, ";") {
		spl := strings.Split(i, "/")
		if len(spl) != 2 {
			log.Fatal("Invalid combo '" + i + "' provided.")
		}
		GOOS := spl[0]
		GOARCH := spl[1]
		log.Println("Compiling...")
		buildargs, err := gosh.Split(*buildcmd)
		if err != nil {
			log.Fatal("Unable to parse command: '"+*buildcmd+"'", err)
		}
		if GOOS == "android" && GOARCH == "all" {
			continue
		}
		compiler.Build(i, *tags, *binname, *builddir+"/bin", ndk, *ldflags, buildargs)
		if GOOS == "linux" && *shouldpkg {
			log.Println("Packaging (deb)...")
			debpackage.Build(i, *binname, *builddir+"/bin", *builddir+"/deb", version)
		}
		if GOOS == "linux" && *appimageit {
			log.Println("Packaging (appimage)...")
			appimage.Package(buildargs, *ldflags, *builddir+"/AppDir", *binname, GOOS, GOARCH, *tags, version, *builddir+"/appimage")
		}
		if GOOS == "darwin" && *shouldpkg {
			log.Println("Packaging (zip)...")
			macpackage.Package(i, *binname, *builddir+"/bin", *builddir+"/mac", version)
		}
		if GOOS == "android" && *apkit {
			log.Println("Packaging (apk)...")
			apkpackage.Package(*binname, *builddir+"/bin", *builddir+"/apk", version, *appurl, sdk, *deltmp, *apktemplate, GOARCH)
		}
	}
}

func acceptLicense() {
	cmd := exec.Command(sdk+"/cmdline-tools/bin/sdkmanager", "--sdk_root="+sdk, "--licenses")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err) //replace with logger, or anything you want
	}
	go func() {
		for {
			time.Sleep(time.Second)
			_, err = io.WriteString(stdin, "y\n")
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()
}

func updateNdk() {
	resp, err := http.Get("https://developer.android.com/studio")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	sdkstra := strings.Split(strings.ReplaceAll(string(body), "\"", ">"), ">")
	sdkstr := grep(".zip", sdkstra)
	sdkstr = grep("commandlinetools-linux", sdkstr)
	link := sdkstr[1]
	b, err := os.ReadFile(path.Join(sdk, "version-sdk"))
	if err != nil || string(b) == link {
		log.Println("current version", string(b), "google's version:", link)
		log.Println("downloading")
		resp, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		out, err := os.Create("/tmp/sdk.zip")
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("writing to file")
		io.Copy(out, resp.Body)
		log.Println("unzipping")
		err = Unzip("/tmp/sdk.zip", path.Join(sdk))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("cleaning up")
		os.Remove("/tmp/sdk.zip")
		os.WriteFile(path.Join(sdk, "version-sdk"), []byte(link), 0755)
	}

	//#!/bin/bash
	//mkdir -p ~/Android/Sdk/ndk/
	log.Println("creating directory", path.Join(sdk, "/ndk"))
	os.MkdirAll(path.Join(sdk, "/ndk"), 0755)
	//latest=$(wget --quiet https://developer.android.com/ndk/downloads/ -O - | tr '>' ">\n" | grep "linux-x86_64.zip" | grep href | tr '"' "\n" | head -2 | tail -1)
	log.Println("fetching latest version number")
	resp, err = http.Get("https://developer.android.com/ndk/downloads/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll", err)
	}
	stra := strings.Split(string(body), ">")
	str := grep("-linux.zip", stra)
	str = strings.Split(strings.Join(str, "\""), "\"")
	link = str[9]
	b, err = os.ReadFile(path.Join(sdk, "version-ndk"))
	if err != nil || string(b) == link {
		log.Println("current version", string(b), "google's version:", link)
		log.Println("downloading")
		resp, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		out, err := os.Create("/tmp/ndk.zip")
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("writing to file")
		io.Copy(out, resp.Body)
		log.Println("unzipping")
		err = Unzip("/tmp/ndk.zip", path.Join(sdk, "/ndk"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("cleaning up")
		os.Remove("/tmp/ndk.zip")
		os.WriteFile(path.Join(sdk, "version-ndk"), []byte(link), 0755)
	}
	// Now imma accept license real quick
}
func grep(search string, in []string) []string {
	var resp []string
	for i := range in {
		if strings.Contains(in[i], search) {
			resp = append(resp, in[i])
		}
	}
	return resp
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

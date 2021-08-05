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
	"strconv"
	"strings"
	"time"

	"git.mrcyjanek.net/mrcyjanek/goprod/apkpackage"
	"git.mrcyjanek.net/mrcyjanek/goprod/compiler"
	"git.mrcyjanek.net/mrcyjanek/goprod/debpackage"
	"git.mrcyjanek.net/mrcyjanek/goprod/macpackage"
)

var (
	combo    = flag.String("combo", "linux/arm;linux/386;linux/arm64;linux/amd64", "Combo that I should serve")
	builddir = flag.String("builddir", "build", "Where should the files get saved.")
	ldflags  = flag.String("ldflags", "", "Things to get passwd by with --ldflags")
	binname  = flag.String("binname", "helloworld", "What is the program name?")
	tags     = flag.String("tags", "goprod", "Tags that are passed to go build command")
	versiona = flag.String("version", "0.0.1", "Version of your program.")
	//TODO: Do **NOT** hardcode the path here, bruh.
	ndka      = flag.String("ndk", "~/Android/Sdk/ndk/android-ndk-r22b/toolchains/llvm/prebuilt/linux-x86_64/bin/", "Path to android toolchain")
	sdkpath   = flag.String("sdkpath", "~/Android/Sdk/", "Path to android Sdk")
	shouldpkg = flag.Bool("package", true, "Should we create a package out of the binary?")
	apkit     = flag.Bool("apkit", true, "Should I create android app?")
	apport    = flag.String("appport", "8081", "What port should I open in native app?")
	deltmp    = flag.Bool("deltmp", true, "Should I delete tmp files?")
)
var ndk string
var sdk string
var version string

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.Parse()
	t := time.Now()
	year := "" + strconv.Itoa(t.Year())
	month := "0" + strconv.Itoa(int(t.Month()))
	month = month[len(month)-2:]
	day := "0" + strconv.Itoa(t.Day())
	day = day[len(day)-2:]
	hour := "0" + strconv.Itoa(t.Hour())
	hour = hour[len(hour)-2:]
	minute := "0" + strconv.Itoa(t.Minute())
	minute = minute[len(minute)-2:]
	version = *versiona + "-" + year + month + day + hour + minute
	usr, err := user.Current()
	if err != nil {
		log.Fatal("user.Current():", err)
	}
	ndk = strings.ReplaceAll(*ndka, "~", usr.HomeDir)
	sdk = strings.ReplaceAll(*sdkpath, "~", usr.HomeDir)
	os.Exit(0)
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
	androidused := false
	for _, i := range strings.Split(*combo, ";") {
		spl := strings.Split(i, "/")
		if len(spl) != 2 {
			log.Fatal("Invalid combo '" + i + "' provided.")
		}
		GOOS := spl[0]
		log.Println("Compiling...")
		compiler.Build(i, *tags, *binname, *builddir+"/bin", ndk, *ldflags)
		if GOOS == "linux" && *shouldpkg {
			log.Println("Packaging...")
			debpackage.Build(i, *binname, *builddir+"/bin", *builddir+"/deb", version)
		}
		if GOOS == "darwin" && *shouldpkg {
			log.Println("Packaging...")
			macpackage.Package(i, *binname, *builddir+"/bin", *builddir+"/mac", version)
		}
		if GOOS == "android" {
			androidused = true
		}
	}
	if *apkit && androidused {
		apkpackage.Package(*binname, *builddir+"/bin", *builddir+"/apk", version, *apport, sdk, *deltmp)
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
		log.Fatal(err)
	}
	stra := strings.Split(string(body), ">")
	str := grep("linux-x86_64.zip", stra)
	str = grep("href", str)
	str = strings.Split(strings.Join(str, "\""), "\"")
	link = str[1]
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

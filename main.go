package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"git.mrcyjanek.net/mrcyjanek/goprod/apkpackage"
	"git.mrcyjanek.net/mrcyjanek/goprod/compiler"
	"git.mrcyjanek.net/mrcyjanek/goprod/debpackage"
	"git.mrcyjanek.net/mrcyjanek/goprod/macpackage"
)

var (
	combo     = flag.String("combo", "linux/arm;linux/386;linux/arm64;linux/amd64", "Combo that I should serve")
	builddir  = flag.String("builddir", "build", "Where should the files get saved.")
	ldflags   = flag.String("ldflags", "", "Things to get passwd by with --ldflags")
	binname   = flag.String("binname", "helloworld", "What is the program name?")
	tags      = flag.String("tags", "goprod", "Tags that are passed to go build command")
	versiona  = flag.String("version", "0.0.1", "Version of your program.")
	ndka      = flag.String("ndk", "~/Android/Sdk/ndk/22.1.7171670/toolchains/llvm/prebuilt/linux-x86_64/bin/", "Path to android toolchain")
	sdkpath   = flag.String("sdkpath", "~/Android/Sdk/", "Path to android Sdk")
	shouldpkg = flag.Bool("package", true, "Should we create a package out of the binary?")
	apkit     = flag.Bool("apkit", true, "Should I create android app?")
	apport    = flag.String("appport", "8081", "What port should I open in native app?")
	deltmp    = flag.Bool("deltmp", true, "Should I delete tmp files?")
)
var ndk string
var version string

func main() {
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
		log.Fatal(err)
	}
	ndk := strings.ReplaceAll(*ndka, "~", usr.HomeDir)
	sdk := strings.ReplaceAll(*sdkpath, "~", usr.HomeDir)
	//os.RemoveAll(*builddir)
	os.MkdirAll(*builddir, 0750)
	log.Println(*combo)
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
	}
	if *apkit {
		apkpackage.Package(*binname, *builddir+"/bin", *builddir+"/apk", version, *apport, sdk, *deltmp)
	}
}

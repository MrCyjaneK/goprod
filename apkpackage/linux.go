// +build linux

package apkpackage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Package(binname string, bindir string, apkdir string, version string, appurl string, sdkpath string, shoulddel bool, apktemplate string, GOARCH string) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	os.MkdirAll(apkdir, 0750)
	workdir, err := ioutil.TempDir(os.TempDir(), "apkbuild")
	if err != nil {
		log.Fatal(err)
	}
	os.Mkdir(workdir, 0770)
	if shoulddel {
		defer func() {
			os.RemoveAll(workdir)
		}()
	}
	log.Println("Building apk in " + workdir)
	copyDir("/usr/share/goprod/android-"+apktemplate, workdir)
	log.Println("Replacing stuff.")
	fileReplace(workdir+"/app/src/main/java/x/x/@appname@/MainActivity.kt", "@appurl@", appurl)
	fileReplace(workdir+"/app/build.gradle", "@version@", version)
	fileReplace(workdir+"/settings.gradle", "@appname@", binname)
	fileReplace(workdir+"/app/build.gradle", "@appname@", binname)
	fileReplace(workdir+"/app/src/androidTest/java/x/x/@appname@/ExampleInstrumentedTest.kt", "@appname@", binname)
	fileReplace(workdir+"/app/src/main/AndroidManifest.xml", "@appname@", binname)
	fileReplace(workdir+"/app/src/main/java/x/x/@appname@/MainActivity.kt", "@appname@", binname)
	fileReplace(workdir+"/app/src/main/java/x/x/@appname@/MainReceiver.kt", "@appname@", binname)
	fileReplace(workdir+"/app/src/main/res/values/strings.xml", "@appname@", binname)
	fileReplace(workdir+"/app/src/main/res/values/themes.xml", "@appname@", binname)
	fileReplace(workdir+"/app/src/main/res/values/themes.xml", "@appname@", binname)
	fileReplace(workdir+"/app/src/main/res/values-night/themes.xml", "@appname@", binname)
	fileReplace(workdir+"/app/src/test/java/x/x/@appname@/ExampleUnitTest.kt", "@appname@", binname)
	log.Println("Renaming stuff.")
	os.Rename(workdir+"/app/src/androidTest/java/x/x/@appname@", workdir+"/app/src/androidTest/java/x/x/"+binname)
	os.Rename(workdir+"/app/src/main/java/x/x/@appname@", workdir+"/app/src/main/java/x/x/"+binname)
	os.Rename(workdir+"/app/src/test/java/x/x/@appname@", workdir+"/app/src/test/java/x/x/"+binname)
	log.Println("Copying binaries.")
	if GOARCH == "amd64" || GOARCH == "all" {
		copyFile(bindir+"/"+binname+"_android_amd64", workdir+"/app/src/main/jniLibs/x86_64/libbin.so")
		copyFile(bindir+"/"+binname+"_android_amd64", workdir+"/app/src/main/resources/lib/x86_64/libbin.so")
	}
	if GOARCH == "386" || GOARCH == "all" {
		copyFile(bindir+"/"+binname+"_android_386", workdir+"/app/src/main/jniLibs/x86/libbin.so")
		copyFile(bindir+"/"+binname+"_android_386", workdir+"/app/src/main/resources/lib/x86/libbin.so")
	}
	if GOARCH == "arm" || GOARCH == "all" {
		copyFile(bindir+"/"+binname+"_android_arm", workdir+"/app/src/main/jniLibs/armeabi-v7a/libbin.so")
		copyFile(bindir+"/"+binname+"_android_arm", workdir+"/app/src/main/resources/lib/armeabi-v7a/libbin.so")
	}
	if GOARCH == "arm64" || GOARCH == "all" {
		copyFile(bindir+"/"+binname+"_android_arm64", workdir+"/app/src/main/jniLibs/arm64-v8a/libbin.so")
		copyFile(bindir+"/"+binname+"_android_arm64", workdir+"/app/src/main/resources/lib/arm64-v8a/libbin.so")
	}
	log.Println("Building yay.")
	wd, err := os.Getwd()
	if err != nil {
		log.Println("err", err)
		os.Exit(2)
	}
	os.Chdir(workdir)
	cmd := exec.Command(workdir+"/gradlew", "build")
	cmd.Env = append(cmd.Env, "ANDROID_SDK_ROOT="+sdkpath)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		log.Println("out:", outbuf.String(), "err:", errbuf.String())
		os.Exit(2)
	}
	os.Chdir(wd)
	log.Println("Copying target app..")
	copyFile(workdir+"/app/build/outputs/apk/debug/app-debug.apk", apkdir+"/"+binname+"."+GOARCH+".apk")
}

func copyFile(src string, dst string) {
	fmt.Println("copyFile", src, "->", dst)
	input, err := ioutil.ReadFile(src)
	if err != nil {
		log.Println(err)
		return
	}

	err = ioutil.WriteFile(dst, input, 0777)
	if err != nil {
		log.Println("Error creating", dst)
		log.Println(err)
		return
	}
}

func fileReplace(filepath string, search string, replace string) {
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println(err)
		return
	}

	target := strings.ReplaceAll(string(input), search, replace)

	ioutil.WriteFile(filepath, []byte(target), 0750)
}

func copyDir(source string, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {
			var data, err1 = ioutil.ReadFile(filepath.Join(source, relPath))
			if err1 != nil {
				return err1
			}
			return ioutil.WriteFile(filepath.Join(destination, relPath), data, 0777)
		}
	})
	return err
}

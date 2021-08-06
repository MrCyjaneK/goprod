// +build linux

package compiler

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func Build(combo string, tags string, binname string, builddir string, ndk string, ldflags string, buildcmd []string) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	spl := strings.Split(combo, "/")
	if len(spl) != 2 {
		log.Fatal("Invalid " + combo + " provided.")
	}
	GOOS := spl[0]
	GOARCH := spl[1]
	log.Println("Building for:", GOOS, GOARCH)
	appendix := ""
	if GOOS == "windows" {
		appendix = ".exe"
	}
	var cmd *exec.Cmd
	if len(buildcmd) == 0 {
		cmd = exec.Command("go", "build", "--ldflags", ldflags, "-o", builddir+"/"+binname+"_"+strings.ToLower(GOOS)+"_"+strings.ToLower(GOARCH)+appendix, "-tags="+tags)
	} else {
		cmd = exec.Command(buildcmd[0])
		if len(buildcmd) > 1 {
			cmd.Args = buildcmd
		}
	}
	// Import env
	out, err := exec.Command("printenv").Output()
	if err != nil {
		log.Fatal(err)
	}
	gopath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		log.Fatal(err)
	}
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		log.Fatal(err)
	}
	envs := strings.Split(string(out), "\n")
	for i := range envs {
		cmd.Env = append(cmd.Env, strings.ReplaceAll(envs[i], "\n", ""))
	}
	cmd.Env = append(cmd.Env, "GOPATH="+strings.ReplaceAll(string(gopath), "\n", ""))
	cmd.Env = append(cmd.Env, "GOROOT="+strings.ReplaceAll(string(goroot), "\n", ""))
	cmd.Env = append(cmd.Env, "GOOS="+GOOS)
	cmd.Env = append(cmd.Env, "GOARCH="+GOARCH)
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	switch GOOS {
	case "android":
		{
			arch := GOARCH
			switch GOARCH {
			case "arm64":
				{
					arch = "aarch64"
					cmd.Env = append(cmd.Env, "CC="+ndk+"/"+arch+"-linux-android21-clang")
					cmd.Env = append(cmd.Env, "CXX="+ndk+"/"+arch+"-linux-android21-clang++")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			case "arm":
				{
					arch = "armv7a"
					cmd.Env = append(cmd.Env, "CC="+ndk+"/"+arch+"-linux-androideabi21-clang")
					cmd.Env = append(cmd.Env, "CXX="+ndk+"/"+arch+"-linux-androideabi21-clang++")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			case "amd64":
				{
					arch = "x86_64"
					cmd.Env = append(cmd.Env, "CC="+ndk+"/"+arch+"-linux-android21-clang")
					cmd.Env = append(cmd.Env, "CXX="+ndk+"/"+arch+"-linux-android21-clang++")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			case "386":
				{
					arch = "i686"
					cmd.Env = append(cmd.Env, "CC="+ndk+"/"+arch+"-linux-android21-clang")
					cmd.Env = append(cmd.Env, "CXX="+ndk+"/"+arch+"-linux-android21-clang++")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			}
		}
	case "linux":
		{
			switch GOARCH {
			case "arm64":
				{
					cmd.Env = append(cmd.Env, "CC=aarch64-linux-gnu-gcc")
					cmd.Env = append(cmd.Env, "CXX=aarch64-linux-gnu-g++")
					cmd.Env = append(cmd.Env, "HOST=aarch64-linux-gnu")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			case "arm":
				{
					cmd.Env = append(cmd.Env, "CC=arm-linux-gnueabihf-gcc")
					cmd.Env = append(cmd.Env, "CXX=arm-linux-gnueabihf-g++")
					cmd.Env = append(cmd.Env, "HOST=arm-linux-gnueabihf")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			case "386":
				{
					cmd.Env = append(cmd.Env, "CC=i686-linux-gnu-gcc")
					cmd.Env = append(cmd.Env, "CXX=i686-linux-gnu-g++")
					cmd.Env = append(cmd.Env, "HOST=i686-linux-gnu")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			case "amd64":
				{
					cmd.Env = append(cmd.Env, "CC=x86_64-linux-gnu-gcc")
					cmd.Env = append(cmd.Env, "CXX=x86_64-linux-gnu-g++")
					cmd.Env = append(cmd.Env, "HOST=x86_64-linux-gnu")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			}
		}
	case "windows":
		{
			switch GOOS {
			case "amd64":
				{
					cmd.Env = append(cmd.Env, "CC=x86_64-w64-mingw32-gcc")
					cmd.Env = append(cmd.Env, "CXX=x86_64-w64-mingw32-g++")
					cmd.Env = append(cmd.Env, "HOST=x86_64-w64-mingw32")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			case "386":
				{
					cmd.Env = append(cmd.Env, "CC=i686-w64-mingw32-gcc")
					cmd.Env = append(cmd.Env, "CXX=i686-w64-mingw32-g++")
					cmd.Env = append(cmd.Env, "HOST=x86_64-w64-mingw32")
					cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
				}
			}
		}
	}
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	log.Println("out:", outbuf.String(), "err:", errbuf.String())
	if err != nil {
		log.Fatal("err:", err)
	}
}

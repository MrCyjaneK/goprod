// +build linux

package debpackage

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Build(combo string, binname string, bindir string, debdir string, version string) {
	spl := strings.Split(combo, "/")
	if len(spl) != 2 {
		log.Fatal("Invalid " + combo + " provided.")
	}
	GOOS := spl[0]
	GOARCH := spl[1]
	if GOOS != "linux" {
		log.Fatal("no")
	}
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
	email, err := exec.Command("git", "config", "user.email").Output()
	if err != nil {
		log.Fatal(err)
	}
	emails := string(email)
	emails = strings.Split(emails, "\n")[0]
	envs := strings.Split(string(out), "\n")
	verout, err := exec.Command("git", "log", "-n", "1").Output()
	if err != nil {
		log.Fatal(err)
	}
	commit := strings.Split(string(verout), " ")[1][0:8]
	arch := "amd64"
	switch GOARCH {
	case "arm":
		{
			arch = "armhf"
		}
	case "386":
		{
			arch = "i386"
		}
	}
	cmd := exec.Command("checkinstall",
		"--type=debian",
		"--install=no",
		"--default",
		"--pkgname="+binname,
		"--pkgversion="+version+"+git"+commit,
		"--arch="+arch,
		"--pakdir="+debdir,
		"--maintainer="+emails,
		"--provides="+binname,
		"--strip=no",
		"--stripso=no")

	for i := range envs {
		cmd.Env = append(cmd.Env, strings.ReplaceAll(envs[i], "\n", ""))
	}
	cmd.Env = append(cmd.Env, "GOPATH="+strings.ReplaceAll(string(gopath), "\n", ""))
	cmd.Env = append(cmd.Env, "GOROOT="+strings.ReplaceAll(string(goroot), "\n", ""))
	cmd.Env = append(cmd.Env, "GOOS="+GOOS)
	cmd.Env = append(cmd.Env, "GOARCH="+GOARCH)
	cmd.Env = append(cmd.Env, "BINNAME="+binname)
	if _, err := os.Stat("Makefile"); err != nil {
		err = ioutil.WriteFile("Makefile", []byte("install:\n"+
			"\tcp build/bin/${BINNAME}_${GOOS}_${GOARCH} /usr/bin/${BINNAME}\n"), 0750)
		if err != nil {
			log.Fatal(err)
		}
	}
	//t := time.Now()
	//log.Fatal("\n", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		log.Fatal("out:", outbuf.String(), "err:", errbuf.String())
	}

}

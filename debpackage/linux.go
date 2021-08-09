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
	log.SetFlags(log.LstdFlags | log.Llongfile)
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
	envs := strings.Split(string(out), "\n")
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
	//TODO: remove `exec` and use go-git?
	email, err := exec.Command("git", "config", "user.email").Output()
	if err != nil {
		log.Println("WARN:", err, "please do `git config user.email \"your@email\"' to configure email address.")
		email = []byte(os.Getenv("EMAIL"))
		if string(email) == "" {
			email = []byte("no-reply@example.com")
		}
	}
	emails := string(email)
	emails = strings.Split(emails, "\n")[0]
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
	case "arm64":
		{
			arch = "arm64"
		}
	}
	cmd := exec.Command("checkinstall",
		"--type=debian",
		"--install=no",
		"--default",
		"--pkgname="+binname,
		"--pkgversion="+version,
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

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		log.Fatal("out:", outbuf.String(), "err:", errbuf.String())
	}

}

package appimage

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Package(buildcmd []string, ldflags string, builddir string, binname string, GOOS string, GOARCH string, tags string, version string, outdir string) {
	// Import env

	var cmd *exec.Cmd
	if len(buildcmd) == 0 {
		cmd = exec.Command("make", "install")
	} else {
		cmd = exec.Command(buildcmd[0])
		if len(buildcmd) > 1 {
			cmd.Args = buildcmd
		}
	}
	out, err := exec.Command("printenv").Output()
	if err != nil {
		log.Fatal(err)
	}
	envs := strings.Split(string(out), "\n")
	for i := range envs {
		cmd.Env = append(cmd.Env, strings.ReplaceAll(envs[i], "\n", ""))
	}
	cmd.Env = append(cmd.Env, "GOOS="+GOOS)
	cmd.Env = append(cmd.Env, "GOARCH="+GOARCH)
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	cmd.Env = append(cmd.Env, "BINNAME="+binname)
	cmd.Env = append(cmd.Env, "DESTDIR="+builddir)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		log.Println("out:", outbuf.String(), "err:", errbuf.String())
		log.Fatal("err:", err)
	}

	var acmd *exec.Cmd
	df := getFile(builddir, ".desktop")
	acmd = exec.Command("appimagetool", "deploy", df)
	for i := range envs {
		acmd.Env = append(acmd.Env, strings.ReplaceAll(envs[i], "\n", ""))
	}
	acmd.Env = append(acmd.Env, "GOOS="+GOOS)
	acmd.Env = append(acmd.Env, "GOARCH="+GOARCH)
	acmd.Env = append(acmd.Env, "ARCH="+GOARCH)
	acmd.Env = append(acmd.Env, "VERSION="+version)
	acmd.Env = append(acmd.Env, "CGO_ENABLED=0")
	acmd.Env = append(acmd.Env, "BINNAME="+binname)
	acmd.Env = append(acmd.Env, "DESTDIR="+builddir)
	var boutbuf, berrbuf bytes.Buffer
	acmd.Stdout = &boutbuf
	acmd.Stderr = &berrbuf
	err = acmd.Run()
	if err != nil {
		log.Println("out:", boutbuf.String(), "err:", berrbuf.String())
		log.Fatal("err:", err)
	}
	filedesktop, err := ioutil.ReadFile(df)
	if err != nil {
		log.Fatal("Unable to read: '", df, "'")
	}
	icons := grep("Icon=", strings.Split(string(filedesktop), "\n"))
	icon := strings.Split(icons[0], "Icon=")[1]
	filec, err := ioutil.ReadFile(getFile(builddir, icon+".png"))
	ext := "png"
	if err != nil {
		filec, err = ioutil.ReadFile(getFile(builddir, icon+".svg"))
		ext = "svg"
		if err != nil {
			filec, err = ioutil.ReadFile(getFile(builddir, icon+".xpm"))
			ext = "xpm"
			if err != nil {
				log.Fatal("Unable to find ", icon, "with suffix .png, .svg, .xpm")
			}
		}
	}
	err = ioutil.WriteFile(builddir+"/"+icon+"."+ext, filec, 0750)
	if err != nil {
		log.Fatal("Unable to write icon to", builddir+"/"+icon+"."+ext)
	}

	acmd = exec.Command("appimagetool", builddir)
	for i := range envs {
		acmd.Env = append(acmd.Env, strings.ReplaceAll(envs[i], "\n", ""))
	}
	acmd.Env = append(acmd.Env, "GOOS="+GOOS)
	acmd.Env = append(acmd.Env, "GOARCH="+GOARCH)
	var arch = "unknown"
	if GOARCH == "amd64" {
		acmd.Env = append(acmd.Env, "ARCH=x86_64")
		arch = "x86_64"
	} else if GOARCH == "386" {
		acmd.Env = append(acmd.Env, "ARCH=i686")
		arch = "i686"
	} else if GOARCH == "arm64" {
		acmd.Env = append(acmd.Env, "ARCH=aarch64")
		arch = "aarch64"
	} else if GOARCH == "arm" {
		acmd.Env = append(acmd.Env, "ARCH=armhf")
		arch = "armhf"
	}
	acmd.Env = append(acmd.Env, "VERSION="+version)
	acmd.Env = append(acmd.Env, "CGO_ENABLED=0")
	acmd.Env = append(acmd.Env, "BINNAME="+binname)
	acmd.Env = append(acmd.Env, "DESTDIR="+builddir)
	var aoutbuf, aerrbuf bytes.Buffer
	acmd.Stdout = &aoutbuf
	acmd.Stderr = &aerrbuf
	err = acmd.Run()
	if err != nil {
		log.Println("out:", aoutbuf.String(), "err:", aerrbuf.String())
		log.Fatal("err:", err)
	}
	defer os.RemoveAll(builddir)
	//fmt.Println(string(filec))
	name := strings.Split(grep("Name=", strings.Split(string(filedesktop), "\n"))[0], "Name=")[1] + "-" + version + "-" + arch + ".AppImage"
	log.Println("AppImage:", name)
	os.MkdirAll(outdir, 0750)
	source, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()
	destination, err := os.Create(outdir + "/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	defer os.Remove(name)
	if err != nil {
		log.Fatal(err)
	}
}

func getFile(root string, endwith string) string {
	var file string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(path, endwith) {
			log.Println("Using", path, "as '", endwith, "' file")
			file = path
			return nil
		}
		return nil
	})
	if err != nil || file == "" {
		log.Fatal("Unable to find any usable .desktop files!", err)
	}
	return file
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

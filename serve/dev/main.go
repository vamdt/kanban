package dev

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type dev struct {
	https      bool
	webpack    bool
	webpackCmd *exec.Cmd
	port       int
}

var Dev *dev = &dev{}

func init() {
	flag.BoolVar(&Dev.webpack, "webpack", true, "start webpack")
}

func dev_static_handle(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	if Dev.webpack && strings.HasPrefix(upath, "/main.js") {
		uri := "http://localhost:"
		if Dev.https {
			uri = "https://localhost:"
		}
		uri = uri + strconv.Itoa(Dev.port) + "/main.js"
		http.Redirect(w, r, uri, 307)
		return
	}

	if strings.HasPrefix(upath, "/bower_components") {
		http.ServeFile(w, r, upath[1:])
		return
	}

	if strings.HasSuffix(upath, "/") {
		upath = upath + "index.html"
	}

	rpath := path.Join("app", upath)
	if _, err := os.Stat(rpath); err == nil {
		http.ServeFile(w, r, rpath)
		return
	}

	rpath = path.Join(".tmp", upath)
	if _, err := os.Stat(rpath); err == nil {
		http.ServeFile(w, r, rpath)
		return
	}

	http.FileServer(http.Dir("dist/static")).ServeHTTP(w, r)
}

func (p *dev) Start(https bool, port string) {
	if p.webpack && p.webpackCmd == nil {
		if iport, err := strconv.Atoi(port[1:]); err == nil {
			p.port = iport + 1
		} else {
			p.port = 31234
		}
		p.https = https
		args := []string{
			"--hot",
			"--inline",
		}
		prot := "http"
		if https {
			args = append(args, "--https")
			prot = "https"
		}
		args = append(args, "--port="+strconv.Itoa(p.port))
		args = append(args, "--output-public-path="+prot+"://localhost"+port+"/")
		args = append(args, "--content-base="+prot+"://localhost"+port+"/")

		p.webpackCmd = exec.Command("webpack-dev-server", args...)
		p.webpackCmd.Stdout = os.Stdout
		p.webpackCmd.Stderr = os.Stderr
		err := p.webpackCmd.Start()
		if err != nil {
			glog.Warningln(err)
		}
	}
}

func (p *dev) Exit() {
	if p.webpackCmd != nil {
		p.webpackCmd.Process.Kill()
		p.webpackCmd = nil
	}
}

func Start(https bool, port string) {
	Dev.Start(https, port)
}

func Exit() {
	Dev.Exit()
}

func (p *dev) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	dev_static_handle(w, req)
}

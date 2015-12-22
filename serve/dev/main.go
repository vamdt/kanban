package dev

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/golang/glog"
)

type dev struct {
	webpack    bool
	webpackCmd *exec.Cmd
}

var Dev *dev = &dev{}

func init() {
	flag.BoolVar(&Dev.webpack, "webpack", true, "start webpack")
}

func dev_static_handle(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	if strings.HasPrefix(upath, "/main.js") {
		http.Redirect(w, r, "https://localhost:9001/main.js", 307)
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

func (p *dev) Start() {
	if p.webpack && p.webpackCmd == nil {
		p.webpackCmd = exec.Command("webpack-dev-server", "--hot --inline", "--output-public-path=https://localhost:3000/", "--content-base=https://localhost:3000/")
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

func Start() {
	Dev.Start()
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

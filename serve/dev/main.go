package dev

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

type Dev struct {
}

func NewDev() *Dev { return &Dev{} }

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

func (p *Dev) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	dev_static_handle(w, req)
}

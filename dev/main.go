package dev

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

type dev struct {
	https      bool
	webpack    bool
	webpackCmd *exec.Cmd
	host       string
	port       int
	open       bool
}

var Dev *dev = &dev{}

func init() {
	flag.BoolVar(&Dev.webpack, "webpack", true, "start webpack in dev mode")
	flag.BoolVar(&Dev.open, "open", false, "open browser in dev mode")
	flag.StringVar(&Dev.host, "host", "", "webpack listen host")
}

func (p *dev) webpackHost(r *http.Request) string {
	host := r.Host
	if i := strings.LastIndex(host, ":"); i > 0 {
		host = host[:i]
	}
	host = host + ":" + strconv.Itoa(p.port)
	uri := "http://"
	if p.https {
		uri = "https://"
	}
	return uri + host
}

func (p *dev) RedirectWebpack(w http.ResponseWriter, r *http.Request) {
	uri = p.webpackHost(r) + r.URL.Path
	http.Redirect(w, r, uri, 302)
}

func dev_static_handle(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	if Dev.webpack && strings.HasPrefix(upath, "/main.js") {
		Dev.RedirectWebpack(w, r)
		return
	}

	if Dev.webpack && strings.Contains(upath, ".hot-update.") {
		Dev.RedirectWebpack(w, r)
		return
	}

	if strings.HasPrefix(upath, "/bower_components") {
		http.ServeFile(w, r, upath[1:])
		return
	}

	paths := []string{"app", ".tmp"}
	for _, v := range paths {
		rpath := path.Join(v, upath)
		if upath == "/" {
			rpath = rpath + "/index.html"
		}
		if _, err := os.Stat(rpath); err == nil {
			if upath == "/" {
				if content, err := ioutil.ReadFile(rpath); err == nil {
					webpack := Dev.webpackHost(r) + "/main.js"
					content = bytes.Replace(content, []byte("src=\"main.js\""), []byte("src=\""+webpack+"\""), 1)
					w.Header().Set("Content-Type", "text/html")
					w.Write(content)
					return
				}
			}
			http.ServeFile(w, r, rpath)
			return
		}
	}

	http.FileServer(http.Dir("dist/static")).ServeHTTP(w, r)
}

func (p *dev) Start(https bool, port string) {
	if iport, err := strconv.Atoi(port[1:]); err == nil {
		p.port = iport + 1
	} else {
		p.port = 31234
	}
	p.https = https
	prot := "http"
	if https {
		prot = "https"
	}

	serve_uri := prot + "://"
	if len(p.host) > 0 {
		serve_uri = serve_uri + p.host
	} else {
		serve_uri = serve_uri + "localhost"
	}
	serve_uri = serve_uri + port + "/"

	if p.webpack && p.webpackCmd == nil {
		args := []string{
			"--hot",
			"--inline",
		}
		if https {
			args = append(args, "--https")
		}
		if len(p.host) > 0 {
			args = append(args, "--host="+p.host)
		}
		args = append(args, "--port="+strconv.Itoa(p.port))
		args = append(args, "--output-public-path="+serve_uri)
		args = append(args, "--content-base="+serve_uri)

		p.webpackCmd = exec.Command("webpack-dev-server", args...)
		p.webpackCmd.Stdout = os.Stdout
		p.webpackCmd.Stderr = os.Stderr
		err := p.webpackCmd.Start()
		if err != nil {
			glog.Warningln(err)
		}
	}

	if p.open {
		go func() {
			time.Sleep(time.Second)
			glog.Infoln("open", serve_uri)
			err := exec.Command("open", serve_uri).Start()
			if err != nil {
				glog.Warning(err)
			}
		}()
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

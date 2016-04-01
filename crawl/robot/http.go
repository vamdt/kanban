package robot

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const maxRetry int = 5

var UA string = "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"

func init() {
	ua := os.Getenv("User_Agent")
	if len(ua) > 0 {
		UA = ua
	}
}

func guess_charset(contype string) string {
	excel := "application/vnd.ms-excel"
	if contype == excel {
		return "GBK"
	}

	charset := strings.ToUpper(contype)
	if i := strings.Index(charset, "CHARSET="); i == -1 {
		return "UTF8"
	} else {
		charset = charset[i+len("CHARSET="):]
		charset = strings.TrimSpace(charset)
		charset = strings.Trim(charset, ";")
	}
	return charset
}

func Http_get(url string, referer *string, tout time.Duration) (body []byte, err error) {
	glog.V(HttpV).Infoln(url)
	var res *http.Response
	for i := 0; i < maxRetry; i++ {
		client := &http.Client{Timeout: tout}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Add("Connection", "keep-alive")
		req.Header.Set("User-Agent", UA)
		if referer != nil {
			req.Header.Set("Referer", *referer)
		}

		res, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		glog.Warningln("http get fail", url, err)
		return
	}
	if res == nil && err == nil {
		err = errors.New("req " + url + " fail")
		return
	}

	defer res.Body.Close()

	contype := res.Header.Get("Content-Type")
	charset := guess_charset(contype)

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
		defer reader.Close()
	default:
		reader = res.Body
	}

	if charset[:2] == "GB" {
		body, err = ioutil.ReadAll(transform.NewReader(reader,
			simplifiedchinese.GBK.NewDecoder()))
	} else {
		body, err = ioutil.ReadAll(reader)
	}

	return body, err
}

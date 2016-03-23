package robot

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
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

func Http_get(url string, referer *string) (res *http.Response, err error) {
	glog.V(HttpV).Infoln(url)
	for i := 0; i < maxRetry; i++ {
		client := &http.Client{Timeout: time.Second * 5}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

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
	}
	if res == nil && err == nil {
		err = errors.New("req " + url + " fail")
	}
	return
}

func Http_get_raw(url string, referer *string) ([]byte, error) {
	resp, err := Http_get(url, referer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Http_get_gbk(url string, referer *string) ([]byte, error) {
	resp, err := Http_get(url, referer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(transform.NewReader(resp.Body,
		simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Download(url string) []byte {
	body, err := Http_get_raw(url, nil)
	if err != nil {
		glog.Warningln("Download fail", url, err)
		return nil
	}
	return body
}

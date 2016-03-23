package robot

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Http_get(url string, referer *string) (res *http.Response, err error) {
	glog.V(HttpV).Infoln(url)
	for i := 0; i < 2; i++ {
		client := &http.Client{Timeout: time.Second * 5}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9")
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
		err = errors.New("req fail")
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

package crawl

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Http_get(url string, referer *string) (*http.Response, error) {
	glog.V(HttpV).Infoln("http get", url)
	client := &http.Client{Timeout: time.Duration(3 * time.Second)}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Warningln("http get fail", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9")
	if referer != nil {
		req.Header.Set("Referer", *referer)
	}

	return client.Do(req)
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
		glog.Warningln("Download fail", err)
		return nil
	}
	return body
}

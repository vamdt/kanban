package crawl

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Http_get(url string, referer *string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.155 Safari/537.36")
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

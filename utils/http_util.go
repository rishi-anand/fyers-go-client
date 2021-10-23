package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

const (
	GET  = "GET"
	POST = "POST"
)

func DoHttpCall(method, url string, body []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		log.Error("Failed to get new http request", err)
		return nil, err
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	reqOut, _ := httputil.DumpRequest(req, true)
	log.Debugln(string(reqOut))
	resp, err := client.Do(req)
	respOut, _ := httputil.DumpResponse(resp, true)
	log.Debugln(string(respOut))
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

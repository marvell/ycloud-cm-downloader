package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

func httpDo(req *http.Request) ([]byte, error) {
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if debugModeFlag {
			reqDump, _ := httputil.DumpRequestOut(req, true)
			dbg("request dump:\n%s", reqDump)

			resDump, _ := httputil.DumpResponse(res, true)
			dbg("response dump:\n%s", resDump)
		}

		return nil, fmt.Errorf("invalid response status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

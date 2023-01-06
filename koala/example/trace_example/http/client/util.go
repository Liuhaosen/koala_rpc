package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpDo(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("statuscode %d, body : %s", resp.StatusCode, body)
	}
	return body, nil
}

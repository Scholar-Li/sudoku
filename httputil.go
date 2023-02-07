package sudoku

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(url string, body interface{}, header map[string]string) (int, []byte, error) {
	return Do(http.MethodPost, url, body, header, make(map[string]interface{}))
}

func Get(url string, query map[string]interface{}, header map[string]string) (int, []byte, error) {
	return Do(http.MethodGet, url, nil, header, query)
}

func Do(method, url string, body interface{}, header map[string]string, query map[string]interface{}) (int, []byte, error) {
	client := &http.Client{}

	byteParams, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	var req *http.Request
	req, err = http.NewRequest(method, url, strings.NewReader(string(byteParams)))
	if err != nil {
		return 0, nil, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = q.Encode()

	req.Close = true

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()

	repBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, repBody, err
}

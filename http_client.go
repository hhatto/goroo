package goroo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpClient struct {
	host string
}

func (h *httpClient) Call(command string, params map[string]string) (*GroongaResult, error) {
	rawurl := fmt.Sprintf("%s://%s", "http", h.host)
	body, err := runCommand(rawurl, command, params)
	if err != nil {
		return nil, err
	}
	return parseResult(body)
}

func runCommand(rawurl, command string, params map[string]string) ([]byte, error) {
	v := url.Values{}
	for value, name := range params {
		v.Set(value, name)
	}
	requestUrl := fmt.Sprintf("%s/d/%s?%s", rawurl, command, v.Encode())
	resp, err := http.Get(requestUrl)
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

func parseResult(body []byte) (*GroongaResult, error) {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	grnInfo := data.([]interface{})
	grnHeader := grnInfo[0].([]interface{})
	result := new(GroongaResult)
	result.RawData = string(body)
	result.Status = int(grnHeader[0].(float64))
	result.StartTime = grnHeader[1].(float64)
	result.ElapsedTime = grnHeader[2].(float64)
	result.Body = grnInfo[1]
	if result.Status != 0 {
		return result, fmt.Errorf("%d - %s", result.Status, grnHeader[3])
	}

	return result, nil
}

func newHttpClient(host string) Client {
	return &httpClient{host}
}

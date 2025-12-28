package goroo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type httpClient struct {
	raw_url string
}

func (h *httpClient) Call(command string, params map[string]string) (*GroongaResult, error) {
	body, err := h.run(h.raw_url, command, params)
	if err != nil {
		return nil, err
	}
	return h.parse(body)
}

func (h *httpClient) run(rawurl, command string, params map[string]string) ([]byte, error) {
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

	return io.ReadAll(resp.Body)
}

func (h *httpClient) parse(body []byte) (*GroongaResult, error) {
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	grnInfo := data.([]any)
	grnHeader := grnInfo[0].([]any)
	result := new(GroongaResult)
	result.RawData = string(body)
	result.Status = int(grnHeader[0].(float64))
	result.StartTime = grnHeader[1].(float64)
	result.ElapsedTime = grnHeader[2].(float64)
	if len(grnInfo) == 2 {
		result.Body = grnInfo[1]
	}
	if result.Status != 0 {
		return result, fmt.Errorf("%d - %s", result.Status, grnHeader[3])
	}

	return result, nil
}

func newHttpClient(host string) Client {
	return &httpClient{host}
}

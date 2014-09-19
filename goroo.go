package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GroongaClient struct {
	Protocol string
	Host     string
	Port     int
}

type GroongaResult struct {
	RawData     string
	Status      int
	StartTime   float64
	ElapsedTime float64
	Body        interface{}
}

func NewGroongaClient(protocol, host string, port int) *GroongaClient {
	client := &GroongaClient{
		Protocol: protocol,
		Host:     host,
		Port:     port,
	}
	return client
}

func (client *GroongaClient) Call(command string, params map[string]string) (result GroongaResult, err error) {
	if len(params) == 0 {
		return result, nil
	}

	v := url.Values{}
	for value, name := range params {
		v.Set(value, name)
	}
	requestUrl := fmt.Sprintf("%s://%s:%d/d/%s?%s",
		client.Protocol, client.Host, client.Port, command, v.Encode())
	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Errorf("http.Get() error: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("response read error: %v", err)
	}

	result, err = setResult(body)
	if err != nil {
		return result, err
	}

	return result, nil
}

func setResult(body []byte) (result GroongaResult, err error) {
	result.RawData = fmt.Sprintf("%s", body)

	var data interface{}
	dec := json.NewDecoder(strings.NewReader(result.RawData))
	dec.Decode(&data)

	grnInfo := data.([]interface{})
	grnHeader := grnInfo[0].([]interface{})
	result.Status = int(grnHeader[0].(float64))
	result.StartTime = grnHeader[1].(float64)
	result.ElapsedTime = grnHeader[2].(float64)
	if len(grnHeader) == 3 {
		// groonga response ok
		result.Body = grnInfo[1]
	} else {
		// groonga response ng
		result.Body = grnHeader[3]
	}

	return result, nil
}

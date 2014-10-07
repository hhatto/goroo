package goroo

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func TestHTTPClient(t *testing.T) {
	doGet = func(string) (*http.Response, error) {
		const body = "[[-22,1412056029.84683,0.000826835632324219,\"already used name was assigned: <Users>\",[[\"grn_obj_register\",\"db.c\",7608]]],false]"
		br := bufio.NewReader(strings.NewReader("HTTP/1.1 200 OK\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
			"\r\n" +
			body))
		resp, _ := http.ReadResponse(br, &http.Request{Method: "GET"})
		return resp, nil
	}

	//client := GroongaClient{Protocol: "http", Host: "localhost", Port: 10041}
	client := NewGroongaClient("http", "localhost", 10041)

	params := map[string]string{
		"name":     "Users",
		"flags":    "TABLE_HASH_KEY",
		"key_type": "ShortText",
	}
	result, _ := client.Call("table_create", params)

	params = map[string]string{
		"table": "Users",
		"name":  "name",
		"flags": "COLUMN_SCALAR",
		"type":  "ShortText",
	}
	result, _ = client.Call("column_create", params)

	doGet = func(string) (*http.Response, error) {
		const body = "[[0,1412056029.84987,9.60826873779297e-05],2]"
		br := bufio.NewReader(strings.NewReader("HTTP/1.1 200 OK\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
			"\r\n" +
			body))
		resp, _ := http.ReadResponse(br, &http.Request{Method: "GET"})
		return resp, nil
	}
	params = map[string]string{
		"table":  "Users",
		"values": "[{\"_key\":\"ken\",\"name\":\"Ken\"},{\"_key\":\"jim\",\"name\":\"Jim\"}]",
	}
	result, _ = client.Call("load", params)

	doGet = func(string) (*http.Response, error) {
		const body = "[[0,1412056029.8505,0.000298976898193359],[[[0],[[\"_id\",\"UInt32\"],[\"_key\",\"ShortText\"],[\"name\",\"ShortText\"]]]]]"
		br := bufio.NewReader(strings.NewReader("HTTP/1.1 200 OK\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
			"\r\n" +
			body))
		resp, _ := http.ReadResponse(br, &http.Request{Method: "GET"})
		return resp, nil
	}
	params = map[string]string{
		"table": "Users",
		"query": "name:@test",
	}
	result, _ = client.Call("select", params)
	if len(result.RawData) == 0 {
		t.Errorf("response body not found")
	}
}

func TestGQTPClientConnectError(t *testing.T) {
	client := NewGroongaClient("gqtp", "localhost", 1)
	params := map[string]string{
		"table": "Users",
		"query": "message:@test",
	}
	_, err := client.Call("select", params)
	if err == nil {
		t.Errorf("invalid sequence")
	}
}

func TestGQTPClientError(t *testing.T) {
	t.Skip("TODO: use mock")
	client := NewGroongaClient("gqtp", "localhost", 10043)
	params := map[string]string{
		"table": "Users",
		"query": "message:@test",
	}
	_, err := client.Call("select", params)
	if err == nil {
		t.Errorf("invalid sequence")
	}
}
func TestGQTPClient(t *testing.T) {
	t.Skip("TODO: use mock")
	client := NewGroongaClient("gqtp", "localhost", 10043)
	params := map[string]string{
		"table": "Users",
		"query": "name:@Jim",
	}
	result, _ := client.Call("select", params)
	if len(result.RawData) == 0 {
		t.Errorf("response body not found")
	}
}

// Benchmarks
func BenchmarkHTTPClient(b *testing.B) {
	client := NewGroongaClient("http", "localhost", 10041)
	params := map[string]string{
		"table": "Users",
		"query": "name:@test",
	}
	for n := 0; n < b.N; n++ {
		result, _ := client.Call("select", params)
		if len(result.RawData) == 0 {
			b.Errorf("response body not found")
		}
	}
}

func BenchmarkGQTPClient(b *testing.B) {
	client := NewGroongaClient("gqtp", "localhost", 10043)
	params := map[string]string{
		"table": "Users",
		"query": "name:@Jim",
	}
	for n := 0; n < b.N; n++ {
		result, _ := client.Call("select", params)
		if len(result.RawData) == 0 {
			b.Errorf("response body not found")
		}
	}
}

package goroo

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func TestHttp_TableList_Empty_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const body = "[[0,1444022807.258,0.0],[[[\"id\",\"UInt32\"],[\"name\",\"ShortText\"],[\"path\",\"ShortText\"],[\"flags\",\"ShortText\"],[\"domain\",\"ShortText\"],[\"range\",\"ShortText\"],[\"default_tokenizer\",\"ShortText\"],[\"normalizer\",\"ShortText\"]]]]"
		fmt.Fprintln(w, body)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}
	client := NewHttpClient(u.Host)
	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Error(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]interface{})) != 1 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestHttp_TableList_Count1_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const body = "[[0,1444024497.318,0.0469999313354492],[[[\"id\",\"UInt32\"],[\"name\",\"ShortText\"],[\"path\",\"ShortText\"],[\"flags\",\"ShortText\"],[\"domain\",\"ShortText\"],[\"range\",\"ShortText\"],[\"default_tokenizer\",\"ShortText\"],[\"normalizer\",\"ShortText\"]],[256,\"TestGQTPClinet\",\"./markdown.db.0000100\",\"TABLE_HASH_KEY|PERSISTENT\",null,null,null,null]]]"
		fmt.Fprintln(w, body)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}
	client := NewHttpClient(u.Host)
	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Error(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]interface{})) != 2 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestHTTPClient(t *testing.T) {
	t.SkipNow()
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

func Test_HttpClient_TableListCommnad(t *testing.T) {
	t.SkipNow()
	client := NewGroongaClient("http", "localhost", 10041)
	result, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Errorf("response body not found")
	}
	if result.Status != 0 {
		t.Errorf("result status not zero.[%d]", result.Status)
	}
	if len(result.RawData) == 0 {
		t.Errorf("response body not found")
	}
}

func Test_GQTPClient_TableCreateCommnad(t *testing.T) {
	t.SkipNow()
	client := NewGroongaClient("gqtp", "localhost", 10043)
	result, err := client.Call("table_create", map[string]string{
		"name": "TestGQTPClinet",
	})
	if err != nil {
		t.Errorf("err is not nil [%s]", err)
	}
	if result.Status != 0 {
		t.Errorf("result status not zero.[%d]", result.Status)
	}
}

func Test_GQTPClient_TableListCommnad(t *testing.T) {
	t.SkipNow()
	client := NewGroongaClient("gqtp", "localhost", 10043)
	result, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Errorf("err is not nil [%s]", err)
	}
	if result.Status != 0 {
		t.Errorf("result status not zero.[%d]", result.Status)
	}
	if len(result.RawData) == 0 {
		t.Errorf("response body not found")
	}
}

func Test_GQTPClient_TableRemoveCommnad(t *testing.T) {
	t.SkipNow()
	client := NewGroongaClient("gqtp", "localhost", 10043)
	result, err := client.Call("table_remove", map[string]string{
		"name": "TestGQTPClinet",
	})
	if err != nil {
		t.Errorf("err is not nil [%s]", err)
	}
	if result.Status != 0 {
		t.Errorf("result status not zero.[%d]", result.Status)
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

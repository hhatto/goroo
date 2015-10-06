package goroo

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func Test_Unsupported_Protocol_IsNil(t *testing.T) {
	client := NewClient("json", "localhost", 10041)
	if client != nil {
		t.Error("return is not error.")
	}
}

func Test_ClinetHttp_TableList_Empty_Success(t *testing.T) {
	const body = "[[0,1444022807.258,0.0],[[[\"id\",\"UInt32\"],[\"name\",\"ShortText\"],[\"path\",\"ShortText\"],[\"flags\",\"ShortText\"],[\"domain\",\"ShortText\"],[\"range\",\"ShortText\"],[\"default_tokenizer\",\"ShortText\"],[\"normalizer\",\"ShortText\"]]]]"
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	schema := u.Scheme
	host, p, _ := net.SplitHostPort(u.Host)
	port, _ := strconv.Atoi(p)
	client := NewClient(schema, host, port)

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

func Test_ClinetGqtp_TableList_Empty_Success(t *testing.T) {
	body := []byte{0xc7, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0xbd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5b, 0x5b, 0x5b, 0x22, 0x69, 0x64, 0x22, 0x2c, 0x22, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x70, 0x61, 0x74, 0x68, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x69, 0x7a, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x5d, 0x5d}
	mock := gqtpMock(body)
	defer mock.Close()

	host, p, _ := net.SplitHostPort(mock.Address)
	port, _ := strconv.Atoi(p)
	client := NewClient("gqtp", host, port)
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
	t.SkipNow()
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
	b.SkipNow()
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

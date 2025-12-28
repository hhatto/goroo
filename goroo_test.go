package goroo

import (
	"net"
	"net/url"
	"strconv"

	"testing"
)

var (
	byteBody = []byte{0xc7, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0,
		0xbd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5b,
		0x5b, 0x5b, 0x22, 0x69, 0x64, 0x22, 0x2c, 0x22, 0x55, 0x49, 0x6e, 0x74,
		0x33, 0x32, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22,
		0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22,
		0x5d, 0x2c, 0x5b, 0x22, 0x70, 0x61, 0x74, 0x68, 0x22, 0x2c, 0x22, 0x53,
		0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b,
		0x22, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f,
		0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64,
		0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72,
		0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x72, 0x61,
		0x6e, 0x67, 0x65, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54,
		0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64, 0x65, 0x66, 0x61,
		0x75, 0x6c, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x69, 0x7a, 0x65,
		0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78,
		0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c,
		0x69, 0x7a, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74,
		0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x5d, 0x5d}
)

const (
	cTEST_BODY = `
[
    [
        0, 
        1444022807.258, 
        0.0
    ], 
    [
        [
            [
                "id", 
                "UInt32"
            ], 
            [
                "name", 
                "ShortText"
            ], 
            [
                "path", 
                "ShortText"
            ], 
            [
                "flags", 
                "ShortText"
            ], 
            [
                "domain", 
                "ShortText"
            ], 
            [
                "range", 
                "ShortText"
            ], 
            [
                "default_tokenizer", 
                "ShortText"
            ], 
            [
                "normalizer", 
                "ShortText"
            ]
        ]
    ]
]`
)

func TestUnsupportedProtocolIsNil(t *testing.T) {
	client := NewClient("json", "localhost", 10041)
	if client != nil {
		t.Error("return is not error.")
	}
}

func TestClinetHttpTableListEmptySuccess(t *testing.T) {
	body := cTEST_BODY
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	schema := u.Scheme
	host, p, _ := net.SplitHostPort(u.Host)
	port, _ := strconv.Atoi(p)
	client := NewClient(schema, host, port)

	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]any)) != 1 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestClinetGqtpTableListEmptySuccess(t *testing.T) {
	body := byteBody
	mock := gqtpMock(body)
	defer mock.Close()

	host, p, _ := net.SplitHostPort(mock.Address)
	port, _ := strconv.Atoi(p)
	client := NewClient("gqtp", host, port)
	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]any)) != 1 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

// Deprecated: It is scheduled to be abolished.
func TestNewGroongaClientUnsupportedProtocolIsNil(t *testing.T) {
	c := NewGroongaClient("json", "localhost", 10041)
	if c.client != nil {
		t.Error("return is not error.")
	}
}

// Deprecated: It is scheduled to be abolished.
func TestNewGroongaClientClinetHttpTableListEmptySuccess(t *testing.T) {
	body := cTEST_BODY
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	schema := u.Scheme
	host, p, _ := net.SplitHostPort(u.Host)
	port, _ := strconv.Atoi(p)
	client := NewGroongaClient(schema, host, port)

	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]any)) != 1 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

// Deprecated: It is scheduled to be abolished.
func TestNewGroongaClientClinetGqtpTableListEmptySuccess(t *testing.T) {
	body := byteBody
	mock := gqtpMock(body)
	defer mock.Close()

	host, p, _ := net.SplitHostPort(mock.Address)
	port, _ := strconv.Atoi(p)
	client := NewGroongaClient("gqtp", host, port)
	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]any)) != 1 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

// Benchmarks
func BenchmarkHTTPClient(b *testing.B) {
	b.SkipNow()
	client := NewClient("http", "localhost", 10041)
	params := map[string]string{
		"table": "Users",
		"query": "name:@test",
	}
	for n := 0; b.N > n; n++ {
		result, _ := client.Call("select", params)
		if len(result.RawData) == 0 {
			b.Errorf("response body not found")
		}
	}
}

func BenchmarkGQTPClient(b *testing.B) {
	b.SkipNow()
	client := NewClient("gqtp", "localhost", 10043)
	params := map[string]string{
		"table": "Users",
		"query": "name:@Jim",
	}
	for n := 0; b.N > n; n++ {
		result, _ := client.Call("select", params)
		if len(result.RawData) == 0 {
			b.Errorf("response body not found")
		}
	}
}

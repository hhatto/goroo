package goroo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func newServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, body)
	}))
}

func TestHttpTableListEmptySuccess(t *testing.T) {
	const body = "[[0,1444022807.258,0.0],[[[\"id\",\"UInt32\"],[\"name\",\"ShortText\"],[\"path\",\"ShortText\"],[\"flags\",\"ShortText\"],[\"domain\",\"ShortText\"],[\"range\",\"ShortText\"],[\"default_tokenizer\",\"ShortText\"],[\"normalizer\",\"ShortText\"]]]]"
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	client := newHttpClient(fmt.Sprintf("%s://%s", u.Scheme, u.Host))

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

func TestHttpTableListCount1Success(t *testing.T) {
	const body = "[[0,1444024497.318,0.0469999313354492],[[[\"id\",\"UInt32\"],[\"name\",\"ShortText\"],[\"path\",\"ShortText\"],[\"flags\",\"ShortText\"],[\"domain\",\"ShortText\"],[\"range\",\"ShortText\"],[\"default_tokenizer\",\"ShortText\"],[\"normalizer\",\"ShortText\"]],[256,\"TestGQTPClinet\",\"./markdown.db.0000100\",\"TABLE_HASH_KEY|PERSISTENT\",null,null,null,null]]]"
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	client := newHttpClient(fmt.Sprintf("%s://%s", u.Scheme, u.Host))

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

func TestHttpColumnCreateUserNameSuccess(t *testing.T) {
	const body = "[[0,1444025635.392,0.00300002098083496],true]"
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	client := newHttpClient(fmt.Sprintf("%s://%s", u.Scheme, u.Host))

	res, err := client.Call("column_create", map[string]string{
		"table": "GQTPTable",
		"name":  "user_name",
		"type":  "ShortText",
	})
	if err != nil {
		t.Error(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if res.Body.(bool) != true {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestHttpColumnCreateUserNameFail(t *testing.T) {
	const body = "[[-22,1444025814.842,0.0,\"already used name was assigned: <user_name>\",[[\"grn_obj_register\",\"db.c\",8966]]],false]"
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	client := newHttpClient(fmt.Sprintf("%s://%s", u.Scheme, u.Host))

	res, err := client.Call("column_create", map[string]string{
		"table": "GQTPTable",
		"name":  "user_name",
		"type":  "ShortText",
	})
	if err == nil {
		t.Errorf("err is nil")
	}
	if res.Status != -22 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if res.Body.(bool) != false {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestHttpSelectTableNotFound(t *testing.T) {
	const body = "[[-22,1444109599.174,0.0,\"invalid table name: <Users>\",[[\"grn_select\",\"proc.c\",1217]]]]"
	ts := newServer(body)
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	client := newHttpClient(fmt.Sprintf("%s://%s", u.Scheme, u.Host))

	res, err := client.Call("column_create", map[string]string{
		"table": "GQTPTable",
		"name":  "user_name",
		"type":  "ShortText",
	})
	if err == nil {
		t.Errorf("err is nil")
	}
	if res.Status != -22 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if res.Body != nil {
		t.Errorf("body is not nil")
	}
}

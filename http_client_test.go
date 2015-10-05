package goroo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

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

func TestHttp_ColumnCreate_UserName_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const body = "[[0,1444025635.392,0.00300002098083496],true]"
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
	if res.Body.(bool) != true {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestHttp_ColumnCreate_UserName_Fail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const body = "[[-22,1444025814.842,0.0,\"already used name was assigned: <user_name>\",[[\"grn_obj_register\",\"db.c\",8966]]],false]"
		fmt.Fprintln(w, body)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}
	client := NewHttpClient(u.Host)
	res, err := client.Call("table_list", map[string]string{})
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

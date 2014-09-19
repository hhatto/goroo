package main

import (
	"fmt"
	"testing"
)

func TestHTTPClient(t *testing.T) {
	//client := GroongaClient{Protocol: "http", Host: "localhost", Port: 10041}
	client := NewGroongaClient("http", "localhost", 10041)

	params := map[string]string{
		"name":     "Users",
		"flags":    "TABLE_HASH_KEY",
		"key_type": "ShortText",
	}
	result, _ := client.Call("table_create", params)
	fmt.Println(result.RawData)
	fmt.Println(result)

	params = map[string]string{
		"table": "Users",
		"name":  "name",
		"flags": "COLUMN_SCALAR",
		"type":  "ShortText",
	}
	result, _ = client.Call("column_create", params)
	fmt.Println(result.RawData)
	fmt.Println(result)

	params = map[string]string{
		"table":  "Users",
		"values": "[{\"_key\":\"ken\",\"name\":\"Ken\"},{\"_key\":\"jim\",\"name\":\"Jim\"}]",
	}
	result, _ = client.Call("load", params)
	fmt.Println(result.RawData)
	fmt.Println(result)

	params = map[string]string{
		"table": "Users",
	}
	result, _ = client.Call("select", params)
	if len(result.RawData) == 0 {
		t.Errorf("response body not found")
	}

	fmt.Println(result.RawData)
	fmt.Println(result)
}

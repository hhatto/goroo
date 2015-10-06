package main

import (
	"fmt"

	"github.com/hhatto/goroo"
)

func main() {
	client := goroo.NewClient("http", "localhost", 10041)
	result, err := client.Call("select", map[string]string{"table": "Users"})
	if err != nil {
		fmt.Println("Call() error:", err)
		return
	}
	fmt.Println(result)
}

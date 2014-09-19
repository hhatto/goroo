# goroo

Yet Another Groonga Client for Go.

## Installation
```
$ go get github.com/hhatto/goroo
```

## Usage

### 1. Start Groonga Server
ex)
```
$ groonga -s -l 8 --log-path ./grn.log --protocol http grn.db
```

### 2. execute client code
```go
package main

import (
	"fmt"

	"github.com/hhatto/goroo"
)

func main() {
	client := goroo.NewGroongaClient("http", "localhost", 10041)
	result, err := client.Call("select", map[string]string{"table": "Users"})
	if err != nil {
		fmt.Println("Call() error:", err)
		return
	}
	fmt.Println(result)
}
```

## TODO
- [ ] GQTP support

## License
MIT

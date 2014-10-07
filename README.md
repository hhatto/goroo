# goroo (ごろう) [![Build Status](https://travis-ci.org/hhatto/goroo.png?branch=master)](https://travis-ci.org/hhatto/goroo)

Yet Another Groonga Client for Go.


## Installation
```
$ go get github.com/hhatto/goroo
```


## Usage

### with HTTP
  1. Start Groonga Server (with HTTP)

    ex)
    ```
    $ groonga -s -l 8 --log-path ./grn.log --protocol http grn.db
    ```

  2. execute client code

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

### with GQTP
  1. Start Groonga Server (with GQTP)

    ex)
    ```
    $ groonga -s -l 8 --log-path ./grn.log --protocol gqtp grn.db
    ```

  2. execute client code

    ```go
    package main

    import (
        "fmt"

        "github.com/hhatto/goroo"
    )

    func main() {
        client := goroo.NewGroongaClient("gqtp", "localhost", 10043)
        result, err := client.Call("select", map[string]string{"table": "Users"})
        if err != nil {
            fmt.Println("Call() error:", err)
            return
        }
        fmt.Println(result)
    }
    ```

## License
MIT

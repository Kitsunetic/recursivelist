# ListRecurrent

List files recurrently using channels

## Installation and usage

Install library using `go get` command.

```bash
go get -u github.com/Kitsunetic/recursivelist
```

```go
package main

import (
    "fmt"
    
    "github.com/Kitsunetic/recursivelist"
)

func main() {
    files, errs, done := recursivelist.RecursiveList("./test")
L:
    for {
        select {
        case file := <-files:
            fmt.Println(file)
        case err := <-errs:
            fmt.Println(err)
        case <-done:
            break L
        }
    }
}
```

The result is

```
/a/a1
/a/a2
/a/b/b1
/a/c/c1
/a/c/c2
/test.go
```

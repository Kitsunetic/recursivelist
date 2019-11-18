# RecursiveList

List files recursively using channels

## Installation and Usage

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

The structure of directory `test` is

```
test
  + a
    + b
      - b1
    + c
      - c1
      - c2
    - a1
    - a2
  - test.go
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

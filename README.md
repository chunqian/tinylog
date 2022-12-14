# tinylog 
tinylog is a dead simple, levelable, colorful logging library.

## use 
```golang
package main

import "C"

import (
  log "github.com/chunqian/tinylog"
)

type Tiny struct {
  name   string
  number int
  keys   []int
  dict   map[string]any
}

func main() {
  t := Tiny{
    "tiny",
    100,
    []int{1, 2, 3, 4, 5, 6},
    map[string]any{
      "start": 123,
      "end":   456,
    },
  }

  msg := [6]int8{'H', 'e', 'l', 'l', 'o', '\x00'}
  msg2 := [4]int8{'G', 'o', '!', '\x00'}
  msg3 := [6]C.char{'H', 'e', 'l', 'l', 'o', '\x00'}
  msg4 := [4]C.char{'G', 'o', '!', '\x00'}

  log.Info("Tiny: {}", t)

  log.Debug("Say: {}, {}", "Hello", "Go!")
  log.Warn("Say: {}, {}", "Hello", "Go!")
  log.Info("Say: {}, {}", "Hello", "Go!")
  log.Error("Say: {}, {}", "Hello", "Go!")
  log.Message("Say: {}, {}", &msg[0], &msg2[0])
  log.Message("Say: {}, {}", unsafe.Pointer(&msg3[0]), unsafe.Pointer(&msg4[0]))
  log.Fatal("Say: {}, {}", "Hello", "Go!")
}

```

## output 
```shell
go build -o stdlog ./test
./stdlog
```
```shell
[INFO] Tiny: main.Tiny{
  name:   "tiny",
  number: 100,
  keys:   {1, 2, 3, 4, 5, 6},
  dict:   {
    "start": int(123),
    "end":   int(456),
  },
}
[DEBUG] Say: Hello, Go!
[WARN] Say: Hello, Go!
[INFO] Say: Hello, Go!
[ERROR] Say: Hello, Go!
[MESSAGE] Say: Hello, Go!
[MESSAGE] Say: Hello, Go!
[FATAL] Say: Hello, Go!
```

## output files 
```golang
//go:build cronolog

package main

import (
  "os"

  log "github.com/chunqian/tinylog"
)

func init() {
  // create writer to writes message to a set of output files
  getwd, _ := os.Getwd()
  writer, _ := log.NewWriter(getwd + "/logs/%Y/%m/%d/test.log")
  log.SetOutput(writer)
}

```
```shell
go build -tags cronolog -o cronolog ./test
./cronolog
```

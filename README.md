# tinylog 
tinylog is a dead simple, levelable, colorful logging library.

## use 
```golang
package main

import (
  log "github.com/chunqian/tinylog"
)

func main() {
  log.Debug("Say: {}, {}", "Hello", "Go!")
  log.Warn("Say: {}, {}", "Hello", "Go!")
  log.Info("Say: {}, {}", "Hello", "Go!")
  log.Error("Say: {}, {}", "Hello", "Go!")
  log.Fatal("Say: {}, {}", "Hello", "Go!")
}
```

## output
```shell
[DEBUG] Say: Hello Go!
[INFO] Say: Hello Go!
[WARN] Say: Hello Go!
[ERROR] Say: Hello Go!
[FATAL] Say: Hello Go!
```

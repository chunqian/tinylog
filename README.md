# tinylog 
tinylog is a dead simple, levelable, colorful logging library.

## use 
```golang
package main

import (
  log "github.com/chunqian/tinylog"
)

func main() {
  log.Debug("Say: {}, {}", "Hello", "Go!") // stdout: [DEBUG] Say: Hello Go!
  log.Warn("Say: {}, {}", "Hello", "Go!") // stdout: [INFO] Say: Hello Go!
  log.Info("Say: {}, {}", "Hello", "Go!") // stdout: [WARN] Say: Hello Go!
  log.Error("Say: {}, {}", "Hello", "Go!") // stdout: [ERROR] Say: Hello Go!
  log.Fatal("Say: {}, {}", "Hello", "Go!") // stdout: [FATAL] Say: Hello Go!
}
```

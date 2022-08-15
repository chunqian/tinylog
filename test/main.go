/**---------------------------------------------------------
 * name: main.go
 * author: shenchunqian
 * created: 2022-08-13
 ---------------------------------------------------------*/

package main

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

	log.Info("Tiny: {}", t)

	log.Debug("Say: {}, {}", "Hello", "Go!")
	log.Warn("Say: {}, {}", "Hello", "Go!")
	log.Info("Say: {}, {}", "Hello", "Go!")
	log.Error("Say: {}, {}", "Hello", "Go!")
	log.Fatal("Say: {}, {}", "Hello", "Go!")
}

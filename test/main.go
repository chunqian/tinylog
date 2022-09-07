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

	msg := [6]int8{'H', 'e', 'l', 'l', 'o', '\x00'}
	msg2 := [4]int8{'G', 'o', '!', '\x00'}

	log.Info("Tiny: {}", t)

	log.Debug("Say: {}, {}", "Hello", "Go!")
	log.Warn("Say: {}, {}", "Hello", "Go!")
	log.Info("Say: {}, {}", "Hello", "Go!")
	log.Error("Say: {}, {}", "Hello", "Go!")
	log.Message("Say: {}, {}", &msg[0], &msg2[0])
	log.Fatal("Say: {}, {}", "Hello", "Go!")
}

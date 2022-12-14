//go:build cronolog

/**---------------------------------------------------------
 * name: cronolog.go
 * author: shenchunqian
 * created: 2022-08-13
 ---------------------------------------------------------*/

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

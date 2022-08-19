/**---------------------------------------------------------
 * name: printf.go
 * author: shenchunqian
 * created: 2022-08-19
 ---------------------------------------------------------*/

package main

import (
	log "github.com/chunqian/tinylog"
	"github.com/chunqian/tinylog/pretty"
)

func printf(zero, level, format *int8, args ...interface{}) int32 {
	goformat := pretty.Gostring(format)
	goLevel := pretty.Gostring(level)

	args = append(args, 0)
	copy(args[1:], args[:])
	args[0] = goformat

	switch goLevel {
	case "DEBUG":
		log.Debug(args...)
	case "INFO":
		log.Info(args...)
	case "WARN":
		log.Warn(args...)
	case "ERROR":
		log.Error(args...)
	case "FATAL":
		log.Fatal(args...)
	case "MESSAGE":
		for i, arg := range args {
			if v, ok := arg.(*int8); ok {
				args[i] = pretty.Gostring(v)
			}
			if v, ok := arg.(*uint8); ok {
				args[i] = pretty.Gostring(v)
			}
		}
		log.Message(args...)
	}
	return 0
}

/**---------------------------------------------------------
 * name: log.go
 * author: shenchunqian
 * created: 2022-08-12
 ---------------------------------------------------------*/

package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"bytes"
	"github.com/chunqian/tinylog/pretty"
	"time"
)

var (
	Prefix     = "[Log]"
	TimeFormat = "06-01-02 15:04:05"

	NonColor   bool
	ShowDepth  bool
	ShowTime   bool
	ShowPrefix bool
	DefaultCallerDepth = 3

	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

func init() {
	if runtime.GOOS == "windows" {
		NonColor = true
	}
	ShowDepth  = false
	ShowTime   = false
	ShowPrefix = false
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Print(level Level, depth int, addNewline bool, args ...interface{}) {
	var buf bytes.Buffer
	top := len(args)
	if top == 1 {
		args = append(args, 0)
		copy(args[1:], args[:])
		args[0] = "{}"
		top = len(args)
	}

	var format string
	var formatSlice []string

	// Add all the string arguments to the buffer
	for i := 0; i < top; i++ {
		var value = args[i]
		if (i == 0) {
			format = value.(string)
			formatSlice = strings.Split(format, "{}")
		}
		if (i > len(formatSlice)) {
			break
		}
		
		if (i > 0) {
			var mStr = ""
			switch value.(type) {
			case string:
				mStr = strings.ReplaceAll(value.(string), "interface {}", "any")
			default:
				mStr = strings.ReplaceAll(fmt.Sprintf("%# v", pretty.Formatter(value)), "interface {}", "any")
			}
			buf.WriteString(mStr)
			buf.WriteString(formatSlice[i])
		} else {
			buf.WriteString(formatSlice[i])
		}
	}
	if addNewline {
		buf.WriteString("\n")
	}

	var depthInfo string
	if ShowDepth {
		if depth == -1 {
			depth = DefaultCallerDepth
		}
		pc, file, line, ok := runtime.Caller(depth)
		if ok {
			// Get caller function name.
			fn := runtime.FuncForPC(pc)
			var fnName string
			if fn == nil {
				fnName = "?()"
			} else {
				fnName = strings.TrimLeft(filepath.Ext(fn.Name()), ".") + "()"
			}
			depthInfo = fmt.Sprintf("[%s:%d %s] ", filepath.Base(file), line, fnName)
		}
	}

	var formatBuf bytes.Buffer
	var selected []any

	if ShowPrefix {
		formatBuf.WriteString("%s ")
		selected = append(selected, Prefix)
	}
	if ShowTime {
		if NonColor {
			formatBuf.WriteString("%s ")
		} else {
			formatBuf.WriteString("\033[36m%s\033[0m ")
		}
		selected = append(selected, time.Now().Format(TimeFormat))
	}

	if NonColor {
		formatBuf.WriteString("[%s] ")
		selected = append(selected, levelFlags[level])
	} else {
		switch level {
		case DEBUG:
			formatBuf.WriteString("[\033[34m%s\033[0m] ")
			selected = append(selected, levelFlags[level])
		case INFO:
			formatBuf.WriteString("[\033[36m%s\033[0m] ")
			selected = append(selected, levelFlags[level])
		case WARNING:
			formatBuf.WriteString("[\033[33m%s\033[0m] ")
			selected = append(selected, levelFlags[level])
		case ERROR:
			formatBuf.WriteString("[\033[31m%s\033[0m] ")
			selected = append(selected, levelFlags[level])
		case FATAL:
			formatBuf.WriteString("[\033[35m%s\033[0m] ")
			selected = append(selected, levelFlags[level])
		default:
			formatBuf.WriteString("[%s] ")
			selected = append(selected, levelFlags[level])
		}
	}

	if ShowDepth {
		formatBuf.WriteString("%s")
		selected = append(selected, depthInfo)
	}

	formatBuf.WriteString("%s")
	selected = append(selected, buf.String())
	formatBuf.WriteString("\n")

	fmt.Printf(formatBuf.String(), selected...)

	if level == FATAL {
		os.Exit(1)
	}

	return
}

func debugD(depth int, args ...interface{}) {
	Print(DEBUG, depth, false, args...)
}

func Debug(args ...interface{}) {
	debugD(-1, args...)
}

func warnD(depth int, args ...interface{}) {
	Print(WARNING, depth, false, args...)
}

func Warn(args ...interface{}) {
	warnD(-1, args...)
}

func infoD(depth int, args ...interface{}) {
	Print(INFO, depth, false, args...)
}

func Info(args ...interface{}) {
	infoD(-1, args...)
}

func errorD(depth int, args ...interface{}) {
	Print(ERROR, depth, false, args...)
}

func Error(args ...interface{}) {
	errorD(-1, args...)
}

func fatalD(depth int, args ...interface{}) {
	Print(FATAL, depth, false, args...)
}

func Fatal(args ...interface{}) {
	fatalD(-1, args...)
}

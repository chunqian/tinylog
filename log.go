/**---------------------------------------------------------
 * name: log.go
 * author: shenchunqian
 * created: 2022-08-12
 ---------------------------------------------------------*/

package log

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/chunqian/tinylog/pretty"
)

var (
	Prefix     = "[Log]"
	TimeFormat = "06-01-02 15:04:05"

	NonColor           bool
	ShowDepth          bool
	ShowTime           bool
	ShowPrefix         bool
	DefaultCallerDepth = 3

	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "MESSAGE"}

	logWriter *Writer
)

func init() {
	if runtime.GOOS == "windows" {
		NonColor = true
	}
	ShowDepth = false
	ShowTime = false
	ShowPrefix = false
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
	MESSAGE
)

func SetOutput(writer *Writer) error {
	logWriter = writer
	return nil
}

func expandToFront(args ...any) []any {
	top := len(args)
	args = append(args, 0)
	copy(args[1:], args[:])
	var fmtStr = ""
	for i := 0; i < top; i++ {
		if i == top-1 {
			fmtStr += "{}"
		} else {
			fmtStr += "{}, "
		}
	}
	args[0] = fmtStr
	return args
}

func expandToEnd(args ...any) []any {
	var count = 0
	switch args[0].(type) {
	case string:
		s := strings.Split(args[0].(string), "{}")
		count = len(s) - 1
	}

	var diff_count = 0
	if len(args)-1 < count {
		diff_count = count - (len(args) - 1)
	}
	for i := 0; i < diff_count; i++ {
		args = append(args, "not found!")
	}
	return args
}

func Print(level Level, depth int, addNewline bool, args ...any) {
	var buf bytes.Buffer

	top := len(args)
	if top == 0 {
		return
	}

	switch args[0].(type) {
	case string:
		matched, _ := regexp.MatchString(`{}`, args[0].(string))
		if !matched {
			args = expandToFront(args...)
		} else {
			args = expandToEnd(args...)
		}
		top = len(args)
	default:
		args = expandToFront(args...)
		top = len(args)
	}

	var format string
	var formatSlice []string

	// Add all the string arguments to the buffer
	for i := 0; i < top; i++ {
		var value = args[i]
		if i == 0 {
			format = value.(string)
			formatSlice = strings.Split(format, "{}")
		}
		if i >= len(formatSlice) {
			break
		}

		if i > 0 {
			var mStr = ""
			if level == MESSAGE {
	      if v, ok := value.(*int8); ok {
	        value = pretty.Gostring(v)
	      }
	      if v, ok := value.(*uint8); ok {
	        value = pretty.Gostring(v)
	      }
			}
			switch value.(type) {
			case string:
				mStr = strings.ReplaceAll(value.(string), "interface {}", "any")
			default:
				mStr = strings.ReplaceAll(fmt.Sprintf("%# v", pretty.Formatter(value)), "interface {}", "any")
			}
			buf.WriteString(mStr)
		}
		buf.WriteString(formatSlice[i])
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
		case INFO:
			formatBuf.WriteString("[\033[36m%s\033[0m] ")
		case WARNING:
			formatBuf.WriteString("[\033[33m%s\033[0m] ")
		case ERROR:
			formatBuf.WriteString("[\033[31m%s\033[0m] ")
		case FATAL:
			formatBuf.WriteString("[\033[35m%s\033[0m] ")
		case MESSAGE:
			formatBuf.WriteString("[\033[34m%s\033[0m] ")
		default:
			formatBuf.WriteString("[%s] ")
		}
		selected = append(selected, levelFlags[level])
	}

	if ShowDepth {
		formatBuf.WriteString("%s")
		selected = append(selected, depthInfo)
	}

	formatBuf.WriteString("%s")
	selected = append(selected, buf.String())
	formatBuf.WriteString("\n")

	if logWriter != nil {
		logWriter.Write([]byte(fmt.Sprintf(formatBuf.String(), selected...)))
	} else {
		fmt.Printf(formatBuf.String(), selected...)
	}

	if level == FATAL {
		os.Exit(1)
	}

	return
}

func debugD(depth int, args ...any) {
	Print(DEBUG, depth, false, args...)
}

func Debug(args ...any) {
	debugD(-1, args...)
}

func warnD(depth int, args ...any) {
	Print(WARNING, depth, false, args...)
}

func Warn(args ...any) {
	warnD(-1, args...)
}

func infoD(depth int, args ...any) {
	Print(INFO, depth, false, args...)
}

func Info(args ...any) {
	infoD(-1, args...)
}

func errorD(depth int, args ...any) {
	Print(ERROR, depth, false, args...)
}

func Error(args ...any) {
	errorD(-1, args...)
}

func fatalD(depth int, args ...any) {
	Print(FATAL, depth, false, args...)
}

func Fatal(args ...any) {
	fatalD(-1, args...)
}

func messageD(depth int, args ...any) {
	Print(MESSAGE, depth, false, args...)
}

func Message(args ...any) {
	messageD(-1, args...)
}

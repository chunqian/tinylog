//go:build cronolog

/**---------------------------------------------------------
 * name: writer.go
 * editor: shenchunqian
 * created: 2022-08-15
 * source: https://github.com/utahta/go-cronowriter
 ---------------------------------------------------------*/

package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/chunqian/tinylog/pretty"
	"github.com/lestrrat-go/strftime"
)

var (
	Prefix     = "[Log]"
	TimeFormat = "06-01-02 15:04:05"

	NonColor           bool
	ShowDepth          bool
	ShowTime           bool
	ShowPrefix         bool
	DefaultCallerDepth = 3

	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

	writer *Writer
)

func init() {
	if runtime.GOOS == "windows" {
		NonColor = true
	}
	ShowDepth = false
	ShowTime = false
	ShowPrefix = false

	getwd, _ := os.Getwd()
	writer, _ = NewWriter(getwd + "/logs/%Y/%m/%d/tiny.log")
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type (
	// A Writer writes message to a set of output files.
	Writer struct {
		pattern *strftime.Strftime // given pattern
		path    string             // current file path
		symlink *strftime.Strftime // symbolic link to current file path
		fp      *os.File           // current file pointer
		loc     *time.Location
		mux     sync.Locker
		init    bool // if true, open the file when New() method is called
	}
)

var (
	_   io.WriteCloser = (*Writer)(nil) // check if object implements interface
	now                = time.Now       // for test
)

// New returns a Writer with the given pattern.
func NewWriter(pattern string) (*Writer, error) {
	p, err := strftime.New(pattern)
	if err != nil {
		return nil, err
	}

	c := &Writer{
		pattern: p,
		path:    "",
		symlink: nil,
		fp:      nil,
		loc:     time.Local,
		mux:     new(sync.Mutex), // default mutex enable
		init:    false,
	}

	if c.init {
		if _, err := c.Write([]byte("")); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Write writes to the file and rotate files automatically based on current date and time.
func (c *Writer) Write(b []byte) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	t := now().In(c.loc)
	path := c.pattern.FormatString(t)

	if c.path != path {
		// close file
		go func(fp *os.File) {
			if fp == nil {
				return
			}
			fp.Close()
		}(c.fp)

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return c.write(nil, err)
		}

		fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return c.write(nil, err)
		}
		c.createSymlink(t, path)

		c.path = path
		c.fp = fp
	}

	return c.write(b, nil)
}

// Path returns the current writing file path.
func (c *Writer) Path() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.path
}

func (c *Writer) createSymlink(t time.Time, path string) {
	if c.symlink == nil {
		return
	}

	symlink := c.symlink.FormatString(t)
	if symlink == path {
		fmt.Println("Can't create symlink. Already file exists.")
		return // ignore error
	}

	if _, err := os.Stat(symlink); err == nil {
		if err := os.Remove(symlink); err != nil {
			fmt.Println(err)
			return // ignore error
		}
	}

	if err := os.Symlink(path, symlink); err != nil {
		fmt.Println(err)
		return // ignore error
	}
}

// Close closes file.
func (c *Writer) Close() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.fp.Close()
}

func (c *Writer) write(b []byte, err error) (int, error) {
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return c.fp.Write(b)
}

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
		if i == 0 {
			format = value.(string)
			formatSlice = strings.Split(format, "{}")
		}
		if i > len(formatSlice) {
			break
		}

		if i > 0 {
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
		case INFO:
			formatBuf.WriteString("[\033[36m%s\033[0m] ")
		case WARNING:
			formatBuf.WriteString("[\033[33m%s\033[0m] ")
		case ERROR:
			formatBuf.WriteString("[\033[31m%s\033[0m] ")
		case FATAL:
			formatBuf.WriteString("[\033[35m%s\033[0m] ")
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

	// fmt.Printf(formatBuf.String(), selected...)
	writer.Write([]byte(fmt.Sprintf(formatBuf.String(), selected...)))

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

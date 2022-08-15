//go:build cronolog

/**---------------------------------------------------------
 * name: writer.go
 * editor: shenchunqian
 * created: 2022-08-15
 * source: https://github.com/utahta/go-cronowriter
 ---------------------------------------------------------*/

package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lestrrat-go/strftime"
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

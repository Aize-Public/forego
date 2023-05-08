package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func Stack(above, max int) []string {
	stack := make([]string, 0, 20)
	for len(stack) < max {
		_, file, line, ok := runtime.Caller(above + 1)
		if !ok {
			return stack
		}
		stack = append(stack, fmt.Sprintf("%s:%d", file, line))
		above++
	}
	return stack
}

type StackFrame struct {
	File string
	Line int
	Func string
}

// Returns the StackFrame of the caller if above is zero.
func Caller(above int) StackFrame {
	rpc := make([]uintptr, 1)
	_ = runtime.Callers(above+2, rpc) // +0 runtime.Caller, +1 this, +2 is the actual caller
	f, _ := runtime.CallersFrames(rpc).Next()
	// we always return something, even if empty
	return StackFrame{
		File: f.File,
		Line: f.Line,
		Func: f.Function,
	}
}

// convert --trimpath like `aize.io/monorepo/...` to full path based on go.mod content
func (this StackFrame) AbsFile() string {
	_, err := os.Stat(this.File)
	if err == nil { // likely this.File is in form of "modname/path/file"
		return this.File
	}
	path := strings.TrimPrefix(this.File, mod.name)
	if path != this.File {
		return mod.path + path
	}

	return ""
}

func (this StackFrame) FileLine() string {
	return fmt.Sprintf("%s:%d", this.File, this.Line)
}

// return only the package.FuncName
func (this StackFrame) ShortFunc() string {
	p := strings.Split(this.Func, "/")
	return p[len(p)-1]
}

func (this StackFrame) FuncName() string {
	p := strings.Split(this.Func, ".")
	return p[len(p)-1]
}

func (this StackFrame) Pkg() string {
	return path.Dir(this.File)
}

// find go.mod, and allow to expand/contract by it
// TODO(oha) we could parse further for local `replace` rules, and allow to expand for any submodule if needed
var mod = func() (out struct {
	name string // the name of the current module in go.mod
	path string // the dir which contains the go.mod file
}) {
	d, _ := os.Getwd()
	for {
		body, _ := os.ReadFile(d + "/go.mod")
		if len(body) > 0 {
			m := regexp.MustCompile(`module +(\S+)`).FindStringSubmatch(string(body))
			if m != nil {
				out.name = m[1]
			}
			out.path = d
			return
		}
		p := filepath.Dir(d)
		if p == d {
			return
		} else {
			d = p
		}
	}
}()

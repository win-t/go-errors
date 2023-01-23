// Provide utility to get stack trace.
package trace

import (
	"fmt"
	"runtime"
	"sync"
)

var pool sync.Pool

// Single location of the trace
type Location struct {
	file  string
	line  int
	func_ string
}

// String representation of Location
func (l *Location) String() string {
	if l.func_ == "" {
		return fmt.Sprintf("%s:%d", l.file, l.line)
	}

	return fmt.Sprintf("%s:%d (%s)", l.file, l.line, l.func_)
}

// The file path that this Location point to
func (l *Location) File() string { return l.file }

// The line that this Location point to
func (l *Location) Line() int { return l.line }

// The path-qualified function that this location point to
func (l *Location) Func() string { return l.func_ }

// Get return list of location of stack trace for calling function.
//
// skip tell Get to skip some trace, 0 is where Get is called.
// deep tell Get how deep the stack trace is.
func Get(skip, deep int) (locations []Location) {
	if deep <= 0 {
		return nil
	}

	if skip < 0 {
		skip = 0
	}
	skip += 2

	dataLen := deep + 10
	var data *[]uintptr
	if tmp, ok := pool.Get().(*[]uintptr); ok {
		if len(*tmp) >= dataLen {
			data = tmp
		}
	}
	if data == nil {
		tmp := make([]uintptr, dataLen)
		data = &tmp
	}

	pc := (*data)[:dataLen]
	pc = pc[:runtime.Callers(skip, pc)]
	if len(pc) == 0 {
		return nil
	}

	locations = make([]Location, 0, len(pc))

	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		if frame.Line != 0 && frame.File != "" &&
			!inPkg(frame.Function, "runtime") &&
			!inPkg(frame.Function, "github.com/win-t/go-errors") {
			locations = append(locations, Location{
				func_: frame.Function,
				file:  frame.File,
				line:  frame.Line,
			})
		}
		if len(locations) == deep {
			break
		}
		if !more {
			break
		}
	}

	pool.Put(data)

	return
}

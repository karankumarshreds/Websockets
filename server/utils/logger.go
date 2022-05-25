package utils

import (
	"log"
	"runtime"
)

// Original from https://stackoverflow.com/a/25927915/1490379
func CustomLogger(msgs ...string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	_, line := f.FileLine(pc[0])
	// file, line := f.FileLine(pc[0])
	// fmt.Printf("%s:%d %s %v\n", file, line, f.Name(), msgs)
	log.Printf("%d: %s %v\n", line, f.Name(), msgs)
}
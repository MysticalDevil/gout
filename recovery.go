package gout

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	// Get the current call stack information
	callers := make([]uintptr, 32)
	n := runtime.Callers(3, callers)

	// Build call stack string
	var builder strings.Builder
	builder.WriteString(message + "\nTraceback:")
	for _, pc := range callers[:n] {
		// Obtain the function information and source code location information corresponding
		// to the call stack frame
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)

		// Format the function information and source code location information
		// into a string and append it to the builder object
		fmt.Fprintf(&builder, "\n\t%s:%d", file, line)
	}

	return builder.String()
}

// Recovery returns a middleware handler that recovers from any panics and writes a 500 error.
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

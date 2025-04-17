package alchemist

import (
	"bytes"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestMortalLogger tests the debug and info calls of the mortalLogger type.
func TestMortalLogger(t *testing.T) {

	var buf bytes.Buffer
	msg := "hello world\n"

	testCases := []struct {
		name   string
		logger mortalLogger
		call   func(mortalLogger, string, ...any) // test info or debug method
		want   string
	}{{
		name:   "info",
		logger: mortalLogger{Logger: log.New(&buf, "", 0)},
		// type.method: this converts the method 'info' to a function that has an
		// additional first parameter: the mortalLogger object to use
		call: mortalLogger.info,
		want: "[INFO] " + msg,
	}, {
		name:   "debug",
		logger: mortalLogger{Logger: log.New(&buf, "", 0)},
		call:   mortalLogger.debug,
	}, {
		name:   "verbose debug",
		logger: mortalLogger{Logger: log.New(&buf, "", 0), verbose: true},
		call:   mortalLogger.debug,
		want:   "[DEBUG] " + msg,
	}, {
		// no nil pointer panic when no logger is available (noop mode)
		name: "nil info",
		call: mortalLogger.info,
	}, {
		name: "nil debug",
		call: mortalLogger.debug,
	}, {
		name:   "nil verbose debug",
		logger: mortalLogger{verbose: true},
		call:   mortalLogger.debug,
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			buf.Reset()

			// call the method by calling the function and providing
			// the mortalLogger object as first parameter.
			c.call(c.logger, "hello %s", "world")
			got := buf.String()
			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("ERROR: got- want+\n%s\n", diff)
			}
		})
	}
}

package alchemist

import "log"

// mortalLogger provides human-readable logging.
type mortalLogger struct {
	*log.Logger
	verbose bool
}

// debug writes a debug message to the logger if verbose is true.
func (m mortalLogger) debug(msg string, args ...any) {
	// this enables using the mortalLogger with its 'zero value'
	// to run in a no-op mode.
	if m.Logger == nil {
		return
	}
	if m.verbose {
		m.Printf("[DEBUG] "+msg, args...)
	}
}

// info writes an info message to the logger.
func (m mortalLogger) info(msg string, args ...any) {
	// this enables using the mortalLogger with its 'zero value'
	// to run in a no-op mode.
	if m.Logger == nil {
		return
	}
	m.Printf("[INFO] "+msg, args...)
}

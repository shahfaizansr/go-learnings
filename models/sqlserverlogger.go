package models

import (
	"sync"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
)

type LogHarbourSQLHooks struct {
	Logger      *logharbour.Logger
	RigelClient *rigel.Rigel
	LogLevel    *LogLevel // Reference to the log level variable, allowing dynamic control.
}

// LogLevel defines a thread-safe structure to hold the current log level.
// This allows us to adjust the logging verbosity at runtime in a thread-safe manner.
type LogLevel struct {
	mu    sync.RWMutex
	level tracelog.LogLevel
}

// Set sets the log level in a thread-safe manner.
func (l *LogLevel) Set(level tracelog.LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Get retrieves the current log level in a thread-safe manner.
func (l *LogLevel) Get() tracelog.LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

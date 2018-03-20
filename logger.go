package aries

import (
	"errors"
	"fmt"
	"log"
	"os"

	"shanhu.io/misc/errcode"
)

// LogPrinter is the interface for printing server logs.
type LogPrinter interface {
	Print(s string)
}

// Logger is a logger for logging server logs
type Logger struct {
	p LogPrinter
}

// NewLogger creates a new logger that prints server
// logs to the given printer.
func NewLogger(p LogPrinter) *Logger {
	return &Logger{p: p}
}

type stdLog struct{}

func (p *stdLog) Print(s string) {
	log.Println(s)
}

var theStdLog = new(stdLog)

// StdLogger returns the logger that logs to the default std log.
func StdLogger() *Logger {
	return &Logger{p: theStdLog}
}

// AltError logs the error and returns an alternative error with code.
func (l *Logger) AltError(err error, code, s string) error {
	if err == nil {
		return nil
	}
	l.p.Print(fmt.Sprintf("%s: %s", s, err))
	return errcode.Add(code, errors.New(s))
}

// AltInternal logs the error and returns an alternative internal error.
func (l *Logger) AltInternal(err error, s string) error {
	return l.AltError(err, errcode.Internal, s)
}

// AltInvalidArg logs the error and returns an alternative invalid arg error.
func (l *Logger) AltInvalidArg(err error, s string) error {
	return l.AltError(err, errcode.InvalidArg, s)
}

// Exit prints the error and exit the service.
func (l *Logger) Exit(err error) error {
	l.p.Print(err.Error())
	os.Exit(1)
	panic("unreachable")
}

// Print prints a message to the logger.
func (l *Logger) Print(s string) {
	l.p.Print(s)
}

// Log logs the message to the logger if the logger is not nil.
func Log(l *Logger, s string) {
	if l == nil {
		return
	}
	l.Print(s)
}

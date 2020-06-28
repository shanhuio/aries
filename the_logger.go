package aries

// TheLogger is the default logger that logs to default golang log.
var TheLogger = StdLogger()

// AltInternal prints error to TheLogger. It is an alias to
// TheLogger.AltInteral
func AltInternal(err error, s string) error {
	return TheLogger.AltInternal(err, s)
}

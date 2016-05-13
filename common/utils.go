package common

import "log"

func notImplemented(msg string) {
	warning(msg + " is not yet implemented on this platform")
}

func unsupportedType(v interface{}) {
	warning("type %T not supported", v)
}

func warning(format string, a ...interface{}) {
	log.Printf("[WARNING] "+format, a...)
}

package engo

import "log"

func notImplemented(msg string) {
	warning(msg + "is not yet implemented on this platform")
}

func unsupportedType() {
	warning("type not supported")
}

func warning(msg string) {
	log.Println("[WARNING] " + msg)
}

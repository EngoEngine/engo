package main

const (
	NAME      = "Game"
	BUNDLE    = NAME + ".app"
	CONTENTS  = BUNDLE + "/Contents"
	EXE       = CONTENTS + "/MacOS"
	RESOURCES = CONTENTS + "/Resources"
)

func main() {
	buildApp()
	buildIcns()
}

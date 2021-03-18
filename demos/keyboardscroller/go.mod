module github.com/EngoEngine/engo/demos/keyboardscroller

go 1.16

replace github.com/EngoEngine/engo => ./../..

replace github.com/EngoEngine/engo/demos/demoutils => ./../demoutils

require (
	github.com/EngoEngine/ecs v1.0.5
	github.com/EngoEngine/engo v0.0.0-00010101000000-000000000000
	github.com/EngoEngine/engo/demos/demoutils v0.0.0-00010101000000-000000000000
)

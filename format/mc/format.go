package mc

import (
	"encoding/json"
	"strings"
)

// https://github.com/egret-labs/egret-docs-en/tree/master/extension/game/movieClip
// http://developer.egret.com/en/github/egret-docs/extension/game/movieClip/index.html
// https://github.com/egret-labs/egret-docs-en/blob/master/tools/TextureMerger/manual/README.md
// http://developer.egret.com/en/github/egret-docs/tools/TextureMerger/manual/index.html

func Unmarshal(data []byte) (MovieClip, error) {
	r := MovieClip{}
	err := json.Unmarshal(data, &r)

	return r, err
}

type MovieClip struct {
	// MovieClip data list, Each attribute in the list represents a MovieClip name
	Mc map[string]Action `json:"mc"`
	// The texture file path corresponding to the data file  (used to help the tool to match the corresponding problem,
	// and the engine will not parse this attribute)
	File string `json:"file"`
	// Texture set data
	Regions map[string]Resource `json:"res"`
}

func (m *MovieClip) MaxXY(needle Frame) (int, int) {
	return maxXY(m.filterByActionNameFrames(needle))
}

func (m *MovieClip) filterByActionNameFrames(needle Frame) []Frame {
	frames := make([]Frame, 0)
	for _, frame := range m.AllFrames() {
		if frame.ActionName() == needle.ActionName() {
			frames = append(frames, frame)
		}
	}

	return frames
}

func (m *MovieClip) AllFrames() []Frame {
	frames := make([]Frame, 0)
	for _, action := range m.Mc {
		for _, frame := range action.Frames {
			frames = append(frames, frame)
		}
	}

	return frames
}

type Action struct {
	// Keyframe data list.
	Frames []Frame `json:"frames"`
	// Frame tag list, [optional attribute]. If there is no frame tag, you can choose not to add this attribute.
	Labels []Label `json:"labels,omitempty"`
	// Frame rate, [optional attribute], the default value of 24, which can be set by the developer through the code.
	FrameRate int `json:"frameRate,omitempty"`
	// Frame events list [optional properties]. If there is no frame action, you can choose not to add this attribute.
	Events []Event `json:"events,omitempty"`
	// Frame script list, [optional properties]. If there is no frame action, you can choose not to add this attribute.
	Scripts []Script `json:"scripts,omitempty"`
}

type Event struct {
	// Event name
	Name string `json:"name"`
	// The frame number where the event is located
	Frame int `json:"frame"`
}

type Script struct {
	// The frame number where the script is located.
	Frame int `json:"frame"`
	// The method name of the script call, supporting six APIs related to animation playback.
	Func string `json:"func"`
	// The frame number where the event is located
	Args []string `json:"args"`
}

type Frame struct {
	// The x coordinate required to be displayed by image, [optional attribute], the default value of 0.
	Y int `json:"y,omitempty"`
	// The y coordinate required to be displayed by image, [optional attribute], the default value of 0.
	X int `json:"x,omitempty"`
	// The number of consecutive frames of the key frame, [optional attribute], the default value of 1.
	Duration int `json:"duration,omitempty"`
	// The image resources needed to be displayed on the key frame, [optional properties],
	// the default value is empty (used for the case of blank frames).
	ResourceName string `json:"res,omitempty"`
}

func (f *Frame) ActionName() string {
	part := strings.Split(f.ResourceName, "_")
	if len(part) < 2 {
		return ""
	}

	return part[0]
}

type Label struct {
	// Tag name
	Name string `json:"name"`
	// The frame number where the label is located
	FrameStart int `json:"frame"`
	FrameEnd   int `json:"end"`
}

// Each attribute in the list represents a resource name
type Resource struct {
	// The y coordinate of the location of the texture set in the resource
	X int `json:"x"`
	// The y coordinate of the location of the texture set in the resource
	Y int `json:"y"`
	// Resource width
	W int `json:"w"`
	// Resource height
	H int `json:"h"`
}

func (r Resource) Centered(x, y int) Resource {
	r.X -= x
	r.Y -= y
	r.W += x
	r.H += y

	return r
}

func maxXY(list []Frame) (int, int) {
	var x, y int
	for _, frame := range list {
		if x < (-frame.X) {
			x = -frame.X
		}

		if y < (-frame.Y) {
			y = -frame.Y
		}
	}

	return x, y
}

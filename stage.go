package eng

type Stage struct {
	batch                     *Batch
	camera                    *Camera
	width, height             float32
	gutterWidth, gutterHeight float32
}

func NewStage(width, height float32, keepAspect bool) *Stage {
	stage := new(Stage)

	if width == 0 {
		width = float32(Width())
	}
	if height == 0 {
		height = float32(Height())
	}

	stage.batch = NewBatch()

	if keepAspect {
		screenWidth := float32(Width())
		screenHeight := float32(Height())
		if screenHeight/screenWidth < height/width {
			toScreenSpace := screenHeight / height
			toViewportSpace := height / screenHeight
			deviceWidth := width * toScreenSpace
			lengthen := (screenWidth - deviceWidth) * toViewportSpace
			stage.width = width + lengthen
			stage.height = height
			stage.gutterWidth = lengthen / 2
			stage.gutterHeight = 0
		} else {
			toScreenSpace := screenWidth / width
			toViewportSpace := width / screenWidth
			deviceHeight := height * toScreenSpace
			lengthen := (screenHeight - deviceHeight) * toViewportSpace
			stage.height = height + lengthen
			stage.width = width
			stage.gutterWidth = 0
			stage.gutterHeight = lengthen / 2
		}
	} else {
		stage.width = width
		stage.height = height
		stage.gutterWidth = 0
		stage.gutterHeight = 0
	}

	stage.camera = NewCamera(stage.width, stage.height)
	stage.camera.Position.X = stage.width / 2
	stage.camera.Position.Y = stage.height / 2

	return stage
}

func (s *Stage) Update() {
	s.camera.Update()
	s.batch.SetProjection(s.camera.Combined)
}

func (s *Stage) Batch() *Batch {
	return s.batch
}

func (s *Stage) ScreenToStage(x, y float32) (float32, float32) {
	tmp.X = x
	tmp.Y = y
	tmp.Z = 1
	s.camera.Unproject(tmp)
	return tmp.X, tmp.Y
}

func (s *Stage) Width() float32 {
	return s.width
}

func (s *Stage) GutterWidth() float32 {
	return s.gutterWidth
}

func (s *Stage) Height() float32 {
	return s.height
}

func (s *Stage) GutterHeight() float32 {
	return s.gutterHeight
}

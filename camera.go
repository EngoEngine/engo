package eng

type Camera struct {
	Zoom           float32
	Position       *Vector
	Direction      *Vector
	Up             *Vector
	Projection     *Matrix
	View           *Matrix
	Combined       *Matrix
	ViewportWidth  float32
	ViewportHeight float32
}

func NewCamera(width, height float32) *Camera {
	camera := new(Camera)
	camera.Zoom = 1
	camera.Position = new(Vector)
	camera.Direction = &Vector{0, 0, -1}
	camera.Up = &Vector{0, 1, 0}
	camera.Projection = NewMatrix()
	camera.View = NewMatrix()
	camera.Combined = NewMatrix()
	camera.ViewportWidth = width
	camera.ViewportHeight = height
	camera.Position.X = width / 2
	camera.Position.Y = height / 2
	camera.Update()
	return camera
}

var tmp = new(Vector)

func (c *Camera) Update() {
	c.Projection.SetToOrtho(c.Zoom*-c.ViewportWidth/2, c.Zoom*c.ViewportWidth/2, c.Zoom*c.ViewportHeight/2, c.Zoom*-c.ViewportHeight/2, 0, 1)
	c.View.SetToLookAt(c.Position, tmp.Set(c.Position).Add(c.Direction), c.Up)
	c.Combined.Set(c.Projection).Mul(c.View)
}

package common

import (
	"fmt"
	"testing"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/stretchr/testify/assert"
)

func TestSpaceComponent_Contains(t *testing.T) {
	space := SpaceComponent{Width: 100, Height: 100}
	pass := []engo.Point{
		engo.Point{X: 10, Y: 10},
		engo.Point{X: 50, Y: 50},
		engo.Point{X: 10, Y: 50},
		engo.Point{X: 99, Y: 99},
	}
	fail := []engo.Point{
		// Totally not within:
		engo.Point{X: -10, Y: -10},
		engo.Point{X: 120, Y: 120},

		// Only one axis within:
		engo.Point{X: 50, Y: 120},
		engo.Point{X: 120, Y: 50},

		// On the edge:
		engo.Point{X: 0, Y: 0},
		engo.Point{X: 0, Y: 50},
		engo.Point{X: 0, Y: 100},
		engo.Point{X: 50, Y: 0},
		engo.Point{X: 50, Y: 100},
		engo.Point{X: 100, Y: 0},
		engo.Point{X: 100, Y: 50},
		engo.Point{X: 100, Y: 100},
	}

	for _, p := range pass {
		assert.True(t, space.Contains(p), fmt.Sprintf("point %v should be within area", p))
	}

	for _, f := range fail {
		assert.False(t, space.Contains(f), fmt.Sprintf("point %v should not be within area", f))
	}
}

func TestSpaceComponent_Contains_Hitboxes(t *testing.T) {
	scs := []SpaceComponent{}
	space0 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 0, Y: 0}}  //AABB
	space1 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 10, Y: 0}} //Triangle
	space1.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 0, Y: 0}, P2: engo.Point{X: 0, Y: 5}},
		engo.Line{P1: engo.Point{X: 0, Y: 5}, P2: engo.Point{X: 5, Y: 5}},
		engo.Line{P1: engo.Point{X: 5, Y: 5}, P2: engo.Point{X: 0, Y: 0}},
	}})
	space2 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 0, Y: 10}} //Multi-Shape
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 2.5, Y: 0}, P2: engo.Point{X: 2.5 + 0.57735026919, Y: 1}},
		engo.Line{P1: engo.Point{X: 2.5 + 0.57735026919, Y: 1}, P2: engo.Point{X: 2.5 - 0.57735026919, Y: 1}},
		engo.Line{P1: engo.Point{X: 2.5 - 0.57735026919, Y: 1}, P2: engo.Point{X: 2.5, Y: 0}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 0, Y: 2.5}, P2: engo.Point{X: 1, Y: 2.5 + 0.57735026919}},
		engo.Line{P1: engo.Point{X: 1, Y: 2.5 + 0.57735026919}, P2: engo.Point{X: 1, Y: 2.5 - 0.57735026919}},
		engo.Line{P1: engo.Point{X: 1, Y: 2.5 - 0.57735026919}, P2: engo.Point{X: 0, Y: 2.5}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 5, Y: 2.5}, P2: engo.Point{X: 4, Y: 2.5 + 0.57735026919}},
		engo.Line{P1: engo.Point{X: 4, Y: 2.5 + 0.57735026919}, P2: engo.Point{X: 4, Y: 2.5 - 0.57735026919}},
		engo.Line{P1: engo.Point{X: 4, Y: 2.5 - 0.57735026919}, P2: engo.Point{X: 5, Y: 2.5}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 0, Y: 5}, P2: engo.Point{X: 1.15470053838, Y: 5}},
		engo.Line{P1: engo.Point{X: 1.15470053838, Y: 5}, P2: engo.Point{X: 0.57735026919, Y: 4}},
		engo.Line{P1: engo.Point{X: 0.57735026919, Y: 4}, P2: engo.Point{X: 0, Y: 5}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 5, Y: 5}, P2: engo.Point{X: 5 - 1.15470053838, Y: 5}},
		engo.Line{P1: engo.Point{X: 5 - 1.15470053838, Y: 5}, P2: engo.Point{X: 5 - 0.57735026919, Y: 5}},
		engo.Line{P1: engo.Point{X: 5 - 0.57735026919, Y: 5}, P2: engo.Point{X: 5, Y: 5}},
	}})
	space3 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 10, Y: 10}} //Circle
	space3.AddShape(Shape{Ellipse: Ellipse{Rx: 2.5, Cx: 2.5, Ry: 2.5, Cy: 2.5}})
	space4 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 20, Y: 0}} //Ellipse
	space4.AddShape(Shape{Ellipse: Ellipse{Rx: 2.5, Cx: 2.5, Ry: 5, Cy: 2.5}})
	space5 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 0, Y: 20}, Rotation: 30} //Ellipse-Rotated
	space5.AddShape(Shape{Ellipse: Ellipse{Rx: 2.5, Cx: 2.5, Ry: 5, Cy: 2.5}})
	scs = append(scs, space0, space1, space2, space3, space4, space5)
	pts := []engo.Point{
		engo.Point{X: 2.5, Y: 2.5},
		engo.Point{X: 12.5, Y: 2.5},
		engo.Point{X: 0.5, Y: 12.5},
		engo.Point{X: 12.5, Y: 12.5},
		engo.Point{X: 22.5, Y: 2.5},
		engo.Point{X: 2.5, Y: 22.5},
	}
	for i := 0; i < len(scs); i++ {
		for j := 0; j < len(pts); j++ {
			if i == j {
				if !scs[i].Contains(pts[j]) {
					println(scs[i].Position.X, scs[i].Position.Y, pts[j].X, pts[j].Y)
					t.Errorf("Space Component %v did not contain point %v, but should.", i, j)
				}
				continue
			}
			if scs[i].Contains(pts[j]) {
				t.Errorf("Space Component %v contained point %v, but should not.", i, j)
			}
		}
	}
}

func TestSpaceComponent_Overlaps(t *testing.T) {
	tol := engo.Point{X: 1e-5, Y: 1e-5}
	scs := []SpaceComponent{}
	space0 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 0, Y: 0}}  //AABB
	space1 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 10, Y: 0}} //Triangle
	space1.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 0, Y: 0}, P2: engo.Point{X: 0, Y: 5}},
		engo.Line{P1: engo.Point{X: 0, Y: 5}, P2: engo.Point{X: 5, Y: 5}},
		engo.Line{P1: engo.Point{X: 5, Y: 5}, P2: engo.Point{X: 0, Y: 0}},
	}})
	space2 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 0, Y: 10}} //Multi-Shape
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 2.5, Y: 0}, P2: engo.Point{X: 2.5 + 0.57735026919, Y: 1}},
		engo.Line{P1: engo.Point{X: 2.5 + 0.57735026919, Y: 1}, P2: engo.Point{X: 2.5 - 0.57735026919, Y: 1}},
		engo.Line{P1: engo.Point{X: 2.5 - 0.57735026919, Y: 1}, P2: engo.Point{X: 2.5, Y: 0}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 0, Y: 2.5}, P2: engo.Point{X: 1, Y: 2.5 + 0.57735026919}},
		engo.Line{P1: engo.Point{X: 1, Y: 2.5 + 0.57735026919}, P2: engo.Point{X: 1, Y: 2.5 - 0.57735026919}},
		engo.Line{P1: engo.Point{X: 1, Y: 2.5 - 0.57735026919}, P2: engo.Point{X: 0, Y: 2.5}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 5, Y: 2.5}, P2: engo.Point{X: 4, Y: 2.5 + 0.57735026919}},
		engo.Line{P1: engo.Point{X: 4, Y: 2.5 + 0.57735026919}, P2: engo.Point{X: 4, Y: 2.5 - 0.57735026919}},
		engo.Line{P1: engo.Point{X: 4, Y: 2.5 - 0.57735026919}, P2: engo.Point{X: 5, Y: 2.5}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 0, Y: 5}, P2: engo.Point{X: 1.15470053838, Y: 5}},
		engo.Line{P1: engo.Point{X: 1.15470053838, Y: 5}, P2: engo.Point{X: 0.57735026919, Y: 4}},
		engo.Line{P1: engo.Point{X: 0.57735026919, Y: 4}, P2: engo.Point{X: 0, Y: 5}},
	}})
	space2.AddShape(Shape{Lines: []engo.Line{
		engo.Line{P1: engo.Point{X: 5, Y: 5}, P2: engo.Point{X: 5 - 1.15470053838, Y: 5}},
		engo.Line{P1: engo.Point{X: 5 - 1.15470053838, Y: 5}, P2: engo.Point{X: 5 - 0.57735026919, Y: 5}},
		engo.Line{P1: engo.Point{X: 5 - 0.57735026919, Y: 5}, P2: engo.Point{X: 5, Y: 5}},
	}})
	space3 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 10, Y: 10}} //Circle
	space3.AddShape(Shape{Ellipse: Ellipse{Rx: 2.5, Cx: 2.5, Ry: 2.5, Cy: 2.5}})
	space4 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 20, Y: 0}} //Ellipse
	space4.AddShape(Shape{Ellipse: Ellipse{Rx: 2.5, Cx: 2.5, Ry: 5, Cy: 2.5}})
	space5 := SpaceComponent{Width: 5, Height: 5, Position: engo.Point{X: 0, Y: 20}, Rotation: 30} //Ellipse-Rotated
	space5.AddShape(Shape{Ellipse: Ellipse{Rx: 2.5, Cx: 2.5, Ry: 5, Cy: 2.5}})
	scs = append(scs, space0, space1, space2, space3, space4, space5)

	//nothing overlaps at first
	for i := 1; i < len(scs); i++ {
		if o, _ := scs[i].Overlaps(scs[0], tol, tol); o {
			t.Errorf("Overlap detected when none of the scs should overlap. Was on Space Component %v", i)
		}
	}

	for i := 1; i < len(scs); i++ {
		scs[0].Position = scs[i].Position
		if o, _ := scs[0].Overlaps(scs[i], tol, tol); !o {
			t.Errorf("SpaceComponents 0 and %v did not overlap but should have", i)
		}
	}

	if o, _ := scs[4].Overlaps(scs[5], tol, tol); o {
		t.Error("SpaceComponents 4 and 5 overlapped but should not have")
	}
	scs[4].Position = scs[5].Position
	if o, _ := scs[4].Overlaps(scs[5], tol, tol); !o {
		t.Error("SpaceComponents 4 and 5 did not overlap but should have")
	}

	scs[4].Position = scs[2].Position
	if o, _ := scs[4].Overlaps(scs[2], tol, tol); !o {
		t.Error("SpaceComponents 4 and 2 did not overlap but should have")
	}

	scs[0].Position = engo.Point{X: 0, Y: 0}
	scs[1].Position = engo.Point{X: 10, Y: 0}
	if o, _ := scs[0].Overlaps(scs[1], tol, tol); o {
		t.Error("SpaceComponents 1 and 0 overlap but should not")
	}
}

func TestSpaceComponent_Corners(t *testing.T) {
	space1 := SpaceComponent{Width: 1, Height: 1}
	exp1 := [4]engo.Point{engo.Point{X: 0, Y: 0}, engo.Point{X: 1, Y: 0}, engo.Point{X: 0, Y: 1}, engo.Point{X: 1, Y: 1}}
	act1 := space1.Corners()
	for i := 0; i < 4; i++ {
		assert.True(t, exp1[i].Equal(act1[i]), fmt.Sprintf("corner %d did not match for rotation %f (got %v expected %v)", i, space1.Rotation, act1[i], exp1[i]))
	}

	space2 := SpaceComponent{Width: 1, Height: 1, Rotation: 90}
	exp2 := [4]engo.Point{engo.Point{X: 0, Y: 0}, engo.Point{X: 0, Y: 1}, engo.Point{X: -1, Y: 0}, engo.Point{X: -1, Y: 1}}
	act2 := space2.Corners()
	for i := 0; i < 4; i++ {
		assert.True(t, exp2[i].Equal(act2[i]), fmt.Sprintf("corner %d did not match for rotation %f (got %v expected %v)", i, space2.Rotation, act2[i], exp2[i]))
	}
}

func TestSpaceComponent_AABB(t *testing.T) {
	space1 := SpaceComponent{Width: 1, Height: 1}
	exp1 := engo.AABB{Min: engo.Point{X: 0, Y: 0}, Max: engo.Point{X: 1, Y: 1}}
	act1 := space1.AABB()
	if !exp1.Min.Equal(act1.Min) || !exp1.Max.Equal(act1.Max) {
		t.Errorf("Space 1's AABB %v did not match expected %v", act1, exp1)
	}

	space2 := SpaceComponent{Width: 1, Height: 1, Rotation: 45}
	exp2 := engo.AABB{Min: engo.Point{X: -0.70710677, Y: 0}, Max: engo.Point{X: 0.70710677, Y: 1.4142135}}
	act2 := space2.AABB()
	if !exp2.Min.Equal(act2.Min) || !exp2.Max.Equal(act2.Max) {
		t.Errorf("Space2's AABB %v did not match expected %v", act2, exp2)
	}
}

const (
	Ball = 1 << iota
	Bat
)

//Test GroupSolid working
func Test_GroupSolid(t *testing.T) {
	//All items in same place, have to collide
	CE := func(m, g CollisionGroup) collisionEntity {
		nb := ecs.NewBasic()
		return collisionEntity{
			BasicEntity: &nb,
			CollisionComponent: &CollisionComponent{
				Main:  m,
				Group: g,
			},
			//All objects in same position
			SpaceComponent: &SpaceComponent{
				Position: engo.Point{X: 10, Y: 10},
				Width:    50,
				Height:   50,
				Rotation: 0,
			},
		}
	}
	ents := []collisionEntity{
		CE(Ball, 0),     //The Ball
		CE(Bat, Ball),   //The Batt
		CE(0, Ball|Bat), //The Wall
		CE(0, 0),        //Ghost Should not collide with anything
	}
	sys := CollisionSystem{
		entities: ents,
		Solids:   Ball, //Only the ball should move as Solid
	}
	sys.Update(0.01)

	if ents[0].Position == ents[1].Position {
		t.Log("Ball should collide Solid")
		t.Fail()
	}

	if ents[3].Collides != 0 {
		t.Log("Ghost should not collide with anything")
		t.Fail()
	}
	for i := 0; i < 2; i++ { //Ball and Bat
		if ents[i].Collides == 0 {
			t.Logf("object %d should collides", i)
			t.Fail()
		}
	}
}

func TestSpaceComponent_Center(t *testing.T) {
	components := []SpaceComponent{
		SpaceComponent{Width: 0, Height: 0},
		SpaceComponent{Width: 100, Height: 100},
		SpaceComponent{Width: 100, Height: 200},
		SpaceComponent{Width: 100, Height: 200, Rotation: 45},
	}
	points := []engo.Point{
		engo.Point{X: 10, Y: 10},
		engo.Point{X: 50, Y: 50},
		engo.Point{X: 10, Y: 50},
		engo.Point{X: 99, Y: 99},
	}

	for _, sc := range components {
		for _, p := range points {
			sc.SetCenter(p)
			c := sc.Center()
			assert.True(t, c.Equal(p), fmt.Sprintf("center %v should be equal to point %v", c, p))
		}
	}
}

func TestShape_Project(t *testing.T) {
	shapes := []Shape{
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 0},
					P2: engo.Point{X: 10, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 0},
					P2: engo.Point{X: 10, Y: 10},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 10},
					P2: engo.Point{X: 0, Y: 0},
				},
			},
		}, //triangle
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 0},
					P2: engo.Point{X: 10, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 0},
					P2: engo.Point{X: 10, Y: 10},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 10},
					P2: engo.Point{X: 0, Y: 10},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 10},
					P2: engo.Point{X: 0, Y: 0},
				},
			},
		}, //square
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 0},
					P2: engo.Point{X: 10, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 0},
					P2: engo.Point{X: 10, Y: 5},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 5},
					P2: engo.Point{X: 0, Y: 5},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 5},
					P2: engo.Point{X: 0, Y: 0},
				},
			},
		}, //rectangle
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 3},
					P2: engo.Point{X: 3, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 3, Y: 0},
					P2: engo.Point{X: 6, Y: 3},
				},
				engo.Line{
					P1: engo.Point{X: 6, Y: 3},
					P2: engo.Point{X: 6, Y: 6},
				},
				engo.Line{
					P1: engo.Point{X: 6, Y: 6},
					P2: engo.Point{X: 0, Y: 6},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 6},
					P2: engo.Point{X: 0, Y: 3},
				},
			},
		}, //pentagon
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 1, Y: 0},
					P2: engo.Point{X: 3, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 3, Y: 0},
					P2: engo.Point{X: 4, Y: 2},
				},
				engo.Line{
					P1: engo.Point{X: 4, Y: 2},
					P2: engo.Point{X: 3, Y: 4},
				},
				engo.Line{
					P1: engo.Point{X: 3, Y: 4},
					P2: engo.Point{X: 1, Y: 4},
				},
				engo.Line{
					P1: engo.Point{X: 1, Y: 4},
					P2: engo.Point{X: 0, Y: 2},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 2},
					P2: engo.Point{X: 1, Y: 0},
				},
			},
		}, //hexagon
		Shape{
			Ellipse: Ellipse{
				Rx: 4,
				Ry: 4,
				Cx: 4,
				Cy: 4,
			},
		}, //circle
		Shape{
			Ellipse: Ellipse{
				Rx: 4,
				Ry: 8,
				Cx: 4,
				Cy: 4,
			},
		}, //ellipse
	}
	pts := []engo.Point{
		engo.Point{X: 1, Y: 0},  //quad1
		engo.Point{X: 0, Y: 1},  //quad2
		engo.Point{X: -1, Y: 0}, //quad3
		engo.Point{X: 0, Y: -1}, //quad4
	}
	exp := [][]float32{
		[]float32{0, 10},                    //0 0
		[]float32{0, 10},                    //0 1
		[]float32{-10, 0},                   //0 2
		[]float32{-10, 0},                   //0 3
		[]float32{0, 10},                    //1 0
		[]float32{0, 10},                    //1 1
		[]float32{-10, 0},                   //1 2
		[]float32{-10, 0},                   //1 3
		[]float32{0, 10},                    //2 0
		[]float32{0, 5},                     //2 1
		[]float32{-10, 0},                   //2 2
		[]float32{-5, 0},                    //2 3
		[]float32{0, 6},                     //3 0
		[]float32{0, 6},                     //3 1
		[]float32{-6, 0},                    //3 2
		[]float32{-6, 0},                    //3 3
		[]float32{0, 4},                     //4 0
		[]float32{0, 4},                     //4 1
		[]float32{-4, 0},                    //4 2
		[]float32{-4, 0},                    //4 3
		[]float32{0.0078930855, 7.992107},   //5 0
		[]float32{0.03154111, 8},            //5 1
		[]float32{-7.992107, -0.0078930855}, //5 2
		[]float32{-8, -0.03154111},          //5 3
		[]float32{0.0078930855, 7.992107},   //6 0
		[]float32{-3.9369178, 12},           //6 1
		[]float32{-7.992107, -0.0078930855}, //6 2
		[]float32{-12, 3.9369178},           //6 3
		[]float32{0.0078930855, 7.992107},   //7 0
		[]float32{0.0078930855, 7.992107},   //7 1
		[]float32{0.0078930855, 7.992107},   //7 2
		[]float32{0.0078930855, 7.992107},   //7 3
	}
	for i, shape := range shapes {
		for j, pt := range pts {
			min, max := shape.Project(pt, SpaceComponent{})
			if !engo.FloatEqual(min, exp[i*len(pts)+j][0]) || !engo.FloatEqual(max, exp[i*len(pts)+j][1]) {
				t.Errorf("Shape projection did not match expected values!\nFor shape %v and point %v\n The min was %v and expected %v, and the max was %v and expected %v", i, j, min, exp[i*len(pts)+j][0], max, exp[i*len(pts)+j][1])
			}
		}
	}
}

func TestSpaceComponentAddShape(t *testing.T) {
	shapes := []Shape{
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 0},
					P2: engo.Point{X: 10, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 0},
					P2: engo.Point{X: 10, Y: 10},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 10},
					P2: engo.Point{X: 0, Y: 0},
				},
			},
		}, //triangle
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 0},
					P2: engo.Point{X: 10, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 0},
					P2: engo.Point{X: 10, Y: 10},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 10},
					P2: engo.Point{X: 0, Y: 10},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 10},
					P2: engo.Point{X: 0, Y: 0},
				},
			},
		}, //square
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 0},
					P2: engo.Point{X: 10, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 0},
					P2: engo.Point{X: 10, Y: 5},
				},
				engo.Line{
					P1: engo.Point{X: 10, Y: 5},
					P2: engo.Point{X: 0, Y: 5},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 5},
					P2: engo.Point{X: 0, Y: 0},
				},
			},
		}, //rectangle
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 0, Y: 3},
					P2: engo.Point{X: 3, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 3, Y: 0},
					P2: engo.Point{X: 6, Y: 3},
				},
				engo.Line{
					P1: engo.Point{X: 6, Y: 3},
					P2: engo.Point{X: 6, Y: 6},
				},
				engo.Line{
					P1: engo.Point{X: 6, Y: 6},
					P2: engo.Point{X: 0, Y: 6},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 6},
					P2: engo.Point{X: 0, Y: 3},
				},
			},
		}, //pentagon
		Shape{
			Lines: []engo.Line{
				engo.Line{
					P1: engo.Point{X: 1, Y: 0},
					P2: engo.Point{X: 3, Y: 0},
				},
				engo.Line{
					P1: engo.Point{X: 3, Y: 0},
					P2: engo.Point{X: 4, Y: 2},
				},
				engo.Line{
					P1: engo.Point{X: 4, Y: 2},
					P2: engo.Point{X: 3, Y: 4},
				},
				engo.Line{
					P1: engo.Point{X: 3, Y: 4},
					P2: engo.Point{X: 1, Y: 4},
				},
				engo.Line{
					P1: engo.Point{X: 1, Y: 4},
					P2: engo.Point{X: 0, Y: 2},
				},
				engo.Line{
					P1: engo.Point{X: 0, Y: 2},
					P2: engo.Point{X: 1, Y: 0},
				},
			},
		}, //hexagon
		Shape{
			Ellipse: Ellipse{
				Rx: 4,
				Ry: 4,
				Cx: 4,
				Cy: 4,
			},
		}, //circle
		Shape{
			Ellipse: Ellipse{
				Rx: 4,
				Ry: 8,
				Cx: 4,
				Cy: 4,
			},
		}, //ellipse
	}
	sc := SpaceComponent{}
	if len(sc.hitboxes) != 0 {
		t.Errorf("SpaceComponent defaulted with hitboxes populated. len was %v", len(sc.hitboxes))
	}
	for i, shape := range shapes {
		sc.AddShape(shape)
		if len(sc.hitboxes) != i+1 {
			t.Errorf("SpaceComponent did not add shape to hitbox. At shape %v, len was %v", i, sc.hitboxes)
		}
	}
}

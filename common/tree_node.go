package core

import (
	"fmt"
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// TreeNodeComponent has local transformation data and heirerachy knowledge
type TreeNodeComponent struct {
	LocPos engo.Point
	LocRot float32
	LocScl engo.Point
}

// TreeNodeEntity affects translation/rotation/scale of child nodes
type TreeNodeEntity struct {
	ecs.BasicEntity
	*common.SpaceComponent
	*common.RenderComponent
	*TreeNodeComponent

	Name string

	children []*TreeNodeEntity
	parent   *TreeNodeEntity
}

// FindByEntityID is a DFS of the Node Tree for a node with the given id
func (n *TreeNodeEntity) FindByEntityID(id uint64) *TreeNodeEntity {

	// Check self
	if n.BasicEntity.ID() == id {
		return n
	}

	// Recursively check child subtrees
	for _, child := range n.children {
		foundInSubtree := child.FindByEntityID(id)
		if foundInSubtree != nil {
			return foundInSubtree
		}
	}

	// Not found
	return nil
}

// AdoptChild ...
func (n *TreeNodeEntity) AdoptChild(child *TreeNodeEntity) {
	n.children = append(n.children, child)
	child.parent = n
}

// AbandonChildEnt ...
func (n *TreeNodeEntity) AbandonChildEnt(entID uint64) {
	var di int = -1
	var child *TreeNodeEntity

	for i, c := range n.children {
		if c.ID() == entID {
			di = i
			child = c
			break
		}
	}

	if di >= 0 {
		n.children = append(n.children[:di], n.children[di+1:]...)
	}

	child.parent = nil
}

func debugstr(id uint64, lvl uint32, str string) {
	// if !(id == 3 || id == 4 || id == 5 || id == 6 || id == 7) {
	if !(id == 5 || id == 6 || id == 7) {
		return
	}

	var pfx string
	var i uint32

	for ; i < lvl; i++ {
		pfx += "    "
	}

	fmt.Printf("%s%s\n", pfx, str)
}

// SubmitLocalsToComponents calculates the node-point scale, rotate, and transform and submits the values to node components
func (n *TreeNodeEntity) SubmitLocalsToComponents(pM *engo.Matrix, pScl engo.Point, pRot float32, lvl uint32) {

	debugstr(n.ID(), lvl, fmt.Sprintf("[ %d:%s ] Submit Locals to Components", n.ID(), n.Name))

	// Calculate new child s, r, t
	cScl := engo.Point{X: pScl.X * n.LocScl.X, Y: pScl.Y * n.LocScl.Y}
	cRot := pRot + n.LocRot
	debugstr(n.ID(), lvl, fmt.Sprintf("  - loc (%.2f, %.2f), rot %.2f deg, scl (%.2f, %.2f)", n.LocPos.X, n.LocPos.Y, n.LocRot, n.LocScl.X, n.LocScl.Y))
	debugstr(n.ID(), lvl, fmt.Sprintf("  - rotation: %.2f deg + %.2f deg = %.2f deg", pRot, n.LocRot, cRot))

	// Calculate offset for sprites
	var dX, dY float32
	if n.SpaceComponent != nil && n.RenderComponent != nil {

		// Calc delta to rotate about the center
		w := float64(n.RenderComponent.Drawable.Width())
		h := float64(n.RenderComponent.Drawable.Height())

		sX, sY := w*float64(cScl.X), h*float64(cScl.Y)
		debugstr(n.ID(), lvl, fmt.Sprintf("  - is a sprite (%.2fpx, %.2fpx) * (%.2f, %.2f) = (%.2fpx, %.2fpx)", w, h, cScl.X, cScl.Y, sX, sY))

		if !engo.FloatEqual(n.LocRot, 0) {
			theta := math.Atan(sY / sX)                        // Original angle
			omega := ((float64(cRot)) * math.Pi / 180) + theta // Rotation plus original angle
			debugstr(n.ID(), lvl, fmt.Sprintf("  - angles tan-1(%.2f/%.2f), %.2f deg + %.2f deg = %.2f deg", sY, sX, theta*(180/math.Pi), cRot, omega*(180/math.Pi)))
			sin, cos := math.Sincos(omega)
			r := math.Sqrt(float64((sX*sX + sY*sY) / 4))
			dX = float32(r * cos)
			dY = float32(r * sin)
			debugstr(n.ID(), lvl, fmt.Sprintf("  - rotated sprite (%.2fpx^2 + %.2fpx^2 = %.2fpx^2), (%.2fdeg) => (%.2fpx, %.2fpx)", sX*0.5, sY*0.5, r, omega*180/math.Pi, dX, dY))
		} else {
			dX = float32(sX / 2)
			dY = float32(sY / 2)
			debugstr(n.ID(), lvl, fmt.Sprintf("  - non rotated sprite  => (%.2fpx, %.2fpx)", dX, dY))
		}
	}


	// Move local
	cPos := engo.Point{X: n.LocPos.X, Y: n.LocPos.Y}
	debugstr(n.ID(), lvl, fmt.Sprintf("  - translate  (befre)=> (%.2fpx, %.2fpx)", n.LocPos.X, n.LocPos.Y))
	cPos.MultiplyMatrixVector(pM).Add(engo.Point{X: -dX, Y: -dY})
	debugstr(n.ID(), lvl, fmt.Sprintf("  - translate  (after)=> (%.2fpx, %.2fpx)", cPos.X, cPos.Y))
	debugstr(n.ID(), lvl, fmt.Sprintf("  - gonna calc the matrix"))

	// Apply scale, rotate, translate, and then apply parent transofrm
	lM := engo.IdentityMatrix().Translate(n.LocPos.X, n.LocPos.Y).Rotate(n.LocRot).Scale(n.LocScl.X, n.LocScl.Y).Multiply(pM)

	// Submit Locals to Components
	if n.RenderComponent != nil && n.SpaceComponent != nil {

		// Get the s, r, t components of the end matrix
		gPosX, gPosY := cPos.X, cPos.Y
		gSclX, gSclY := cScl.X, cScl.Y

		// Set global s, r, t
		n.SpaceComponent.Position = engo.Point{X: gPosX, Y: gPosY}
		n.RenderComponent.Scale = engo.Point{X: gSclX, Y: gSclY}
		n.SpaceComponent.Rotation = cRot
	}

	// Instruct child nodes to submit calculated values to components
	for _, c := range n.children {
		c.SubmitLocalsToComponents(lM, cScl, cRot, lvl+1)
	}

}

// String ...
func (n *TreeNodeEntity) String() string {
	return n.string(0)
}

func (n *TreeNodeEntity) string(depth uint32) string {
	pfx := ""
	for d := uint32(0); d < depth; d++ {
		pfx += "--"
	}

	pfx += fmt.Sprintf("%s> [ %d ] %s\n", pfx, n.ID(), n.Name)
	for _, c := range n.children {
		pfx += c.string(depth + 1)
	}

	return pfx
}

// TreeNodeSystem ...
type TreeNodeSystem struct {
	root *TreeNodeEntity
}

// New ...
func (s *TreeNodeSystem) New(w *ecs.World) {
	// Set root to new component at the origin
	rootEnt := ecs.NewBasic()
	s.root = &TreeNodeEntity{
		Name:        "root",
		BasicEntity: rootEnt,
		TreeNodeComponent: &TreeNodeComponent{
			LocPos: engo.Point{X: 0, Y: 0},
			LocRot: 0,
			LocScl: engo.Point{X: 1, Y: 1},
		},
	}

	// Print the new tree
	fmt.Println(s.root)
	fmt.Println()
}

// MoveNodeToParent moves the node matching ent to node matching parent
func (s *TreeNodeSystem) MoveNodeToParent(entID, parentID uint64) {
	node := s.root.FindByEntityID(entID)
	prnt := s.root.FindByEntityID(parentID)
	if node == nil || prnt == nil {
		panic(fmt.Sprintf("Could not parent (%d) to child (%d).  One or both not found.\n%s\n", parentID, entID, s.root))
	}

	node.parent.AbandonChildEnt(entID)
	prnt.AdoptChild(node)

	// Print the new tree
	fmt.Println(s.root)
	fmt.Println()
}

// Update applies SRT (Scale-Rotate-Translate) Transforms to full tree, passing down parent SRT Transform matrix to children
func (s *TreeNodeSystem) Update(dt float32) {
	s.root.SubmitLocalsToComponents(engo.IdentityMatrix(), engo.Point{X: 1, Y: 1}, 0, 0)
	fmt.Println()
}

// Remove ...
func (s *TreeNodeSystem) Remove(e ecs.BasicEntity) {
	// Get node being removed
	n := s.root.FindByEntityID(e.ID())
	if n == nil {
		return
	}

	// If it has children, we don't care.  Children will just be GC'd along with node

	// Remove it from the parent's slice
	if n.parent != nil {
		n.parent.AbandonChildEnt(e.ID())
	}
}

// GetRoot ...
func (s *TreeNodeSystem) GetRoot() *TreeNodeEntity {
	return s.root
}

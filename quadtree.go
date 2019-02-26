package engo

import (
	"sync"
)

var (
	quadtreeNodePool *sync.Pool
	nodeDataPool     *sync.Pool
)

const minQuadtreeCellSize = 0.01

func init() {
	quadtreeNodePool = &sync.Pool{
		New: func() interface{} {
			return new(quadtreeNode)
		},
	}
	nodeDataPool = &sync.Pool{
		New: func() interface{} {
			return new(quadtreeNodeData)
		},
	}
}

func aabbOverlaps(a, b AABB) bool {
	// a is left of b
	if a.Max.X < b.Min.X {
		return false
	}

	// a is right of b
	if a.Min.X > b.Max.X {
		return false
	}

	// a is above b
	if a.Max.Y < b.Min.Y {
		return false
	}

	// a is below b
	if a.Min.Y > b.Max.Y {
		return false
	}

	// The two overlap
	return true
}

func aabbWidth(x AABB) float32 {
	return x.Max.X - x.Min.X
}

func aabbHeight(x AABB) float32 {
	return x.Max.Y - x.Min.Y
}

func aabbRect(x, y, width, height float32) AABB {
	return AABB{
		Min: Point{
			X: x,
			Y: y,
		},
		Max: Point{
			X: x + width,
			Y: y + height,
		},
	}
}

type quadtreeNodeData struct {
	Value AABBer
	AABB  AABB
}

type quadtreeNode struct {
	Bounds   AABB
	Level    int
	Objects  []*quadtreeNodeData
	hasNodes bool
	Nodes    [4]*quadtreeNode
	Tree     *Quadtree
}

// Quadtree implementation which can store AABBer values
type Quadtree struct {
	MaxObjects int // Maximum objects a node can hold before splitting into 4 subnodes
	MaxLevels  int // Total max levels inside root Quadtree
	root       *quadtreeNode
	usePool    bool
	Total      int
}

func calcMaxLevel(width, height float32) int {
	res := 0
	for width > minQuadtreeCellSize && height > minQuadtreeCellSize {
		res++
		width, height = width/2, height/2
	}
	return res
}

// NewQuadtree creates a new quadtree for the given bounds.
// When setting usePool to true, the internal values will be taken from a sync.Pool which reduces the allocation overhead.
// maxObjects tells the tree how many objects should be stored within a level before the quadtree cell is split.
func NewQuadtree(bounds AABB, usePool bool, maxObjects int) *Quadtree {
	qt := &Quadtree{MaxObjects: maxObjects, usePool: usePool}
	qt.root = qt.newNode(bounds, 0)
	qt.MaxLevels = calcMaxLevel(aabbWidth(bounds), aabbHeight(bounds))
	return qt
}

// Destroy frees the nodes if the Quadtree uses the node pool
func (qt *Quadtree) Destroy() {
	qt.freeQuadtreeNode(qt.root)
	qt.root = nil
}

func (qt *Quadtree) newNode(bounds AABB, level int) (node *quadtreeNode) {
	if qt.usePool {
		node = quadtreeNodePool.Get().(*quadtreeNode)
	} else {
		node = new(quadtreeNode)
	}

	node.Tree = qt
	node.Bounds = bounds
	node.Level = level
	return node
}

func (qt *Quadtree) newQuadtreeNodeData(item AABBer, r AABB) *quadtreeNodeData {
	if qt.usePool {
		d := nodeDataPool.Get().(*quadtreeNodeData)
		d.AABB = r
		d.Value = item
		return d
	}
	return &quadtreeNodeData{item, r}
}

func (qt *Quadtree) freeQuadtreeNodeData(n *quadtreeNodeData) {
	if !qt.usePool {
		return
	}
	nodeDataPool.Put(n)
}

func (qt *Quadtree) freeQuadtreeNode(n *quadtreeNode) {
	if !qt.usePool {
		return
	}
	if n.hasNodes {
		for i, child := range n.Nodes {
			qt.freeQuadtreeNode(child)
			n.Nodes[i] = nil
		}
	}

	if n.Objects != nil {
		for _, o := range n.Objects {
			qt.freeQuadtreeNodeData(o)
		}
	}
	n.Objects = nil
	n.Tree = nil
	n.hasNodes = false
	quadtreeNodePool.Put(n)
}

// split - split the node into 4 subnodes
func (qt *quadtreeNode) split() {
	if qt.hasNodes {
		return
	}
	qt.hasNodes = true

	nextLevel := qt.Level + 1
	subWidth := aabbWidth(qt.Bounds) / 2
	subHeight := aabbHeight(qt.Bounds) / 2
	x := qt.Bounds.Min.X
	y := qt.Bounds.Min.Y

	//top right node (0)
	qt.Nodes[0] = qt.Tree.newNode(aabbRect(x+subWidth, y, subWidth, subHeight), nextLevel)

	//top left node (1)
	qt.Nodes[1] = qt.Tree.newNode(aabbRect(x, y, subWidth, subHeight), nextLevel)

	//bottom left node (2)
	qt.Nodes[2] = qt.Tree.newNode(aabbRect(x, y+subHeight, subWidth, subHeight), nextLevel)

	//bottom right node (3)
	qt.Nodes[3] = qt.Tree.newNode(aabbRect(x+subWidth, y+subHeight, subWidth, subHeight), nextLevel)
}

func (qt *quadtreeNode) isEmpty() bool {
	return len(qt.Objects) == 0 && !qt.hasNodes
}

func (qt *quadtreeNode) unsplit() {
	for i := 0; i < 4; i++ {
		if !qt.Nodes[i].isEmpty() {
			return
		}
	}
	for i := 0; i < 4; i++ {
		qt.Tree.freeQuadtreeNode(qt.Nodes[i])
		qt.Nodes[i] = nil
	}
	qt.hasNodes = false
}

// getIndex - Determine which quadrant the object belongs to (0-3)
func (qt *quadtreeNode) getIndex(pRect AABB) int {
	horzMidpoint := qt.Bounds.Min.X + (aabbWidth(qt.Bounds) / 2)
	vertMidpoint := qt.Bounds.Min.Y + (aabbHeight(qt.Bounds) / 2)

	//pRect can completely fit within the top quadrants
	topQuadrant := (pRect.Min.Y < vertMidpoint) && (pRect.Max.Y < vertMidpoint)

	//pRect can completely fit within the bottom quadrants
	bottomQuadrant := (pRect.Min.Y > vertMidpoint)

	//pRect can completely fit within the left quadrants
	if (pRect.Min.X < horzMidpoint) && (pRect.Max.X < horzMidpoint) {
		if topQuadrant {
			return 1
		} else if bottomQuadrant {
			return 2
		}
	} else if pRect.Min.X > horzMidpoint {
		//pRect can completely fit within the right quadrants
		if topQuadrant {
			return 0
		} else if bottomQuadrant {
			return 3
		}
	}

	return -1 // index of the subnode (0-3), or -1 if pRect cannot completely fit within a subnode and is part of the parent node
}

// Insert inserts the given item to the quadtree
func (qt *Quadtree) Insert(item AABBer) {
	qt.Total++
	pRect := item.AABB()
	qt.root.Insert(qt.newQuadtreeNodeData(item, pRect))
}

func (qt *quadtreeNode) Insert(item *quadtreeNodeData) {
	if qt.hasNodes {
		index := qt.getIndex(item.AABB)

		if index != -1 {
			qt.Nodes[index].Insert(item)
			return
		}
	}

	// If we don't subnodes within the Quadtree
	qt.Objects = append(qt.Objects, item)

	// If total objects is greater than max objects and level is less than max levels
	if (len(qt.Objects) > qt.Tree.MaxObjects) && (qt.Tree.MaxLevels <= 0 || qt.Level < qt.Tree.MaxLevels) {
		// split if we don't already have subnodes
		if !qt.hasNodes {
			qt.split()
		}

		// Add all objects to there corresponding subNodes
		for i := 0; i < len(qt.Objects); {
			object := qt.Objects[i] // Get the object out of the slice
			bounds := object.AABB
			index := qt.getIndex(bounds)
			if index != -1 {
				qt.Objects = append(qt.Objects[:i], qt.Objects[i+1:]...) // Remove the object from the slice
				qt.Nodes[index].Insert(object)
			} else {
				i++
			}
		}
	}
}

func (qt *quadtreeNode) Remove(item AABBer, pRect AABB) {
	if qt.hasNodes {
		index := qt.getIndex(pRect)
		if index != -1 {
			qt.Nodes[index].Remove(item, pRect)
			qt.unsplit()
			return
		}
	}
	for i := 0; i < len(qt.Objects); i++ {
		if qt.Objects[i].Value == item {
			qt.Tree.freeQuadtreeNodeData(qt.Objects[i])
			qt.Objects = append(qt.Objects[:i], qt.Objects[i+1:]...) // Remove the object from the slice
			return
		}
	}
}

// Remove removes the given item from the quadtree
func (qt *Quadtree) Remove(item AABBer) {
	bounds := item.AABB()
	qt.root.Remove(item, bounds)
}

// Retrieve returns all objects that could collide with the given bounding box
func (qt *quadtreeNode) Retrieve(pRect AABB) []AABBer {
	index := qt.getIndex(pRect)

	// Array with all detected objects
	result := make([]AABBer, len(qt.Objects))
	for i, o := range qt.Objects {
		result[i] = o.Value
	}

	//if we have subnodes ...
	if qt.hasNodes {
		//if pRect fits into a subnode ..
		if index != -1 {
			result = append(result, qt.Nodes[index].Retrieve(pRect)...)
		} else {
			//if pRect does not fit into a subnode, check it against all subnodes
			for i := 0; i < 4; i++ {
				result = append(result, qt.Nodes[i].Retrieve(pRect)...)
			}
		}
	}

	return result

}

// Retrieve returns all objects that could collide with the given bounding box and passing the given filter function.
func (qt *Quadtree) Retrieve(find AABB, filter func(aabb AABBer) bool) []AABBer {
	var foundIntersections []AABBer

	potentials := qt.root.Retrieve(find)
	for _, p := range potentials {
		if aabbOverlaps(find, p.AABB()) && (filter == nil || filter(p)) {
			foundIntersections = append(foundIntersections, p)
		}
	}

	return foundIntersections
}

//Clear removes all items from the quadtree
func (qt *Quadtree) Clear() {
	bounds := qt.root.Bounds
	qt.freeQuadtreeNode(qt.root)
	qt.root = qt.newNode(bounds, 0)
	qt.Total = 0
}

package fauxgl

func NewTreeForMesh(mesh *Mesh, depth int) *Node {
	boxes := make([]Box, len(mesh.Triangles))
	for i, t := range mesh.Triangles {
		boxes[i] = t.BoundingBox()
	}
	node := NewNode(boxes)
	node.Split(depth)
	return node
}

type Node struct {
	Box   Box
	Boxes []Box
	Axis  Axis
	Point float64
	Left  *Node
	Right *Node
}

func NewNode(boxes []Box) *Node {
	box := BoxForBoxes(boxes)
	return &Node{box, boxes, AxisNone, 0, nil, nil}
}

func (node *Node) Leaves(maxDepth int) []Box {
	var result []Box
	if maxDepth == 0 || (node.Left == nil && node.Right == nil) {
		return []Box{node.Box}
	}
	if node.Left != nil {
		result = append(result, node.Left.Leaves(maxDepth-1)...)
	}
	if node.Right != nil {
		result = append(result, node.Right.Leaves(maxDepth-1)...)
	}
	return result
}

func (node *Node) PartitionScore(axis Axis, point float64, side bool) float64 {
	// var left, right Box
	// for _, box := range node.Boxes {
	// 	l, r := box.Partition(axis, point)
	// 	if l && r {
	// 		if side {
	// 			left = left.Extend(box)
	// 		} else {
	// 			right = right.Extend(box)
	// 		}
	// 	} else if l {
	// 		left = left.Extend(box)
	// 	} else if r {
	// 		right = right.Extend(box)
	// 	}
	// }
	l, r := node.Partition(axis, point, side)
	left := BoxForBoxes(l)
	right := BoxForBoxes(r)
	return left.Volume() + right.Volume() - left.Intersection(right).Volume()
}

func (node *Node) Partition(axis Axis, point float64, side bool) (left, right []Box) {
	for _, box := range node.Boxes {
		l, r := box.Partition(axis, point)
		if l && r {
			if side {
				left = append(left, box)
			} else {
				right = append(right, box)
			}
		} else if l {
			left = append(left, box)
		} else if r {
			right = append(right, box)
		}
	}
	// return
	if side {
		outer := BoxForBoxes(left)
		a := right[:0]
		for _, box := range right {
			if outer.ContainsBox(box) {
				left = append(left, box)
			} else {
				a = append(a, box)
			}
		}
		right = a
	} else {
		outer := BoxForBoxes(right)
		a := left[:0]
		for _, box := range left {
			if outer.ContainsBox(box) {
				right = append(right, box)
			} else {
				a = append(a, box)
			}
		}
		left = a
	}
	return
}

func (node *Node) Split(depth int) {
	if depth == 0 {
		return
	}
	box := node.Box
	best := box.Volume()
	bestAxis := AxisNone
	bestPoint := 0.0
	bestSide := false
	const N = 16
	for s := 0; s < 2; s++ {
		side := s == 1
		for i := 1; i < N; i++ {
			p := float64(i) / N
			x := box.Min.X + (box.Max.X-box.Min.X)*p
			y := box.Min.Y + (box.Max.Y-box.Min.Y)*p
			z := box.Min.Z + (box.Max.Z-box.Min.Z)*p
			sx := node.PartitionScore(AxisX, x, side)
			if sx < best {
				best = sx
				bestAxis = AxisX
				bestPoint = x
				bestSide = side
			}
			sy := node.PartitionScore(AxisY, y, side)
			if sy < best {
				best = sy
				bestAxis = AxisY
				bestPoint = y
				bestSide = side
			}
			sz := node.PartitionScore(AxisZ, z, side)
			if sz < best {
				best = sz
				bestAxis = AxisZ
				bestPoint = z
				bestSide = side
			}
		}
	}
	if bestAxis == AxisNone {
		return
	}
	l, r := node.Partition(bestAxis, bestPoint, bestSide)
	node.Axis = bestAxis
	node.Point = bestPoint
	node.Left = NewNode(l)
	node.Right = NewNode(r)
	node.Left.Split(depth - 1)
	node.Right.Split(depth - 1)
	node.Boxes = nil // only needed at leaf nodes
}

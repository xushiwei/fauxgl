package fauxgl

import "math"

var EmptyBox = Box{}

type Box struct {
	Min, Max Vector
}

func BoxForBoxes(boxes []Box) Box {
	if len(boxes) == 0 {
		return EmptyBox
	}
	x0 := boxes[0].Min.X
	y0 := boxes[0].Min.Y
	z0 := boxes[0].Min.Z
	x1 := boxes[0].Max.X
	y1 := boxes[0].Max.Y
	z1 := boxes[0].Max.Z
	for _, box := range boxes {
		x0 = math.Min(x0, box.Min.X)
		y0 = math.Min(y0, box.Min.Y)
		z0 = math.Min(z0, box.Min.Z)
		x1 = math.Max(x1, box.Max.X)
		y1 = math.Max(y1, box.Max.Y)
		z1 = math.Max(z1, box.Max.Z)
	}
	return Box{Vector{x0, y0, z0}, Vector{x1, y1, z1}}
}

func (a Box) Outline() []*Line {
	x0 := a.Min.X
	y0 := a.Min.Y
	z0 := a.Min.Z
	x1 := a.Max.X
	y1 := a.Max.Y
	z1 := a.Max.Z
	return []*Line{
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x0, y0, z1}),
		NewLineForPoints(Vector{x0, y1, z0}, Vector{x0, y1, z1}),
		NewLineForPoints(Vector{x1, y0, z0}, Vector{x1, y0, z1}),
		NewLineForPoints(Vector{x1, y1, z0}, Vector{x1, y1, z1}),
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x0, y1, z0}),
		NewLineForPoints(Vector{x0, y0, z1}, Vector{x0, y1, z1}),
		NewLineForPoints(Vector{x1, y0, z0}, Vector{x1, y1, z0}),
		NewLineForPoints(Vector{x1, y0, z1}, Vector{x1, y1, z1}),
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x1, y0, z0}),
		NewLineForPoints(Vector{x0, y1, z0}, Vector{x1, y1, z0}),
		NewLineForPoints(Vector{x0, y0, z1}, Vector{x1, y0, z1}),
		NewLineForPoints(Vector{x0, y1, z1}, Vector{x1, y1, z1}),
	}
}

func (a Box) Volume() float64 {
	s := a.Size()
	return s.X * s.Y * s.Z
}

func (a Box) Anchor(anchor Vector) Vector {
	return a.Min.Add(a.Size().Mul(anchor))
}

func (a Box) Center() Vector {
	return a.Anchor(Vector{0.5, 0.5, 0.5})
}

func (a Box) Size() Vector {
	return a.Max.Sub(a.Min)
}

func (a Box) Extend(b Box) Box {
	if a == EmptyBox {
		return b
	}
	return Box{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

func (a Box) Contains(b Vector) bool {
	return a.Min.X <= b.X && a.Max.X >= b.X &&
		a.Min.Y <= b.Y && a.Max.Y >= b.Y &&
		a.Min.Z <= b.Z && a.Max.Z >= b.Z
}

func (a Box) ContainsBox(b Box) bool {
	return a.Min.X <= b.Min.X && a.Max.X >= b.Max.X &&
		a.Min.Y <= b.Min.Y && a.Max.Y >= b.Min.Y &&
		a.Min.Z <= b.Min.Z && a.Max.Z >= b.Min.Z
}

func (a Box) Intersects(b Box) bool {
	return !(a.Min.X > b.Max.X || a.Max.X < b.Min.X || a.Min.Y > b.Max.Y ||
		a.Max.Y < b.Min.Y || a.Min.Z > b.Max.Z || a.Max.Z < b.Min.Z)
}

func (a Box) Intersection(b Box) Box {
	min := a.Min.Max(b.Min)
	max := a.Max.Min(b.Max)
	return Box{min, max}
}

func (a Box) Partition(axis Axis, point float64) (left, right bool) {
	switch axis {
	case AxisX:
		left = a.Min.X <= point
		right = a.Max.X >= point
	case AxisY:
		left = a.Min.Y <= point
		right = a.Max.Y >= point
	case AxisZ:
		left = a.Min.Z <= point
		right = a.Max.Z >= point
	}
	return
}

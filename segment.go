package snake

import (
	"math"
)

type Segment interface {
	Intersect(o Segment) (Vec, bool)
}

type Vec struct {
	X, Y float64
}

func (p Vec) Less(o Vec) bool {
	return p.X < o.X && p.Y < o.Y
}

func (p Vec) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Vec) Minus(o Vec) Vec {
	return Vec{
		X: p.X - o.X,
		Y: p.Y - o.Y,
	}
}

func (p Vec) Plus(o Vec) Vec {
	return Vec{
		X: p.X + o.X,
		Y: p.Y + o.Y,
	}
}

func (p Vec) Multiply(t float64) Vec {
	return Vec{
		X: p.X * t,
		Y: p.Y * t,
	}
}

func (p Vec) Distance(o Vec) float64 {
	return o.Minus(p).Length()
}

func (p Vec) CrossProduct(o Vec) float64 {
	return p.X*o.Y - p.Y*o.X
}

func (p Vec) DotProduct(o Vec) float64 {
	return p.X*o.X + p.Y*o.Y
}

type Line struct {
	A, B Vec
}

func (l Line) Dimensions() Vec {
	return l.B.Minus(l.A)
}

func (l Line) Length() float64 {
	return l.Dimensions().Length()
}

func (l Line) CrossProduct() float64 {
	return l.B.CrossProduct(l.A)
}

func (l Line) ContainsPoint(p Vec) bool {
	return (l.A.Distance(p) + l.B.Distance(p)) == l.Length()
}

func (l Line) IntersectLine(o Line) (Vec, bool) {
	if l.A == o.A || l.A == o.B {
		return l.A, true
	}
	if l.B == o.A || l.B == o.B {
		return l.B, true
	}

	da, db := l.Dimensions(), o.Dimensions()
	c := da.CrossProduct(db)

	dab := l.B.Minus(l.A)
	if c == 0 {
		return Vec{}, false
	}

	t := dab.CrossProduct(db) / c
	u := dab.CrossProduct(da) / c
	if (c == 0) || (t < 0 || t > 1) || (u < 0 || u > 1) {
		return Vec{}, false
	}

	return l.A.Plus(da.Multiply((t + u) / 2.0)), true
}

func (l Line) IntersectArc(a Arc) (Vec, bool) {
	return a.IntersectLine(l)
}

func (l Line) Intersect(o Segment) (Vec, bool) {
	switch v := o.(type) {
	case Line:
		return l.IntersectLine(v)
	case Arc:
		return l.IntersectArc(v)
	default:
		panic("unexpected segment type")
	}
}

type Arc struct {
	C       Vec
	R, S, D float64
}

func (a Arc) Point(t float64) Vec {
	return a.C.Plus(Vec{
		X: a.R * math.Cos(t),
		Y: a.R * math.Sin(t),
	})
}

func (a Arc) ContainsAngle(t float64) bool {
	return a.Interval().Contains(t)
}

func (a Arc) Interval() Interval {
	return Interval{
		S: math.Min(a.S, a.S+a.D),
		E: math.Max(a.S+a.D, a.S),
	}
}

func (a Arc) IntersectLine(o Line) (Vec, bool) {
	d := o.Dimensions()
	p := math.Atan(d.X / d.Y)
	q := math.Acos(-(a.C.X + a.C.Y + o.CrossProduct()) / (a.R * d.Length()))

	if t := p - q; a.ContainsAngle(t) {
		return a.Point(t), true
	}
	if t := p + q; a.ContainsAngle(t) {
		return a.Point(t), true
	}

	return Vec{}, false
}

func (a Arc) IntersectArc(o Arc) (Vec, bool) {
	d := a.C.Distance(o.C)
	// arcs don't intersect
	if a.R+o.R > d {
		return Vec{}, false
	}
	// arcs are within each other
	if math.Abs(a.R-o.R) > d {
		return Vec{}, false
	}
	// arcs overlap
	if d == 0 && a.R == o.R {
		if i := a.Interval().IntersectInterval(o.Interval()); !i.Valid() {
			return a.Point(i.S), true
		}
		return Vec{}, false
	}
	// arcs intersect once
	if a.R+o.R == d {
		return a.C.Plus(o.C.Minus(a.C).Multiply(d / a.R)), true
	}
	// arcs intersect twice
	h := (a.R*a.R - o.R*o.R + d*d) / (2 * d)
	p, i := math.Acos(h/a.R), a.Interval()
	if i.Contains(a.S - p) {
		return a.Point(a.S - p), true
	}
	if i.Contains(a.S + p) {
		return a.Point(a.S + p), true
	}
	return Vec{}, false
}

func (a Arc) Intersect(o Segment) (Vec, bool) {
	switch v := o.(type) {
	case Line:
		return a.IntersectLine(v)
	case Arc:
		return a.IntersectArc(v)
	default:
		panic("unexpected segment type")
	}
}

type Interval struct {
	S, E float64
}

func (i Interval) Contains(t float64) bool {
	return t >= i.S && t <= i.E
}

func (i Interval) ContainsInterval(o Interval) bool {
	return i.IntersectInterval(o).Valid()
}

func (i Interval) Valid() bool {
	return i.S <= i.E
}

func (i Interval) IntersectInterval(o Interval) Interval {
	return Interval{
		S: math.Max(i.S, o.S),
		E: math.Min(i.E, o.E),
	}
}

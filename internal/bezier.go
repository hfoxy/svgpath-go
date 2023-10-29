package internal

import "math"

type BezierType string

const (
	BezierTypeCubic     BezierType = "cubic"
	BezierTypeQuadratic BezierType = "quadratic"
)

type Bezier struct {
	a          Point
	b          Point
	c          Point
	d          Point
	bezierType BezierType
	length     float64
}

func NewBezier(a Point, b Point, c Point, d Point) Bezier {
	r := Bezier{
		a: a,
		b: b,
		c: c,
	}

	if d.blank {
		r.bezierType = BezierTypeCubic
		r.d = d
	} else {
		r.bezierType = BezierTypeQuadratic
		r.d = Point{X: 0, Y: 0}
	}

	r.length = r.GetArcLength(
		[]float64{a.X, b.X, c.X, d.X},
		[]float64{a.Y, b.Y, c.Y, d.Y},
		1,
	)

	return r
}

func (b Bezier) GetArcLength(xs []float64, ys []float64, t float64) float64 {
	if b.bezierType == BezierTypeCubic {
		return getCubicArcLength(xs, ys, t)
	}

	return getQuadraticArcLength(xs, ys, t)
}

func (b Bezier) GetPoint(xs []float64, ys []float64, t float64) Point {
	if b.bezierType == BezierTypeCubic {
		return cubicPoint(xs, ys, t)
	}

	return quadraticPoint(xs, ys, t)
}

func (b Bezier) GetDerivative(xs []float64, ys []float64, t float64) Point {
	if b.bezierType == BezierTypeCubic {
		return cubicDerivative(xs, ys, t)
	}

	return quadraticDerivative(xs, ys, t)
}

func (b Bezier) GetTotalLength() float64 {
	return b.length
}

func (b Bezier) GetPointAtLength(pos float64) Point {
	xs := []float64{b.a.X, b.b.X, b.c.X, b.d.X}
	ys := []float64{b.a.Y, b.b.Y, b.c.Y, b.d.Y}
	t := t2length(pos, b.length, func(t float64) float64 {
		return b.GetArcLength(xs, ys, t)
	})

	return b.GetPoint(xs, ys, t)
}

func (b Bezier) GetTangentAtLength(pos float64) Point {
	xs := []float64{b.a.X, b.b.X, b.c.X, b.d.X}
	ys := []float64{b.a.Y, b.b.Y, b.c.Y, b.d.Y}
	t := t2length(pos, b.length, func(t float64) float64 {
		return b.GetArcLength(xs, ys, t)
	})

	derivative := b.GetDerivative(xs, ys, t)
	mdl := math.Sqrt(derivative.X*derivative.X + derivative.Y*derivative.Y)

	var tangent Point
	if mdl > 0 {
		tangent = Point{X: derivative.X / mdl, Y: derivative.Y / mdl}
	} else {
		tangent = Point{X: 0, Y: 0}
	}

	return tangent
}

func (b Bezier) GetPropertiesAtLength(pos float64) PointProperties {
	xs := []float64{b.a.X, b.b.X, b.c.X, b.d.X}
	ys := []float64{b.a.Y, b.b.Y, b.c.Y, b.d.Y}
	t := t2length(pos, b.length, func(t float64) float64 {
		return b.GetArcLength(xs, ys, t)
	})

	derivative := b.GetDerivative(xs, ys, t)
	mdl := math.Sqrt(derivative.X*derivative.X + derivative.Y*derivative.Y)

	var tangent Point
	if mdl > 0 {
		tangent = Point{X: derivative.X / mdl, Y: derivative.Y / mdl}
	} else {
		tangent = Point{X: 0, Y: 0}
	}

	point := b.GetPoint(xs, ys, t)
	return PointProperties{
		X:        point.X,
		Y:        point.Y,
		TangentX: tangent.X,
		TangentY: tangent.Y,
	}
}
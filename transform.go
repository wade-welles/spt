package spt

import (
	"encoding/gob"
	"math"
)

func init() {
	gob.Register(SDFTransform{})
	gob.Register(SDFScale{})
	gob.Register(SDFDistort{})
	gob.Register(SDFMirror{})
	gob.Register(SDFUnion{})
	gob.Register(SDFDifference{})
	gob.Register(SDFIntersection{})
}

type Matrix44 struct {
	X00, X01, X02, X03 float64
	X10, X11, X12, X13 float64
	X20, X21, X22, X23 float64
	X30, X31, X32, X33 float64
}

func (a Matrix44) MulVec3(b Vec3) Vec3 {
	x := a.X00*b.X + a.X01*b.Y + a.X02*b.Z + a.X03
	y := a.X10*b.X + a.X11*b.Y + a.X12*b.Z + a.X13
	z := a.X20*b.X + a.X21*b.Y + a.X22*b.Z + a.X23
	return Vec3{x, y, z}
}

func (a Matrix44) Determinant() float64 {
	return (a.X00*a.X11*a.X22*a.X33 - a.X00*a.X11*a.X23*a.X32 +
		a.X00*a.X12*a.X23*a.X31 - a.X00*a.X12*a.X21*a.X33 +
		a.X00*a.X13*a.X21*a.X32 - a.X00*a.X13*a.X22*a.X31 -
		a.X01*a.X12*a.X23*a.X30 + a.X01*a.X12*a.X20*a.X33 -
		a.X01*a.X13*a.X20*a.X32 + a.X01*a.X13*a.X22*a.X30 -
		a.X01*a.X10*a.X22*a.X33 + a.X01*a.X10*a.X23*a.X32 +
		a.X02*a.X13*a.X20*a.X31 - a.X02*a.X13*a.X21*a.X30 +
		a.X02*a.X10*a.X21*a.X33 - a.X02*a.X10*a.X23*a.X31 +
		a.X02*a.X11*a.X23*a.X30 - a.X02*a.X11*a.X20*a.X33 -
		a.X03*a.X10*a.X21*a.X32 + a.X03*a.X10*a.X22*a.X31 -
		a.X03*a.X11*a.X22*a.X30 + a.X03*a.X11*a.X20*a.X32 -
		a.X03*a.X12*a.X20*a.X31 + a.X03*a.X12*a.X21*a.X30)
}

func (a Matrix44) Inverse() Matrix44 {
	m := Matrix44{}
	d := a.Determinant()
	m.X00 = (a.X12*a.X23*a.X31 - a.X13*a.X22*a.X31 + a.X13*a.X21*a.X32 - a.X11*a.X23*a.X32 - a.X12*a.X21*a.X33 + a.X11*a.X22*a.X33) / d
	m.X01 = (a.X03*a.X22*a.X31 - a.X02*a.X23*a.X31 - a.X03*a.X21*a.X32 + a.X01*a.X23*a.X32 + a.X02*a.X21*a.X33 - a.X01*a.X22*a.X33) / d
	m.X02 = (a.X02*a.X13*a.X31 - a.X03*a.X12*a.X31 + a.X03*a.X11*a.X32 - a.X01*a.X13*a.X32 - a.X02*a.X11*a.X33 + a.X01*a.X12*a.X33) / d
	m.X03 = (a.X03*a.X12*a.X21 - a.X02*a.X13*a.X21 - a.X03*a.X11*a.X22 + a.X01*a.X13*a.X22 + a.X02*a.X11*a.X23 - a.X01*a.X12*a.X23) / d
	m.X10 = (a.X13*a.X22*a.X30 - a.X12*a.X23*a.X30 - a.X13*a.X20*a.X32 + a.X10*a.X23*a.X32 + a.X12*a.X20*a.X33 - a.X10*a.X22*a.X33) / d
	m.X11 = (a.X02*a.X23*a.X30 - a.X03*a.X22*a.X30 + a.X03*a.X20*a.X32 - a.X00*a.X23*a.X32 - a.X02*a.X20*a.X33 + a.X00*a.X22*a.X33) / d
	m.X12 = (a.X03*a.X12*a.X30 - a.X02*a.X13*a.X30 - a.X03*a.X10*a.X32 + a.X00*a.X13*a.X32 + a.X02*a.X10*a.X33 - a.X00*a.X12*a.X33) / d
	m.X13 = (a.X02*a.X13*a.X20 - a.X03*a.X12*a.X20 + a.X03*a.X10*a.X22 - a.X00*a.X13*a.X22 - a.X02*a.X10*a.X23 + a.X00*a.X12*a.X23) / d
	m.X20 = (a.X11*a.X23*a.X30 - a.X13*a.X21*a.X30 + a.X13*a.X20*a.X31 - a.X10*a.X23*a.X31 - a.X11*a.X20*a.X33 + a.X10*a.X21*a.X33) / d
	m.X21 = (a.X03*a.X21*a.X30 - a.X01*a.X23*a.X30 - a.X03*a.X20*a.X31 + a.X00*a.X23*a.X31 + a.X01*a.X20*a.X33 - a.X00*a.X21*a.X33) / d
	m.X22 = (a.X01*a.X13*a.X30 - a.X03*a.X11*a.X30 + a.X03*a.X10*a.X31 - a.X00*a.X13*a.X31 - a.X01*a.X10*a.X33 + a.X00*a.X11*a.X33) / d
	m.X23 = (a.X03*a.X11*a.X20 - a.X01*a.X13*a.X20 - a.X03*a.X10*a.X21 + a.X00*a.X13*a.X21 + a.X01*a.X10*a.X23 - a.X00*a.X11*a.X23) / d
	m.X30 = (a.X12*a.X21*a.X30 - a.X11*a.X22*a.X30 - a.X12*a.X20*a.X31 + a.X10*a.X22*a.X31 + a.X11*a.X20*a.X32 - a.X10*a.X21*a.X32) / d
	m.X31 = (a.X01*a.X22*a.X30 - a.X02*a.X21*a.X30 + a.X02*a.X20*a.X31 - a.X00*a.X22*a.X31 - a.X01*a.X20*a.X32 + a.X00*a.X21*a.X32) / d
	m.X32 = (a.X02*a.X11*a.X30 - a.X01*a.X12*a.X30 - a.X02*a.X10*a.X31 + a.X00*a.X12*a.X31 + a.X01*a.X10*a.X32 - a.X00*a.X11*a.X32) / d
	m.X33 = (a.X01*a.X12*a.X20 - a.X02*a.X11*a.X20 + a.X02*a.X10*a.X21 - a.X00*a.X12*a.X21 - a.X01*a.X10*a.X22 + a.X00*a.X11*a.X22) / d
	return m
}

func Translation(v Vec3) Matrix44 {
	return Matrix44{
		1, 0, 0, v.X,
		0, 1, 0, v.Y,
		0, 0, 1, v.Z,
		0, 0, 0, 1,
	}
}

func Rotation(v Vec3, a float64) Matrix44 {
	a *= math.Pi / 180.0
	v = v.Unit()
	s := math.Sin(a)
	c := math.Cos(a)
	m := 1 - c
	return Matrix44{
		m*v.X*v.X + c, m*v.X*v.Y + v.Z*s, m*v.Z*v.X - v.Y*s, 0,
		m*v.X*v.Y - v.Z*s, m*v.Y*v.Y + c, m*v.Y*v.Z + v.X*s, 0,
		m*v.Z*v.X + v.Y*s, m*v.Y*v.Z - v.X*s, m*v.Z*v.Z + c, 0,
		0, 0, 0, 1,
	}
}

type SDFTransform struct {
	SDF3
	M Matrix44
	I Matrix44
}

func (s SDFTransform) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	m := s.I
	return func(p Vec3) float64 {
		// micro-optimization alert! ... but this is such a hot spot that manually
		// inlining MulVec3 gives a 9% speed up on my Ryzen for spt_test.go
		x := m.X00*p.X + m.X01*p.Y + m.X02*p.Z + m.X03
		y := m.X10*p.X + m.X11*p.Y + m.X12*p.Z + m.X13
		z := m.X20*p.X + m.X21*p.Y + m.X22*p.Z + m.X23
		return sdf(Vec3{x, y, z})
		//return sdf(m.MulVec3(p))
	}
}

func (s SDFTransform) Sphere() (Vec3, float64) {
	center, radius := s.SDF3.Sphere()
	return s.M.MulVec3(center), radius
}

func Translate(v Vec3, sdf SDF3) SDF3 {
	m := Translation(v)
	return SDFTransform{sdf, m, m.Inverse()}
}

func TranslateX(n float64, sdf SDF3) SDF3 {
	return Translate(V3(n, 0, 0), sdf)
}

func TranslateY(n float64, sdf SDF3) SDF3 {
	return Translate(V3(0, n, 0), sdf)
}

func TranslateZ(n float64, sdf SDF3) SDF3 {
	return Translate(V3(0, 0, n), sdf)
}

func Rotate(v Vec3, deg float64, sdf SDF3) SDF3 {
	m := Rotation(v, deg)
	return SDFTransform{sdf, m, m.Inverse()}
}

func RotateX(deg float64, sdf SDF3) SDF3 {
	return Rotate(V3(1, 0, 0), deg, sdf)
}

func RotateY(deg float64, sdf SDF3) SDF3 {
	return Rotate(V3(0, 1, 0), deg, sdf)
}

func RotateZ(deg float64, sdf SDF3) SDF3 {
	return Rotate(V3(0, 0, 1), deg, sdf)
}

type SDFScale struct {
	SDF3
	Factor float64
}

func (s SDFScale) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	return func(pos Vec3) float64 {
		return sdf(pos.Scale(1.0/s.Factor)) * s.Factor
	}
}

func (s SDFScale) Sphere() (Vec3, float64) {
	center, radius := s.SDF3.Sphere()
	return center, radius * s.Factor
}

func Scale(factor float64, sdf SDF3) SDF3 {
	return SDFScale{sdf, factor}
}

// non-uniform scaling
type SDFDistort struct {
	SDF3
	Factor Vec3
}

func (s SDFDistort) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	return func(pos Vec3) float64 {
		return sdf(pos.Div(s.Factor)) * min(s.Factor.X, min(s.Factor.Y, s.Factor.Z))
	}
}

func (s SDFDistort) Sphere() (Vec3, float64) {
	center, radius := s.SDF3.Sphere()
	return center, radius * max(max(s.Factor.X, s.Factor.Y), s.Factor.Z)
}

func Distort(factor Vec3, sdf SDF3) SDF3 {
	return SDFDistort{sdf, factor}
}

type SDFMirror struct {
	SDF3
	Mul Vec3
}

func (s SDFMirror) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	return func(pos Vec3) float64 {
		return sdf(pos.Mul(s.Mul))
	}
}

func (s SDFMirror) Sphere() (Vec3, float64) {
	center, radius := s.SDF3.Sphere()
	return center.Mul(s.Mul), radius
}

func Mirror(mul Vec3, sdf SDF3) SDF3 {
	return SDFMirror{sdf, mul}
}

func MirrorX(sdf SDF3) SDF3 {
	return Mirror(V3(-1, 1, 1), sdf)
}

func MirrorY(sdf SDF3) SDF3 {
	return Mirror(V3(1, -1, 1), sdf)
}

func MirrorZ(sdf SDF3) SDF3 {
	return Mirror(V3(1, 1, -1), sdf)
}

func itemsBoundingSphere(items []SDF3) (Vec3, float64) {
	centers := Zero3
	mradius := 0.0
	points := []Vec3{}
	for _, item := range items {
		center, radius := item.Sphere()
		mradius = max(mradius, radius)
		centers = centers.Add(center)
		points = append(points, center.Add(Vec3{radius, 0, 0}))
		points = append(points, center.Add(Vec3{-radius, 0, 0}))
		points = append(points, center.Add(Vec3{0, radius, 0}))
		points = append(points, center.Add(Vec3{0, -radius, 0}))
		points = append(points, center.Add(Vec3{0, 0, radius}))
		points = append(points, center.Add(Vec3{0, 0, -radius}))
	}

	center := centers.Scale(1.0 / float64(len(items)))
	radius := mradius

	for {
		encompass := true
		for _, p := range points {
			d := len3(sub3(p, center)) - radius
			if d > 0 {
				encompass = false
				radius += 1.0
			}
		}
		if encompass {
			break
		}
	}

	return center, radius
}

type SDFUnion struct {
	Items []SDF3
}

func (s SDFUnion) SDF() func(Vec3) float64 {
	var items []func(Vec3) float64
	var spheres []sphere
	for _, item := range s.Items {
		items = append(items, item.SDF())
		c, r := item.Sphere()
		spheres = append(spheres, sphere{c, r})
	}
	return func(pos Vec3) float64 {
		var dist float64
		for i, sdf := range items {
			if i > 0 {
				bd := spheres[i].distance(pos)
				if bd > dist {
					continue
				}
			}
			d := sdf(pos)
			if i == 0 || d < dist {
				dist = d
			}
		}
		return dist
	}
}

func (s SDFUnion) Sphere() (Vec3, float64) {
	center, radius := itemsBoundingSphere(s.Items)
	return center, radius
}

func Union(items ...SDF3) SDF3 {
	return SDFUnion{items}
}

type SDFDifference struct {
	Items []SDF3
}

func (s SDFDifference) SDF() func(Vec3) float64 {
	var items []func(Vec3) float64
	var spheres []sphere
	for _, item := range s.Items {
		items = append(items, item.SDF())
		c, r := item.Sphere()
		spheres = append(spheres, sphere{c, r})
	}
	return func(pos Vec3) float64 {
		var dist float64
		for i, sdf := range items {
			if i > 0 {
				bd := spheres[i].distance(pos)
				if -bd < dist {
					continue
				}
			}
			d := sdf(pos)
			if i == 0 {
				dist = d
			} else if -d > dist {
				dist = -d
			}
		}
		return dist
	}
}

func (s SDFDifference) Sphere() (Vec3, float64) {
	center, radius := itemsBoundingSphere(s.Items)
	return center, radius
}

func Difference(items ...SDF3) SDF3 {
	return SDFDifference{items}
}

type SDFIntersection struct {
	Items []SDF3
}

func (s SDFIntersection) SDF() func(Vec3) float64 {
	var items []func(Vec3) float64
	var spheres []sphere
	for _, item := range s.Items {
		items = append(items, item.SDF())
		c, r := item.Sphere()
		spheres = append(spheres, sphere{c, r})
	}
	return func(pos Vec3) float64 {
		var dist float64
		for i, sdf := range items {
			if i > 0 {
				bd := spheres[i].distance(pos)
				if bd > dist {
					continue
				}
			}
			d := sdf(pos)
			if i == 0 || d > dist {
				dist = d
			}
		}
		return dist
	}
}

func (s SDFIntersection) Sphere() (Vec3, float64) {
	center, radius := itemsBoundingSphere(s.Items)
	return center, radius
}

func Intersection(items ...SDF3) SDF3 {
	return SDFIntersection{items}
}

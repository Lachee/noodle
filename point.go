package noodle

import (
	"math"
	"unsafe"
)

//Point a 2 dimensional point
type Point struct {
	X int
	Y int
}

//NewPoint creates a new point with defined components
func NewPoint(x, y int) Point { return Point{X: x, Y: y} }

//NewPointd creates a new point with double precesion float (float64)
func NewPointd(x, y float64) Point { return Point{X: int(x), Y: int(y)} }

//NewPointi creates a new point
func NewPointf(x, y float32) Point { return Point{X: int(x), Y: int(y)} }

//NewPointZero creates a point with all components equaling 0
func NewPointZero() Point { return Point{X: 0, Y: 0} }

//NewPointUp creates a normalized point pointing up
func NewPointUp() Point { return Point{X: 0, Y: 1} }

//NewPointRight creates a normalized point pointing right
func NewPointRight() Point { return Point{X: 1, Y: 0} }

//NewPointOne creates a point that is completely 1
func NewPointOne() Point { return Point{X: 1, Y: 1} }

//DecomposePointer the point into a slice of floats
func (v Point) DecomposePointer() *[2]float32 { return (*[2]float32)(unsafe.Pointer(&v)) }

//Decompose the Vector into a new slice of floats.
func (v Point) Decompose() []float32 { return []float32{float32(v.X), float32(v.Y)} }

//Decomposei the point into a new slice of ints
func (v Point) Decomposei() []int { return []int{v.X, v.Y} }

//Add two points (v1 + v2)
func (v Point) Add(v2 Point) Point {
	return Point{X: v.X + v2.X, Y: v.Y + v2.Y}
}

//Subtract two points (v1 - v2)
func (v Point) Subtract(v2 Point) Point {
	return Point{X: v.X - v2.X, Y: v.Y - v2.Y}
}

//Length of the point
func (v Point) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X) + float64(v.Y*v.Y)))
}

//SqrLength is the squared length of the point
func (v Point) SqrLength() float32 {
	return float32(float64(v.X*v.X) + float64(v.Y*v.Y))
}

//Dot of the point
func (v Point) Dot(v2 Point) int {
	return v.X*v2.X + v.Y*v2.Y
}

//Perpendicular to this point
func (v Point) Perpendicular() Point {
	return Point{X: v.Y, Y: v.X}
}

//RotateByRadians rotates the point in radians. Use Deg2Rad to convert degress into radians.
func (v Point) RotateByRadians(radians float32) Point {
	sin := float32(math.Sin(float64(radians)))
	cos := float32(math.Cos(float64(radians)))
	return Point{
		X: int((cos * float32(v.X)) - (sin * float32(v.Y))),
		Y: int((sin * float32(v.X)) + (cos * float32(v.Y))),
	}
}

//Angle the point creates with another point
func (v Point) Angle(v2 Point) float32 {
	result := float32(math.Atan2(float64(v2.Y-v.Y), float64(v2.X-v.X))) * Rad2Deg
	if result < 0 {
		result += 360
	}
	return result
}

//Scale the point (v * scale)
func (v Point) Scale(scale float32) Point {
	return Point{X: int(float32(v.X) * scale), Y: int(float32(v.Y) * scale)}
}

//Multiply a point by another point
func (v Point) Multiply(v2 Point) Point {
	return Point{X: v.X * v2.X, Y: v.Y * v2.Y}
}

//Negate or Inverts a point
func (v Point) Negate() Point {
	return Point{X: -v.X, Y: -v.Y}
}

//Divide  a point by a value ( v1.x / d, v1.y / d )
func (v Point) Divide(d float32) Point {
	return Point{X: int(float32(v.X) / d), Y: int(float32(v.Y) / d)}
}

//DivideV a point by another point: ( v1.x / v2.x, v1.y / v2.y )
func (v Point) DivideV(v2 Point) Point {
	return Point{X: v.X / v2.X, Y: v.Y / v2.Y}
}

//Normalize a point
func (v Point) Normalize() Point {
	len := v.Length()
	if len == 0 {
		return v
	}
	return v.Divide(len)
}

//Lerp a point towards another point
func (v Point) Lerp(target Point, amount float32) Point {
	return Point{
		X: int(float32(v.X) + amount*float32(target.X-v.X)),
		Y: int(float32(v.Y) + amount*float32(target.Y-v.Y)),
	}
}

//Distance between two points
func (v Point) Distance(v2 Point) float32 {
	d := v2.Subtract(v)
	return d.Length()
}

//Reflect a point. The mirror normal can be invisioned as a mirror perpendicular to the surface that is hit.
func (v Point) Reflect(mirrorNormal Point) Point {
	return v.Add(v.Scale(float32(-2 * v.Dot(mirrorNormal))))
}

//Min value for each pair of components
func (v Point) Min(v2 Point) Point {
	return Point{
		X: int(math.Min(float64(v.X), float64(v2.X))),
		Y: int(math.Min(float64(v.Y), float64(v2.Y))),
	}
}

//Max value for each pair of components
func (v Point) Max(v2 Point) Point {
	return Point{
		X: int(math.Max(float64(v.X), float64(v2.X))),
		Y: int(math.Max(float64(v.Y), float64(v2.Y))),
	}
}

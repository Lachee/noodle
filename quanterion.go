package noodle

import (
	"log"
	"math"
	"unsafe"
)

/*
Copyright (c) 2020 Lachee
Copyright ©2013 The go-gl Authors. All rights reserved.
*/

//https://github.com/raysan5/raylib/blob/master/src/raymath.h
//https://github.com/go-gl/mathgl/blob/master/mgl32/quat.go

//Quaternion A represntation of rotations that does not suffer from gimbal lock
type Quaternion Vector4

//NewQuaternionIdentity creates a Quaternion Identity (a blank quaternion)
func NewQuaternionIdentity() Quaternion { return Quaternion{X: 0, Y: 0, Z: 0, W: 1} }

//NewQuaternionEuler creates a quaternion from Euler Angles
func NewQuaternionEuler(theta Vector3) Quaternion {
	var q Quaternion
	cosz2 := float32(math.Cos(0.5 * float64(theta.Z)))
	cosy2 := float32(math.Cos(0.5 * float64(theta.Y)))
	cosx2 := float32(math.Cos(0.5 * float64(theta.X)))

	sinz2 := float32(math.Sin(0.5 * float64(theta.Z)))
	siny2 := float32(math.Sin(0.5 * float64(theta.Y)))
	sinx2 := float32(math.Sin(0.5 * float64(theta.X)))

	// and now compute Quaternion
	q.W = cosz2*cosy2*cosx2 + sinz2*siny2*sinx2
	q.X = cosz2*cosy2*sinx2 - sinz2*siny2*cosx2
	q.Y = cosz2*siny2*cosx2 + sinz2*cosy2*sinx2
	q.Z = sinz2*cosy2*cosx2 - cosz2*siny2*sinx2
	return q
}

//Add two quaternions (q1 + q2)
func (q Quaternion) Add(q2 Quaternion) Quaternion {
	return Quaternion{X: q.X + q2.X, Y: q.Y + q2.Y, Z: q.Z + q2.Z, W: q.W + q2.W}
}

//Subtract two quaternions (q1 - q2)
func (q Quaternion) Subtract(q2 Quaternion) Quaternion {
	return Quaternion{X: q.X - q2.X, Y: q.Y - q2.Y, Z: q.Z - q2.Z, W: q.W - q2.W}
}

//Multiply two quaternions together
func (q Quaternion) Multiply(q2 Quaternion) Quaternion {
	return Quaternion{
		W: q.W*q2.W - (q.X*q2.X + q.Y*q2.Y + q.Z*q2.Z),
		X: q.Y*q2.Z - q.Z*q2.Y + q.W*q2.X + q.X*q2.W,
		Y: q.Z*q2.X - q.X*q2.Z + q.W*q2.Y + q.Y*q2.W,
		Z: q.X*q2.Y - q.Y*q2.X + q.W*q2.Z + q.Z*q2.W,
	}
}

//Conjugate the quaternion (invert)
func (q Quaternion) Conjugate() Quaternion {
	return Quaternion{X: -q.X, Y: -q.Y, Z: -q.Z, W: -q.W}
}

//Invert the quaternion
func (q Quaternion) Invert() Quaternion {
	sqrl := q.SqrLength()
	return Quaternion{X: -q.X / sqrl, Y: -q.Y / sqrl, Z: -q.Z / sqrl, W: -q.W / sqrl}
}

//Scale the quaternion (q * scale)
func (q Quaternion) Scale(scale float32) Quaternion {
	return Quaternion{X: q.X * scale, Y: q.Y * scale, Z: q.Z * scale, W: q.W * scale}
}

//Length of the quaternion
func (q Quaternion) Length() float32 {
	return float32(math.Sqrt(float64(q.X*q.X) + float64(q.Y*q.Y) + float64(q.Z*q.Z) + float64(q.W*q.W)))
}

//SqrLength of the quaternion
func (q Quaternion) SqrLength() float32 {
	return float32((float64(q.X*q.X) + float64(q.Y*q.Y) + float64(q.Z*q.Z) + float64(q.W*q.W)))
}

//Lerp lineraly interpolates between 2 points
func (q Quaternion) Lerp(q2 Quaternion, t float32) Quaternion {
	//TODO: Implement Lerp and Slerp
	return q
}

//Slerp spherically interpolates between 2 points
func (q Quaternion) Slerp(q2 Quaternion, t float32) Quaternion {
	//TODO: Implement Lerp and Slerp
	return q
}

//Rotate a vector by the quaternion
func (q Quaternion) Rotate(v Vector3) Vector3 {
	cong := q.Conjugate()
	r := q.Multiply(Quaternion{v.X, v.Y, v.Z, 0}).Multiply(cong)
	return Vector3{r.X, r.Y, r.Z}
}

//ToMatrix obsolete, use QuaternionRotation instead. This would just give isomorphic matrix
func (q Quaternion) ToMatrix() Matrix {
	log.Println("ToMatrix() isnt supported. Use NewMatrixRotation or NewMatrixRotationIsomorphic")
	return NewMatrixRotation(q)
}

//Decompose the Vector into a new slice of floats.
func (q Quaternion) Decompose() []float32 { return []float32{q.X, q.Y, q.Z, q.W} }

//DecomposePointer the vector into a slice of floats
func (q Quaternion) DecomposePointer() *[4]float32 { return (*[4]float32)(unsafe.Pointer(&q)) }

/*

//NewQuaternionLookRotation looks at a point
// https://answers.unity.com/questions/467614/what-is-the-source-code-of-quaternionlookrotation.html
func newQuaternionLookRotation(forward, up Vector3) Quaternion {
	var quaternion = NewQuaternionIdentity()

	v := forward.Normalize()
	v2 := up.CrossProduct(v).Normalize()
	v3 := v.CrossProduct(v2)

	var m00 = v2.X
	var m01 = v2.Y
	var m02 = v2.Z
	var m10 = v3.X
	var m11 = v3.Y
	var m12 = v3.Z
	var m20 = v.X
	var m21 = v.Y
	var m22 = v.Z

	num8 := (m00 + m11) + m22

	if num8 > 0 {
		num := float32(math.Sqrt(float64(num8) + 1.0))
		quaternion.W = num * 0.5
		num = 0.5 / num
		quaternion.X = (m12 - m21) * num
		quaternion.Y = (m20 - m02) * num
		quaternion.Z = (m01 - m10) * num
		return quaternion
	}

	if (m00 >= m11) && (m00 >= m22) {
		var num7 = float32(math.Sqrt(float64(((1.0 + m00) - m11) - m22)))
		var num4 = 0.5 / num7
		quaternion.X = 0.5 * num7
		quaternion.Y = (m01 + m10) * num4
		quaternion.Z = (m02 + m20) * num4
		quaternion.W = (m12 - m21) * num4
		return quaternion
	}

	if m11 > m22 {
		var num6 = float32(math.Sqrt(float64(((1 + m11) - m00) - m22)))
		var num3 = 0.5 / num6
		quaternion.X = (m10 + m01) * num3
		quaternion.Y = 0.5 * num6
		quaternion.Z = (m21 + m12) * num3
		quaternion.W = (m20 - m02) * num3
		return quaternion
	}

	var num5 = float32(math.Sqrt(float64(((1 + m22) - m00) - m11)))
	var num2 = 0.5 / num5
	quaternion.X = (m20 + m02) * num2
	quaternion.Y = (m21 + m12) * num2
	quaternion.Z = 0.5 * num5
	quaternion.W = (m01 - m10) * num2
	return quaternion
}
*/

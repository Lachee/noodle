package noodle

import (
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

// NewQuaternionAxis creates an angle from an axis and an angle relative to that axis in radians.
func NewQuaternionAxis(axis Vector3, angle float32) Quaternion {
	// angle = (float32(math.Pi) * angle) / 180.0

	c, s := float32(math.Cos(float64(angle/2))), float32(math.Sin(float64(angle/2)))
	return Quaternion{axis.X * s, axis.Y * s, axis.Z * s, c}
}

// NewQuaternionAngle calculates the rotation between two vectors
func NewQuaternionAngle(start, dest Vector3) Quaternion {
	const epsilon = float32(0.001)

	// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-17-quaternions/#I_need_an_equivalent_of_gluLookAt__How_do_I_orient_an_object_towards_a_point__
	// https://github.com/g-truc/glm/blob/0.9.5/glm/gtx/quaternion.inl#L225
	// https://bitbucket.org/sinbad/ogre/src/d2ef494c4a2f5d6e2f0f17d3bfb9fd936d5423bb/OgreMain/include/OgreVector3.h?at=default#cl-654

	start = start.Normalize()
	dest = dest.Normalize()

	cosTheta := start.Dot(dest)
	if cosTheta < -1.0+epsilon {
		// special case when vectors in opposite directions:
		// there is no "ideal" rotation axis
		// So guess one; any will do as long as it's perpendicular to start
		axis := Vector3{1, 0, 0}.Cross(start)
		if axis.Dot(axis) < epsilon {
			// bad luck, they were parallel, try again!
			axis = Vector3{0, 1, 0}.Cross(start)
		}

		return NewQuaternionAxis(axis.Normalize(), math.Pi)
	}

	axis := start.Cross(dest)
	s := float32(math.Sqrt(float64(1.0+cosTheta) * 2.0))
	v := axis.Scale(1.0 / s)
	return Quaternion{v.X, v.Y, v.Z, s * 0.5}
}

//NewQuaternionLookAt creates a rotation from the eye vector to the center vector.
func NewQuaternionLookAt(eye, center, up Vector3) Quaternion {
	direction := center.Subtract(eye).Normalize()
	rotDir := NewQuaternionAngle(Vector3{0, 0, 1}, direction)

	//Force the up
	right := direction.Cross(up)
	up = right.Cross(direction)

	// Recompute up so that it's perpendicular to the direction
	// You can skip that part if you really want to force up
	upCur := rotDir.Rotate(Vector3{0, 1, 0})
	rotUp := NewQuaternionAngle(upCur, up)

	rotTarget := rotUp.Multiply(rotDir)
	return rotTarget.Inverse()
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
	q1V := Vector3{q.X, q.Y, q.Z}
	q2V := Vector3{q2.X, q2.Y, q2.Z}
	v := q1V.Cross(q2V).Add(q2V.Scale(q.W)).Add(q1V.Scale(q2.W))
	return Quaternion{v.X, v.Y, v.Z, q.W*q2.W - q1V.Dot(q2V)}
}

//Conjugate the quaternion (invert)
func (q Quaternion) Conjugate() Quaternion {
	return Quaternion{X: -q.X, Y: -q.Y, Z: -q.Z, W: q.W}
}

//Inverse the quaternion
func (q Quaternion) Inverse() Quaternion {
	return q.Conjugate().Scale(1 / q.Dot(q))
}

// Rotate a vector by the rotation this quaternion represents.
// This will result in a 3D vector. Strictly speaking, this is
// equivalent to q1.v.q* where the "."" is quaternion multiplication and v is interpreted
// as a quaternion with W 0 and V v. In code:
// q1.Mul(Quat{0,v}).Mul(q1.Conjugate()), and
// then retrieving the imaginary (vector) part.
//
// In practice, we hand-compute this in the general case and simplify
// to save a few operations.
func (q Quaternion) Rotate(v Vector3) Vector3 {
	q1V := Vector3{q.X, q.Y, q.Z}
	cross := q1V.Cross(v)
	return v.Add(cross.Scale(2 * q.W)).Add(q1V.Scale(2).Cross(cross))
}

//ToMatrix is an alias of NewMatrixRotation
func (q Quaternion) ToMatrix() Matrix { return NewMatrixRotation(q) }

// Dot product between two quaternions, equivalent to if this was a Vec4.
func (q Quaternion) Dot(q2 Quaternion) float32 {
	return q.W*q2.W + q.X*q2.X + q.Y*q2.Y + q.Z*q2.Z
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

// Normalize the quaternion, returning its versor (unit quaternion).
// This is the same as normalizing it as a Vec4.
func (q Quaternion) Normalize() Quaternion {
	length := q.Length()
	if length == 0 {
		length = 1
	}

	ilength := 1 / length
	return q.Scale(ilength)
}

//Decompose the Vector into a new slice of floats.
func (q Quaternion) Decompose() []float32 { return []float32{q.X, q.Y, q.Z, q.W} }

//DecomposePointer the vector into a slice of floats
func (q Quaternion) DecomposePointer() *[4]float32 { return (*[4]float32)(unsafe.Pointer(&q)) }

// QuatSlerp is *S*pherical *L*inear Int*erp*olation, a method of interpolating
// between two quaternions. This always takes the straightest path on the sphere between
// the two quaternions, and maintains constant velocity.
//
// However, it's expensive and QuatSlerp(q1,q2) is not the same as QuatSlerp(q2,q1)
func QuaternionSlerp(q1, q2 Quaternion, amount float32) Quaternion {
	q1, q2 = q1.Normalize(), q2.Normalize()
	dot := q1.Dot(q2)

	// If the inputs are too close for comfort, linearly interpolate and normalize the result.
	if dot > 0.9995 {
		return QuaternionLerp(q1, q2, amount).Normalize()
	}

	// This is here for precision errors, I'm perfectly aware that *technically* the dot is bound [-1,1], but since Acos will freak out if it's not (even if it's just a liiiiitle bit over due to normal error) we need to clamp it
	dot = Clamp32(dot, -1, 1)

	theta := float32(math.Acos(float64(dot))) * amount
	c, s := float32(math.Cos(float64(theta))), float32(math.Sin(float64(theta)))
	rel := q2.Subtract(q1.Scale(dot)).Normalize()

	return q1.Scale(c).Add(rel.Scale(s))
}

// QuatLerp is a *L*inear Int*erp*olation between two Quaternions, cheap and simple.
//
// Not excessively useful, but uses can be found.
func QuaternionLerp(q1, q2 Quaternion, amount float32) Quaternion {
	return q1.Add(q2.Subtract(q1).Scale(amount))
}

/*

//NewQuaternionLookRotation looks at a point
// https://answers.unity.com/questions/467614/what-is-the-source-code-of-quaternionlookrotation.html
func newQuaternionLookRotation(forward, up Vector3) Quaternion {
	var quaternion = NewQuaternionIdentity()

	v := forward.Normalize()
	v2 := up.Cross(v).Normalize()
	v3 := v.Cross(v2)

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

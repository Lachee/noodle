package noodle

import "math"

/*
Copyright (c) 2020 Lachee
Copyright ©2013 The go-gl Authors. All rights reserved.
*/

//https://github.com/raysan5/raylib/blob/master/src/raymath.h
//https://github.com/go-gl/mathgl/blob/master/mgl32/quat.go

//Quaternion A represntation of rotations that does not suffer from gimbal lock
type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

//NewQuaternionIdentity creates a Quaternion Identity (a blank quaternion)
func NewQuaternionIdentity() Quaternion { return Quaternion{X: 0, Y: 0, Z: 0, W: 1} }

//newQuaternionVector3ToVector3 creates a quaternion that is the angle between 2 vectors
func newQuaternionVector3ToVector3(from, too Vector3) Quaternion {
	cos2theta := from.DotProduct(too)
	cross := from.CrossProduct(too)
	return Quaternion{X: cross.X, Y: cross.Y, Z: cross.Z, W: 1 + cos2theta}.Normalize()
}

//NewQuaternionBetweenVectors creates a rotation between the two vectors.
func NewQuaternionBetweenVectors(start, dest Vector3) Quaternion {
	//https://github.com/go-gl/mathgl/blob/master/mgl32/quat.go#L431
	start = start.Normalize()
	dest = dest.Normalize()
	epsilon := float32(0.001)
	cosTheta := start.DotProduct(dest)

	if cosTheta < -1.0+epsilon {
		// special case when vectors in opposite directions:
		// there is no "ideal" rotation axis
		// So guess one; any will do as long as it's perpendicular to start
		axis := Vector3{1, 0, 0}.CrossProduct(start)
		if axis.DotProduct(axis) < epsilon {
			// bad luck, they were parallel, try again!
			axis = Vector3{0, 1, 0}.CrossProduct(start)
		}

		return NewQuaternionAxisAngle(axis.Normalize(), PI)
	}

	axis := start.CrossProduct(dest)
	s := float32(math.Sqrt(float64(1.0+cosTheta) * 2.0))

	comps := axis.Scale(1.0 / s)
	return Quaternion{
		comps.X,
		comps.Y,
		comps.Z,
		s * 0.5,
	}
}

//NewQuaternionAxisAngle creates a quaternion from an axis and its rotation
func NewQuaternionAxisAngle(axis Vector3, angle float32) Quaternion {

	//If the axis isn't 0, we will just extream it a bit.
	//if axis.Length() != 0 {
	//	angle *= 0.5
	//}

	axis = axis.Normalize()
	cosres := float32(math.Cos(float64(angle)))
	sinres := float32(math.Sin(float64(angle)))
	return Quaternion{X: axis.X * sinres, Y: axis.Y * sinres, Z: axis.Z * sinres, W: cosres}.Normalize()
}

//NewQuaternionEuler creates a quaternion from euler angles (roll, yaw, pitch)
func NewQuaternionEuler(euler Vector3) Quaternion {
	x0 := float32(math.Cos(float64(euler.X * 0.5)))
	x1 := float32(math.Sin(float64(euler.X * 0.5)))
	y0 := float32(math.Cos(float64(euler.Y * 0.5)))
	y1 := float32(math.Sin(float64(euler.Y * 0.5)))
	z0 := float32(math.Cos(float64(euler.Z * 0.5)))
	z1 := float32(math.Sin(float64(euler.Z * 0.5)))
	return Quaternion{
		X: x1*y0*z0 - x0*y1*z1,
		Y: x0*y1*z0 + x1*y0*z1,
		Z: x0*y0*z1 - x1*y1*z0,
		W: x0*y0*z0 + x1*y1*z1,
	}
}

// NewQuaternionMatrix converts a pure rotation matrix into a quaternion
func NewQuaternionMatrix(matrix Matrix) Quaternion {
	m := matrix.DecomposePointer()

	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToQuaternion/index.htm

	if tr := m[0] + m[5] + m[10]; tr > 0 {
		s := float32(0.5 / math.Sqrt(float64(tr+1.0)))
		return Quaternion{
			(m[6] - m[9]) * s,
			(m[8] - m[2]) * s,
			(m[1] - m[4]) * s,
			0.25 / s,
		}
	}

	if (m[0] > m[5]) && (m[0] > m[10]) {
		s := float32(2.0 * math.Sqrt(float64(1.0+m[0]-m[5]-m[10])))
		return Quaternion{
			0.25 * s,
			(m[4] + m[1]) / s,
			(m[8] + m[2]) / s,
			(m[6] - m[9]) / s,
		}
	}

	if m[5] > m[10] {
		s := float32(2.0 * math.Sqrt(float64(1.0+m[5]-m[0]-m[10])))
		return Quaternion{
			(m[4] + m[1]) / s,
			0.25 * s,
			(m[9] + m[6]) / s,
			(m[8] - m[2]) / s,
		}

	}

	s := float32(2.0 * math.Sqrt(float64(1.0+m[10]-m[0]-m[5])))
	return Quaternion{
		(m[8] + m[2]) / s,
		(m[9] + m[6]) / s,
		0.25 * s,
		(m[1] - m[4]) / s,
	}
}

//NewQuaternionLookAt creates a rotation from the eye to the center, with the give up.
func NewQuaternionLookAt(eye, center, up Vector3) Quaternion {
	// https://github.com/go-gl/mathgl/blob/master/mgl32/quat.go#L406
	// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-17-quaternions/#I_need_an_equivalent_of_gluLookAt__How_do_I_orient_an_object_towards_a_point__
	// https://bitbucket.org/sinbad/ogre/src/d2ef494c4a2f5d6e2f0f17d3bfb9fd936d5423bb/OgreMain/src/OgreCamera.cpp?at=default#cl-161

	direction := center.Subtract(eye).Normalize()

	// Find the rotation between the front of the object (that we assume towards Z-,
	// but this depends on your model) and the desired direction
	rotDir := NewQuaternionBetweenVectors(Vector3{0, 0, -1}, direction)

	// Recompute up so that it's perpendicular to the direction
	// You can skip that part if you really want to force up
	//right := direction.Cross(up)
	//up = right.Cross(direction)

	// Because of the 1rst rotation, the up is probably completely screwed up.
	// Find the rotation between the "up" of the rotated object, and the desired up
	upCur := rotDir.Rotate(Vector3{0, 1, 0})
	rotUp := NewQuaternionBetweenVectors(upCur, up)

	rotTarget := rotUp.Multiply(rotDir) // remember, in reverse order.
	return rotTarget.Invert()           // camera rotation should be inversed!
}

//Invert a quaternions components
func (q Quaternion) Invert() Quaternion {
	length := q.SqrLength()
	if length != 0 {
		i := 1 / length
		return Quaternion{
			X: q.X * -i,
			Y: q.Y * -i,
			Z: q.Z * -i,
			W: q.W * i,
		}
	}
	return q
}

//Decompose the quaternion into a slice of floats
func (q Quaternion) Decompose() []float32 { return []float32{q.X, q.Y, q.Z, q.W} }

//Length of the quaternion
func (q Quaternion) Length() float32 {
	return float32(math.Sqrt(float64(q.X*q.X) + float64(q.Y*q.Y) + float64(q.Z*q.Z) + float64(q.W*q.W)))
}

//SqrLength is the squared length of the quaternion
func (q Quaternion) SqrLength() float32 {
	return float32(float64(q.X*q.X) + float64(q.Y*q.Y) + float64(q.Z*q.Z) + float64(q.W*q.W))
}

//Scale the quaternion (v * scale)
func (q Quaternion) Scale(scale float32) Quaternion {
	return Quaternion{X: q.X * scale, Y: q.Y * scale, Z: q.Z * scale, W: q.W * scale}
}

//Normalize a quaternion
func (q Quaternion) Normalize() Quaternion {
	length := q.Length()
	if length == 0 {
		length = 1
	}

	ilength := 1 / length
	return q.Scale(ilength)
}

//Multiply two Quaternion together, doing queraternion mathmatics
func (q Quaternion) Multiply(q2 Quaternion) Quaternion {
	return Quaternion{
		X: q.X*q2.W + q.W*q2.X + q.Y*q2.Z - q.Z*q2.Y,
		Y: q.Y*q2.W + q.W*q2.Y + q.Z*q2.X - q.X*q2.Z,
		Z: q.Z*q2.W + q.W*q2.Z + q.X*q2.Y - q.Y*q2.X,
		W: q.W*q2.W - q.X*q2.X - q.Y*q2.Y - q.Z*q2.Z,
	}
}

//Lerp a vector towards another vector
func (q Quaternion) Lerp(target Quaternion, amount float32) Quaternion {
	return Quaternion{
		X: q.X + amount*(target.X-q.X),
		Y: q.Y + amount*(target.Y-q.Y),
		Z: q.Z + amount*(target.Z-q.Z),
		W: q.W + amount*(target.W-q.W),
	}
}

//Nlerp slerp-optimized interpolation between two quaternions
func (q Quaternion) Nlerp(target Quaternion, amount float32) Quaternion {
	return q.Lerp(target, amount).Normalize()
}

//Slerp Spherically Lerped
func (q Quaternion) Slerp(q2 Quaternion, amount float32) Quaternion {
	cosHalfTheta := q.X*q2.X + q.Y*q2.Y + q.Z*q2.Z + q.W*q2.W
	if math.Abs((float64(cosHalfTheta))) >= 1 {
		return q
	}

	if cosHalfTheta > 0.95 {
		return q.Nlerp(q2, amount)
	}

	halfTheta := float32(math.Acos(float64(cosHalfTheta)))
	sinHalfTheta := float32(math.Sqrt(float64(1 - cosHalfTheta*cosHalfTheta)))

	if math.Abs(float64(sinHalfTheta)) < 0.001 {
		return Quaternion{
			X: q.X*0.5 + q.X*0.5,
			Y: q.Y*0.5 + q.Y*0.5,
			Z: q.Z*0.5 + q.Z*0.5,
			W: q.W*0.5 + q.W*0.5,
		}
	}

	ratioA := float32(math.Sin(float64((1-amount)*halfTheta)) / float64(sinHalfTheta))
	ratioB := float32(math.Sin(float64(amount*halfTheta)) / float64(sinHalfTheta))

	return Quaternion{
		X: q.X*ratioA + q.X*ratioB,
		Y: q.Y*ratioA + q.Y*ratioB,
		Z: q.Z*ratioA + q.Z*ratioB,
		W: q.W*ratioA + q.W*ratioB,
	}
}

// Rotate a vector by the rotation this quaternion represents.
func (q Quaternion) Rotate(v Vector3) Vector3 {
	qv := Vector3{q.X, q.Y, q.Z}
	cross := qv.CrossProduct(v)
	return v.Add(cross.Scale(2 * q.W)).Add(qv.Scale(2).CrossProduct(cross))
}

//ToAxisAngle returns the rotation angle and axis for a given quaternion
func (q Quaternion) ToAxisAngle() (Vector3, float32) {

	var den float32
	var resAngle float32
	var resAxis Vector3

	if math.Abs(float64(q.W)) > 1 {
		q = q.Normalize()
	}

	resAxis = Vector3{0, 0, 0}
	resAngle = 2 * float32(math.Atan(float64(q.W)))
	den = float32(math.Sqrt(float64(1 - q.W*q.W)))
	if den > 0.0001 {
		resAxis.X = q.X / den
		resAxis.Y = q.Y / den
		resAxis.Z = q.Z / den
	} else {
		resAxis.X = 1
	}

	return resAxis, resAngle
}

//ToEuler turns the quaternion into equivalent euler angles (roll, putch, yaw). Values are returned in Degrees
func (q Quaternion) ToEuler() Vector3 {
	x0 := 2 * (q.W*q.X + q.Y*q.Z)
	x1 := 1 - 2*(q.X*q.X+q.Y*q.Y)
	y0 := Clamp(float64(2*(q.W*q.Y-q.Z*q.X)), -1, 1)
	z0 := 2 * (q.W*q.Z + q.X*q.Y)
	z1 := 1 - 2*(q.Y*q.Y+q.Z*q.Z)

	return Vector3{
		X: float32(math.Atan2(float64(x0), float64(x1))) * Rad2Deg,
		Y: float32(math.Asin(y0)) * Rad2Deg,
		Z: float32(math.Atan2(float64(z0), float64(z1))) * Rad2Deg,
	}
}

//ToMatrix turns the quaternion into a matrix representation
func (q Quaternion) ToMatrix() Matrix {
	w, x, y, z := q.W, q.X, q.Y, q.Z
	return Matrix{
		1 - 2*y*y - 2*z*z, 2*x*y + 2*w*z, 2*x*z - 2*w*y, 0,
		2*x*y - 2*w*z, 1 - 2*x*x - 2*z*z, 2*y*z + 2*w*x, 0,
		2*x*z + 2*w*y, 2*y*z - 2*w*x, 1 - 2*x*x - 2*y*y, 0,
		0, 0, 0, 1,
	}
}

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

package noodle

import (
	"log"
	"math"
)

//Sources:
//https://github.com/mrdoob/three.js/blob/dev/src/math/Quaternion.js
//https://github.com/raysan5/raylib/blob/master/src/raymath.h
//https://github.com/go-gl/mathgl/blob/master/mgl32/quat.go

//EulerOrder the ordering of the euler angles
type EulerOrder uint8

const (
	//OrderXYZ of the quaternion
	OrderXYZ EulerOrder = iota
	//OrderYXZ of the quaternion
	OrderYXZ
	//OrderZXY of the quaternion
	OrderZXY
	//OrderZYX of the quaternion
	OrderZYX
	//OrderYZX of the quaternion
	OrderYZX
	//OrderXZY of the quaternion
	OrderXZY
)

//QuaternionOrdering defines the Euler Order that will be used when reading in and out Euler Angles. Modifying this will change the behaviour of many functions.
var QuaternionOrdering = OrderXYZ

//Quaternion A represntation of rotations that does not suffer from gimbal lock
//Based of THREE.JS implementation
//https://github.com/mrdoob/three.js/blob/dev/src/math/Quaternion.js
type Quaternion Vector4

//NewQuaternionIdentity creates a Quaternion Identity (a blank quaternion)
func NewQuaternionIdentity() Quaternion { return Quaternion{X: 0, Y: 0, Z: 0, W: 1} }

//NewQuaternionEulerAngle creates a new quaternion from the euler and the current QuaternionOrdering
func NewQuaternionEulerAngle(euler Vector3) Quaternion {
	x, y, z := euler.X, euler.Y, euler.Z

	c1 := float32(math.Cos(float64(x / 2)))
	c2 := float32(math.Cos(float64(y / 2)))
	c3 := float32(math.Cos(float64(z / 2)))

	s1 := float32(math.Sin(float64(x / 2)))
	s2 := float32(math.Sin(float64(y / 2)))
	s3 := float32(math.Sin(float64(z / 2)))

	result := Quaternion{0, 0, 0, 1}
	switch QuaternionOrdering {
	case OrderXYZ:
		result.X = s1*c2*c3 + c1*s2*s3
		result.Y = c1*s2*c3 - s1*c2*s3
		result.Z = c1*c2*s3 + s1*s2*c3
		result.W = c1*c2*c3 - s1*s2*s3
		break
	case OrderYXZ:
		result.X = s1*c2*c3 + c1*s2*s3
		result.Y = c1*s2*c3 - s1*c2*s3
		result.Z = c1*c2*s3 - s1*s2*c3
		result.W = c1*c2*c3 + s1*s2*s3
		break
	case OrderZXY:
		result.X = s1*c2*c3 - c1*s2*s3
		result.Y = c1*s2*c3 + s1*c2*s3
		result.Z = c1*c2*s3 + s1*s2*c3
		result.W = c1*c2*c3 - s1*s2*s3
		break
	case OrderZYX:
		result.X = s1*c2*c3 - c1*s2*s3
		result.Y = c1*s2*c3 + s1*c2*s3
		result.Z = c1*c2*s3 - s1*s2*c3
		result.W = c1*c2*c3 + s1*s2*s3
		break
	case OrderYZX:
		result.X = s1*c2*c3 + c1*s2*s3
		result.Y = c1*s2*c3 + s1*c2*s3
		result.Z = c1*c2*s3 - s1*s2*c3
		result.W = c1*c2*c3 - s1*s2*s3
		break
	case OrderXZY:
		result.X = s1*c2*c3 - c1*s2*s3
		result.Y = c1*s2*c3 - s1*c2*s3
		result.Z = c1*c2*s3 + s1*s2*c3
		result.W = c1*c2*c3 + s1*s2*s3
		break
	default:
		log.Fatalln("invalid quaternion ordering supplied.")
		break
	}

	return result
}

//NewQuaternionAxisAngle creates a new quaternion that is the axis rotated by the angle (in radians)
func NewQuaternionAxisAngle(axis Vector3, angle float32) Quaternion {
	s := float32(math.Sin(float64(angle / 2)))
	return Quaternion{
		X: axis.X * s,
		Y: axis.Y * s,
		Z: axis.Z * s,
		W: float32(math.Cos(float64(angle))),
	}
}

//NewQuaternionRotationMatrix creates a quaternion from the rotation component of a given matrix
func NewQuaternionRotationMatrix(matrix Matrix) Quaternion {
	te := matrix.DecomposePointer()
	m11 := te[0]
	m21 := te[1]
	m31 := te[2]
	m12 := te[4]
	m22 := te[5]
	m32 := te[6]
	m13 := te[8]
	m23 := te[9]
	m33 := te[10]
	trace := m11 + m22 + m33

	result := Quaternion{0, 0, 0, 1}
	if trace > 0 {

		s := 0.5 / float32(math.Sqrt(float64(trace+1.0)))

		result.W = 0.25 / s
		result.X = (m32 - m23) * s
		result.Y = (m13 - m31) * s
		result.Z = (m21 - m12) * s

	} else if m11 > m22 && m11 > m33 {

		s := 2.0 * float32(math.Sqrt(float64(1.0+m11-m22-m33)))

		result.W = (m32 - m23) / s
		result.X = 0.25 * s
		result.Y = (m12 + m21) / s
		result.Z = (m13 + m31) / s

	} else if m22 > m33 {

		s := 2.0 * float32(math.Sqrt(float64(1.0+m22-m11-m33)))

		result.W = (m13 - m31) / s
		result.X = (m12 + m21) / s
		result.Y = 0.25 * s
		result.Z = (m23 + m32) / s

	} else {

		s := 2.0 * float32(math.Sqrt(float64(1.0+m33-m11-m22)))

		result.W = (m21 - m12) / s
		result.X = (m13 + m31) / s
		result.Y = (m23 + m32) / s
		result.Z = 0.25 * s
	}

	return result
}

//NewQuaternionUnitVectors creates a quaternion that is the angle between the two unit vectors.
func NewQuaternionUnitVectors(from, to Vector3) Quaternion {
	const EPS = 0.000001

	r := from.Dot(to) + 1
	result := Quaternion{0, 0, 0, 1}
	if r < EPS {
		r = 0
		if math.Abs(float64(from.X)) > math.Abs(float64(to.Z)) {

			result.X = -from.Y
			result.Y = from.X
			result.Z = 0
			result.W = r

		} else {

			result.X = 0
			result.Y = -from.Z
			result.Z = from.Y
			result.W = r

		}

	} else {
		result.X = from.Y*to.Z - from.Z*to.Y
		result.Y = from.Z*to.X - from.X*to.Z
		result.Z = from.X*to.Y - from.Y*to.X
		result.W = r

	}
	log.Fatalln("not implemented")
	return result
	//TODO: Implement Normalise
	//return result.Normalize()
}

//- angleTo
//- roatateTwoards

//Inverse the quaternion
func (q Quaternion) Inverse() Quaternion {
	//TODO: Implement this. It assumes its a unit length, but i dont trust it.
	return q.Conjugate()
}

//Conjugate the quaternion
func (q Quaternion) Conjugate() Quaternion {
	return Quaternion{-q.X, -q.Y, -q.Z, q.W}
}

//Length of the Quaternion
func (q Quaternion) Length() float32 {
	return float32(math.Sqrt(float64(q.X*q.X) + float64(q.Y*q.Y) + float64(q.Z*q.Z) + float64(q.W*q.W)))
}

//SqrLength is the squared length of the Quaternion
func (q Quaternion) SqrLength() float32 {
	return float32(float64(q.X*q.X) + float64(q.Y*q.Y) + float64(q.Z*q.Z) + float64(q.W*q.W))
}

//Dot of the Quaternion
func (q Quaternion) Dot(q2 Quaternion) float32 {
	return q.X*q2.X + q.Y*q2.Y + q.Z*q2.Z + q.W*q2.W
}

//Normalize a Quaternion
func (q Quaternion) Normalize() Quaternion {
	l := q.Length()
	if l == 0 {
		return Quaternion{0, 0, 0, 1}
	}

	l = 1 / l
	return Quaternion{
		q.X * l,
		q.Y * l,
		q.Z * l,
		q.W * l,
	}
}

//Multiply two quaternions toegether
func (q Quaternion) Multiply(q2 Quaternion) Quaternion {
	qax, qay, qaz, qaw := q.X, q.Y, q.Z, q.W
	qbx, qby, qbz, qbw := q2.X, q2.Y, q2.Z, q2.W

	result := Quaternion{0, 0, 0, 1}
	result.X = qax*qbw + qaw*qbx + qay*qbz - qaz*qby
	result.Y = qay*qbw + qaw*qby + qaz*qbx - qax*qbz
	result.Z = qaz*qbw + qaw*qbz + qax*qby - qay*qbx
	result.W = qaw*qbw - qax*qbx - qay*qby - qaz*qbz
	return result
}

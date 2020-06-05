package noodle

import "math"

//https://github.com/raysan5/raylib/blob/master/src/raymath.h

//Quaternion A represntation of rotations that does not suffer from gimbal lock
type Quaternion struct {
	X float64
	Y float64
	Z float64
	W float64
}

//NewQuaternionIdentity creates a Quaternion Identity (a blank quaternion)
func NewQuaternionIdentity() Quaternion { return Quaternion{X: 0, Y: 0, Z: 0, W: 1} }

//NewQuaternionVector3ToVector3 creates a quaternion that is the angle between 2 vectors
func NewQuaternionVector3ToVector3(from, too Vector3) Quaternion {
	cos2theta := from.DotProduct(too)
	cross := from.CrossProduct(too)
	return Quaternion{X: cross.X, Y: cross.Y, Z: cross.Z, W: 1 + cos2theta}.Normalize()
}

//NewQuaternionFromAxisAngle creates a quaternion from an axis and its rotation
func NewQuaternionFromAxisAngle(axis Vector3, angle float64) Quaternion {

	if axis.Length() != 0 {
		angle *= 0.5
	}

	axis = axis.Normalize()

	sinres := math.Sin(angle)
	cosres := math.Cos(angle)

	return Quaternion{X: axis.X * sinres, Y: axis.Y * sinres, Z: axis.Z * sinres, W: cosres}.Normalize()
}

//NewQuaternionFromMatrix creates a Quaternion from a rotation matrix
func NewQuaternionFromMatrix(mat Matrix) Quaternion {
	var s float64
	var invS float64
	trace := mat.Trace()

	if trace > 0 {
		s = math.Sqrt(float64(trace+1)) * 2
		invS = 1 / s
		return Quaternion{
			X: float64(mat.M6-mat.M9) * invS,
			Y: float64(mat.M8-mat.M2) * invS,
			Z: float64(mat.M1-mat.M4) * invS,
			W: s * 0.25,
		}
	}

	m00 := mat.M0
	m11 := mat.M5
	m22 := mat.M10

	if m00 > m11 && m00 > m22 {
		s = math.Sqrt(float64(1+m00-m11-m22)) * 2
		invS = 1 / s
		return Quaternion{
			X: s * 0.25,
			Y: float64(mat.M4-mat.M1) * invS,
			Z: float64(mat.M8-mat.M2) * invS,
			W: float64(mat.M6-mat.M9) * invS,
		}
	} else if m11 > m22 {
		s = math.Sqrt(float64(1+m11-m00-m22)) * 2
		invS = 1 / s
		return Quaternion{
			X: float64(mat.M4-mat.M1) * invS,
			Y: s * 0.25,
			Z: float64(mat.M9-mat.M6) * invS,
			W: float64(mat.M8-mat.M2) * invS,
		}
	}

	s = math.Sqrt(float64(1+m22-m00-m11)) * 2
	invS = 1 / s
	return Quaternion{
		X: float64(mat.M8-mat.M2) * invS,
		Y: float64(mat.M9-mat.M6) * invS,
		Z: s * 0.25,
		W: float64(mat.M1-mat.M4) * invS,
	}
}

//NewQuaternionFromEuler creates a quaternion from euler angles (roll, yaw, pitch)
func NewQuaternionFromEuler(euler Vector3) Quaternion {
	x0 := math.Cos(euler.X * 0.5)
	x1 := math.Sin(euler.X * 0.5)
	y0 := math.Cos(euler.Y * 0.5)
	y1 := math.Sin(euler.Y * 0.5)
	z0 := math.Cos(euler.Z * 0.5)
	z1 := math.Sin(euler.Z * 0.5)
	return Quaternion{
		X: x1*y0*z0 - x0*y1*z1,
		Y: x0*y1*z0 + x1*y0*z1,
		Z: x0*y0*z1 - x1*y1*z0,
		W: x0*y0*z0 + x1*y1*z1,
	}
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
func (q Quaternion) Decompose() []float64 { return []float64{q.X, q.Y, q.Z, q.W} }

//Length of the quaternion
func (q Quaternion) Length() float64 {
	return math.Sqrt((q.X * q.X) + (q.Y * q.Y) + (q.Z * q.Z) + (q.W * q.W))
}

//SqrLength is the squared length of the quaternion
func (q Quaternion) SqrLength() float64 {
	return (q.X * q.X) + (q.Y * q.Y) + (q.Z * q.Z) + (q.W * q.W)
}

//Scale the quaternion (v * scale)
func (q Quaternion) Scale(scale float64) Quaternion {
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
func (q Quaternion) Lerp(target Quaternion, amount float64) Quaternion {
	return Quaternion{
		X: q.X + amount*(target.X-q.X),
		Y: q.Y + amount*(target.Y-q.Y),
		Z: q.Z + amount*(target.Z-q.Z),
		W: q.W + amount*(target.W-q.W),
	}
}

//Nlerp slerp-optimized interpolation between two quaternions
func (q Quaternion) Nlerp(target Quaternion, amount float64) Quaternion {
	return q.Lerp(target, amount).Normalize()
}

//Slerp Spherically Lerped
func (q Quaternion) Slerp(q2 Quaternion, amount float64) Quaternion {
	cosHalfTheta := q.X*q2.X + q.Y*q2.Y + q.Z*q2.Z + q.W*q2.W
	if math.Abs(cosHalfTheta) >= 1 {
		return q
	}

	if cosHalfTheta > 0.95 {
		return q.Nlerp(q2, amount)
	}

	halfTheta := math.Acos(cosHalfTheta)
	sinHalfTheta := math.Sqrt(1 - cosHalfTheta*cosHalfTheta)

	if math.Abs(sinHalfTheta) < 0.001 {
		return Quaternion{
			X: q.X*0.5 + q.X*0.5,
			Y: q.Y*0.5 + q.Y*0.5,
			Z: q.Z*0.5 + q.Z*0.5,
			W: q.W*0.5 + q.W*0.5,
		}
	}

	ratioA := math.Sin(((1 - amount) * halfTheta)) / (sinHalfTheta)
	ratioB := math.Sin((amount * halfTheta)) / (sinHalfTheta)

	return Quaternion{
		X: q.X*ratioA + q.X*ratioB,
		Y: q.Y*ratioA + q.Y*ratioB,
		Z: q.Z*ratioA + q.Z*ratioB,
		W: q.W*ratioA + q.W*ratioB,
	}
}

//ToMatrix converts the quaternion into a rotation matrix
func (q Quaternion) ToMatrix() Matrix {
	return NewMatrixFromQuaternion(q)
}

//ToAxisAngle returns the rotation angle and axis for a given quaternion
func (q Quaternion) ToAxisAngle() (Vector3, float64) {

	var den float64
	var resAngle float64
	var resAxis Vector3

	if math.Abs(q.W) > 1 {
		q = q.Normalize()
	}

	resAxis = Vector3{0, 0, 0}
	resAngle = 2 * math.Atan(q.W)
	den = math.Sqrt(1 - q.W*q.W)
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
		X: math.Atan2(x0, x1) * Rad2Deg,
		Y: math.Asin(y0) * Rad2Deg,
		Z: math.Atan2(z0, z1) * Rad2Deg,
	}
}

//Transform a quaternion, given a transformation matrix
func (q Quaternion) Transform(mat Matrix) Quaternion {
	X := float32(q.X)
	Y := float32(q.Y)
	Z := float32(q.Z)
	W := float32(q.W)

	return Quaternion{
		X: float64(mat.M0*X + mat.M4*Y + mat.M8*Z + mat.M12*W),
		Y: float64(mat.M1*X + mat.M5*Y + mat.M9*Z + mat.M13*W),
		Z: float64(mat.M2*X + mat.M6*Y + mat.M10*Z + mat.M14*W),
		W: float64(mat.M3*X + mat.M7*Y + mat.M11*Z + mat.M15*W),
	}
}

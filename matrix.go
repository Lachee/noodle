package noodle

import (
	"math"
	"unsafe"
)

//Matrix A representation of a 4 x 4 matrix
type Matrix struct {
	M0, M1, M2, M3     float32
	M4, M5, M6, M7     float32
	M8, M9, M10, M11   float32
	M12, M13, M14, M15 float32
}

//NewMatrix creates a identity
func NewMatrix() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

//NewMatrixTranslate creates a new translate matrix
func NewMatrixTranslate(v Vector3) Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v.X, v.Y, v.Z, 1}
}

//NewMatrixScale creates a new scale matrix
func NewMatrixScale(s Vector3) Matrix {
	return Matrix{
		s.X, 0, 0, 0,
		0, s.Y, 0, 0,
		0, 0, s.Z, 0,
		0, 0, 0, 1,
	}
}

// NewMatrixRotation returns the homogeneous 3D rotation matrix corresponding to the quaternion.
func NewMatrixRotation(q1 Quaternion) Matrix {
	w, x, y, z := q1.W, q1.X, q1.Y, q1.Z
	return Matrix{
		1 - 2*y*y - 2*z*z, 2*x*y + 2*w*z, 2*x*z - 2*w*y, 0,
		2*x*y - 2*w*z, 1 - 2*x*x - 2*z*z, 2*y*z + 2*w*x, 0,
		2*x*z + 2*w*y, 2*y*z - 2*w*x, 1 - 2*x*x - 2*y*y, 0,
		0, 0, 0, 1,
	}
}

//NewMatrixRotationX creates a matrix which rotates around the X axis
func NewMatrixRotationX(radians float32) Matrix {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))
	return Matrix{
		1, 0, 0, 0,
		0, c, s, 0,
		0, -s, c, 0,
		0, 0, 0, 1,
	}
}

//NewMatrixRotationY creates a matrix which rotates around the Y axis
func NewMatrixRotationY(radians float32) Matrix {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))
	return Matrix{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

//NewMatrixRotationZ creates a matrix which rotates around the Z axis
func NewMatrixRotationZ(radians float32) Matrix {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))
	return Matrix{
		c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

//NewMatrixPerspective creates a perspective projection matrix. FOVY is in degrees
func NewMatrixPerspective(fovy, aspect, near, far float32) Matrix {
	f := float32(math.Tan((math.Pi * 0.5) - (0.5 * float64(fovy*Deg2Rad))))
	rangeInv := 1.0 / (near - far)
	return Matrix{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (near + far) * rangeInv, -1,
		0, 0, near * far * rangeInv * 2, 0,
	}
}

//NewMatrixOrtho creates a orthographic projection
func NewMatrixOrtho(left, right, bottom, top, near, far float32) Matrix {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)
	return Matrix{float32(2. / rml), 0, 0, 0, 0, float32(2. / tmb), 0, 0, 0, 0, float32(-2. / fmn), 0, float32(-(right + left) / rml), float32(-(top + bottom) / tmb), float32(-(far + near) / fmn), 1}
}

// NewMatrixFrustum generates a Frustum Matrix.
func NewMatrixFrustum(left, right, bottom, top, near, far float32) Matrix {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)
	A, B, C, D := (right+left)/rml, (top+bottom)/tmb, -(far+near)/fmn, -(2*far*near)/fmn

	return Matrix{float32((2. * near) / rml), 0, 0, 0, 0, float32((2. * near) / tmb), 0, 0, float32(A), float32(B), float32(C), -1, 0, 0, float32(D), 0}
}

//NewMatrixLookAt creates a matrix to look at a target
func NewMatrixLookAt(eye, target, up Vector3) Matrix {
	var zAxis = eye.Subtract(target).Normalize() //Local Z Direction
	var xAxis = up.Cross(zAxis).Normalize()      //Local X Direction
	var yAxis = zAxis.Cross(xAxis).Normalize()   //Local Y Direction
	return Matrix{
		xAxis.X, xAxis.Y, xAxis.Z, 0,
		yAxis.X, yAxis.Y, yAxis.Z, 0,
		zAxis.X, zAxis.Y, zAxis.Z, 0,
		eye.X, eye.Y, eye.Z, 1, //Translate it
	}
}

//Trace of the matrix (sum of values along diagonal)
func (m Matrix) Trace() float32 {
	return m.M0 + m.M5 + m.M10 + m.M15
}

//Add two matrices
func (m Matrix) Add(right Matrix) Matrix {
	m1 := m.DecomposePointer()
	m2 := right.DecomposePointer()
	return Matrix{m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2], m1[3] + m2[3], m1[4] + m2[4], m1[5] + m2[5], m1[6] + m2[6], m1[7] + m2[7], m1[8] + m2[8], m1[9] + m2[9], m1[10] + m2[10], m1[11] + m2[11], m1[12] + m2[12], m1[13] + m2[13], m1[14] + m2[14], m1[15] + m2[15]}
}

//Subtract two matrices
func (m Matrix) Subtract(right Matrix) Matrix {
	m1 := m.DecomposePointer()
	m2 := right.DecomposePointer()
	return Matrix{m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2], m1[3] - m2[3], m1[4] - m2[4], m1[5] - m2[5], m1[6] - m2[6], m1[7] - m2[7], m1[8] - m2[8], m1[9] - m2[9], m1[10] - m2[10], m1[11] - m2[11], m1[12] - m2[12], m1[13] - m2[13], m1[14] - m2[14], m1[15] - m2[15]}
}

//Multiplyf scales all components of the matrix by c
func (m Matrix) Multiplyf(c float32) Matrix {
	m1 := m.DecomposePointer()
	return Matrix{m1[0] * c, m1[1] * c, m1[2] * c, m1[3] * c, m1[4] * c, m1[5] * c, m1[6] * c, m1[7] * c, m1[8] * c, m1[9] * c, m1[10] * c, m1[11] * c, m1[12] * c, m1[13] * c, m1[14] * c, m1[15] * c}
}

//Multiply 2 matrixs together
func (m Matrix) Multiply(right Matrix) Matrix {
	result := Matrix{}
	a := m.DecomposePointer()
	b := right.DecomposePointer()
	dst := result.DecomposePointer()

	var b00 = b[0*4+0]
	var b01 = b[0*4+1]
	var b02 = b[0*4+2]
	var b03 = b[0*4+3]
	var b10 = b[1*4+0]
	var b11 = b[1*4+1]
	var b12 = b[1*4+2]
	var b13 = b[1*4+3]
	var b20 = b[2*4+0]
	var b21 = b[2*4+1]
	var b22 = b[2*4+2]
	var b23 = b[2*4+3]
	var b30 = b[3*4+0]
	var b31 = b[3*4+1]
	var b32 = b[3*4+2]
	var b33 = b[3*4+3]
	var a00 = a[0*4+0]
	var a01 = a[0*4+1]
	var a02 = a[0*4+2]
	var a03 = a[0*4+3]
	var a10 = a[1*4+0]
	var a11 = a[1*4+1]
	var a12 = a[1*4+2]
	var a13 = a[1*4+3]
	var a20 = a[2*4+0]
	var a21 = a[2*4+1]
	var a22 = a[2*4+2]
	var a23 = a[2*4+3]
	var a30 = a[3*4+0]
	var a31 = a[3*4+1]
	var a32 = a[3*4+2]
	var a33 = a[3*4+3]
	dst[0] = b00*a00 + b01*a10 + b02*a20 + b03*a30
	dst[1] = b00*a01 + b01*a11 + b02*a21 + b03*a31
	dst[2] = b00*a02 + b01*a12 + b02*a22 + b03*a32
	dst[3] = b00*a03 + b01*a13 + b02*a23 + b03*a33
	dst[4] = b10*a00 + b11*a10 + b12*a20 + b13*a30
	dst[5] = b10*a01 + b11*a11 + b12*a21 + b13*a31
	dst[6] = b10*a02 + b11*a12 + b12*a22 + b13*a32
	dst[7] = b10*a03 + b11*a13 + b12*a23 + b13*a33
	dst[8] = b20*a00 + b21*a10 + b22*a20 + b23*a30
	dst[9] = b20*a01 + b21*a11 + b22*a21 + b23*a31
	dst[10] = b20*a02 + b21*a12 + b22*a22 + b23*a32
	dst[11] = b20*a03 + b21*a13 + b22*a23 + b23*a33
	dst[12] = b30*a00 + b31*a10 + b32*a20 + b33*a30
	dst[13] = b30*a01 + b31*a11 + b32*a21 + b33*a31
	dst[14] = b30*a02 + b31*a12 + b32*a22 + b33*a32
	dst[15] = b30*a03 + b31*a13 + b32*a23 + b33*a33
	return result
}

//MultiplyVector4 multiplies a vector with the matrix (m * v)
func (m Matrix) MultiplyVector4(v Vector4) Vector4 {
	m1 := m.DecomposePointer()
	m2 := v.DecomposePointer()
	return Vector4{
		m1[0]*m2[0] + m1[4]*m2[1] + m1[8]*m2[2] + m1[12]*m2[3],
		m1[1]*m2[0] + m1[5]*m2[1] + m1[9]*m2[2] + m1[13]*m2[3],
		m1[2]*m2[0] + m1[6]*m2[1] + m1[10]*m2[2] + m1[14]*m2[3],
		m1[3]*m2[0] + m1[7]*m2[1] + m1[11]*m2[2] + m1[15]*m2[3],
	}
}

//Translate the given matrix by the given vector. m * v
func (m Matrix) Translate(v Vector3) Matrix {
	return m.Multiply(NewMatrixTranslate(v))
}

//RotateX rotates the matrix
func (m Matrix) RotateX(radians float32) Matrix {
	return m.Multiply(NewMatrixRotationX(radians))
}

//RotateY rotates the matrix
func (m Matrix) RotateY(radians float32) Matrix {
	return m.Multiply(NewMatrixRotationY(radians))
}

//RotateZ rotates the matrix
func (m Matrix) RotateZ(radians float32) Matrix {
	return m.Multiply(NewMatrixRotationZ(radians))
}

//Scale the given matrix
func (m Matrix) Scale(s Vector3) Matrix {
	return m.Multiply(NewMatrixScale(s))
}

//Inverse the matrix
func (m Matrix) Inverse() Matrix {
	result := Matrix{}
	dst := result.DecomposePointer()

	mm := m.DecomposePointer()
	var m00 = mm[0*4+0]
	var m01 = mm[0*4+1]
	var m02 = mm[0*4+2]
	var m03 = mm[0*4+3]
	var m10 = mm[1*4+0]
	var m11 = mm[1*4+1]
	var m12 = mm[1*4+2]
	var m13 = mm[1*4+3]
	var m20 = mm[2*4+0]
	var m21 = mm[2*4+1]
	var m22 = mm[2*4+2]
	var m23 = mm[2*4+3]
	var m30 = mm[3*4+0]
	var m31 = mm[3*4+1]
	var m32 = mm[3*4+2]
	var m33 = mm[3*4+3]
	var tmp_0 = m22 * m33
	var tmp_1 = m32 * m23
	var tmp_2 = m12 * m33
	var tmp_3 = m32 * m13
	var tmp_4 = m12 * m23
	var tmp_5 = m22 * m13
	var tmp_6 = m02 * m33
	var tmp_7 = m32 * m03
	var tmp_8 = m02 * m23
	var tmp_9 = m22 * m03
	var tmp_10 = m02 * m13
	var tmp_11 = m12 * m03
	var tmp_12 = m20 * m31
	var tmp_13 = m30 * m21
	var tmp_14 = m10 * m31
	var tmp_15 = m30 * m11
	var tmp_16 = m10 * m21
	var tmp_17 = m20 * m11
	var tmp_18 = m00 * m31
	var tmp_19 = m30 * m01
	var tmp_20 = m00 * m21
	var tmp_21 = m20 * m01
	var tmp_22 = m00 * m11
	var tmp_23 = m10 * m01

	var t0 = (tmp_0*m11 + tmp_3*m21 + tmp_4*m31) -
		(tmp_1*m11 + tmp_2*m21 + tmp_5*m31)
	var t1 = (tmp_1*m01 + tmp_6*m21 + tmp_9*m31) -
		(tmp_0*m01 + tmp_7*m21 + tmp_8*m31)
	var t2 = (tmp_2*m01 + tmp_7*m11 + tmp_10*m31) -
		(tmp_3*m01 + tmp_6*m11 + tmp_11*m31)
	var t3 = (tmp_5*m01 + tmp_8*m11 + tmp_11*m21) -
		(tmp_4*m01 + tmp_9*m11 + tmp_10*m21)

	var d = 1.0 / (m00*t0 + m10*t1 + m20*t2 + m30*t3)

	dst[0] = d * t0
	dst[1] = d * t1
	dst[2] = d * t2
	dst[3] = d * t3
	dst[4] = d * ((tmp_1*m10 + tmp_2*m20 + tmp_5*m30) -
		(tmp_0*m10 + tmp_3*m20 + tmp_4*m30))
	dst[5] = d * ((tmp_0*m00 + tmp_7*m20 + tmp_8*m30) -
		(tmp_1*m00 + tmp_6*m20 + tmp_9*m30))
	dst[6] = d * ((tmp_3*m00 + tmp_6*m10 + tmp_11*m30) -
		(tmp_2*m00 + tmp_7*m10 + tmp_10*m30))
	dst[7] = d * ((tmp_4*m00 + tmp_9*m10 + tmp_10*m20) -
		(tmp_5*m00 + tmp_8*m10 + tmp_11*m20))
	dst[8] = d * ((tmp_12*m13 + tmp_15*m23 + tmp_16*m33) -
		(tmp_13*m13 + tmp_14*m23 + tmp_17*m33))
	dst[9] = d * ((tmp_13*m03 + tmp_18*m23 + tmp_21*m33) -
		(tmp_12*m03 + tmp_19*m23 + tmp_20*m33))
	dst[10] = d * ((tmp_14*m03 + tmp_19*m13 + tmp_22*m33) -
		(tmp_15*m03 + tmp_18*m13 + tmp_23*m33))
	dst[11] = d * ((tmp_17*m03 + tmp_20*m13 + tmp_23*m23) -
		(tmp_16*m03 + tmp_21*m13 + tmp_22*m23))
	dst[12] = d * ((tmp_14*m22 + tmp_17*m32 + tmp_13*m12) -
		(tmp_16*m32 + tmp_12*m12 + tmp_15*m22))
	dst[13] = d * ((tmp_20*m32 + tmp_12*m02 + tmp_19*m22) -
		(tmp_18*m22 + tmp_21*m32 + tmp_13*m02))
	dst[14] = d * ((tmp_18*m12 + tmp_23*m32 + tmp_15*m02) -
		(tmp_22*m32 + tmp_14*m02 + tmp_19*m12))
	dst[15] = d * ((tmp_22*m22 + tmp_16*m02 + tmp_21*m12) -
		(tmp_20*m12 + tmp_23*m22 + tmp_17*m02))

	return result
}

//Decompose turns a matrix into an slice of floats
func (m Matrix) Decompose() []float32 {
	return []float32{m.M0, m.M1, m.M2, m.M3, m.M4, m.M5, m.M6, m.M7, m.M8, m.M9, m.M10, m.M11, m.M12, m.M13, m.M14, m.M15}
}

//DecomposePointer is an unsafe Decompose. Instead of the values being copied, a pointer to the matrix is cast into a float array pointer and returned.
func (m *Matrix) DecomposePointer() *[16]float32 {
	return (*[16]float32)(unsafe.Pointer(m))
}

package noodle

import (
	"math"
	"unsafe"
)

//Matrix A representation of a 4 x 4 matrix
type Matrix struct {
	M00 float32
	M10 float32
	M20 float32
	M30 float32

	M01 float32
	M11 float32
	M21 float32
	M31 float32

	M02 float32
	M12 float32
	M22 float32
	M32 float32

	M03 float32
	M13 float32
	M23 float32
	M33 float32
}

//NewMatrix creates a identity
func NewMatrix() Matrix {
	return Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

//NewMatrixTranslate creates a new translate matrix
func NewMatrixTranslate(v Vector3) Matrix {
	var r Matrix
	r.M00 = 1
	r.M11 = 1
	r.M22 = 1
	r.M33 = 1
	r.M30 = v.X
	r.M31 = v.Y
	r.M32 = v.Z
	return r
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

//NewMatrixRotation creates a new rotation matrix
func NewMatrixRotation(q Quaternion) Matrix {
	//https://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToMatrix/index.htm

	var invs, tmp1, tmp2 float32
	var m00, m11, m22, m10, m01, m20, m02, m21, m12 float32

	sqx := q.X * q.X
	sqy := q.Y * q.Y
	sqz := q.Z * q.Z
	sqw := q.W * q.W

	// invs (inverse square length) is only required if quaternion is not already normalised
	invs = 1.0 / (sqx + sqy + sqz + sqw)
	m00 = (sqx - sqy - sqz + sqw) * invs // since sqw + sqx + sqy + sqz =1/invs*invs
	m11 = (-sqx + sqy - sqz + sqw) * invs
	m22 = (-sqx - sqy + sqz + sqw) * invs

	tmp1 = q.X * q.Y
	tmp2 = q.Z * q.W
	m10 = 2.0 * (tmp1 + tmp2) * invs
	m01 = 2.0 * (tmp1 - tmp2) * invs

	tmp1 = q.X * q.Z
	tmp2 = q.Y * q.W
	m20 = 2.0 * (tmp1 - tmp2) * invs
	m02 = 2.0 * (tmp1 + tmp2) * invs

	tmp1 = q.Y * q.Z
	tmp2 = q.X * q.W
	m21 = 2.0 * (tmp1 + tmp2) * invs
	m12 = 2.0 * (tmp1 - tmp2) * invs
	return Matrix{
		M00: m00,
		M11: m11,
		M22: m22,
		M10: m10,
		M01: m01,
		M20: m20,
		M02: m02,
		M21: m21,
		M12: m12,
	}
}

//NewMatrixPerspective creates a perspective projection matrix. FOVY is in degrees
func NewMatrixPerspective(fovy, aspectRatio, near, far float32) Matrix {
	fieldOfViewInRadians := float64(fovy) * (math.Pi / 180.0)
	f := float32(1.0 / math.Tan(fieldOfViewInRadians/2.0))
	rangeInv := 1.0 / (near - far)

	return Matrix{
		f / aspectRatio, 0, 0, 0,
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

	var r Matrix
	r.M00 = float32((2. * near) / rml)
	r.M11 = float32((2. * near) / tmb)
	r.M20 = A
	r.M21 = B
	r.M22 = C
	r.M23 = -1
	r.M32 = D
	return r
}

//NewMatrixLookAt creates a matrix to look at a target
func NewMatrixLookAt(eye, target, up Vector3) Matrix {
	f := target.Subtract(eye).Normalize()
	s := f.CrossProduct(up.Normalize()).Normalize()
	u := s.CrossProduct(f)
	matrix := Matrix{
		M00: float32(s.X),
		M10: float32(u.X),
		M20: float32(-f.X),
		M30: 0,

		M01: float32(s.Y),
		M11: float32(u.Y),
		M21: float32(-f.Y),
		M31: 0,

		M02: float32(s.Z),
		M12: float32(u.Z),
		M22: float32(-f.Z),
		M32: 0,

		M03: 0,
		M13: 0,
		M23: 0,
		M33: 1,
	}

	//	return M.Mul4(Translate3D(float32(-eye[0]), float32(-eye[1]), float32(-eye[2])))
	return matrix.Multiply(NewMatrixTranslate(Vector3{-eye.X, -eye.Y, -eye.Z}))
}

//Trace of the matrix (sum of values along diagonal)
func (m Matrix) Trace() float32 {
	return m.M00 + m.M11 + m.M22 + m.M33
}

//Add two matrices
func (m Matrix) Add(right Matrix) Matrix {
	return Matrix{
		m.M00 + right.M00, m.M01 + right.M01, m.M02 + right.M02, m.M03 + right.M03,
		m.M10 + right.M10, m.M11 + right.M11, m.M12 + right.M12, m.M13 + right.M13,
		m.M20 + right.M20, m.M21 + right.M21, m.M22 + right.M22, m.M23 + right.M23,
		m.M30 + right.M30, m.M31 + right.M31, m.M32 + right.M32, m.M33 + right.M33,
	}
}

//Subtract two matrices
func (m Matrix) Subtract(right Matrix) Matrix {
	return Matrix{
		m.M00 - right.M00, m.M01 - right.M01, m.M02 - right.M02, m.M03 - right.M03,
		m.M10 - right.M10, m.M11 - right.M11, m.M12 - right.M12, m.M13 - right.M13,
		m.M20 - right.M20, m.M21 - right.M21, m.M22 - right.M22, m.M23 - right.M23,
		m.M30 - right.M30, m.M31 - right.M31, m.M32 - right.M32, m.M33 - right.M33,
	}
}

//Multiply two matrix together. Note that order matters.
func (m Matrix) Multiply(m2 Matrix) Matrix {
	var r Matrix
	r.M00 = m.M00*m2.M00 + m.M01*m2.M10 + m.M02*m2.M20 + m.M03*m2.M30
	r.M01 = m.M00*m2.M01 + m.M01*m2.M11 + m.M02*m2.M21 + m.M03*m2.M31
	r.M02 = m.M00*m2.M02 + m.M01*m2.M12 + m.M02*m2.M22 + m.M03*m2.M32
	r.M03 = m.M00*m2.M03 + m.M01*m2.M13 + m.M02*m2.M23 + m.M03*m2.M33

	r.M10 = m.M10*m2.M00 + m.M11*m2.M10 + m.M12*m2.M20 + m.M13*m2.M30
	r.M11 = m.M10*m2.M01 + m.M11*m2.M11 + m.M12*m2.M21 + m.M13*m2.M31
	r.M12 = m.M10*m2.M02 + m.M11*m2.M12 + m.M12*m2.M22 + m.M13*m2.M32
	r.M13 = m.M10*m2.M03 + m.M11*m2.M13 + m.M12*m2.M23 + m.M13*m2.M33

	r.M20 = m.M20*m2.M00 + m.M21*m2.M10 + m.M22*m2.M20 + m.M23*m2.M30
	r.M21 = m.M20*m2.M01 + m.M21*m2.M11 + m.M22*m2.M21 + m.M23*m2.M31
	r.M22 = m.M20*m2.M02 + m.M21*m2.M12 + m.M22*m2.M22 + m.M23*m2.M32
	r.M23 = m.M20*m2.M03 + m.M21*m2.M13 + m.M22*m2.M23 + m.M23*m2.M33

	r.M30 = m.M30*m2.M00 + m.M31*m2.M10 + m.M32*m2.M20 + m.M33*m2.M30
	r.M31 = m.M30*m2.M01 + m.M31*m2.M11 + m.M32*m2.M21 + m.M33*m2.M31
	r.M32 = m.M30*m2.M02 + m.M31*m2.M12 + m.M32*m2.M22 + m.M33*m2.M32
	r.M33 = m.M30*m2.M03 + m.M31*m2.M13 + m.M32*m2.M23 + m.M33*m2.M33
	return r
}

//MultiplyVector3 multiplies a vector with the matrix (m * v)
func (m Matrix) MultiplyVector3(v Vector3) Vector3 {
	var r Vector3

	fInvW := 1.0 / (m.M30*v.X + m.M31*v.Y + m.M32*v.Z + m.M33)
	r.X = (m.M00*v.X + m.M01*v.Y + m.M02*v.Z + m.M03) * fInvW
	r.Y = (m.M10*v.X + m.M11*v.Y + m.M12*v.Z + m.M13) * fInvW
	r.Z = (m.M20*v.X + m.M21*v.Y + m.M22*v.Z + m.M23) * fInvW
	return r
}

//MultiplyVector4 multiplies a vector with the matrix (m * v)
func (m Matrix) MultiplyVector4(v Vector4) Vector4 {
	var r Vector4

	r.X = m.M00*v.X + m.M01*v.Y + m.M02*v.Z + m.M03*v.W
	r.Y = m.M10*v.X + m.M11*v.Y + m.M12*v.Z + m.M13*v.W
	r.Z = m.M20*v.X + m.M21*v.Y + m.M22*v.Z + m.M23*v.W
	r.W = m.M30*v.X + m.M31*v.Y + m.M32*v.Z + m.M33*v.W
	return r
}

//Translation gets the translation component
func (m Matrix) Translation() Vector3 {
	return Vector3{m.M03, m.M13, m.M23}
}

//Decompose turns a matrix into an slice of floats
func (m Matrix) Decompose() []float32 {
	return []float32{m.M00, m.M10, m.M20, m.M30, m.M01, m.M11, m.M21, m.M31, m.M02, m.M12, m.M22, m.M32, m.M03, m.M13, m.M23, m.M33}
}

//DecomposePointer is an unsafe Decompose. Instead of the values being copied, a pointer to the matrix is cast into a float array pointer and returned.
func (m *Matrix) DecomposePointer() *[16]float32 {
	return (*[16]float32)(unsafe.Pointer(m))
}

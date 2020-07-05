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
	return Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
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

//NewMatrixPerspective creates a perspective projection matrix. FOVY is in degrees
func NewMatrixPerspective(fovy, aspect, near, far float32) Matrix {

	nmf := near - far
	rangeInv := 1.0 / nmf
	f := float32(math.Tan(math.Pi*0.5 - 0.5*float64(fovy*Deg2Rad)))
	return Matrix{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (near + far) * rangeInv, -1,
		0, 0, near * far * rangeInv * 2, 0,
	}

	/*
		//Identitcal
		fieldOfViewInRadians := float64(fovy) * (math.Pi / 180.0)
		nmf, f := near-far, float32(1./math.Tan(fieldOfViewInRadians/2.0))
		return Matrix{float32(f / aspect), 0, 0, 0, 0, float32(f), 0, 0, 0, 0, float32((near + far) / nmf), -1, 0, 0, float32((2. * far * near) / nmf), 0}
	*/
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
	f := target.Subtract(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)

	M := Matrix{
		s.X, u.X, -f.X, 0,
		s.Y, u.Y, -f.Y, 0,
		s.Z, u.Z, -f.Z, 0,
		0, 0, 0, 1,
	}

	return M.Multiply(NewMatrixTranslate(Vector3{float32(-eye.X), float32(-eye.Y), float32(-eye.Z)}))
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

//Scale a matrix
func (m Matrix) Scale(c float32) Matrix {
	m1 := m.DecomposePointer()
	return Matrix{m1[0] * c, m1[1] * c, m1[2] * c, m1[3] * c, m1[4] * c, m1[5] * c, m1[6] * c, m1[7] * c, m1[8] * c, m1[9] * c, m1[10] * c, m1[11] * c, m1[12] * c, m1[13] * c, m1[14] * c, m1[15] * c}
}

//Multiply 2 matrixs together
func (m Matrix) Multiply(right Matrix) Matrix {
	m1 := m.DecomposePointer()
	m2 := m.DecomposePointer()
	return Matrix{
		m1[0]*m2[0] + m1[4]*m2[1] + m1[8]*m2[2] + m1[12]*m2[3],
		m1[1]*m2[0] + m1[5]*m2[1] + m1[9]*m2[2] + m1[13]*m2[3],
		m1[2]*m2[0] + m1[6]*m2[1] + m1[10]*m2[2] + m1[14]*m2[3],
		m1[3]*m2[0] + m1[7]*m2[1] + m1[11]*m2[2] + m1[15]*m2[3],
		m1[0]*m2[4] + m1[4]*m2[5] + m1[8]*m2[6] + m1[12]*m2[7],
		m1[1]*m2[4] + m1[5]*m2[5] + m1[9]*m2[6] + m1[13]*m2[7],
		m1[2]*m2[4] + m1[6]*m2[5] + m1[10]*m2[6] + m1[14]*m2[7],
		m1[3]*m2[4] + m1[7]*m2[5] + m1[11]*m2[6] + m1[15]*m2[7],
		m1[0]*m2[8] + m1[4]*m2[9] + m1[8]*m2[10] + m1[12]*m2[11],
		m1[1]*m2[8] + m1[5]*m2[9] + m1[9]*m2[10] + m1[13]*m2[11],
		m1[2]*m2[8] + m1[6]*m2[9] + m1[10]*m2[10] + m1[14]*m2[11],
		m1[3]*m2[8] + m1[7]*m2[9] + m1[11]*m2[10] + m1[15]*m2[11],
		m1[0]*m2[12] + m1[4]*m2[13] + m1[8]*m2[14] + m1[12]*m2[15],
		m1[1]*m2[12] + m1[5]*m2[13] + m1[9]*m2[14] + m1[13]*m2[15],
		m1[2]*m2[12] + m1[6]*m2[13] + m1[10]*m2[14] + m1[14]*m2[15],
		m1[3]*m2[12] + m1[7]*m2[13] + m1[11]*m2[14] + m1[15]*m2[15],
	}
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

//TransformCoordinate multiplies a 3D vector by a transformation given by
// the homogeneous 4D matrix m, applying any translation.
// If this transformation is non-affine, it will project this
// vector onto the plane w=1 before returning the result.
func (m Matrix) TransformCoordinate(v Vector3) Vector3 {
	t := v.ToVector4()
	t = m.MultiplyVector4(t)
	t = t.Scale(1 / t.W)
	return Vector3{t.X, t.Y, t.Z}
}

//Decompose turns a matrix into an slice of floats
func (m Matrix) Decompose() []float32 {
	return []float32{m.M0, m.M1, m.M2, m.M3, m.M4, m.M5, m.M6, m.M7, m.M8, m.M9, m.M10, m.M11, m.M12, m.M13, m.M14, m.M15}
}

//DecomposePointer is an unsafe Decompose. Instead of the values being copied, a pointer to the matrix is cast into a float array pointer and returned.
func (m *Matrix) DecomposePointer() *[16]float32 {
	return (*[16]float32)(unsafe.Pointer(m))
}

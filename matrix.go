package noodle

import (
	"math"
	"unsafe"
)

//Matrix A representation of a 4 x 4 matrix
type Matrix struct {
	//0   4   8   12
	//1   5   9   13
	//2   6   10  14
	//3   7   11  15
	//0 4 8 12    1 5 9 13    2 6 10 14    3 7 11 15
	/*
		M0  float32
		M4  float32
		M8  float32
		M12 float32

		M1  float32
		M5  float32
		M9  float32
		M13 float32

		M2  float32
		M6  float32
		M10 float32
		M14 float32

		M3  float32
		M7  float32
		M11 float32
		M15 float32
	*/

	//0    1    2    3
	//4    5    6    7
	//8    9    10   11
	//12   13   14   15

	M0 float32 //M0 is 0,0
	M1 float32 //M1 is 1,0
	M2 float32 //M2 is 2,0
	M3 float32 //M3 is 3,0

	M4 float32 //M4 is 0,1
	M5 float32 //M5 is 1,1
	M6 float32 //M6 is 2,1
	M7 float32 //M7 is 3,1

	M8  float32 //M8 is 0,2
	M9  float32 //M9 is 1,2
	M10 float32 //M10 is 2,2
	M11 float32 //M11 is 3,2

	M12 float32 //M12 is 0,3
	M13 float32 //M13 is 1,3
	M14 float32 //M14 is 2,3
	M15 float32 //M15 is 3,3
}

func newMatrixFromPointer(ptr unsafe.Pointer) Matrix { return *(*Matrix)(ptr) }

//NewMatrixQuaternion creates a new rotation matrix from a quaternion
// https://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToMatrix/index.htm
func NewMatrixQuaternion(q Quaternion) Matrix {
	sqw := q.W * q.W
	sqx := q.X * q.X
	sqy := q.Y * q.Y
	sqz := q.Z * q.Z

	invs := 1 / (sqx + sqy + sqz + sqw)
	m00 := (sqx - sqy - sqz + sqw) * invs // since sqw + sqx + sqy + sqz =1/invs*invs
	m11 := (-sqx + sqy - sqz + sqw) * invs
	m22 := (-sqx - sqy + sqz + sqw) * invs

	tmp1 := q.X * q.Y
	tmp2 := q.Z * q.W
	m10 := 2.0 * (tmp1 + tmp2) * invs
	m01 := 2.0 * (tmp1 - tmp2) * invs

	tmp1 = q.X * q.Z
	tmp2 = q.Y * q.W
	m20 := 2.0 * (tmp1 - tmp2) * invs
	m02 := 2.0 * (tmp1 + tmp2) * invs

	tmp1 = q.Y * q.Z
	tmp2 = q.X * q.W
	m21 := 2.0 * (tmp1 + tmp2) * invs
	m12 := 2.0 * (tmp1 - tmp2) * invs
	return Matrix{
		m00, m01, m02, 0,
		m10, m11, m12, 0,
		m20, m21, m22, 0,
		0, 0, 0, 1,
	}
}

//NewMatrixIdentity creates a identity
func NewMatrixIdentity() Matrix {
	return Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

//NewMatrixTranslate creates a blank translation matrix from vector
func NewMatrixTranslate(v Vector3) Matrix {
	return Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, v.X, v.Y, v.Z, 1}
}

//NewMatrixTranslate32 creates a blank translation matrix
func NewMatrixTranslate32(x, y, z float32) Matrix {
	return Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, x, y, z, 1}
}

//NewMatrixTranslate64 creates a blank translation matrix
func NewMatrixTranslate64(x, y, z float64) Matrix {
	return NewMatrixTranslate32(float32(x), float32(y), float32(z))
}

//NewMatrixRotate creates a rotation matrix based of the axis and radians
func NewMatrixRotate(axis Vector3, radians float32) Matrix {
	x := float32(axis.X)
	y := float32(axis.Y)
	z := float32(axis.Z)
	length := float32(axis.Length())
	if length != 1 && length != 0 {
		length = 1 / length
		x *= length
		y *= length
		z *= length
	}

	sinres := float32(math.Sin(float64(radians)))
	cosres := float32(math.Cos(float64(radians)))
	t := 1 - cosres
	return Matrix{
		M0: x*x*t + cosres,
		M1: y*x*t + z*sinres,
		M2: z*x*t - y*sinres,
		M3: 0,

		M4: x*y*t - z*sinres,
		M5: y*y*t + cosres,
		M6: z*y*t + x*sinres,
		M7: 0,

		M8:  x*z*t + y*sinres,
		M9:  y*z*t - x*sinres,
		M10: z*z*t + cosres,
		M11: 0,

		M12: 0,
		M13: 0,
		M14: 0,
		M15: 1,
	}
}

//NewMatrixRotateXYZ new xyz-rotation matrix (in radians)
func NewMatrixRotateXYZ(radians Vector3) Matrix {
	cosz := float32(math.Cos(float64(-radians.Z)))
	sinz := float32(math.Sin(float64(-radians.Z)))
	cosy := float32(math.Cos(float64(-radians.Y)))
	siny := float32(math.Sin(float64(-radians.Y)))
	cosx := float32(math.Cos(float64(-radians.X)))
	sinx := float32(math.Sin(float64(-radians.X)))
	result := NewMatrixIdentity()
	result.M0 = cosz * cosy
	result.M4 = (cosz * siny * sinx) - (sinz * cosx)
	result.M8 = (cosz * siny * cosx) + (sinz * sinx)
	result.M1 = sinz * cosy
	result.M5 = (sinz * siny * sinx) + (cosz * cosx)
	result.M9 = (sinz * siny * cosx) - (cosz * sinx)
	result.M2 = -siny
	result.M6 = cosy * sinx
	result.M10 = cosy * cosx
	return result

}

//NewMatrixRotateX creates a new matrix that is rotated
func NewMatrixRotateX(radians float32) Matrix {
	result := NewMatrixIdentity()
	cosres := float32(math.Cos(float64(radians)))
	sinres := float32(math.Sin(float64(radians)))
	result.M5 = cosres
	result.M6 = -sinres
	result.M9 = sinres
	result.M10 = cosres
	return result
}

//NewMatrixRotateY creates a new matrix that is rotated
func NewMatrixRotateY(radians float32) Matrix {
	result := NewMatrixIdentity()
	cosres := float32(math.Cos(float64(radians)))
	sinres := float32(math.Sin(float64(radians)))
	result.M0 = cosres
	result.M2 = sinres
	result.M8 = -sinres
	result.M10 = cosres
	return result

}

//NewMatrixRotateZ creates a new matrix that is rotated
func NewMatrixRotateZ(radians float32) Matrix {
	result := NewMatrixIdentity()
	cosres := float32(math.Cos(float64(radians)))
	sinres := float32(math.Sin(float64(radians)))
	result.M0 = cosres
	result.M1 = -sinres
	result.M4 = sinres
	result.M5 = cosres
	return result

}

//NewMatrixRotateAxis creates a rotation matrix from the axis
func NewMatrixRotateAxis(xAxis, yAxis, zAxis Vector3) Matrix {
	return Matrix{
		M0:  xAxis.X,
		M1:  yAxis.X,
		M2:  zAxis.X,
		M3:  0,
		M4:  xAxis.Y,
		M5:  yAxis.Y,
		M6:  zAxis.Y,
		M7:  0,
		M8:  xAxis.Z,
		M9:  yAxis.Z,
		M10: zAxis.Z,
		M11: 0,
		M12: 0,
		M13: 0,
		M14: 0,
		M15: 1,
	}
}

//NewMatrixScale creates a new scalling matrix
func NewMatrixScale(scale Vector3) Matrix {
	return Matrix{
		M0: float32(scale.X), M1: 0, M2: 0, M3: 0,
		M4: 0, M5: float32(scale.Y), M6: 0, M7: 0,
		M8: 0, M9: 0, M10: float32(scale.Z), M11: 0,
		M12: 0, M13: 0, M14: 0, M15: 1,
	}
}

//NewMatrixPerspective creates a perspective projection matrix. FOVY is in degrees
func NewMatrixPerspective(fovy, aspect, near, far float32) Matrix {
	fovy = fovy * Deg2Rad
	nmf, f := near-far, float32(1./math.Tan(float64(fovy)/2.0))
	return Matrix{float32(f / aspect), 0, 0, 0, 0, float32(f), 0, 0, 0, 0, float32((near + far) / nmf), -1, 0, 0, float32((2. * far * near) / nmf), 0}
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
	s := f.CrossProduct(up.Normalize()).Normalize()
	u := s.CrossProduct(f)
	matrix := Matrix{
		M0: float32(s.X),
		M1: float32(u.X),
		M2: float32(-f.X),
		M3: 0,

		M4: float32(s.Y),
		M5: float32(u.Y),
		M6: float32(-f.Y),
		M7: 0,

		M8:  float32(s.Z),
		M9:  float32(u.Z),
		M10: float32(-f.Z),
		M11: 0,

		M12: 0,
		M13: 0,
		M14: 0,
		M15: 1,
	}

	//	return M.Mul4(Translate3D(float32(-eye[0]), float32(-eye[1]), float32(-eye[2])))
	return matrix.Multiply(NewMatrixTranslate32(-eye.X, -eye.Y, -eye.Z))
}

/*
//NewMatrixTransform creates a new matrix based off a transform
func NewMatrixTransform(transform Transform) Matrix {
	return NewMatrixTranslate(transform.Position).Multiply(NewMatrixQuaternion(transform.Rotation)).Multiply(NewMatrixScale(transform.Scale))
}
*/

//GetTranslation gets the matrix translation
func (m Matrix) GetTranslation() Vector3 {
	return Vector3{m.M12, m.M13, m.M14}
}

//Trace of the matrix (sum of values along diagonal)
func (m Matrix) Trace() float32 {
	return m.M0 + m.M5 + m.M10 + m.M15
}

//Detrimant of the matrix
func (m Matrix) Detrimant() float32 {
	// Cache the matrix values (speed optimization)
	a00 := m.M0
	a01 := m.M1
	a02 := m.M2
	a03 := m.M3
	a10 := m.M4
	a11 := m.M5
	a12 := m.M6
	a13 := m.M7
	a20 := m.M8
	a21 := m.M9
	a22 := m.M10
	a23 := m.M11
	a30 := m.M12
	a31 := m.M13
	a32 := m.M14
	a33 := m.M15

	return a30*a21*a12*a03 - a20*a31*a12*a03 - a30*a11*a22*a03 + a10*a31*a22*a03 +
		a20*a11*a32*a03 - a10*a21*a32*a03 - a30*a21*a02*a13 + a20*a31*a02*a13 +
		a30*a01*a22*a13 - a00*a31*a22*a13 - a20*a01*a32*a13 + a00*a21*a32*a13 +
		a30*a11*a02*a23 - a10*a31*a02*a23 - a30*a01*a12*a23 + a00*a31*a12*a23 +
		a10*a01*a32*a23 - a00*a11*a32*a23 - a20*a11*a02*a33 + a10*a21*a02*a33 +
		a20*a01*a12*a33 - a00*a21*a12*a33 - a10*a01*a22*a33 + a00*a11*a22*a33
}

//Transpose the matrix
func (m Matrix) Transpose() Matrix {
	return Matrix{
		M0:  m.M0,
		M1:  m.M4,
		M2:  m.M8,
		M3:  m.M12,
		M4:  m.M1,
		M5:  m.M5,
		M6:  m.M9,
		M7:  m.M13,
		M8:  m.M2,
		M9:  m.M6,
		M10: m.M10,
		M11: m.M14,
		M12: m.M3,
		M13: m.M7,
		M14: m.M11,
		M15: m.M15,
	}
}

//Invert the matrix
func (m Matrix) Invert() Matrix {
	a00 := m.M0
	a01 := m.M1
	a02 := m.M2
	a03 := m.M3
	a10 := m.M4
	a11 := m.M5
	a12 := m.M6
	a13 := m.M7
	a20 := m.M8
	a21 := m.M9
	a22 := m.M10
	a23 := m.M11
	a30 := m.M12
	a31 := m.M13
	a32 := m.M14
	a33 := m.M15

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	// Calculate the invert determinant (inlined to avoid double-caching)
	invDet := 1 / (b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06)
	return Matrix{
		M0:  (a11*b11 - a12*b10 + a13*b09) * invDet,
		M1:  (-a01*b11 + a02*b10 - a03*b09) * invDet,
		M2:  (a31*b05 - a32*b04 + a33*b03) * invDet,
		M3:  (-a21*b05 + a22*b04 - a23*b03) * invDet,
		M4:  (-a10*b11 + a12*b08 - a13*b07) * invDet,
		M5:  (a00*b11 - a02*b08 + a03*b07) * invDet,
		M6:  (-a30*b05 + a32*b02 - a33*b01) * invDet,
		M7:  (a20*b05 - a22*b02 + a23*b01) * invDet,
		M8:  (a10*b10 - a11*b08 + a13*b06) * invDet,
		M9:  (-a00*b10 + a01*b08 - a03*b06) * invDet,
		M10: (a30*b04 - a31*b02 + a33*b00) * invDet,
		M11: (-a20*b04 + a21*b02 - a23*b00) * invDet,
		M12: (-a10*b09 + a11*b07 - a12*b06) * invDet,
		M13: (a00*b09 - a01*b07 + a02*b06) * invDet,
		M14: (-a30*b03 + a31*b01 - a32*b00) * invDet,
		M15: (a20*b03 - a21*b01 + a22*b00) * invDet,
	}
}

//Normalize calcuates the normal of the matrix
func (m Matrix) Normalize() Matrix {
	det := m.Detrimant()
	return Matrix{
		M0:  m.M0 / det,
		M1:  m.M1 / det,
		M2:  m.M2 / det,
		M3:  m.M3 / det,
		M4:  m.M4 / det,
		M5:  m.M5 / det,
		M6:  m.M6 / det,
		M7:  m.M7 / det,
		M8:  m.M8 / det,
		M9:  m.M9 / det,
		M10: m.M10 / det,
		M11: m.M11 / det,
		M12: m.M12 / det,
		M13: m.M13 / det,
		M14: m.M14 / det,
		M15: m.M15 / det,
	}
}

//Add two matrices
func (m Matrix) Add(right Matrix) Matrix {
	return Matrix{
		M0:  m.M0 + right.M0,
		M1:  m.M1 + right.M1,
		M2:  m.M2 + right.M2,
		M3:  m.M3 + right.M3,
		M4:  m.M4 + right.M4,
		M5:  m.M5 + right.M5,
		M6:  m.M6 + right.M6,
		M7:  m.M7 + right.M7,
		M8:  m.M8 + right.M8,
		M9:  m.M9 + right.M9,
		M10: m.M10 + right.M10,
		M11: m.M11 + right.M11,
		M12: m.M12 + right.M12,
		M13: m.M13 + right.M13,
		M14: m.M14 + right.M14,
		M15: m.M15 + right.M15,
	}
}

//Subtract two matrices
func (m Matrix) Subtract(right Matrix) Matrix {
	return Matrix{
		M0:  m.M0 - right.M0,
		M1:  m.M1 - right.M1,
		M2:  m.M2 - right.M2,
		M3:  m.M3 - right.M3,
		M4:  m.M4 - right.M4,
		M5:  m.M5 - right.M5,
		M6:  m.M6 - right.M6,
		M7:  m.M7 - right.M7,
		M8:  m.M8 - right.M8,
		M9:  m.M9 - right.M9,
		M10: m.M10 - right.M10,
		M11: m.M11 - right.M11,
		M12: m.M12 - right.M12,
		M13: m.M13 - right.M13,
		M14: m.M14 - right.M14,
		M15: m.M15 - right.M15,
	}
}

//Multiply two matrix together. Note that order matters.
func (m Matrix) Multiply(right Matrix) Matrix {
	m1 := m.DecomposePointer()
	m2 := right.DecomposePointer()
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

//Decompose turns a matrix into an slice of floats
func (m Matrix) Decompose() []float32 {
	return []float32{
		m.M0, m.M1, m.M2, m.M3,
		m.M4, m.M5, m.M6, m.M7,
		m.M8, m.M9, m.M10, m.M11,
		m.M12, m.M13, m.M14, m.M15,
	}
}

//DecomposePointer is an unsafe Decompose. Instead of the values being copied, a pointer to the matrix is cast into a float array pointer and returned.
func (m *Matrix) DecomposePointer() *[16]float32 {
	return (*[16]float32)(unsafe.Pointer(m))
}

package main

import (
	"log"

	n "github.com/lachee/noodle"
)

//AxisTestApp shows off the 3D capabilities
type AxisTestApp struct {
	positionBuffer n.WebGLBuffer
	colorBuffer    n.WebGLBuffer

	shader            *n.Shader
	positionAttribute n.WebGLAttributeLocation
	colorAttribute    n.WebGLAttributeLocation
	matrixLocation    n.WebGLUniformLocation
}

//Start is called by the noodle engine when ready
func (app *AxisTestApp) Start() bool {

	//Load the cube shader
	shader, err := n.LoadShaderFromURL("resources/shader/axis.vert", "resources/shader/axis.frag")
	if err != nil {
		log.Fatalln("Failed to load the shaders! ", err)
		return false
	}
	app.shader = shader
	app.positionAttribute = shader.GetAttribLocation("a_position")
	app.colorAttribute = shader.GetAttribLocation("a_color")
	app.matrixLocation = shader.GetUniformLocation("u_matrix")

	// Create vertex buffer
	app.positionBuffer = n.GL.CreateBuffer()
	n.GL.BindBuffer(n.GlArrayBuffer, app.positionBuffer)
	app.setGeometry()

	app.colorBuffer = n.GL.CreateBuffer()
	n.GL.BindBuffer(n.GlArrayBuffer, app.colorBuffer)
	app.setColors()

	return true
}

//Update occurs once a frame
func (app *AxisTestApp) Update(dt float32) {
}

//Render occurs when the screen needs updating
func (app *AxisTestApp) Render() {

	translation := Vector3{-150, 0, -360}
	rotation := Vector3{190, 40, 320}
	scale := Vector3{1, 1, 1}

	//Clearing and Shader
	n.GL.Clear(n.GlColorBufferBit | n.GlDepthBufferBit)
	n.GL.Enable(n.GlCullFace)
	n.GL.Enable(n.GlDepthTest)
	app.shader.Use()

	//Buffers
	n.GL.EnableVertexAttribArray(app.positionAttribute)                        //Turn on the position attribute
	n.GL.BindBuffer(n.GlArrayBuffer, app.positionBuffer)                       //Bind the buffer
	n.GL.VertexAttribPointer(app.positionAttribute, 3, n.GlFloat, false, 0, 0) //Tell the position attribute how to get the data out of the buffer

	n.GL.EnableVertexAttribArray(app.colorAttribute)
	n.GL.BindBuffer(n.GlArrayBuffer, app.colorBuffer)
	n.GL.VertexAttribPointer(app.colorAttribute, 3, n.GlUnsignedByte, true, 0, 0)

	//Matrix
	matrix := n.NewMatrixPerspective(60, n.GL.AspectRatio(), 1, 2000)
	matrix = matrix.Translate(translation)
	matrix = matrix.RotateX(rotation.X * n.Deg2Rad)
	matrix = matrix.RotateY(rotation.Y * n.Deg2Rad)
	matrix = matrix.RotateZ(rotation.Z * n.Deg2Rad)
	matrix = matrix.Scale(scale)
	n.GL.UniformMatrix4fv(app.matrixLocation, matrix)

	//Drawing
	n.GL.DrawArrays(n.GlTriangles, 0, 16*6)
}

//setGeometry fills the buffer with the F values
func (app *AxisTestApp) setGeometry() {
	n.GL.BufferData(n.GlArrayBuffer, []float32{
		// left column front
		0, 0, 0,
		0, 150, 0,
		30, 0, 0,
		0, 150, 0,
		30, 150, 0,
		30, 0, 0,

		// top rung front
		30, 0, 0,
		30, 30, 0,
		100, 0, 0,
		30, 30, 0,
		100, 30, 0,
		100, 0, 0,

		// middle rung front
		30, 60, 0,
		30, 90, 0,
		67, 60, 0,
		30, 90, 0,
		67, 90, 0,
		67, 60, 0,

		// left column back
		0, 0, 30,
		30, 0, 30,
		0, 150, 30,
		0, 150, 30,
		30, 0, 30,
		30, 150, 30,

		// top rung back
		30, 0, 30,
		100, 0, 30,
		30, 30, 30,
		30, 30, 30,
		100, 0, 30,
		100, 30, 30,

		// middle rung back
		30, 60, 30,
		67, 60, 30,
		30, 90, 30,
		30, 90, 30,
		67, 60, 30,
		67, 90, 30,

		// top
		0, 0, 0,
		100, 0, 0,
		100, 0, 30,
		0, 0, 0,
		100, 0, 30,
		0, 0, 30,

		// top rung right
		100, 0, 0,
		100, 30, 0,
		100, 30, 30,
		100, 0, 0,
		100, 30, 30,
		100, 0, 30,

		// under top rung
		30, 30, 0,
		30, 30, 30,
		100, 30, 30,
		30, 30, 0,
		100, 30, 30,
		100, 30, 0,

		// between top rung and middle
		30, 30, 0,
		30, 60, 30,
		30, 30, 30,
		30, 30, 0,
		30, 60, 0,
		30, 60, 30,

		// top of middle rung
		30, 60, 0,
		67, 60, 30,
		30, 60, 30,
		30, 60, 0,
		67, 60, 0,
		67, 60, 30,

		// right of middle rung
		67, 60, 0,
		67, 90, 30,
		67, 60, 30,
		67, 60, 0,
		67, 90, 0,
		67, 90, 30,

		// bottom of middle rung.
		30, 90, 0,
		30, 90, 30,
		67, 90, 30,
		30, 90, 0,
		67, 90, 30,
		67, 90, 0,

		// right of bottom
		30, 90, 0,
		30, 150, 30,
		30, 90, 30,
		30, 90, 0,
		30, 150, 0,
		30, 150, 30,

		// bottom
		0, 150, 0,
		0, 150, 30,
		30, 150, 30,
		0, 150, 0,
		30, 150, 30,
		30, 150, 0,

		// left side
		0, 0, 0,
		0, 0, 30,
		0, 150, 30,
		0, 0, 0,
		0, 150, 30,
		0, 150, 0,
	}, n.GlStaticDraw)
}

//setColors Fill the buffer with colors for the 'F'.
func (app *AxisTestApp) setColors() {
	n.GL.BufferData(n.GlArrayBuffer, []uint8{
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,

		// top rung front
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,

		// middle rung front
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,
		200, 70, 120,

		// left column back
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,

		// top rung back
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,

		// middle rung back
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,
		80, 70, 200,

		// top
		70, 200, 210,
		70, 200, 210,
		70, 200, 210,
		70, 200, 210,
		70, 200, 210,
		70, 200, 210,

		// top rung right
		200, 200, 70,
		200, 200, 70,
		200, 200, 70,
		200, 200, 70,
		200, 200, 70,
		200, 200, 70,

		// under top rung
		210, 100, 70,
		210, 100, 70,
		210, 100, 70,
		210, 100, 70,
		210, 100, 70,
		210, 100, 70,

		// between top rung and middle
		210, 160, 70,
		210, 160, 70,
		210, 160, 70,
		210, 160, 70,
		210, 160, 70,
		210, 160, 70,

		// top of middle rung
		70, 180, 210,
		70, 180, 210,
		70, 180, 210,
		70, 180, 210,
		70, 180, 210,
		70, 180, 210,

		// right of middle rung
		100, 70, 210,
		100, 70, 210,
		100, 70, 210,
		100, 70, 210,
		100, 70, 210,
		100, 70, 210,

		// bottom of middle rung.
		76, 210, 100,
		76, 210, 100,
		76, 210, 100,
		76, 210, 100,
		76, 210, 100,
		76, 210, 100,

		// right of bottom
		140, 210, 80,
		140, 210, 80,
		140, 210, 80,
		140, 210, 80,
		140, 210, 80,
		140, 210, 80,

		// bottom
		90, 130, 110,
		90, 130, 110,
		90, 130, 110,
		90, 130, 110,
		90, 130, 110,
		90, 130, 110,

		// left side
		160, 160, 220,
		160, 160, 220,
		160, 160, 220,
		160, 160, 220,
		160, 160, 220,
		160, 160, 220,
	}, n.GlStaticDraw)
}

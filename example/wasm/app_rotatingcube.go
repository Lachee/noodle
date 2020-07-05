package main

import (
	"image"
	"image/color"
	"log"

	n "github.com/lachee/noodle"
)

//RotatingCubeApp shows off the 3D capabilities
type RotatingCubeApp struct {
	vertexBuffer n.WebGLBuffer
	colorBuffer  n.WebGLBuffer
	indexBuffer  n.WebGLBuffer
	uvBuffer     n.WebGLBuffer

	shader *n.Shader

	uProjMatrixLoc  n.WebGLUniformLocation
	uViewMatrixLoc  n.WebGLUniformLocation
	uModelMatrixLoc n.WebGLUniformLocation
	uSamplerLoc     n.WebGLUniformLocation

	projMatrix  Matrix
	viewMatrix  Matrix
	modelMatrix Matrix
	moveMatrix  Matrix

	rotation float32
	texture  *n.Texture
}

func (app *RotatingCubeApp) prepareImage() (*n.Image, error) {

	//Size of the image
	const width = 255
	const height = 255

	//Create the image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			color := color.RGBA{uint8(x), uint8(y), 128, 0xff}
			img.Set(x, y, color)
		}
	}

	return n.LoadImageRGBA(img)
	//return n.LoadImage("resources/moomin.png") // The image URL
}

//Start is called by the noodle engine when ready
func (app *RotatingCubeApp) Start() bool {

	// Create vertex buffer
	app.vertexBuffer = n.GL.NewBuffer(n.GlArrayBuffer, rotCubeVerts, n.GlStaticDraw)
	app.colorBuffer = n.GL.NewBuffer(n.GlArrayBuffer, rotCubeColours, n.GlStaticDraw)
	app.indexBuffer = n.GL.NewBuffer(n.GlElementArrayBuffer, rotCubeTris, n.GlStaticDraw)
	app.uvBuffer = n.GL.NewBuffer(n.GlArrayBuffer, rotCubeUV, n.GlStaticDraw)

	// == Load the cube image and the shaders
	image, err := app.prepareImage()
	if err != nil {
		log.Fatalln("Failed to load image", err)
		return false
	}
	app.texture = n.NewTexture(image)

	//Load the cube shader
	shader, err := n.LoadShaderFromURL("resources/shader/rotatingCube.vert", "resources/shader/rotatingCube.frag")
	if err != nil {
		log.Fatalln("Failed to load the shaders! ", err)
		return false
	}
	app.shader = shader

	// == Link the shaders up
	// Associate attributes to vertex shader
	app.uProjMatrixLoc = app.shader.GetUniformLocation("Pmatrix")
	app.uViewMatrixLoc = app.shader.GetUniformLocation("Vmatrix")
	app.uModelMatrixLoc = app.shader.GetUniformLocation("Mmatrix")
	app.uSamplerLoc = app.shader.GetUniformLocation("uSampler")

	app.shader.BindVertexData("position", n.GlArrayBuffer, app.vertexBuffer, 3, n.GlFloat, false, 0, 0)
	app.shader.BindVertexData("color", n.GlArrayBuffer, app.colorBuffer, 3, n.GlFloat, false, 0, 0)
	app.shader.BindVertexData("textureCoord", n.GlArrayBuffer, app.uvBuffer, 2, n.GlFloat, false, 0, 0)
	app.shader.Use()

	// == Set WeebGL properties
	//n.GL.ClearColor(0.5, 0.5, 0.5, 0.9)
	n.GL.ClearDepth(1)
	n.GL.Viewport(0, 0, n.GL.Width(), n.GL.Height())
	n.GL.DepthFunc(n.GlLEqual)

	// == Create Matrixes
	// Generate and apply projection matrix
	//app.projMatrix = n.NewMatrixPerspective(45.0, n.GL.AspectRatio(), 1, 100.0)
	//n.GL.UniformMatrix4fv(app.uProjMatrixLoc, app.projMatrix)

	// Generate and apply view matrix
	app.viewMatrix = n.NewMatrixLookAt(Vector3{3.0, 3.0, 3.0}, Vector3{0, 0, 0}, Vector3{0, 1, 0})
	n.GL.UniformMatrix4fv(app.uViewMatrixLoc, app.viewMatrix)

	//Update the texture shit
	// Activate the text0, tell the texture to bind, then tell the
	//  sampler that it should be on 0
	app.texture.SetSampler(app.uSamplerLoc, 0)
	return true
}

//Update occurs once a frame
func (app *RotatingCubeApp) Update(dt float32) {
	app.rotation = app.rotation + dt/1
	app.rotation = 1

	app.moveMatrix = n.NewMatrixTranslate(Vector3{0, 0, 0})
	/*

		//Update the move matrix
		axis := Vector3{0, 1, 1}
		angle := float32(0.5)

			var movMatrix Matrix
			rota := n.NewQuaternionAxis(axis, angle)
			rotb := n.NewQuaternionAxis(n.Vector3{0, 0, 1}, 0.3*app.rotation)
			rotc := n.NewQuaternionAxis(n.Vector3{1, 0, 0}, 0.2*app.rotation)

			movMatrix = n.NewMatrixRotation(rota)
			movMatrix = movMatrix.Multiply(n.NewMatrixRotation(rotb))
			movMatrix = movMatrix.Multiply(n.NewMatrixRotation(rotc))

			app.moveMatrix = movMatrix
	*/

}

//Render occurs when the screen needs updating
func (app *RotatingCubeApp) Render() {

	//Set the clear colour
	n.GL.ClearColor(n.White)

	//Force the projection matrix to update this frame.
	app.projMatrix = n.NewMatrixPerspective(45.0, n.GL.AspectRatio(), 1, 100.0)
	n.GL.UniformMatrix4fv(app.uProjMatrixLoc, app.projMatrix)

	//Update matrix
	n.GL.UniformMatrix4fv(app.uModelMatrixLoc, app.moveMatrix)

	//Clear
	n.GL.BindBuffer(n.GlElementArrayBuffer, app.indexBuffer)
	n.GL.Enable(n.GlDepthTest)
	n.GL.Clear(n.GlColorBufferBit | n.GlDepthBufferBit)

	// Draw the cube
	n.GL.DrawElements(n.GlTriangles, len(rotCubeTris), n.GlUnsignedShort, 0)
}

var rotCubeVerts = []Vector3{
	Vector3{-1, -1, -1}, Vector3{1, -1, -1}, Vector3{1, 1, -1}, Vector3{-1, 1, -1},
	Vector3{-1, -1, 1}, Vector3{1, -1, 1}, Vector3{1, 1, 1}, Vector3{-1, 1, 1},
	Vector3{-1, -1, -1}, Vector3{-1, 1, -1}, Vector3{-1, 1, 1}, Vector3{-1, -1, 1},
	Vector3{1, -1, -1}, Vector3{1, 1, -1}, Vector3{1, 1, 1}, Vector3{1, -1, 1},
	Vector3{-1, -1, -1}, Vector3{-1, -1, 1}, Vector3{1, -1, 1}, Vector3{1, -1, -1},
	Vector3{-1, 1, -1}, Vector3{-1, 1, 1}, Vector3{1, 1, 1}, Vector3{1, 1, -1},
}
var rotCubeUV = []Vector2{
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
}
var rotCubeColours = []float32{
	5, 3, 7, 5, 3, 7, 5, 3, 7, 5, 3, 7,
	1, 1, 3, 1, 1, 3, 1, 1, 3, 1, 1, 3,
	0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1,
	1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0,
	1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0,
	0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0,
}
var rotCubeTris = []uint16{
	0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7,
	8, 9, 10, 8, 10, 11, 12, 13, 14, 12, 14, 15,
	16, 17, 18, 16, 18, 19, 20, 21, 22, 20, 22, 23,
}

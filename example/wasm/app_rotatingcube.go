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

	shader      *n.Shader
	uMatrixLoc  n.WebGLUniformLocation
	uSamplerLoc n.WebGLUniformLocation

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

	//return n.LoadImageRGBA(img)
	return n.LoadImage("resources/direction.png") // The image URL
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
	shader, err := n.LoadShaderFromURL("resources/shader/flat.vert", "resources/shader/flat.frag")
	if err != nil {
		log.Fatalln("Failed to load the shaders! ", err)
		return false
	}
	app.shader = shader

	// == Link the shaders up
	// Associate attributes to vertex shader
	app.uMatrixLoc = app.shader.GetUniformLocation("n_matrix")
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

	//Update the texture shit
	// Activate the text0, tell the texture to bind, then tell the
	//  sampler that it should be on 0
	app.texture.SetSampler(app.uSamplerLoc, 0)
	return true
}

//Update occurs once a frame
func (app *RotatingCubeApp) Update(dt float32) {
}

//Render occurs when the screen needs updating
func (app *RotatingCubeApp) Render() {

	const radius = float32(10)

	//Set the clear colour
	n.GL.ClearColor(n.White)

	//Prepare the projection
	projectionMatrix := n.NewMatrixPerspective(90, n.GL.AspectRatio(), 1, 2000.0)

	//PRepare the camera
	//cameraAngleRadians := float32(n.GetFrameTime()) * n.PI * 0.5
	//cameraMatrix := n.NewMatrixRotationY(cameraAngleRadians)
	//cameraMatrix = cameraMatrix.Translate(Vector3{0, 0, radius * 1.5})

	//cameraMatrix := n.NewMatrixTranslate(Vector3{0, 0, radius})

	cameraMatrix := n.NewMatrixTranslate(Vector3{0, 0, -3})

	//Create the new matrix
	viewMatrix := cameraMatrix //cameraMatrix.Inverse()
	viewProjectionMatrix := projectionMatrix.Multiply(viewMatrix)

	//Create the model matrix
	modelMatrix := n.NewMatrixTranslate(Vector3{0, 0, 0})
	modelMatrix = modelMatrix.RotateX((float32)(n.GetFrameTime()))
	modelMatrix = modelMatrix.RotateY((float32)(n.GetFrameTime()))

	//Set the Unfiform
	viewProjectionModelMatrix := viewProjectionMatrix.Multiply(modelMatrix)
	n.GL.UniformMatrix4fv(app.uMatrixLoc, viewProjectionModelMatrix)

	//Clear
	n.GL.BindBuffer(n.GlElementArrayBuffer, app.indexBuffer)
	n.GL.Enable(n.GlDepthTest)
	n.GL.Clear(n.GlColorBufferBit | n.GlDepthBufferBit)

	// Draw the cube
	n.GL.DrawElements(n.GlTriangles, len(rotCubeTris), n.GlUnsignedShort, 0)
}

var rotCubeVerts = []Vector3{
	//Vector3{-1, -1, -1}, Vector3{1, -1, -1}, Vector3{1, 1, -1}, Vector3{-1, 1, -1}, // Back Face
	Vector3{-1, -1, 1}, Vector3{1, -1, 1}, Vector3{1, 1, 1}, Vector3{-1, 1, 1}, // Front Face
	//Vector3{-1, -1, -1}, Vector3{-1, 1, -1}, Vector3{-1, 1, 1}, Vector3{-1, -1, 1}, // Left Face
	Vector3{1, -1, -1}, Vector3{1, 1, -1}, Vector3{1, 1, 1}, Vector3{1, -1, 1}, // Right Face
	//Vector3{-1, -1, -1}, Vector3{-1, -1, 1}, Vector3{1, -1, 1}, Vector3{1, -1, -1}, // Bottom Face
	Vector3{-1, 1, -1}, Vector3{-1, 1, 1}, Vector3{1, 1, 1}, Vector3{1, 1, -1}, //Top Face
}
var rotCubeUV = []Vector2{
	//Vector2{0.75, 0.25}, Vector2{1, 0.25}, Vector2{1, 0.0}, Vector2{0.75, 0.0}, //Back
	Vector2{0.5, 0.25}, Vector2{0.75, 0.25}, Vector2{0.75, 0.0}, Vector2{0.5, 0.0}, //Front
	//Vector2{0.25, 0.25}, Vector2{0.5, 0.25}, Vector2{0.5, 0}, Vector2{0.25, 0}, //Left
	Vector2{0, 0.25}, Vector2{0.25, 0.25}, Vector2{0.25, 0}, Vector2{0, 0}, //Right
	//Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0}, // Bottom Face
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0}, // Top Face
}
var rotCubeColours = []float32{

	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,

	/*
		5, 3, 7, 5, 3, 7, 5, 3, 7, 5, 3, 7,
		1, 1, 3, 1, 1, 3, 1, 1, 3, 1, 1, 3,
		0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1,
		0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1,
		1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0,
		1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0,
	*/
}
var rotCubeTris = []uint16{
	0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7,
	8, 9, 10, 8, 10, 11, 12, 13, 14, 12, 14, 15,
	16, 17, 18, 16, 18, 19, 20, 21, 22, 20, 22, 23,
}

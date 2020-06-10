package main

import (
	"log"
	"math"

	n "github.com/lachee/noodle"
)

type NineSliceApp struct {
	vertexBuffer n.WebGLBuffer
	colorBuffer  n.WebGLBuffer
	indexBuffer  n.WebGLBuffer
	uvBuffer     n.WebGLBuffer

	shader *n.Shader

	uProjMatrixLoc  n.WebGLUniformLocation
	uViewMatrixLoc  n.WebGLUniformLocation
	uModelMatrixLoc n.WebGLUniformLocation
	uSamplerLoc     n.WebGLUniformLocation
	uDimensionLoc   n.WebGLUniformLocation
	uBorderLoc      n.WebGLUniformLocation

	projMatrix  Matrix
	viewMatrix  Matrix
	modelMatrix Matrix
	moveMatrix  Matrix

	texture       *n.Texture //The texture of the slice
	textureBorder int        //The size of the slice border (in pixels)

	scale     float32
	clip      Vector2 //Absolute clip space of the image (in pixels)
	border    Vector2 //Relative border of the texture
	dimension Vector2 //Translated size of the geometry
}

func (app *NineSliceApp) Start() bool {

	// Create vertex buffer
	app.vertexBuffer = n.GL.NewBuffer(n.GlArrayBuffer, rotCubeVerts, n.GlStaticDraw)
	app.colorBuffer = n.GL.NewBuffer(n.GlArrayBuffer, rotCubeColours, n.GlStaticDraw)
	app.indexBuffer = n.GL.NewBuffer(n.GlElementArrayBuffer, rotCubeTris, n.GlStaticDraw)
	app.uvBuffer = n.GL.NewBuffer(n.GlArrayBuffer, rotCubeUV, n.GlStaticDraw)

	// == Load the cube image and the shaders
	app.textureBorder = 16                             // Border Size of the image
	image, err := n.LoadImage("resources/tilefat.png") // The image URL
	if err != nil {
		log.Fatalln("Failed to load image", err)
		return false
	}
	app.texture = n.NewTexture(image)
	app.clip = n.NewVector2i(app.texture.Width(), app.texture.Height())
	app.border = n.NewVector2(float32(app.textureBorder)/app.clip.X, float32(app.textureBorder)/app.clip.Y)

	//Load the cube shader
	shader, err := n.LoadShaderFromURL("resources/shader/nineSlice.vert", "resources/shader/nineSlice.frag")
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
	app.uDimensionLoc = app.shader.GetUniformLocation("uDimensions")
	app.uBorderLoc = app.shader.GetUniformLocation("uBorder")

	app.shader.BindVertexData("position", n.GlArrayBuffer, app.vertexBuffer, 3, n.GlFloat, false, 0, 0)
	app.shader.BindVertexData("color", n.GlArrayBuffer, app.colorBuffer, 3, n.GlFloat, false, 0, 0)
	app.shader.BindVertexData("textureCoord", n.GlArrayBuffer, app.uvBuffer, 2, n.GlFloat, false, 0, 0)
	app.shader.Use()

	// == Set WeebGL properties
	n.GL.ClearColor(0.5, 0.5, 0.5, 0.9)
	n.GL.ClearDepth(1)
	n.GL.Viewport(0, 0, n.Width(), n.Height())
	n.GL.DepthFunc(n.GlLEqual)

	// == Create Matrixes
	// Generate and apply projection matrix
	app.projMatrix = n.NewMatrixPerspective(45.0, float64(n.Width())/float64(n.Height()), 1, 100.0)
	n.GL.UniformMatrix4fv(app.uProjMatrixLoc, app.projMatrix)

	// Generate and apply view matrix
	app.viewMatrix = n.NewMatrixLookAt(Vector3{0.0, 0, -10.0}, Vector3{0, 0, 0}, Vector3{0, 1, 0})
	n.GL.UniformMatrix4fv(app.uViewMatrixLoc, app.viewMatrix)

	//Update the texture shit
	// Activate the text0, tell the texture to bind, then tell the
	//  sampler that it should be on 0
	app.texture.SetSampler(app.uSamplerLoc, 0)
	return true
}

//Update occurs once a frame
func (app *NineSliceApp) Update(dt float32) {
	speed := 0.0015
	scaleX := (float32(math.Sin(n.GetFrameTime()*speed)/2) + 1) * 3
	scaleY := (float32(math.Cos(n.GetFrameTime()*speed)/2) + 1) * 3
	scale := Vector2{scaleX, scaleY}

	//Size of the geometry
	box := app.clip.Multiply(scale)
	app.dimension = n.NewVector2(float32(app.textureBorder)/box.X, float32(app.textureBorder)/box.Y)

	//Update the move matrix
	movMatrix := n.NewMatrixRotate(n.NewVector3Up(), n.PI*2)
	movMatrix = movMatrix.Multiply(n.NewMatrixScale(scale.ToVector3()))
	app.moveMatrix = movMatrix
}

//Render occurs when the screen needs updating
func (app *NineSliceApp) Render() {

	//Update matrix
	n.GL.UniformMatrix4fv(app.uModelMatrixLoc, app.moveMatrix)
	n.GL.Uniform2v(app.uBorderLoc, app.border)
	n.GL.Uniform2v(app.uDimensionLoc, app.dimension)

	//Clear
	n.GL.BindBuffer(n.GlElementArrayBuffer, app.indexBuffer)
	n.GL.Enable(n.GlDepthTest)
	n.GL.Clear(n.GlColorBufferBit | n.GlDepthBufferBit)

	// Draw the cube
	n.GL.DrawElements(n.GlTriangles, len(rotCubeTris), n.GlUnsignedShort, 0)
}

var ninePlaneVerts = []Vector3{
	Vector3{0, 0, 0}, Vector3{0, 1, 0},
	Vector3{1, 1, 0}, Vector3{1, 0, 0},
}
var ninePlaneUV = []Vector2{
	Vector2{0, 0}, Vector2{0, 1},
	Vector2{1, 1}, Vector2{1, 0},
}
var ninePlaneTris = []uint16{0, 1, 2, 0, 2, 3}

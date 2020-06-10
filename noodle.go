// +build js,wasm
package noodle

/* Here is a list of interesting Readings
https://blog.scottlogic.com/2019/11/18/drawing-lines-with-webgl.html - Line Rendering
https://github.com/davidwparker/programmingtil-webgl/tree/master/0016-single-line - More Lines
https://mattdesl.svbtle.com/drawing-lines-is-hard - Even more lines
https://github.com/mattdesl/webgl-lines/blob/master/expanded/frag.glsl - A line shader
https://stdiopt.github.io/gowasm-experiments/rainbow-mouse/ - Example Go WASM
https://www.gamedev.net/forums/topic/696879-glsl-9-slicing/ - 9 Slice
https://github.com/Lachee/engi/blob/master/SpriteRenderer.go - Engi SpriteRenderer Renderer
*/

import (
	"log"
	"math"
	"syscall/js"
)

var (
	//GL gives direct access to the WebGL component of the canvas.
	GL              *WebGL
	document        js.Value
	canvas          js.Value
	frameRenderFunc js.Func
	input           *Input
	app             Application
	width           int
	height          int
	texture         *Texture
	frameTime       float64
	deltaTime       float64
)

//GetFrameTime returns the time the last frame was rendered
func GetFrameTime() float64 { return frameTime }

//GetDeltaTime returns a high accuracy difference in time between the last frame and the current one.
func GetDeltaTime() float64 { return deltaTime }

//DT returns a less accurate version of GetDeltaTime, for all your 32bit mathmatic needs.
func DT() float32 { return float32(deltaTime) }

//GetInput returns the current input manager
func GetInput() *Input {
	return input
}

//Width gets the width of the screen
func Width() int {
	return canvas.Get("width").Int()
}

//Height gets the width of the screen
func Height() int {
	return canvas.Get("height").Int()
}

//Initialize sets up the Noodle renderer
func Initialize(application Application) {
	app = application

	input = newInput()
	document = js.Global().Get("document")
	canvas = document.Call("getElementById", "gocanvas")

	//Set the width and height of the canvas to conver the entire screen
	width = document.Get("body").Get("clientWidth").Int()
	height = document.Get("body").Get("clientHeight").Int()
	canvas.Set("width", width)
	canvas.Set("height", height)

	done := make(chan struct{}, 0)

	//Get the GL context
	context := canvas.Call("getContext", "webgl")
	if context.IsUndefined() {
		context = canvas.Call("getContext", "experimental-webgl")
	}
	if context.IsUndefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

	//Create a new GL instance
	GL = newWebGL(context)

	//log.Println("GlTexture", GlTexture0, GlTexture1, GlTexture2, GlTexture31, GlTexture30, GlTexture29)
	//log.Println("Literal32", 0x84C0, 0x84C1, 0x84C2, 0x84df, 0x84de, 0x84dd)

	//Define the texture
	//image, err := LoadImage("resources/moomin.png")
	//image, err := LoadImage("resources/firefox.svg")
	image, err := LoadImage("resources/tilefat.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
		return
	}

	texture = NewTexture(image)

	//Record canvas events
	onMouseChangeEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		rect := canvas.Call("getBoundingClientRect")
		x := evt.Get("clientX").Int() - rect.Get("left").Int()
		y := evt.Get("clientY").Int() - rect.Get("top").Int()
		input.setMousePosition(x, y)
		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mousemove", onMouseChangeEvent)

	//Record canvas events
	onMouseUpEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		input.setMouseUp(button)
		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mouseup", onMouseUpEvent)

	//Record canvas events
	onMouseDownEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		input.setMouseDown(button)
		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mousedown", onMouseDownEvent)

	//Setup the animation frame
	renderSetup()
	frameRenderFunc = js.FuncOf(onRequestAnimationFrame)
	defer frameRenderFunc.Release()
	js.Global().Call("requestAnimationFrame", frameRenderFunc)
	<-done
}

//RequestRedraw requests for a new animation frame
func RequestRedraw() {
	js.Global().Call("requestAnimationFrame", frameRenderFunc)
}

//onRequestAnimationFrame callback for animations
func onRequestAnimationFrame(this js.Value, args []js.Value) interface{} {
	//Setupt he time
	time := args[0].Float()
	deltaTime = time - frameTime
	frameTime = time

	input.update()
	if render() {
		RequestRedraw()
	}
	return nil
}

var (
	vertexBuffer WebGLBuffer
	colorBuffer  WebGLBuffer
	indexBuffer  WebGLBuffer
	uvBuffer     WebGLBuffer
	basicShader  *Shader

	uProjMatrixLoc  WebGLUniformLocation
	uViewMatrixLoc  WebGLUniformLocation
	uModelMatrixLoc WebGLUniformLocation
	uSamplerLoc     WebGLUniformLocation
	uDimensionLoc   WebGLUniformLocation
	uBorderLoc      WebGLUniformLocation

	projMatrix  Matrix
	viewMatrix  Matrix
	modelMatrix Matrix

	rotation float32
)

func renderSetup() {

	// Create vertex buffer
	//vertexBuffer := GL.NewBuffer(GlArrayBuffer, verticesNative, GlStaticDraw)
	vertexBuffer = GL.NewBuffer(GlArrayBuffer, planeVertices, GlStaticDraw)

	// Create color buffer
	colorBuffer = GL.NewBuffer(GlArrayBuffer, colorsNative, GlStaticDraw)

	// Create index buffer
	indexBuffer = GL.NewBuffer(GlElementArrayBuffer, planeIndices, GlStaticDraw)

	// Create uv buffer
	uvBuffer = GL.NewBuffer(GlArrayBuffer, planeUV, GlStaticDraw)

	//=== SHADER
	basicShader, err := LoadShader(newvertShaderCode, newfragShaderCode)
	if err != nil {
		log.Fatalln("Failed to load the shaders! ", err)
		return
	}

	// Associate attributes to vertex shader
	uProjMatrixLoc = basicShader.GetUniformLocation("Pmatrix")
	uViewMatrixLoc = basicShader.GetUniformLocation("Vmatrix")
	uModelMatrixLoc = basicShader.GetUniformLocation("Mmatrix")
	uSamplerLoc = basicShader.GetUniformLocation("uSampler")
	uDimensionLoc = basicShader.GetUniformLocation("u_dimensions")
	uBorderLoc = basicShader.GetUniformLocation("u_border")

	//Bind the data we have
	basicShader.BindVertexData("position", GlArrayBuffer, vertexBuffer, 3, GlFloat, false, 0, 0)
	basicShader.BindVertexData("color", GlArrayBuffer, colorBuffer, 3, GlFloat, false, 0, 0)
	basicShader.BindVertexData("textureCoord", GlArrayBuffer, uvBuffer, 2, GlFloat, false, 0, 0)
	basicShader.Use()

	// Set WeebGL properties
	GL.ClearColor(0.5, 0.5, 0.5, 0.9)
	GL.ClearDepth(1)
	GL.Viewport(0, 0, width, height)
	GL.DepthFunc(GlLEqual)

	//// Create Matrixes ////
	// Generate and apply projection matrix
	projMatrix = NewMatrixPerspective(45.0, float64(width)/float64(height), 1, 100.0)
	//projMatrix = NewMatrixOrtho(left, right, bottom, top, 1, 100.0)
	GL.UniformMatrix4fv(uProjMatrixLoc, projMatrix)

	// Generate and apply view matrix
	viewMatrix = NewMatrixLookAt(NewVector3(3.0, 0.0, 3.0), Vector3{0, 1.5, 0}, NewVector3Up())
	GL.UniformMatrix4fv(uViewMatrixLoc, viewMatrix)

	//Update the texture shit
	// Activate the text0, tell the texture to bind, then tell the
	//  sampler that it should be on 0
	GL.ActiveTexture(GlTexture0)
	texture.Bind()
	GL.Uniform1i(uSamplerLoc, 0)
	//texture.SetSampler(uSamplerLoc, 0)
}

//render outpouts the render
func render() bool {

	// Bind to element array for draw function
	GL.BindBuffer(GlElementArrayBuffer, indexBuffer)

	// Calculate rotation rate
	rotation = rotation + DT()/500
	diff := float32(math.Sin(frameTime*0.005)/2) + 1
	scale := diff * 3

	clip := NewVector2i(texture.Width(), texture.Height())
	box := NewVector2(scale*clip.X, scale*clip.Y)

	var slice float32 = 15.0
	border := NewVector2(slice/clip.X, slice/clip.Y)
	dimension := NewVector2(slice/box.X, slice/box.Y)

	GL.Uniform2f(uBorderLoc, border.X, border.Y)
	GL.Uniform2f(uDimensionLoc, dimension.X, dimension.Y)

	// Do new model matrix calculations
	movMatrix := NewMatrixRotate(NewVector3Up(), 0.5)
	movMatrix = movMatrix.Multiply(NewMatrixScale(box.Scale(0.01).ToVector3()))
	//movMatrix = movMatrix.Multiply(NewMatrixRotate(NewVector3Forward(), 0.3*rotation))
	//movMatrix = movMatrix.Multiply(NewMatrixRotate(NewVector3Right(), 0.2*rotation))
	GL.UniformMatrix4fv(uModelMatrixLoc, movMatrix)

	// Clear the screen
	GL.Enable(GlDepthTest)
	GL.Clear(GlColorBufferBit | GlDepthBufferBit)

	// Draw the cube
	GL.DrawElements(GlTriangles, len(planeIndices), GlUnsignedShort, 0)

	return true
}

var planeVertices = []Vector3{
	Vector3{0, 0, 0}, Vector3{0, 1, 0},
	Vector3{1, 1, 0}, Vector3{1, 0, 0},
}
var planeUV = []Vector2{
	Vector2{0, 0}, Vector2{0, 1},
	Vector2{1, 1}, Vector2{1, 0},
}
var planeIndices = []uint16{0, 1, 2, 0, 2, 3}

//// BUFFERS + SHADERS ////
// Shamelessly copied from https://www.tutorialspoint.com/webgl/webgl_cube_rotation.htm //
var verticesNativeV = []Vector3{
	Vector3{-1, -1, -1}, Vector3{1, -1, -1}, Vector3{1, 1, -1}, Vector3{-1, 1, -1},
	Vector3{-1, -1, 1}, Vector3{1, -1, 1}, Vector3{1, 1, 1}, Vector3{-1, 1, 1},
	Vector3{-1, -1, -1}, Vector3{-1, 1, -1}, Vector3{-1, 1, 1}, Vector3{-1, -1, 1},
	Vector3{1, -1, -1}, Vector3{1, 1, -1}, Vector3{1, 1, 1}, Vector3{1, -1, 1},
	Vector3{-1, -1, -1}, Vector3{-1, -1, 1}, Vector3{1, -1, 1}, Vector3{1, -1, -1},
	Vector3{-1, 1, -1}, Vector3{-1, 1, 1}, Vector3{1, 1, 1}, Vector3{1, 1, -1},
}
var uvNativeV = []Vector2{
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
	Vector2{0.0, 0.0}, Vector2{1.0, 0.0}, Vector2{1.0, 1.0}, Vector2{0.0, 1.0},
}
var colorsNative = []float32{
	5, 3, 7, 5, 3, 7, 5, 3, 7, 5, 3, 7,
	1, 1, 3, 1, 1, 3, 1, 1, 3, 1, 1, 3,
	0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1,
	1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0,
	1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0,
	0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0,
}
var indicesNative = []uint16{
	0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7,
	8, 9, 10, 8, 10, 11, 12, 13, 14, 12, 14, 15,
	16, 17, 18, 16, 18, 19, 20, 21, 22, 20, 22, 23,
}

const vertShaderCode = `
attribute vec3 position;
attribute vec2 textureCoord;

uniform mat4 Pmatrix;
uniform mat4 Vmatrix;
uniform mat4 Mmatrix;

attribute vec3 color;

varying vec3 vColor;
varying highp vec2 vTextureCoord;

void main(void) {
	gl_Position = Pmatrix * Vmatrix * Mmatrix * vec4(position, 1.);
	vColor = color;
	vTextureCoord = textureCoord;
}
`
const fragShaderCode = `
precision mediump float;
varying vec3 vColor;

varying highp vec2 vTextureCoord;
uniform sampler2D uSampler;

void main(void) {
	gl_FragColor = vec4(vColor, 1.) * texture2D(uSampler, vTextureCoord);
}
`

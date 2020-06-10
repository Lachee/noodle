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
	"syscall/js"
)

var (
	GL              *WebGL
	document        js.Value
	canvas          js.Value
	frameRenderFunc js.Func
	input           *Input
	app             Application
	width           int
	height          int
	texture         *Texture
)

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

	//Define the texture
	image, err := LoadImage("resources/moomin.png")
	//image, err := LoadImage("resources/firefox.svg")
	//image, err := LoadImage("resources/tile.png")
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

	//=== BUFFER
	// Create vertex buffer
	//vertexBuffer := GL.NewBuffer(GlArrayBuffer, verticesNative, GlStaticDraw)
	vertexBuffer := GL.NewBuffer(GlArrayBuffer, verticesNativeV, GlStaticDraw)

	// Create color buffer
	colorBuffer := GL.NewBuffer(GlArrayBuffer, colorsNative, GlStaticDraw)

	// Create index buffer
	indexBuffer := GL.NewBuffer(GlElementArrayBuffer, indicesNative, GlStaticDraw)

	// Create uv buffer
	uvBuffer := GL.NewBuffer(GlArrayBuffer, uvNativeV, GlStaticDraw)

	//=== SHADER
	basicShader := LoadShader(vertShaderCode, fragShaderCode)

	// Associate attributes to vertex shader
	PositionMatrix := basicShader.GetUniformLocation("Pmatrix")
	ViewMatrix := basicShader.GetUniformLocation("Vmatrix")
	ModelMatrix := basicShader.GetUniformLocation("Mmatrix")
	Sampler := basicShader.GetUniformLocation("uSampler")

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
	ratio := float64(width) / float64(height)

	// Generate and apply projection matrix
	projMatrix := NewMatrixPerspective(45.0, ratio, 1, 100.0)
	GL.UniformMatrix4fv(PositionMatrix, projMatrix)

	// Generate and apply view matrix
	viewMatrix := NewMatrixLookAt(NewVector3(3.0, 3.0, 3.0), NewVector3Zero(), NewVector3Up())
	GL.UniformMatrix4fv(ViewMatrix, viewMatrix)

	//// Drawing the Cube ////
	movMatrix := NewMatrixIdentity()
	var tmark float32
	var rotation = float32(0)

	//Update the texture shit
	// Activate the text0, tell the texture to bind, then tell the
	//  sampler that it should be on 0
	GL.ActiveTexture(GlTexture0)
	texture.Bind()
	GL.Uniform1i(Sampler, 0)

	// Bind to element array for draw function
	GL.BindBuffer(GlElementArrayBuffer, indexBuffer)
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		// Calculate rotation rate
		now := float32(args[0].Float())
		tdiff := now - tmark
		tmark = now
		rotation = rotation + float32(tdiff)/500

		// Do new model matrix calculations
		movMatrix = NewMatrixRotate(NewVector3Up(), 0.5*rotation)
		movMatrix = movMatrix.Multiply(NewMatrixRotate(NewVector3Forward(), 0.3*rotation))
		movMatrix = movMatrix.Multiply(NewMatrixRotate(NewVector3Right(), 0.2*rotation))
		GL.UniformMatrix4fv(ModelMatrix, movMatrix)

		// Clear the screen
		GL.Enable(GlDepthTest)
		GL.Clear(GlColorBufferBit | GlDepthBufferBit)

		// Draw the cube
		GL.DrawElements(GlTriangles, len(indicesNative), GlUnsignedShort, 0)

		// Call next frame
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)
	//*/
	/*
		//Setup the animation frame
		renderSetup()
		frameRenderFunc = js.FuncOf(onRequestAnimationFrame)
		defer frameRenderFunc.Release()
		js.Global().Call("requestAnimationFrame", frameRenderFunc)
	*/
	<-done
}

//RequestRedraw requests for a new animation frame
func RequestRedraw() {
	js.Global().Call("requestAnimationFrame", frameRenderFunc)
}

//onRequestAnimationFrame callback for animations
func onRequestAnimationFrame(this js.Value, args []js.Value) interface{} {
	input.update()
	if render() {
		RequestRedraw()
	}
	return nil
}

func renderSetup() {

}

//render outpouts the render
func render() bool {
	return true
}

var planeVertices = []Vector3{}

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

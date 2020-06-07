// +build js,wasm
package noodle

/* Here is a list of interesting Readings
https://blog.scottlogic.com/2019/11/18/drawing-lines-with-webgl.html - Line Rendering
https://github.com/davidwparker/programmingtil-webgl/tree/master/0016-single-line - More Lines
https://mattdesl.svbtle.com/drawing-lines-is-hard - Even more lines
https://github.com/mattdesl/webgl-lines/blob/master/expanded/frag.glsl - A line shader
https://stdiopt.github.io/gowasm-experiments/rainbow-mouse/ - Example Go WASM
https://www.gamedev.net/forums/topic/696879-glsl-9-slicing/ - 9 Slice
https://github.com/Lachee/engi/blob/master/batch.go - Engi Batch Renderer
*/

import (
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
)

//GetInput returns the current input manager
func GetInput() *Input {
	return input
}

//Initialize sets up the Noodle renderer
func Initialize(application Application) {
	app = application

	input = newInput()
	document = js.Global().Get("document")
	canvas = document.Call("getElementById", "gocanvas")
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
	vertexBuffer := GL.NewBuffer(GlArrayBuffer, verticesNative, GlStaticDraw)

	// Create color buffer
	colorBuffer := GL.NewBuffer(GlArrayBuffer, colorsNative, GlStaticDraw)

	// Create index buffer
	indexBuffer := GL.NewBuffer(GlElementArrayBuffer, indicesNative, GlStaticDraw)

	//=== SHADER
	basicShader := LoadShader(vertShaderCode, fragShaderCode)

	/*
		// Create a vertex shader object
		vertShader := GL.NewShader(GlVertexShader, vertShaderCode)

		// Create fragment shader object
		fragShader := GL.NewShader(GlFragmentShader, fragShaderCode)
	*/

	// Create a shader program object to store
	// the combined shader program
	//shaderProgram := GL.NewProgram([]WebGLShader{vertShader, fragShader})

	// Associate attributes to vertex shader
	//PositionMatrix := GL.GetUniformLocation(shaderProgram, "Pmatrix")
	//ViewMatrix := GL.GetUniformLocation(shaderProgram, "Vmatrix")
	//ModelMatrix := GL.GetUniformLocation(shaderProgram, "Mmatrix")
	PositionMatrix := basicShader.GetUniformLocation("Pmatrix")
	ViewMatrix := basicShader.GetUniformLocation("Vmatrix")
	ModelMatrix := basicShader.GetUniformLocation("Mmatrix")

	//shaderProgram := basicShader.GetProgram()

	//GL.Call("bindBuffer", glTypes.ArrayBuffer, vertexBuffer)
	basicShader.BindVertexData("position", GlArrayBuffer, vertexBuffer, 3, GlFloat, false, 0, 0)
	basicShader.BindVertexData("color", GlArrayBuffer, colorBuffer, 3, GlFloat, false, 0, 0)
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
	var renderFrame js.Func
	var tmark float32
	var rotation = float32(0)

	// Bind to element array for draw function
	GL.BindBuffer(GlElementArrayBuffer, indexBuffer)

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
	//Setup the animation frame
	//frameRenderFunc = js.FuncOf(onRequestAnimationFrame)
	//defer frameRenderFunc.Release()
	//js.Global().Call("requestAnimationFrame", frameRenderFunc)

	<-done
}

func onRequestAnimationFrame(this js.Value, args []js.Value) interface{} {
	input.update()

	if input.GetButtonDown(0) {
		println("Button Was Down!")
	}

	if input.GetButtonUp(0) {
		println("Button Was Up!", input.GetMouseX(), input.GetMouseY())
	}

	js.Global().Call("requestAnimationFrame", frameRenderFunc)
	return nil
}

//// BUFFERS + SHADERS ////
// Shamelessly copied from https://www.tutorialspoint.com/webgl/webgl_cube_rotation.htm //
var verticesNative = []float32{
	-1, -1, -1, 1, -1, -1, 1, 1, -1, -1, 1, -1,
	-1, -1, 1, 1, -1, 1, 1, 1, 1, -1, 1, 1,
	-1, -1, -1, -1, 1, -1, -1, 1, 1, -1, -1, 1,
	1, -1, -1, 1, 1, -1, 1, 1, 1, 1, -1, 1,
	-1, -1, -1, -1, -1, 1, 1, -1, 1, 1, -1, -1,
	-1, 1, -1, -1, 1, 1, 1, 1, 1, 1, 1, -1,
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
uniform mat4 Pmatrix;
uniform mat4 Vmatrix;
uniform mat4 Mmatrix;
attribute vec3 color;
varying vec3 vColor;
void main(void) {
	gl_Position = Pmatrix*Vmatrix*Mmatrix*vec4(position, 1.);
	vColor = color;
}
`
const fragShaderCode = `
precision mediump float;
varying vec3 vColor;
void main(void) {
	gl_FragColor = vec4(vColor, 1.);
}
`

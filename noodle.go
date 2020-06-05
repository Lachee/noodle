// +build js,wasm
package noodle

import (
	"syscall/js"
	"unsafe"
)

var (
	document        js.Value
	canvas          js.Value
	gl              js.Value
	glTypes         GLTypes
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
	gl = canvas.Call("getContext", "webgl")
	if gl.IsUndefined() {
		gl = canvas.Call("getContext", "experimental-webgl")
	}
	if gl.IsUndefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

	//Prepare the GL Types
	glTypes.new(gl)

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

	// Convert buffers to JS TypedArrays
	var colors = sliceToTypedArray(colorsNative)
	var vertices = sliceToTypedArray(verticesNative)
	var indices = sliceToTypedArray(indicesNative)

	//=== BUFFER
	// Create vertex buffer
	vertexBuffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", glTypes.ArrayBuffer, vertexBuffer)
	gl.Call("bufferData", glTypes.ArrayBuffer, vertices, glTypes.StaticDraw)

	// Create color buffer
	colorBuffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", glTypes.ArrayBuffer, colorBuffer)
	gl.Call("bufferData", glTypes.ArrayBuffer, colors, glTypes.StaticDraw)

	// Create index buffer
	indexBuffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", glTypes.ElementArrayBuffer, indexBuffer)
	gl.Call("bufferData", glTypes.ElementArrayBuffer, indices, glTypes.StaticDraw)

	//=== SHADER
	// Create a vertex shader object
	vertShader := gl.Call("createShader", glTypes.VertexShader)
	gl.Call("shaderSource", vertShader, vertShaderCode)
	gl.Call("compileShader", vertShader)

	// Create fragment shader object
	fragShader := gl.Call("createShader", glTypes.FragmentShader)
	gl.Call("shaderSource", fragShader, fragShaderCode)
	gl.Call("compileShader", fragShader)

	// Create a shader program object to store
	// the combined shader program
	shaderProgram := gl.Call("createProgram")
	gl.Call("attachShader", shaderProgram, vertShader)
	gl.Call("attachShader", shaderProgram, fragShader)
	gl.Call("linkProgram", shaderProgram)

	// Associate attributes to vertex shader
	PositionMatrix := gl.Call("getUniformLocation", shaderProgram, "Pmatrix")
	ViewMatrix := gl.Call("getUniformLocation", shaderProgram, "Vmatrix")
	ModelMatrix := gl.Call("getUniformLocation", shaderProgram, "Mmatrix")

	gl.Call("bindBuffer", glTypes.ArrayBuffer, vertexBuffer)
	position := gl.Call("getAttribLocation", shaderProgram, "position")
	gl.Call("vertexAttribPointer", position, 3, glTypes.Float, false, 0, 0)
	gl.Call("enableVertexAttribArray", position)

	gl.Call("bindBuffer", glTypes.ArrayBuffer, colorBuffer)
	color := gl.Call("getAttribLocation", shaderProgram, "color")
	gl.Call("vertexAttribPointer", color, 3, glTypes.Float, false, 0, 0)
	gl.Call("enableVertexAttribArray", color)

	gl.Call("useProgram", shaderProgram)

	// Set WeebGL properties
	gl.Call("clearColor", 0.5, 0.5, 0.5, 0.9) // Color the screen is cleared to
	gl.Call("clearDepth", 1.0)                // Z value that is set to the Depth buffer every frame
	gl.Call("viewport", 0, 0, width, height)  // Viewport size
	gl.Call("depthFunc", glTypes.LEqual)

	//// Create Matrixes ////
	ratio := float64(width) / float64(height)

	// Generate and apply projection matrix
	//projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(ratio), 1, 100.0)
	projMatrix := NewMatrixPerspective(45.0, ratio, 1, 100.0)
	var projMatrixBuffer *[16]float32
	//projMatrixBuffer = projMatrix.ToBuffer()
	projMatrixBuffer = (*[16]float32)(unsafe.Pointer(&projMatrix))
	typedProjMatrixBuffer := sliceToTypedArray([]float32((*projMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", PositionMatrix, false, typedProjMatrixBuffer)

	// Generate and apply view matrix
	viewMatrix := NewMatrixLookAt(NewVector3(3.0, 3.0, 3.0), NewVector3Zero(), NewVector3Up())
	var viewMatrixBuffer *[16]float32
	viewMatrixBuffer = viewMatrix.DecomposePointer()
	typedViewMatrixBuffer := sliceToTypedArray([]float32((*viewMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", ViewMatrix, false, typedViewMatrixBuffer)

	//// Drawing the Cube ////
	movMatrix := NewMatrixIdentity() // mgl32.Ident4()
	//movMatrix := mgl32.Ident4()
	var renderFrame js.Func
	var tmark float32
	var rotation = float32(0)

	// Bind to element array for draw function
	gl.Call("bindBuffer", glTypes.ElementArrayBuffer, indexBuffer)

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Calculate rotation rate
		now := float32(args[0].Float())
		tdiff := now - tmark
		tmark = now
		rotation = rotation + float32(tdiff)/500

		var modelMatrixBuffer *[16]float32

		// Do new model matrix calculations
		movMatrix = NewMatrixRotate(NewVector3Up(), 0.5*rotation)
		movMatrix = movMatrix.Multiply(NewMatrixRotate(NewVector3Forward(), 0.3*rotation))
		movMatrix = movMatrix.Multiply(NewMatrixRotate(NewVector3Right(), 0.2*rotation))
		modelMatrixBuffer = movMatrix.DecomposePointer()

		//movMatrix = mgl32.HomogRotate3DX(0.5 * rotation)
		//movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DY(0.3 * rotation))
		//movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DZ(0.2 * rotation))
		//modelMatrixBuffer = (*[16]float32)(unsafe.Pointer(&movMatrix))

		typedModelMatrixBuffer := sliceToTypedArray([]float32((*modelMatrixBuffer)[:]))

		// Apply the model matrix
		gl.Call("uniformMatrix4fv", ModelMatrix, false, typedModelMatrixBuffer)

		// Clear the screen
		gl.Call("enable", glTypes.DepthTest)
		gl.Call("clear", glTypes.ColorBufferBit)
		gl.Call("clear", glTypes.DepthBufferBit)

		// Draw the cube
		gl.Call("drawElements", glTypes.Triangles, len(indicesNative), glTypes.UnsignedShort, 0)

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

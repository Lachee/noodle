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
	//GL gives direct access to the WebGL component of the canvas.
	GL *WebGL

	document        js.Value
	canvas          js.Value
	frameRenderFunc js.Func
	m_input         *InputHandler
	app             Application
	width           int
	height          int
	texture         *Texture
	frameTime       float64
	deltaTime       float64

	//AlwaysDraw continously draws
	AlwaysDraw = true
)

//GetFrameTime returns the time the last frame was rendered
func GetFrameTime() float64 { return frameTime }

//GetDeltaTime returns a high accuracy difference in time between the last frame and the current one.
func GetDeltaTime() float64 { return deltaTime }

//DT returns a less accurate version of GetDeltaTime, for all your 32bit mathmatic needs.
func DT() float32 { return float32(deltaTime) }

//Input returns the current input handler
func Input() *InputHandler {
	return m_input
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

	m_input = newInput()
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

	//Setup the animation frame
	if !app.Start() {
		log.Println("Failed to start the application")
		return
	}

	//Record canvas events
	onMouseChangeEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		rect := canvas.Call("getBoundingClientRect")
		x := evt.Get("clientX").Int() - rect.Get("left").Int()
		y := evt.Get("clientY").Int() - rect.Get("top").Int()
		m_input.setMousePosition(x, y)
		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mousemove", onMouseChangeEvent)

	//Record canvas events
	onMouseUpEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		m_input.setMouseUp(button)
		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mouseup", onMouseUpEvent)

	//Record canvas events
	onMouseDownEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		m_input.setMouseDown(button)
		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mousedown", onMouseDownEvent)

	//Begin rendering
	GL.Viewport(0, 0, Width(), Height())
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

	//Update the input
	m_input.update()

	//Call update on the Application
	app.Update(float32(deltaTime))

	//Render everything
	app.Render()

	//If we need to draw again, then do so
	if AlwaysDraw {
		RequestRedraw()
	}

	//Return nil to JS
	return nil
}

//Package noodle is a WebGL game engine, designed for low level access for fast and efficent 3D applications
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
	inputHandler         *InputHandler
	app             Application

	width  int
	height int

	texture   *Texture
	frameTime float64
	deltaTime float64

	awaiter chan int

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
	return inputHandler
}

//Width gets the width of the screen
func Width() int {
	return width
}

//Height gets the width of the screen
func Height() int {
	return height
}

//Run setups the WebGL context and runs the application. It is blocking and returns an exit code if Exit() is ever called.
func Run(application Application) int {
	app = application

	inputHandler = newInput()
	document = js.Global().Get("document")
	canvas = document.Call("getElementById", "gocanvas")

	//Set the width and height of the canvas to conver the entire screen
	width = document.Get("body").Get("clientWidth").Int()
	height = document.Get("body").Get("clientHeight").Int()
	SetCanvasSize(width, height)
	awaiter = make(chan int, 0)

	//Get the GL context
	contextOptions := js.Global().Get("JSON").Call("parse", "{ \"desynchronized\": true }")
	context := canvas.Call("getContext", "webgl", contextOptions)
	if context.IsUndefined() {
		context = canvas.Call("getContext", "experimental-webgl", contextOptions)
	}
	if context.IsUndefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return 0
	}

	//Create a new GL instance
	GL = newWebGL(context)

	//Setup the animation frame
	if !app.Start() {
		log.Println("Failed to start the application")
		return 0
	}

	//Record canvas events
	onMouseChangeEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		x := evt.Get("offsetX").Int()
		y := evt.Get("offsetY").Int()
		inputHandler.setMousePosition(x, y)

		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mousemove", onMouseChangeEvent, false)

	//Record canvas events
	onMouseUpEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		inputHandler.setMouseUp(button)

		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mouseup", onMouseUpEvent)

	//Record canvas events
	onMouseDownEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		inputHandler.setMouseDown(button)

		return nil
	})
	defer onMouseChangeEvent.Release()
	canvas.Call("addEventListener", "mousedown", onMouseDownEvent)

	//Begin rendering
	GL.Viewport(0, 0, width, height)
	frameRenderFunc = js.FuncOf(onRequestAnimationFrame)
	defer frameRenderFunc.Release()
	js.Global().Call("requestAnimationFrame", frameRenderFunc)
	return <-awaiter
}

//SetCanvasSize the size of the canvas
func SetCanvasSize(w, h int) {
	canvas.Set("width", w)
	canvas.Set("height", h)
	width = canvas.Get("width").Int()
	height = canvas.Get("height").Int()
	log.Println("resized canvas")
}

//RequestRedraw requests for a new animation frame
func RequestRedraw() {
	js.Global().Call("requestAnimationFrame", frameRenderFunc)
}

//Exit the application
func Exit() {
	awaiter <- 1
}

//onRequestAnimationFrame callback for animations
func onRequestAnimationFrame(this js.Value, args []js.Value) interface{} {
	//Setupt he time
	time := args[0].Float()
	deltaTime = time - frameTime
	frameTime = time

	//Update the input
	inputHandler.update()

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

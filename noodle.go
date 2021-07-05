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
	errBox          js.Value
	hasErrBox       = false
	frameRenderFunc js.Func
	inputHandler    *InputHandler
	app             Application

	texture    *Texture
	frameTime  float64
	deltaTime  float64
	frameCount int64

	awaiter chan int

	//DebugDraw causes renderers to display debug information
	DebugDraw = false
	//DebugDrawLoops is a less efficent way of drawing, but preserves the box representation when drawing
	DebugDrawLoops = true

	//AlwaysDraw continously draws
	AlwaysDraw = true

	antiAlias = true
)

//GetFrameTime returns the time the last frame was rendered
func GetFrameTime() float64 { return frameTime }

//GetDeltaTime returns a high accuracy difference in time between the last frame and the current one.
func GetDeltaTime() float64 { return deltaTime }

//GetFrameCount returns the current frame
func GetFrameCount() int64 { return frameCount }

//DT returns a less accurate version of GetDeltaTime, for all your 32bit mathmatic needs.
func DT() float32 { return float32(deltaTime) }

//Canvas gets the current canvas object. It's not recommended to deal with canvas directly
func Canvas() *js.Value {
	//TODO: Return a 'dom' element of the canvas
	return &canvas
}

//Input returns the current input handler
func Input() *InputHandler {
	return inputHandler
}

//SetAntiAlias sets the anti-alias state
func SetAntiAlias(state bool) {
	antiAlias = state
	setGLContext()
}

//setGLContext sets the context reference
func setGLContext() {
	GL = newWebGL(antiAlias, canvas)
}

//Run setups the WebGL context and runs the application. It is blocking and returns an exit code if Exit() is ever called.
func Run(application Application, canvasSelector string) int {
	app = application

	//Prepare the everything
	document = js.Global().Get("document")
	canvas = document.Call("querySelector", canvasSelector)
	canvas.Get("classList").Call("add", "noodle-canvas")
	setGLContext()

	inputHandler = newInput()

	//Prepare the awaiter
	awaiter = make(chan int, 0)

	//Setup the animation frame
	if !app.Start() {
		Error("Failed to start the application", nil)
		return 0
	}

	//Cursor Moved
	onMouseChangeEvent := AddEventListener("mousemove", func(this js.Value, args []js.Value) interface{} {
		evt := args[0]

		bounding := GL.BoundingBox()
		x := float32(evt.Get("pageX").Float()) - bounding.X
		y := float32(evt.Get("pageY").Float()) - bounding.Y
		inputHandler.setMousePosition(int(x), int(y))
		RequestRedraw()
		return nil
	})
	defer onMouseChangeEvent.Release()

	//Mouse Up
	onMouseUpEvent := AddEventListener("mouseup", func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		inputHandler.setMouseUp(button)
		return nil
	})
	defer onMouseUpEvent.Release()

	//Mouse Down
	onMouseDownEvent := AddEventListener("mousedown", func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		button := evt.Get("button").Int()
		inputHandler.setMouseDown(button)
		return nil
	})
	defer onMouseDownEvent.Release()

	//Mouse Scroll
	onMouseScrollEvent := AddEventListener("wheel", func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		scroll := evt.Get("deltaY").Float()
		inputHandler.setMouseScroll(float32(scroll))
		return nil
	})
	defer onMouseScrollEvent.Release()

	//Key Down
	onKeyDownEvent := AddEventListener("keydown", func(this js.Value, args []js.Value) interface{} {
		//Get the event and ditch repeated keys
		evt := args[0]
		if evt.Get("repeat").Bool() {
			return nil
		}
		key := evt.Get("keyCode").Int()
		inputHandler.setKeyDown(key)
		return nil
	})
	defer onKeyDownEvent.Release()

	//Key Up
	onKeyUpEvent := AddEventListener("keyup", func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		key := evt.Get("keyCode").Int()
		inputHandler.setKeyUp(key)
		return nil
	})
	defer onKeyUpEvent.Release()

	//Initial resize
	width, height := GL.Resize()
	GL.Viewport(0, 0, width, height)

	//Request a animation frame.
	frameRenderFunc = js.FuncOf(onRequestAnimationFrame)
	defer frameRenderFunc.Release()
	js.Global().Call("requestAnimationFrame", frameRenderFunc)
	return <-awaiter
}

//AddEventListener adds a new event listener to the document. It will return a JS function that needs to be Released() when its no longer required.
func AddEventListener(event string, fn func(this js.Value, args []js.Value) interface{}) js.Func {
	jsfunc := js.FuncOf(fn)
	document.Call("addEventListener", event, jsfunc)
	return jsfunc
}

//RequestRedraw requests for a new animation frame
func RequestRedraw() {
	js.Global().Call("requestAnimationFrame", frameRenderFunc)
}

//Exit the application
func Exit() {
	log.Println("Exiting...")
	awaiter <- 1
}

//onRequestAnimationFrame callback for animations
func onRequestAnimationFrame(this js.Value, args []js.Value) interface{} {
	//Setupt he time
	time := args[0].Float() / 1000
	deltaTime = time - frameTime
	frameTime = time
	frameCount++

	//Update the input
	inputHandler.update()

	//Call update on the Application
	app.Update(float32(deltaTime))

	//Prepare the view port and then render everything
	width, height := GL.Resize()
	GL.Viewport(0, 0, width, height)

	//Clear the canvas
	app.Render()

	//If we need to draw again, then do so
	if AlwaysDraw {
		RequestRedraw()
	}

	//Return nil to JS
	return nil
}

//Error logs the error message and additionally shows on screen
func Error(message string, err error) {
	log.Println(message, err)

	if !hasErrBox {
		hasErrBox = true

		// Create container
		errBox = document.Call("createElement", "div")
		errBox.Get("classList").Call("add", "noodle-errors")
		document.Get("body").Call("appendChild", errBox)

		// Set element style
		bounding := GL.BoundingBox()
		style := errBox.Get("style")
		style.Set("position", "absolute")
		style.Set("left", bounding.X)
		style.Set("top", bounding.Y)
		style.Set("width", bounding.Width)
		style.Set("height", bounding.Height)
		style.Set("overflow", "scroll")
	}

	// Append the error
	errorElement := document.Call("createElement", "div")
	errorElement.Get("classList").Call("add", "noodle-error")

	errorMessage := document.Call("createElement", "div")
	errorMessage.Get("classList").Call("add", "noodle-error-message")
	errorMessage.Set("innerText", message)

	errorStack := document.Call("createElement", "div")
	errorStack.Get("classList").Call("add", "noodle-error-stack")
	errorStack.Set("innerText", err.Error())

	errorElement.Call("appendChild", errorMessage)
	errorElement.Call("appendChild", errorStack)
	errBox.Call("appendChild", errorElement)
}

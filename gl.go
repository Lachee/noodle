package noodle

import (
	"errors"
	"log"
	"syscall/js"
)

//WebGLBuffer the js representation of a buffer
type WebGLBuffer = js.Value

//WebGLShader is a JS representation of a shader
type WebGLShader = js.Value

//WebGLShaderProgram is a JS representation of a shader program
type WebGLShaderProgram = js.Value

//WebGLUniformLocation is a JS representation of a uniform location
type WebGLUniformLocation = js.Value

//WebGLAttributeLocation is a representation of a attribute location
type WebGLAttributeLocation = int

//WebGLTexture is a JS representation of a texture
type WebGLTexture = js.Value

//WebGL is the base class that wraps GL functionality.
type WebGL struct {
	context                   js.Value
	canvas                    js.Value
	targetWidth, targetHeight int
}

func newWebGL(canvas js.Value) *WebGL {

	//Get the GL context
	contextOptions := js.Global().Get("JSON").Call("parse", "{ \"desynchronized\": true }")
	context := canvas.Call("getContext", "webgl", contextOptions)
	if context.IsUndefined() {
		context = canvas.Call("getContext", "experimental-webgl", contextOptions)
	}

	//Verify we actually have a context
	if context.IsUndefined() {
		js.Global().Call("alert", "browser might not support webgl")
		log.Fatalln("Context is undefined. Browser doesn't support WebGL!")
	}

	return &WebGL{
		canvas:       canvas,
		context:      context,
		targetWidth:  -1,
		targetHeight: -1,
	}
}

//NewBuffer creates, binds and sets the data of a new buffer
func (gl *WebGL) NewBuffer(target GLEnum, data interface{}, usage GLEnum) WebGLBuffer {
	buffer := gl.CreateBuffer()
	gl.BindBuffer(target, buffer)
	gl.BufferData(target, data, usage)
	return buffer
}

//CreateBuffer creates a WebGLBuffer object.
func (gl *WebGL) CreateBuffer() WebGLBuffer {
	return WebGLBuffer(gl.context.Call("createBuffer"))
}

//BindBuffer binds a given WebGLBuffer to a target.
func (gl *WebGL) BindBuffer(target GLEnum, buffer WebGLBuffer) {
	gl.context.Call("bindBuffer", target, buffer)
}

//BufferData sets the data of a buffer
func (gl *WebGL) BufferData(target GLEnum, data interface{}, usage GLEnum) {
	values := sliceToTypedArray(data)
	gl.context.Call("bufferData", target, values, usage)
}

//BufferSubData updates a subset of a buffer object's data store.
func (gl *WebGL) BufferSubData(target GLEnum, offset int, data interface{}) {
	values := sliceToTypedArray(data)
	gl.context.Call("bufferSubData", target, offset, values)
}

//CreateShader creates a new WebGLShader
func (gl *WebGL) CreateShader(shaderType GLEnum) WebGLShader {
	return gl.context.Call("createShader", shaderType)
}

//ShaderSource sets the shader source code
func (gl *WebGL) ShaderSource(shader WebGLShader, source string) {
	gl.context.Call("shaderSource", shader, source)
}

//CompileShader compiles the shader
func (gl *WebGL) CompileShader(shader WebGLShader) error {
	gl.context.Call("compileShader", shader)

	if !gl.GetShaderParameter(shader, GlCompileStatus).Bool() {
		err := errors.New(gl.GetShaderInfoLog(shader))
		reportError("Failed to compile shader", err)
		return err
	}

	return nil
}

//GetShaderParameter returns information about the given shader.
func (gl *WebGL) GetShaderParameter(shader WebGLShader, param GLEnum) js.Value {
	return gl.context.Call("getShaderParameter", shader, param)
}

//GetShaderInfoLog returns the information log for the specified WebGLShader object. It contains warnings, debugging and compile information.
func (gl *WebGL) GetShaderInfoLog(shader WebGLShader) string {
	return gl.context.Call("getShaderInfoLog", shader).String()
}

//NewShader creates, sources and compiles a new shader
func (gl *WebGL) NewShader(shaderType GLEnum, sourceCode string) (WebGLShader, error) {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, sourceCode)
	err := gl.CompileShader(shader)
	return shader, err
}

//DeleteShader marks a given WebGLShader object for deletion. It will then be deleted whenever the shader is no longer in use.
func (gl *WebGL) DeleteShader(shader WebGLShader) {
	gl.context.Call("deleteShader", shader)
}

//CreateProgram creates a new webgl shader program
func (gl *WebGL) CreateProgram() WebGLShaderProgram {
	return gl.context.Call("createProgram")
}

//AttachShader attaches a shader to the program
func (gl *WebGL) AttachShader(shaderProgram WebGLShaderProgram, shader WebGLShader) {
	gl.context.Call("attachShader", shaderProgram, shader)
}

//LinkProgram inks a given WebGLProgram, completing the process of preparing the GPU code for the program's fragment and vertex shaders.
func (gl *WebGL) LinkProgram(shaderProgram WebGLShaderProgram) error {
	gl.context.Call("linkProgram", shaderProgram)

	if !gl.GetProgramParameter(shaderProgram, GlLinkStatus).Bool() {
		err := errors.New(gl.GetProgramInfoLog(shaderProgram))
		return err
	}

	return nil
}

//UseProgram tells webgl to start using this program
func (gl *WebGL) UseProgram(shaderProgram WebGLShaderProgram) {
	gl.context.Call("useProgram", shaderProgram)
}

//GetProgramParameter returns information about the given program.
func (gl *WebGL) GetProgramParameter(shaderProgram WebGLShaderProgram, param GLEnum) js.Value {
	return gl.context.Call("getProgramParameter", shaderProgram, param)
}

//GetProgramInfoLog returns the information log for the specified WebGLProgram object. It contains errors that occurred during failed linking or validation of WebGLProgram objects.
func (gl *WebGL) GetProgramInfoLog(shaderProgram WebGLShaderProgram) string {
	return gl.context.Call("getProgramInfoLog", shaderProgram).String()
}

//NewProgram creates a new webgl shader program with some shaders and links it
func (gl *WebGL) NewProgram(shaders []WebGLShader) (WebGLShaderProgram, error) {
	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}

	err := gl.LinkProgram(program)
	return program, err
}

//GetUniformLocation returns the location of a specific uniform variable which is part of a given WebGLProgram.
func (gl *WebGL) GetUniformLocation(shaderProgram WebGLShaderProgram, location string) WebGLUniformLocation {
	return gl.context.Call("getUniformLocation", shaderProgram, location)
}

//GetAttribLocation gets a location of an attribute
func (gl *WebGL) GetAttribLocation(shaderProgram WebGLShaderProgram, attribute string) WebGLAttributeLocation {
	return gl.context.Call("getAttribLocation", shaderProgram, attribute).Int()
}

//VertexAttribPointer binds the buffer currently bound to gl.ARRAY_BUFFER to a generic vertex attribute of the current vertex buffer object and specifies its layout.
func (gl *WebGL) VertexAttribPointer(position WebGLAttributeLocation, size int, valueType GLEnum, normalized bool, stride int, offset int) {
	gl.context.Call("vertexAttribPointer", position, size, valueType, normalized, stride, offset)
}

//EnableVertexAttribArray turns on the generic vertex attribute array at the specified index into the list of attribute arrays.
func (gl *WebGL) EnableVertexAttribArray(position WebGLAttributeLocation) {
	gl.context.Call("enableVertexAttribArray", position)
}

//ClearColor sets the colour the screen will be cleared to
func (gl *WebGL) ClearColor(color Color) {
	gl.context.Call("clearColor", color.R, color.G, color.B, color.A)
}

//ClearDepth sets the z value that is set to the depth buffer every frame
func (gl *WebGL) ClearDepth(depth float64) {
	gl.context.Call("clearDepth", float32(depth))
}

//Viewport sets the viewport, which specifies the affine transformation of x and y from normalized device coordinates to window coordinates.
func (gl *WebGL) Viewport(x, y, width, height int) {
	gl.context.Call("viewport", x, y, width, height)
}

//Resize updates the size of the canvas to be contained completely within its parent
func (gl *WebGL) Resize() (int, int) {

	width := gl.canvas.Get("width").Int()
	height := gl.canvas.Get("height").Int()
	displayWidth := gl.canvas.Get("clientWidth").Int()
	displayHeight := gl.canvas.Get("clientHeight").Int()

	//If we manually set it, then do so.
	if (gl.targetWidth | gl.targetHeight) > 0 {
		displayWidth = gl.targetWidth
		displayHeight = gl.targetHeight
	}

	if width != displayWidth || height != displayHeight {
		gl.canvas.Set("width", displayWidth)
		gl.canvas.Set("height", displayHeight)
		return displayWidth, displayHeight
	}

	return width, height
}

//Width of the canvas in pixels
func (gl *WebGL) Width() int { return gl.canvas.Get("width").Int() }

//Height of the canvas in pixels
func (gl *WebGL) Height() int { return gl.canvas.Get("height").Int() }

//AspectRatio of the canvas
func (gl *WebGL) AspectRatio() float32 { return float32(gl.Width()) / float32(gl.Height()) }

//Size is the width and height of the canvas in pixels
func (gl *WebGL) Size() Vector2 { return NewVector2i(gl.Width(), gl.Height()) }

//DepthFunc specifies a function that compares incoming pixel depth to the current depth buffer value.
func (gl *WebGL) DepthFunc(function GLEnum) {
	gl.context.Call("depthFunc", function)
}

//BlendFunc specifies a function that compares incoming pixel depth to the current blending.
func (gl *WebGL) BlendFunc(sFactor GLEnum, gFactor GLEnum) {
	gl.context.Call("blendFunc", sFactor, gFactor)
}

//Enable enables a option
func (gl *WebGL) Enable(option GLEnum) {
	gl.context.Call("enable", option)
}

//Disable disables a option
func (gl *WebGL) Disable(option GLEnum) {
	gl.context.Call("disable", option)
}

//Clear empties the buffers
func (gl *WebGL) Clear(option GLEnum) {
	gl.context.Call("clear", option)
}

//DrawElements renders primitives from array data.
func (gl *WebGL) DrawElements(mode GLEnum, count int, valueType GLEnum, offset int) {
	gl.context.Call("drawElements", mode, count, valueType, offset)
}

//DrawArrays renders primitives from array data.
func (gl *WebGL) DrawArrays(mode GLEnum, first int, count int) {
	gl.context.Call("drawArrays", mode, first, count)
}

//CreateTexture creates a new texture on the GPU
func (gl *WebGL) CreateTexture() WebGLTexture {
	return gl.context.Call("createTexture")
}

//BindTexture binds a given WebGLTexture to a target (binding point).
func (gl *WebGL) BindTexture(target GLEnum, texture WebGLTexture) {
	gl.context.Call("bindTexture", target, texture)
}

//UnbindTexture unbinds the target texture. Alias of bindTexture(target, nil) as WebGLTexture cannot be nil
func (gl *WebGL) UnbindTexture(target GLEnum) {
	gl.context.Call("bindTexture", target, nil)
}

//ActiveTexture tells WebGL what texture state will be now modified
func (gl *WebGL) ActiveTexture(target GLEnum) {
	gl.context.Call("activeTexture", target)
}

//TexImage2D specifies a 2D image
func (gl *WebGL) TexImage2D(target GLEnum, level int, internalFormat GLEnum, format GLEnum, texelType GLEnum, pixels interface{}) {
	gl.context.Call("texImage2D", target, level, internalFormat, format, texelType, pixels)
}

//GenerateMipmap creats the Mipmap for a texture
func (gl *WebGL) GenerateMipmap(target GLEnum) {
	gl.context.Call("generateMipmap", target)
}

//TexParameteri set texture parameters
func (gl *WebGL) TexParameteri(target GLEnum, param GLEnum, value int) {
	gl.context.Call("texParameteri", target, param, value)
}

//TexParameterf set texture parameters
func (gl *WebGL) TexParameterf(target GLEnum, param GLEnum, value float64) {
	gl.context.Call("texParameterf", target, param, value)
}

//=== Uniform Setting

//Uniform1f specifies values of uniform variables
func (gl *WebGL) Uniform1f(location WebGLUniformLocation, value float32) {
	gl.context.Call("uniform1f", location, value)
}

//Uniform1fv specifies values of uniform variables
func (gl *WebGL) Uniform1fv(location WebGLUniformLocation, value []float32) {
	slice := sliceToTypedArray(value)
	gl.context.Call("uniform1fv", location, slice)
}

//Uniform1i specifies values of uniform variables
func (gl *WebGL) Uniform1i(location WebGLUniformLocation, value int) {
	gl.context.Call("uniform1i", location, value)
}

//Uniform1iv specifies values of uniform variables
func (gl *WebGL) Uniform1iv(location WebGLUniformLocation, value []int) {
	slice := sliceToTypedArray(value)
	gl.context.Call("uniform1iv", location, slice)
}

//Uniform2f specifies values of uniform variables
func (gl *WebGL) Uniform2f(location WebGLUniformLocation, value, value2 float32) {
	gl.context.Call("uniform2f", location, value, value2)
}

//Uniform2fv specifies values of uniform variables
func (gl *WebGL) Uniform2fv(location WebGLUniformLocation, value []float32) {
	slice := sliceToTypedArray(value)
	gl.context.Call("Uniform2fv", location, slice)
}

//Uniform2i specifies values of uniform variables
func (gl *WebGL) Uniform2i(location WebGLUniformLocation, value, value2 int) {
	gl.context.Call("uniform2i", location, value, value2)
}

//Uniform2iv specifies values of uniform variables
func (gl *WebGL) Uniform2iv(location WebGLUniformLocation, value []int) {
	slice := sliceToTypedArray(value)
	gl.context.Call("uniform2iv", location, slice)
}

//Uniform2v is an alias of Uniform2fv but with Vector support
func (gl *WebGL) Uniform2v(location WebGLUniformLocation, value Vector2) {
	tmp := sliceToTypedArray([]float32((*value.DecomposePointer())[:]))
	gl.context.Call("uniform2fv", location, tmp)
}

//UniformMatrix4fv specify matrix values for uniform variables.
// Transpose is excluded as it always has to be false anyways.
func (gl *WebGL) UniformMatrix4fv(location WebGLUniformLocation, matrix Matrix) {
	buffer := matrix.DecomposePointer()
	typedBuffer := sliceToTypedArray([]float32((*buffer)[:]))
	gl.context.Call("uniformMatrix4fv", location, false, typedBuffer)
}

//Call the internal context and reutrns the JS value
func (gl *WebGL) Call(m string, args ...interface{}) js.Value {
	//Prepare a list of converted values
	convts := make([]interface{}, len(args))

	//Iterate over the values
	for i, a := range args {

		//Convert
		val, isJsValue := a.(js.Value)
		if isJsValue {
			convts[i] = val
		} else {
			//We failed!
			//log.Println("Failed to convert a value to a JS equiv in the call wrapper", i, a)
			//log.Fatalln("^^^ GLCall Failed", m, args)
			convts[i] = a
		}
	}

	//Pass it to the context
	//log.Println("GLCall", m, convts)
	return gl.context.Call(m, convts...)
}

//IsUndefined checks if the context is undefined
func (gl *WebGL) IsUndefined() bool {
	return gl.context.IsUndefined()
}

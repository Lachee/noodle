package noodle

import (
	"syscall/js"
)

//WebGLBuffer the js representation of a buffer
type WebGLBuffer = js.Value
type WebGLShader = js.Value
type WebGLShaderProgram = js.Value
type WebGLUniformLocation = js.Value

//GL is a helper class that wraps webgl
type GL struct {
	context js.Value
}

func newGL(context js.Value) *GL {
	return &GL{
		context: context,
	}
}

//NewBuffer creates, binds and sets the data of a new buffer
func (gl *GL) NewBuffer(target GLEnum, data interface{}, usage GLEnum) WebGLBuffer {
	buffer := gl.CreateBuffer()
	gl.BindBuffer(target, buffer)
	gl.BufferData(target, data, usage)
	return buffer
}

//CreateBuffer creates a WebGLBuffer object.
func (gl *GL) CreateBuffer() WebGLBuffer {
	return WebGLBuffer(gl.context.Call("createBuffer"))
}

//BindBuffer binds a given WebGLBuffer to a target.
func (gl *GL) BindBuffer(target GLEnum, buffer WebGLBuffer) {
	gl.context.Call("bindBuffer", target, buffer)
}

//BufferData sets the data of a buffer
func (gl *GL) BufferData(target GLEnum, data interface{}, usage GLEnum) {
	gl.context.Call("bufferData", target, data, usage)
}

//CreateShader creates a new WebGLShader
func (gl *GL) CreateShader(shaderType GLEnum) WebGLShader {
	return gl.context.Call("createShader", shaderType)
}

//ShaderSource sets the shader source code
func (gl *GL) ShaderSource(shader WebGLShader, source string) {
	gl.context.Call("shaderSource", shader, source)
}

//CompileShader compiles the shader
func (gl *GL) CompileShader(shader WebGLShader) {
	gl.context.Call("compileShader", shader)
}

//NewShader creates, sources and compiles a new shader
func (gl *GL) NewShader(shaderType GLEnum, sourceCode string) WebGLShader {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, sourceCode)
	gl.CompileShader(shader)
	return shader
}

//CreateProgram creates a new webgl shader program
func (gl *GL) CreateProgram() WebGLShaderProgram {
	return gl.context.Call("createProgram")
}

//AttachShader attaches a shader to the program
func (gl *GL) AttachShader(shaderProgram WebGLShaderProgram, shader WebGLShader) {
	gl.context.Call("attachShader", shaderProgram, shader)
}

//LinkProgram inks a given WebGLProgram, completing the process of preparing the GPU code for the program's fragment and vertex shaders.
func (gl *GL) LinkProgram(shaderProgram WebGLShaderProgram) {
	gl.context.Call("linkProgram", shaderProgram)
}

//UseProgram tells webgl to start using this program
func (gl *GL) UseProgram(shaderProgram WebGLShaderProgram) {
	gl.context.Call("useProgram", shaderProgram)
}

//NewProgram creates a new webgl shader program with some shaders and links it
func (gl *GL) NewProgram(shaders []WebGLShader) WebGLShaderProgram {
	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)
	return program
}

//GetUniformLocation returns the location of a specific uniform variable which is part of a given WebGLProgram.
func (gl *GL) GetUniformLocation(shaderProgram WebGLShaderProgram, location string) WebGLUniformLocation {
	return gl.context.Call("getUniformLocation", shaderProgram, location)
}

//GetAttribLocation gets a location of an attribute
func (gl *GL) GetAttribLocation(shaderProgram WebGLShaderProgram, attribute string) int {
	return gl.context.Call("getAttribLocation", shaderProgram, attribute).Int()
}

//VertexAttribPointer binds the buffer currently bound to gl.ARRAY_BUFFER to a generic vertex attribute of the current vertex buffer object and specifies its layout.
func (gl *GL) VertexAttribPointer(position int, size int, valueType GLEnum, normalized bool, stride int, offset int) {
	gl.context.Call("vertexAttribPointer", position, size, valueType, normalized, stride, offset)
}

//EnableVertexAttribArray turns on the generic vertex attribute array at the specified index into the list of attribute arrays.
func (gl *GL) EnableVertexAttribArray(position int) {
	gl.context.Call("enableVertexAttribArray", position)
}

//ClearColor sets the colour the screen will be cleared to
func (gl *GL) ClearColor(r, g, b, a float64) {
	gl.context.Call("clearColor", float32(r), float32(g), float32(b), float32(a))
}

//ClearDepth sets the z value that is set to the depth buffer every frame
func (gl *GL) ClearDepth(depth float64) {
	gl.context.Call("clearDepth", float32(depth))
}

//Viewport sets the viewport, which specifies the affine transformation of x and y from normalized device coordinates to window coordinates.
func (gl *GL) Viewport(x, y, width, height int) {
	gl.context.Call("viewport", x, y, width, height)
}

//DepthFunc specifies a function that compares incoming pixel depth to the current depth buffer value.
func (gl *GL) DepthFunc(function GLEnum) {
	gl.context.Call("depthFunc", function)
}

//UniformMatrix4fv specify matrix values for uniform variables.
func (gl *GL) UniformMatrix4fv(location WebGLUniformLocation, matrix Matrix) {
	buffer := matrix.DecomposePointer()
	typedBuffer := sliceToTypedArray([]float32((*buffer)[:]))
	gl.Call("uniformMatrix4fv", location, false, typedBuffer)
}

//Enable enables a option
func (gl *GL) Enable(option GLEnum) {
	gl.context.Call("enable", option)
}

//Disable disables a option
func (gl *GL) Disable(option GLEnum) {
	gl.context.Call("disable", option)
}

//Clear empties the buffers
func (gl *GL) Clear(option GLEnum) {
	gl.context.Call("clear", option)
}

//DrawElements renders primitives from array data.
func (gl *GL) DrawElements(mode GLEnum, count int, valueType GLEnum, offset int) {
	gl.context.Call("drawElements", mode, count, valueType, offset)
}

//Call the internal context and reutrns the JS value
func (gl *GL) Call(m string, args ...interface{}) js.Value {
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
func (gl *GL) IsUndefined() bool {
	return gl.context.IsUndefined()
}

package noodle

import (
	"syscall/js"
)

//WebGLBuffer the js representation of a buffer
type WebGLBuffer = js.Value
type WebGLShader = js.Value
type WebGLShaderProgram = js.Value

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

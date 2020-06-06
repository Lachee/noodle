package noodle

import "syscall/js"

type GL struct {
	types   *GLTypes
	context js.Value
}

func newGL(context js.Value) *GL {
	return &GL{
		types:   newGLTypes(context),
		context: context,
	}
}

//CreateBuffer creates a WebGLBuffer object.
func (gl *GL) CreateBuffer() js.Value {
	return gl.Call("createBuffer")
}

//Call the internal context and reutrns the JS value
func (gl *GL) Call(m string, args ...interface{}) js.Value {
	return gl.context.Call(m, args...)
}

//IsUndefined checks if the context is undefined
func (gl *GL) IsUndefined() bool {
	return gl.context.IsUndefined()
}

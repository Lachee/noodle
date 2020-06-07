// +build js,wasm
//Borrowed from https://github.com/bobcob7/wasm-rotating-cube/blob/master/gltypes/gltypes.go
package noodle

import (
	"log"
	"syscall/js"
)

//GLTypesValue is a value in the GLTypes
//type GLEnum = js.Value

// GLTypes provides WebGL bindings.
type GLEnumCollection struct {
	StaticDraw         GLEnum
	ArrayBuffer        GLEnum
	ElementArrayBuffer GLEnum
	VertexShader       GLEnum
	FragmentShader     GLEnum
	Float              GLEnum
	DepthTest          GLEnum
	ColorBufferBit     GLEnum
	DepthBufferBit     GLEnum
	Triangles          GLEnum
	UnsignedShort      GLEnum
	LEqual             GLEnum
	LineLoop           GLEnum
}

func newGLEnumCollection(context js.Value) *GLEnumCollection {
	log.Println("Create Collection")
	var gltypes = &GLEnumCollection{}
	gltypes.find(context)
	return gltypes
}

// New grabs the WebGL bindings from a GL context.
func (collection *GLEnumCollection) find(gl js.Value) {
	/*
		collection.StaticDraw = GLEnum(gl.Get("STATIC_DRAW"))
		collection.ArrayBuffer = GLEnum(gl.Get("ARRAY_BUFFER"))
		collection.ElementArrayBuffer = GLEnum(gl.Get("ELEMENT_ARRAY_BUFFER"))
		collection.VertexShader = GLEnum(gl.Get("VERTEX_SHADER"))
		collection.FragmentShader = GLEnum(gl.Get("FRAGMENT_SHADER"))
		collection.Float = GLEnum(gl.Get("FLOAT"))
		collection.DepthTest = GLEnum(gl.Get("DEPTH_TEST"))
		collection.ColorBufferBit = GLEnum(gl.Get("COLOR_BUFFER_BIT"))
		collection.Triangles = GLEnum(gl.Get("TRIANGLES"))
		collection.UnsignedShort = GLEnum(gl.Get("UNSIGNED_SHORT"))
		collection.LEqual = GLEnum(gl.Get("LEQUAL"))
		collection.DepthBufferBit = GLEnum(gl.Get("DEPTH_BUFFER_BIT"))
		collection.LineLoop = GLEnum(gl.Get("LINE_LOOP"))
	*/
}

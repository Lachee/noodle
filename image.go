package noodle

import (
	"errors"
	"syscall/js"
)

type Image struct {
	data js.Value
}

//LoadImage loads a new image
func LoadImage(url string) (*Image, error) {
	ch := make(chan error, 1)
	img := js.Global().Get("Image").New()

	//Prepare the events
	loadEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() { ch <- nil }()
		return nil
	})
	defer loadEvent.Release()
	img.Call("addEventListener", "load", loadEvent)

	errorEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() { ch <- errors.New("Failed to load image") }()
		return nil
	})
	defer errorEvent.Release()
	img.Call("addEventListener", "error", errorEvent)

	//Set the source
	img.Set("src", url)

	//Wait for the source to load
	err := <-ch
	if err != nil {
		return nil, err
	}

	//Finish
	return &Image{img}, nil
}

//Data gets the JS value
func (i *Image) Data() js.Value {
	return i.data
}

//Width gets the width in pixels
func (i *Image) Width() int {
	return i.data.Get("width").Int()
}

//Height gets the height in pixels
func (i *Image) Height() int {
	return i.data.Get("height").Int()
}

//IsPowerOf2 checks if the image is a square power
func (i *Image) IsPowerOf2() bool {
	w := i.Width()
	h := i.Height()
	return ((w & (w - 1)) == 0) && ((h & (h - 1)) == 0)
}

type Texture struct {
	target  GLEnum
	level   int
	format  GLEnum
	texture WebGLTexture
}

//NewTexture a new Texture from the image
func NewTexture(image *Image) *Texture {
	webglTexture := GL.CreateTexture()
	tex := &Texture{
		target:  GlTexture2D,
		level:   0,
		format:  GlRGBA,
		texture: webglTexture,
	}

	tex.SetImage(image)
	return tex
}

//Data gets the internal JS reference
func (tex *Texture) Data() WebGLTexture {
	return tex.texture
}

//SetImage sets the texture's image
func (tex *Texture) SetImage(image *Image) {
	GL.BindTexture(tex.target, tex.texture)
	GL.TexImage2D(tex.target, tex.level, tex.format, tex.format, GlUnsignedByte, image.Data())

	//Generate mips
	if image.IsPowerOf2() {
		GL.GenerateMipmap(tex.target)
	} else {
		//Turn of mips, not square
		GL.TexParameteri(tex.target, GlTextureWrapS, GlClampToEdge)
		GL.TexParameteri(tex.target, GlTextureWrapT, GlClampToEdge)
		GL.TexParameteri(tex.target, GlTextureMinFilter, GlLinear)
	}
}

//Bind tells GL to use this texture
func (tex *Texture) Bind() {
	GL.BindTexture(tex.target, tex.texture)
}

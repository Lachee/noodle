package noodle

import (
	"errors"
	"syscall/js"
)

//UVTile interface provides methods for sprites
type UVTile interface {
	//Texture returns the texture
	Texture() *Texture

	//Width is the size in pixels of the source image
	Width() int

	//Height is the size in pixels of the source image
	Height() int

	//Slice returns uv extremes
	Slice() (float32, float32, float32, float32)
}

type Sprite struct {
	Source    *Texture
	Rectangle Rectangle
}

func (spr *Sprite) Texture() *Texture { return spr.Source }
func (spr *Sprite) Width() int        { return int(spr.Rectangle.X) }
func (spr *Sprite) Height() int       { return int(spr.Rectangle.Y) }
func (spr *Sprite) Slice() (float32, float32, float32, float32) {
	invTexWidth := 1.0 / float32(spr.Source.Width())
	invTexHeight := 1.0 / float32(spr.Source.Height())

	u := spr.Rectangle.X * invTexWidth
	v := spr.Rectangle.Y * invTexHeight
	u2 := (spr.Rectangle.X + spr.Rectangle.Width) * invTexWidth
	v2 := (spr.Rectangle.Y + spr.Rectangle.Height) * invTexHeight
	return u, v, u2, v2
}

//Image is a CPU image
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

//Texture is a GPU image
type Texture struct {
	target    GLEnum
	level     int
	format    GLEnum
	texture   WebGLTexture
	width     int
	height    int
	noMipMaps bool
}

//NewTexture a new Texture from the image
func NewTexture(image *Image) *Texture {
	webglTexture := GL.CreateTexture()
	tex := &Texture{
		target:    GlTexture2D,
		level:     0,
		format:    GlRGBA,
		texture:   webglTexture,
		width:     image.Width(),
		height:    image.Height(),
		noMipMaps: true,
	}

	tex.SetImage(image)
	return tex
}

//Width gets the width of the texture
func (tex *Texture) Width() int { return tex.width }

//Height gets the width of the texture
func (tex *Texture) Height() int { return tex.height }

//Texture returns this texture. Exists for compatability with the UVTile
func (tex *Texture) Texture() *Texture { return tex }

//Slice slices the texture into a UV region
func (tex *Texture) Slice() (float64, float64, float64, float64) {
	return 0.0, 0.0, 1.0, 1.0
}

//CreateSprite creates a sprite
func (tex *Texture) CreateSprite(region Rectangle) *Sprite {
	return &Sprite{
		Source:    tex,
		Rectangle: region,
	}
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
	if !tex.noMipMaps && image.IsPowerOf2() {
		GL.GenerateMipmap(tex.target)
	} else {
		//Turn of mips, not square
		GL.TexParameteri(tex.target, GlTextureWrapS, GlClampToEdge)
		GL.TexParameteri(tex.target, GlTextureWrapT, GlClampToEdge)
		GL.TexParameteri(tex.target, GlTextureMinFilter, GlNearest)
	}
}

//Bind tells GL to use this texture
func (tex *Texture) Bind() {
	GL.BindTexture(tex.target, tex.texture)
}

//SetSampler updates the texture sampler
func (tex *Texture) SetSampler(sampler WebGLUniformLocation, textureIndex int) {
	GL.ActiveTexture(GlTexture0 + textureIndex)
	GL.BindTexture(tex.target, tex.texture)
	GL.Uniform1i(sampler, textureIndex)
}

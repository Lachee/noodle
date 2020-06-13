package noodle

import (
	"errors"
	"image"
	"runtime"
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

	//Slice returns uv extremes, UV Min and UV Max
	Slice() (Vector2, Vector2)
}

//Sprite is a basic implementation of a UVTile. Its a single texture with a rectangle slice.
type Sprite struct {
	Source    *Texture
	Rectangle Rectangle
}

//Texture gets the current sprit texture
func (spr *Sprite) Texture() *Texture { return spr.Source }

//Width gets the sprites width in pixels
func (spr *Sprite) Width() int        { return int(spr.Rectangle.Width) }

//Height gets the sprites height in pixels
func (spr *Sprite) Height() int       { return int(spr.Rectangle.Height) }

//Slice generates a UV slice for the SpriteRenderer
func (spr *Sprite) Slice() (Vector2, Vector2) {
	invTexWidth := 1.0 / float32(spr.Source.Width())
	invTexHeight := 1.0 / float32(spr.Source.Height())

	u := spr.Rectangle.X * invTexWidth
	v := spr.Rectangle.Y * invTexHeight
	u2 := (spr.Rectangle.X + spr.Rectangle.Width) * invTexWidth
	v2 := (spr.Rectangle.Y + spr.Rectangle.Height) * invTexHeight
	return Vector2{u, v}, Vector2{u2, v2}
}

//NewSprite a new sprite
func NewSprite(source *Texture, rectangle Rectangle) *Sprite {
	return &Sprite{source, rectangle}
}

//Image is a CPU image
type Image struct {
	data   js.Value
	format GLEnum
	width  int
	height int
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

	width := img.Get("width").Int()
	height := img.Get("height").Int()

	//Finish
	return &Image{img, GlRGBA, width, height}, nil
}

//LoadImageRGBA loads a go RGBA image
func LoadImageRGBA(rgba *image.RGBA) (*Image, error) {
	//Get the pixels and convert it into a Uint8ClampedArray
	s := rgba.Pix
	a := js.Global().Get("Uint8Array").New(len(s))
	js.CopyBytesToJS(a, sliceToByteSlice(s))
	runtime.KeepAlive(s)
	buf := a.Get("buffer")
	ac := js.Global().Get("Uint8ClampedArray").New(buf, a.Get("byteOffset"), a.Get("byteLength"))

	//Create the image data
	bounds := rgba.Bounds()
	imageData := js.Global().Get("ImageData").New(ac, bounds.Dx(), bounds.Dy())

	//Return the final image
	return &Image{imageData, GlRGBA, bounds.Dx(), bounds.Dy()}, nil
}

//Data gets the JS value
func (i *Image) Data() js.Value {
	return i.data
}

//Width gets the width in pixels
func (i *Image) Width() int {
	return i.width
}

//Height gets the height in pixels
func (i *Image) Height() int {
	return i.height
}

//IsPowerOf2 checks if the image is a square power
func (i *Image) IsPowerOf2() bool {
	w := i.Width()
	h := i.Height()
	return ((w & (w - 1)) == 0) && ((h & (h - 1)) == 0)
}

//CreateTexture creates a new texture
func (i *Image) CreateTexture() *Texture {
	return NewTexture(i)
}

//TextureFilter  is the filter to use on a texture
type TextureFilter = GLEnum

//TextureWrap is how a texture should be wrapped.
type TextureWrap = GLEnum

const (
	//TextureFilterLinear linear filtering
	TextureFilterLinear TextureFilter = GlLinear
	//TextureFilterNearest nearest neighbour filtering
	TextureFilterNearest = GlNearest
	//TextureFilterNearestMipmapNearest nearest neighbour filtering
	TextureFilterNearestMipmapNearest = GlNearestMipmapNearest
	//TextureFilterLinearMipmapNearest linear filtering
	TextureFilterLinearMipmapNearest = GlLinearMipmapNearest
	//TextureFilterNearestMipmapLinear nearest neighbour filtering
	TextureFilterNearestMipmapLinear = GlNearestMipmapLinear
	//TextureFilterLinearMipmapLinear linear filtering
	TextureFilterLinearMipmapLinear = GlLinearMipmapLinear
)

const (
	//TextureWrapRepeat repeats the texture
	TextureWrapRepeat TextureWrap = GlRepeat
	//TextureWrapClampToEdge clamps the texture
	TextureWrapClampToEdge = GlClampToEdge
	//TextureWrapMirroredRepeat mirrors the texture
	TextureWrapMirroredRepeat = GlMirroredRepeat
)

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

//SetImage copies the data from the Image into the texture, setting initial filtering.
func (tex *Texture) SetImage(image *Image) {

	//Update the formatting
	tex.format = image.format

	//Setup the texture
	GL.BindTexture(tex.target, tex.texture)
	GL.TexImage2D(tex.target, tex.level, tex.format, tex.format, GlUnsignedByte, image.Data())

	//Generate mips
	if !tex.noMipMaps && image.IsPowerOf2() {
		GL.GenerateMipmap(tex.target)
	} else {
		//Turn of mips, not square
		GL.TexParameteri(tex.target, GlTextureWrapS, GlClampToEdge)
		GL.TexParameteri(tex.target, GlTextureWrapT, GlClampToEdge)
		GL.TexParameteri(tex.target, GlTextureMinFilter, GlLinear)
		GL.TexParameteri(tex.target, GlTextureMagFilter, GlLinear)
	}
}

//SetFilter binds the texture and sets the filtering
func (tex *Texture) SetFilter(filter TextureFilter) {
	GL.BindTexture(tex.target, tex.texture)
	GL.TexParameteri(tex.target, GlTextureMinFilter, filter)
	GL.TexParameteri(tex.target, GlTextureMagFilter, filter)
}

//SetWrap binds the texture and sets how the texture will be wrapped
func (tex *Texture) SetWrap(wrap TextureWrap) {
	GL.BindTexture(tex.target, tex.texture)
	GL.TexParameteri(tex.target, GlTextureWrapS, wrap)
	GL.TexParameteri(tex.target, GlTextureWrapT, wrap)
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

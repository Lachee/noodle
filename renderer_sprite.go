package noodle

import (
	"log"
	"math"
)

//Based Heavily from https://github.com/ajhager/engi/blob/master/b.go

const SpriteRendererSize = 10000

type SpriteRenderer struct {
	shader       *Shader
	inPosition   int
	inColor      int
	inTexCoords  int
	ufProjection WebGLUniformLocation
	ufvDimension WebGLUniformLocation
	ufvBorder    WebGLUniformLocation

	projX float64
	projY float64

	vertices     []float32
	indices      []uint16
	vertexBuffer WebGLBuffer
	indexBuffer  WebGLBuffer

	drawing     bool
	lastTexture *Texture
	index       int
}

func NewSpriteRenderer() *SpriteRenderer {
	b := &SpriteRenderer{}

	b.projX = float64(Width()) / 2
	b.projY = float64(Height()) / 2

	//Prepare the shader
	//b.shader = LoadShader(SpriteRendererVertCode, SpriteRendererFragCode)
	var shaderError error
	b.shader, shaderError = LoadShader(SpriteRendererVertCode, NineSliceRendererFragCode)
	if shaderError != nil {
		log.Fatalln("Failed to compile batch shader!", shaderError)
		return nil
	}

	b.inPosition = b.shader.GetAttribLocation("in_Position")
	b.inColor = b.shader.GetAttribLocation("in_Color")
	b.inTexCoords = b.shader.GetAttribLocation("in_TexCoords")
	b.ufProjection = b.shader.GetUniformLocation("uf_Projection")
	b.ufvDimension = b.shader.GetUniformLocation("u_dimensions")
	b.ufvBorder = b.shader.GetUniformLocation("u_border")

	//Prepare the verticies
	b.vertices = make([]float32, 20*SpriteRendererSize)
	b.indices = make([]uint16, 6*SpriteRendererSize)

	//Link all the indicies with the verticies, forming the quads
	for i, j := 0, 0; i < SpriteRendererSize*6; i, j = i+6, j+4 {
		b.indices[i+0] = uint16(j + 0)
		b.indices[i+1] = uint16(j + 1)
		b.indices[i+2] = uint16(j + 2)
		b.indices[i+3] = uint16(j + 0)
		b.indices[i+4] = uint16(j + 2)
		b.indices[i+5] = uint16(j + 3)
	}

	//Create the buffers

	b.indexBuffer = GL.CreateBuffer()
	GL.BindBuffer(GlElementArrayBuffer, b.indexBuffer)
	GL.BufferData(GlElementArrayBuffer, b.indices, GlStaticDraw)

	b.vertexBuffer = GL.CreateBuffer()
	GL.BindBuffer(GlArrayBuffer, b.vertexBuffer)
	GL.BufferData(GlArrayBuffer, b.vertices, GlDynamicDraw)

	//Enable them
	GL.EnableVertexAttribArray(b.inPosition)
	GL.EnableVertexAttribArray(b.inTexCoords)
	GL.EnableVertexAttribArray(b.inColor)

	//Bind the attributes
	GL.VertexAttribPointer(b.inPosition, 2, GlFloat, false, 20, 0)
	GL.VertexAttribPointer(b.inTexCoords, 2, GlFloat, false, 20, 8)
	GL.VertexAttribPointer(b.inColor, 4, GlUnsignedByte, true, 20, 16)

	return b
}

//Begin starts a SpriteRenderer
func (b *SpriteRenderer) Begin() {
	if b.drawing {
		log.Fatal("b.End() must be called first")
	}

	b.drawing = true
	GL.Enable(GlBlend)
	GL.BlendFunc(GlSrcColor, GlOneMinusSrcAlpha)
	GL.UseProgram(b.shader.GetProgram())
}

//End finalises a SpriteRenderer
func (b *SpriteRenderer) End() {
	if !b.drawing {
		log.Fatal("b.Begin() must be called first")
	}

	if b.index > 0 {
		b.flush()
	}

	b.drawing = false
	b.lastTexture = nil
}

//flush pushes the texture to GL
func (b *SpriteRenderer) flush() {
	if b.lastTexture == nil {
		return
	}

	//Bind the previous texture
	b.lastTexture.Bind()

	//Set the projection X and Y
	GL.Uniform2f(b.ufProjection, float32(b.projX), float32(b.projY))
	GL.Uniform2f(b.ufvBorder, 0.1, 0.1)
	GL.Uniform2f(b.ufvDimension, 1, 1)

	//Draw the buffer
	GL.BufferSubData(GlArrayBuffer, 0, b.vertices)
	GL.DrawElements(GlTriangles, 6*b.index, GlUnsignedShort, 0)

	b.index = 0
}

//Draw a particular texture
func (b *SpriteRenderer) Draw(r UVTile, position, origin, scale Vector2, rotation float64, color uint32, transparency float64) {
	if !b.drawing {
		log.Fatal("b.Begin() must be called first")
	}

	//Get the sprites image. If it has changed we need to flush ourselves then update the last image
	spriteImage := r.Texture()
	if spriteImage != b.lastTexture {
		if b.lastTexture != nil {
			b.flush()
		}
		b.lastTexture = spriteImage
	}

	x := position.X
	y := position.Y
	originX := origin.X
	originY := origin.Y
	scaleX := scale.X
	scaleY := scale.Y

	w := float32(r.Width())
	h := float32(r.Height())

	x -= originX * w
	y -= originY * h

	originX = w * originX
	originY = h * originY

	worldOriginX := x + originX
	worldOriginY := y + originY
	fx := -originX
	fy := -originY
	fx2 := w - originX
	fy2 := h - originY

	if scaleX != 1 || scaleY != 1 {
		fx *= scaleX
		fy *= scaleY
		fx2 *= scaleX
		fy2 *= scaleY
	}

	p1x := fx
	p1y := fy
	p2x := fx
	p2y := fy2
	p3x := fx2
	p3y := fy2
	p4x := fx2
	p4y := fy

	var x1 float32
	var y1 float32
	var x2 float32
	var y2 float32
	var x3 float32
	var y3 float32
	var x4 float32
	var y4 float32

	if rotation != 0 {
		rot := rotation * (math.Pi / 180.0)

		cos := float32(math.Cos(rot))
		sin := float32(math.Sin(rot))

		x1 = cos*p1x - sin*p1y
		y1 = sin*p1x + cos*p1y

		x2 = cos*p2x - sin*p2y
		y2 = sin*p2x + cos*p2y

		x3 = cos*p3x - sin*p3y
		y3 = sin*p3x + cos*p3y

		x4 = x1 + (x3 - x2)
		y4 = y3 - (y2 - y1)
	} else {
		x1 = p1x
		y1 = p1y

		x2 = p2x
		y2 = p2y

		x3 = p3x
		y3 = p3y

		x4 = p4x
		y4 = p4y
	}

	x1 += worldOriginX
	y1 += worldOriginY
	x2 += worldOriginX
	y2 += worldOriginY
	x3 += worldOriginX
	y3 += worldOriginY
	x4 += worldOriginX
	y4 += worldOriginY

	red := (color >> 16) & 0xFF
	green := ((color >> 8) & 0xFF) << 8
	blue := (color & 0xFF) << 16
	alpha := uint32(transparency*255.0) << 24
	tint := math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)

	idx := b.index * 20

	u, v, u2, v2 := r.Slice()

	b.vertices[idx+0] = float32(x1)
	b.vertices[idx+1] = float32(y1)
	b.vertices[idx+2] = float32(u)
	b.vertices[idx+3] = float32(v)
	b.vertices[idx+4] = float32(tint)

	b.vertices[idx+5] = float32(x4)
	b.vertices[idx+6] = float32(y4)
	b.vertices[idx+7] = float32(u2)
	b.vertices[idx+8] = float32(v)
	b.vertices[idx+9] = float32(tint)

	b.vertices[idx+10] = float32(x3)
	b.vertices[idx+11] = float32(y3)
	b.vertices[idx+12] = float32(u2)
	b.vertices[idx+13] = float32(v2)
	b.vertices[idx+14] = float32(tint)

	b.vertices[idx+15] = float32(x2)
	b.vertices[idx+16] = float32(y2)
	b.vertices[idx+17] = float32(u)
	b.vertices[idx+18] = float32(v2)
	b.vertices[idx+19] = float32(tint)

	//log.Println("x1, y1, u, v, tint", x1, y1, u, v)

	//increment the item index
	b.index += 1

	//We have reached the max, we should flush
	if b.index >= SpriteRendererSize {
		b.flush()
	}
}

var SpriteRendererVertCode = `
attribute vec2 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;
uniform vec2 uf_Projection;
varying vec4 var_Color;
varying vec2 var_TexCoords;
const vec2 center = vec2(-1.0, 1.0);
void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
	gl_Position = vec4(in_Position.x / uf_Projection.x + center.x,
										 in_Position.y / -uf_Projection.y + center.y,
										 0.0, 1.0);
}`

var SpriteRendererFragCode = `
#ifdef GL_ES
#define LOWP lowp
	precision mediump float;
#else
#define LOWP
#endif
varying vec4 var_Color;
varying vec2 var_TexCoords;
uniform sampler2D uf_Texture;
void main (void) {
  gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords);
}`

var NineSliceRendererFragCode = `

varying vec4 var_Color;
varying vec2 var_TexCoords;
uniform sampler2D uf_Texture;

uniform vec2 u_dimensions;
uniform vec2 u_border;

float map(float value, float originalMin, float originalMax, float newMin, float newMax) {
    return (value - originalMin) / (originalMax - originalMin) * (newMax - newMin) + newMin;
}

// Helper function, because WET code is bad code
// Takes in the coordinate on the current axis and the borders
float processAxis(float coord, float textureBorder, float windowBorder) {
    if (coord < windowBorder)
        return map(coord, 0, windowBorder, 0, textureBorder) ;
    if (coord < 1 - windowBorder)
        return map(coord,  windowBorder, 1 - windowBorder, textureBorder, 1 - textureBorder);
    return map(coord, 1 - windowBorder, 1, 1 - textureBorder, 1);
}

void main(void) {
		vec2 newUV = vec2(
         processAxis(var_TexCoords.x, u_border.x, u_dimensions.x),
         processAxis(var_TexCoords.y, u_border.y, u_dimensions.y)
     );
     gl_FragColor = texture2D(tex, newUV);
}
`

package noodle

import (
	"log"
	"math"
)

//Based Heavily from https://github.com/ajhager/engi/blob/master/b.go

const batchMaxSize = 10000

//SpriteRenderer renders UVTiles in a batched manner
type SpriteRenderer struct {
	shader       *Shader
	inPosition   int
	inColor      int
	inTexCoords  int
	ufProjection WebGLUniformLocation
	ufvDimension WebGLUniformLocation
	ufvBorder    WebGLUniformLocation

	projX float32
	projY float32

	vertices     []float32
	indices      []uint16
	vertexBuffer WebGLBuffer
	indexBuffer  WebGLBuffer

	drawing     bool
	lastTexture *Texture
	index       int
}

//NewSpriteRenderer creates a new sprite renderer
func NewSpriteRenderer() *SpriteRenderer {
	b := &SpriteRenderer{}

	//Prepare the shader
	var shaderError error
	b.shader, shaderError = LoadShader(spriteRendererVertCode, spriteRendererFragCode)
	if shaderError != nil {
		log.Fatalln("Failed to compile batch shader!", shaderError)
		return nil
	}

	b.inPosition = b.shader.GetAttribLocation("in_Position")
	b.inColor = b.shader.GetAttribLocation("in_Color")
	b.inTexCoords = b.shader.GetAttribLocation("in_TexCoords")
	b.ufProjection = b.shader.GetUniformLocation("uf_Projection")

	//Prepare the verticies
	b.vertices = make([]float32, 20*batchMaxSize)
	b.indices = make([]uint16, 6*batchMaxSize)

	//Link all the indicies with the verticies, forming the quads
	for i, j := 0, 0; i < batchMaxSize*6; i, j = i+6, j+4 {
		b.indices[i+0] = uint16(j + 0)
		b.indices[i+1] = uint16(j + 1)
		b.indices[i+2] = uint16(j + 2)
		b.indices[i+3] = uint16(j + 0)
		b.indices[i+4] = uint16(j + 2)
		b.indices[i+5] = uint16(j + 3)
	}

	//Create the buffers

	b.indexBuffer = GL.CreateBuffer()
	b.vertexBuffer = GL.CreateBuffer()

	GL.BindBuffer(GlElementArrayBuffer, b.indexBuffer)
	GL.BufferData(GlElementArrayBuffer, b.indices, GlStaticDraw)

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

	b.projX = float32(Width()) / 2.0
	b.projY = float32(Height()) / 2.0

	GL.Enable(GlBlend)
	GL.BlendFunc(GlSrcColor, GlOneMinusSrcAlpha)
	return b
}

//Begin starts a SpriteRenderer
func (b *SpriteRenderer) Begin() *SpriteRenderer {
	if b.drawing {
		log.Fatal("b.End() must be called first")
	}

	b.drawing = true
	GL.UseProgram(b.shader.GetProgram())

	//Set the projection X and Y
	b.projX = float32(Width()) / 2.0
	b.projY = float32(Height()) / 2.0
	GL.Uniform2f(b.ufProjection, b.projX, b.projY)

	return b
}

//End finalises a SpriteRenderer
func (b *SpriteRenderer) End() *SpriteRenderer {
	if !b.drawing {
		log.Fatal("b.Begin() must be called first")
	}

	if b.index > 0 {
		b.flush()
	}

	b.drawing = false
	//b.lastTexture = nil
	return b
}

//flush pushes the texture to GL
func (b *SpriteRenderer) flush() {
	if b.lastTexture == nil {
		return
	}

	//Bind the previous texture
	b.lastTexture.Bind()

	//Draw the buffer
	GL.BufferSubData(GlArrayBuffer, 0, b.vertices)
	GL.DrawElements(GlTriangles, 6*b.index, GlUnsignedShort, 0)

	//Reset the index
	b.index = 0
}

//Draw a particular texture
func (b *SpriteRenderer) Draw(r UVTile, origin Vector2, transform Transform2D, color uint32, transparency64 float64) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	if r.Texture() != b.lastTexture {
		if b.lastTexture != nil {
			b.flush()
		}
		b.lastTexture = r.Texture()
	}

	x := transform.Position.X
	y := transform.Position.Y
	originX := origin.X
	originY := origin.Y
	scaleX := transform.Scale.X
	scaleY := transform.Scale.Y
	rotation := transform.Rotation
	transparency := float32(transparency64)

	x -= originX * float32(r.Width())
	y -= originY * float32(r.Height())

	originX = float32(r.Width()) * originX
	originY = float32(r.Height()) * originY

	worldOriginX := x + originX
	worldOriginY := y + originY
	fx := -originX
	fy := -originY
	fx2 := float32(r.Width()) - originX
	fy2 := float32(r.Height()) - originY

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
		rot := float64(rotation * (math.Pi / 180.0))

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

	min, max := r.Slice()
	u := min.X
	v := min.Y
	u2 := max.X
	v2 := max.Y

	b.vertices[idx+0] = x1
	b.vertices[idx+1] = y1
	b.vertices[idx+2] = u
	b.vertices[idx+3] = v
	b.vertices[idx+4] = tint

	b.vertices[idx+5] = x4
	b.vertices[idx+6] = y4
	b.vertices[idx+7] = u2
	b.vertices[idx+8] = v
	b.vertices[idx+9] = tint

	b.vertices[idx+10] = x3
	b.vertices[idx+11] = y3
	b.vertices[idx+12] = u2
	b.vertices[idx+13] = v2
	b.vertices[idx+14] = tint

	b.vertices[idx+15] = x2
	b.vertices[idx+16] = y2
	b.vertices[idx+17] = u
	b.vertices[idx+18] = v2
	b.vertices[idx+19] = tint

	b.index++

	if b.index >= batchMaxSize {
		b.flush()
	}
}

var spriteRendererVertCode = `
attribute vec2 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;
uniform vec2 uf_Projection;
varying vec4 var_Color;
varying vec2 var_TexCoords;
const vec2 center = vec2(-1.0, 1.0);
void main() { var_Color = in_Color; var_TexCoords = in_TexCoords;	gl_Position = vec4(in_Position.x / uf_Projection.x + center.x, in_Position.y / -uf_Projection.y + center.y,  0.0, 1.0); }`

var spriteRendererFragCode = `
precision mediump float;
varying vec4 var_Color;
varying vec2 var_TexCoords;
uniform sampler2D uf_Texture;
void main (void) { gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords); }`

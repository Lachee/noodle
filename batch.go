package noodle

import (
	"log"
)

//Based Heavily from https://github.com/ajhager/engi/blob/master/batch.go

const batchSize = 10000

type TextureSlice interface {
	Texture() *Texture
	Width() float64
	Height() float64
	View() (float64, float64, float64, float64)
}

type Batch struct {
	shader       *Shader
	inPosition   int
	inColor      int
	inTexCoords  int
	ufProjection WebGLUniformLocation

	vertices     []float32
	indices      []uint16
	vertexBuffer WebGLBuffer
	indexBuffer  WebGLBuffer

	drawing     bool
	lastTexture *Texture
}

func NewBatch() *Batch {
	batch := &Batch{}

	//Prepare the shader
	batch.shader = LoadShader(batchVertCode, batchFragCode)
	batch.inPosition = batch.shader.GetAttribLocation("in_Position")
	batch.inColor = batch.shader.GetAttribLocation("in_Color")
	batch.inTexCoords = batch.shader.GetAttribLocation("in_TexCoords")
	batch.ufProjection = batch.shader.GetUniformLocation("uf_Projection")

	//Prepare the verticies
	batch.vertices = make([]float32, 20*batchSize)
	batch.indices = make([]uint16, 6*batchSize)

	//Link all the indicies with the verticies, forming the quads
	for i, j := 0, 0; i < batchSize*6; i, j = i+6, j+4 {
		batch.indices[i+0] = uint16(j + 0)
		batch.indices[i+1] = uint16(j + 1)
		batch.indices[i+2] = uint16(j + 2)
		batch.indices[i+3] = uint16(j + 0)
		batch.indices[i+4] = uint16(j + 2)
		batch.indices[i+5] = uint16(j + 3)
	}

	//Create the buffers

	batch.vertexBuffer = GL.NewBuffer(GlArrayBuffer, batch.vertices, GlDynamicDraw)
	batch.indexBuffer = GL.NewBuffer(GlElementArrayBuffer, batch.indices, GlStaticDraw)

	//Bind the attributes
	GL.VertexAttribPointer(batch.inPosition, 2, GlFloat, false, 20, 0)
	GL.VertexAttribPointer(batch.inTexCoords, 2, GlFloat, false, 20, 8)
	GL.VertexAttribPointer(batch.inColor, 4, GlUnsignedByte, true, 20, 16)

	//Enable them
	GL.EnableVertexAttribArray(batch.inPosition)
	GL.EnableVertexAttribArray(batch.inTexCoords)
	GL.EnableVertexAttribArray(batch.inColor)
	return batch
}

func (b *Batch) Begin() {
	if b.drawing {
		log.Fatal("Batch.End() must be called first")
	}

	b.drawing = true
	GL.Enable(GlBlend)
	GL.BlendFunc(GlSrcColor, GlOneMinusSrcAlpha)
	GL.UseProgram(b.shader.GetProgram())
}

func (b *Batch) End() {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	b.drawing = false
	b.lastTexture = nil
}

/*

//flush pushes the texture to GL
func (b *Batch) flush() {
	if b.lastImage == nil {
		return
	}

	gl.BindTexture(gl.TEXTURE_2D, b.lastTexture)

	gl.Uniform2f(b.ufProjection, b.projX, b.projY)

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, b.vertices)
	gl.DrawElements(gl.TRIANGLES, 6*b.index, gl.UNSIGNED_SHORT, 0)

	b.index = 0
}


func (b *Batch) Draw(r Sprite, x, y, originX, originY, scaleX, scaleY, rotation float64, color uint32, transparency float64) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	//Get the sprites image. If it has changed we need to flush ourselves then update the last image
	spriteImage := r.Image()
	if spriteImage != b.lastImage {
		if b.lastImage != nil {
			b.flush()
		}
		b.lastImage = spriteImage
	}

	x -= originX * r.Width()
	y -= originY * r.Height()

	originX = r.Width() * originX
	originY = r.Height() * originY

	worldOriginX := x + originX
	worldOriginY := y + originY
	fx := float32(-originX)
	fy := float32(-originY)
	fx2 := float32(r.Width() - originX)
	fy2 := float32(r.Height() - originY)

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

	u, v, u2, v2 := r.View()

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

	b.index += 1

	if b.index >= size {
		b.flush()
	}
}
*/

var batchVertCode = `
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

var batchFragCode = `
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

package noodle

import (
	"log"
	"math"
)

//uiRendererVertexLength is how many bytes are in each "vertex" element.
const uiRendererVertexLength = 44

//UIRenderer renders UVTiles in a batched manner
type UIRenderer struct {
	Zoom float32

	shader *Shader

	inPosition    WebGLAttributeLocation
	inTexCoords   WebGLAttributeLocation
	inSliceCoords WebGLAttributeLocation
	inDimension   WebGLAttributeLocation
	inColor       WebGLAttributeLocation

	uProjection WebGLUniformLocation
	uSampler    WebGLUniformLocation
	uBorder     WebGLUniformLocation

	vertices     []float32
	indices      []uint16
	vertexBuffer WebGLBuffer
	indexBuffer  WebGLBuffer

	drawing bool
	index   int

	sprite  *SliceSprite
	texture *Texture
}

//NewUIRenderer creates a new sprite renderer
func NewUIRenderer() *UIRenderer {
	b := &UIRenderer{}

	b.Zoom = 2.0

	//Prepare the shader
	var shaderError error
	b.shader, shaderError = LoadShaderFromURL("resources/shader/ui.vert", "resources/shader/ui.frag")
	if shaderError != nil {
		log.Fatalln("Failed to compile batch shader!", shaderError)
		return nil
	}

	b.inPosition = b.shader.GetAttribLocation("position")
	b.inTexCoords = b.shader.GetAttribLocation("texcoords")
	b.inSliceCoords = b.shader.GetAttribLocation("slicecoords")
	b.inDimension = b.shader.GetAttribLocation("dimension")
	b.inColor = b.shader.GetAttribLocation("color")

	b.uProjection = b.shader.GetUniformLocation("uProjection")
	b.uSampler = b.shader.GetUniformLocation("uSampler")
	b.uBorder = b.shader.GetUniformLocation("uBorder")

	//Prepare the verticies
	b.vertices = make([]float32, uiRendererVertexLength*batchMaxSize)
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
	GL.EnableVertexAttribArray(b.inSliceCoords)
	GL.EnableVertexAttribArray(b.inDimension)
	GL.EnableVertexAttribArray(b.inColor)

	//Bind the attributes
	//these numbers represent the values in the vertices buffer we are pushing, not the shader memory
	GL.VertexAttribPointer(b.inPosition, 2, GlFloat, false, uiRendererVertexLength, 0)
	GL.VertexAttribPointer(b.inTexCoords, 4, GlFloat, false, uiRendererVertexLength, 8)
	GL.VertexAttribPointer(b.inSliceCoords, 2, GlFloat, false, uiRendererVertexLength, 24)
	GL.VertexAttribPointer(b.inDimension, 2, GlFloat, false, uiRendererVertexLength, 32)
	GL.VertexAttribPointer(b.inColor, 4, GlUnsignedByte, true, uiRendererVertexLength, 40)
	return b
}

//Screen2UISpace converts the screen cooridinates to UI coords
func (b *UIRenderer) Screen2UISpace(screen Vector2) Vector2 {
	return screen.Scale(b.Zoom)
}

//Begin starts a UIRenderer
func (b *UIRenderer) Begin() *UIRenderer {
	if b.drawing {
		log.Fatal("b.End() must be called first")
	}

	b.drawing = true
	GL.UseProgram(b.shader.GetProgram())

	//Set the projection X and Y
	projX := float32(Width()) / b.Zoom
	projY := float32(Height()) / b.Zoom
	GL.Uniform2f(b.uProjection, projX, projY)

	//Clear the cache
	b.texture = nil
	b.sprite = nil

	//Enable the blending
	GL.Enable(GlBlend)
	GL.BlendFunc(GlOne, GlOneMinusSrcAlpha)
	return b
}

//End finalises a UIRenderer
func (b *UIRenderer) End() *UIRenderer {
	if !b.drawing {
		log.Fatal("b.Begin() must be called first")
	}

	if b.index > 0 {
		b.flush()
	}

	b.drawing = false
	return b
}

//flush pushes the texture to GL
func (b *UIRenderer) flush() {

	//Bind the texture
	if b.texture != b.sprite.Texture() {
		b.texture = b.sprite.Texture()
	}

	if b.texture != nil {
		b.texture.SetSampler(b.uSampler, 0)
	}

	//Update the border
	border := b.sprite.relativeBorder()
	GL.Uniform2v(b.uBorder, border)

	//Draw the buffer
	GL.BufferSubData(GlArrayBuffer, 0, b.vertices)
	GL.DrawElements(GlTriangles, 6*b.index, GlUnsignedShort, 0)

	//Reset the index
	b.index = 0
}

//SetSprite sets teh current sprite
func (b *UIRenderer) SetSprite(sprite *SliceSprite) {
	//Flush previous sprites
	if sprite != b.sprite && b.sprite != nil {
		b.flush()
	}

	//Setup the new sprite
	b.sprite = sprite
}

//Draw a particular texture
func (b *UIRenderer) Draw(rect Rectangle, color Color) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	//Prepare position and scale
	x := rect.X
	y := rect.Y
	scaleX := rect.Width
	scaleY := rect.Height

	//Unsused additional details that can be used for rotations and stuff
	rotation := float32(0)
	originX := float32(0)
	originY := float32(0)

	//Setup the origins
	x -= originX * float32(b.sprite.Width())
	y -= originY * float32(b.sprite.Height())

	originX = float32(b.sprite.Width()) * originX
	originY = float32(b.sprite.Height()) * originY

	worldOriginX := x + originX
	worldOriginY := y + originY
	fx := -originX
	fy := -originY
	fx2 := float32(b.sprite.Width()) - originX
	fy2 := float32(b.sprite.Height()) - originY

	//Update the scale
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

	//Prepare the tin
	tint := color.ToTint()

	//Get the min/max
	min, max := b.sprite.Slice()
	u := min.X
	v := min.Y
	u2 := max.X
	v2 := max.Y

	//Get the dimension
	dimensions := b.sprite.dimension(rect.Size())
	dx := dimensions.X
	dy := dimensions.Y

	//Setup the verts
	idx := b.index * uiRendererVertexLength
	b.vertices[idx+0] = x1
	b.vertices[idx+1] = y1
	b.vertices[idx+2] = u
	b.vertices[idx+3] = v
	b.vertices[idx+4] = u2
	b.vertices[idx+5] = v2
	b.vertices[idx+6] = 0
	b.vertices[idx+7] = 0
	b.vertices[idx+8] = dx
	b.vertices[idx+9] = dy
	b.vertices[idx+10] = tint

	b.vertices[idx+11] = x4
	b.vertices[idx+12] = y4
	b.vertices[idx+13] = u
	b.vertices[idx+14] = v
	b.vertices[idx+15] = u2
	b.vertices[idx+16] = v2
	b.vertices[idx+17] = 1
	b.vertices[idx+18] = 0
	b.vertices[idx+19] = dx
	b.vertices[idx+20] = dy
	b.vertices[idx+21] = tint

	b.vertices[idx+22] = x3
	b.vertices[idx+23] = y3
	b.vertices[idx+24] = u
	b.vertices[idx+25] = v
	b.vertices[idx+26] = u2
	b.vertices[idx+27] = v2
	b.vertices[idx+28] = 1
	b.vertices[idx+29] = 1
	b.vertices[idx+30] = dx
	b.vertices[idx+31] = dy
	b.vertices[idx+32] = tint

	b.vertices[idx+33] = x2
	b.vertices[idx+34] = y2
	b.vertices[idx+35] = u
	b.vertices[idx+36] = v
	b.vertices[idx+37] = u2
	b.vertices[idx+38] = v2
	b.vertices[idx+39] = 0
	b.vertices[idx+40] = 1
	b.vertices[idx+41] = dx
	b.vertices[idx+42] = dy
	b.vertices[idx+43] = tint

	//todo: region code hered

	b.index++

	if b.index >= batchMaxSize {
		b.flush()
	}
}

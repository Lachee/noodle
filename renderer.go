package noodle

/*
type Renderer interface {
	Setup()
	Render()
}

type TriangleRenderer struct {
	vertices []float32
}

func NewTriangleRenderer() *TriangleRenderer {
	renderer = &TriangleRenderer{}
	renderer.vertices = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}

	return renderer
}

func (renderer *TriangleRenderer) Setup() {

}
*/
/*

type CubeRenderer struct {
	verticesNative []float32
	colorsNative   []float32
	indicesNative  []float32
}

func NewCubeRenderer() *CubeRenderer {
	var cube = &CubeRenderer{}
	cube.verticesNative = []float32{
		-1, -1, -1, 1, -1, -1, 1, 1, -1, -1, 1, -1,
		-1, -1, 1, 1, -1, 1, 1, 1, 1, -1, 1, 1,
		-1, -1, -1, -1, 1, -1, -1, 1, 1, -1, -1, 1,
		1, -1, -1, 1, 1, -1, 1, 1, 1, 1, -1, 1,
		-1, -1, -1, -1, -1, 1, 1, -1, 1, 1, -1, -1,
		-1, 1, -1, -1, 1, 1, 1, 1, 1, 1, 1, -1,
	}
	cube.colorsNative = []float32{
		5, 3, 7, 5, 3, 7, 5, 3, 7, 5, 3, 7,
		1, 1, 3, 1, 1, 3, 1, 1, 3, 1, 1, 3,
		0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1,
		1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0,
		1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0,
		0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0,
	}
	cube.indicesNative = []uint16{
		0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7,
		8, 9, 10, 8, 10, 11, 12, 13, 14, 12, 14, 15,
		16, 17, 18, 16, 18, 19, 20, 21, 22, 20, 22, 23,
	}

	return cube
}

func (renderer *CubeRenderer) PreRender(gl *js.Value) {

}

func (renderer *CubeRenderer) Render(gl *js.Value) {

}
*/

attribute vec2 position;	//position
attribute vec4 texcoords;	//uv
attribute vec2 slicecoords;	//sprite size
attribute vec2 dimension;	//size of the box
attribute vec4 color;

uniform vec2 uProjection;

varying vec4 vTexCoords;
varying vec2 vSliceCoords;
varying vec2 vDimension;
varying vec4 vColor;

const vec2 center = vec2(-1.0, 1.0);

void main() { 
	vTexCoords = texcoords;
	vSliceCoords = slicecoords;	
	vDimension = dimension;
	vColor = color;

	gl_Position = vec4(position.x / uProjection.x + center.x, position.y / -uProjection.y + center.y,  0.0, 1.0); 
}
attribute vec2 position;
attribute vec2 texcoords;
attribute vec2 dimension;
attribute vec4 color;

uniform vec2 uProjection;
varying vec2 vTexCoords;
varying vec2 vDimension;
varying vec4 vColor;

const vec2 center = vec2(-1.0, 1.0);

void main() { 
	vTexCoords = texcoords;
	vColor = color;	
	vDimension = dimension;
	gl_Position = vec4(position.x / uProjection.x + center.x, position.y / -uProjection.y + center.y,  0.0, 1.0); 
}
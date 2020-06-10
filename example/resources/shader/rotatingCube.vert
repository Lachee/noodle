attribute vec3 position;
attribute vec2 textureCoord;

uniform mat4 Pmatrix;
uniform mat4 Vmatrix;
uniform mat4 Mmatrix;

attribute vec3 color;

varying vec3 vColor;
varying highp vec2 vTextureCoord;

void main(void) {
	gl_Position = Pmatrix * Vmatrix * Mmatrix * vec4(position, 1.);
	vColor = color;
	vTextureCoord = textureCoord;
}

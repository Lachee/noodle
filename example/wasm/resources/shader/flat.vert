attribute vec3 position;
attribute vec2 textureCoord;

uniform mat4 n_matrix;

attribute vec3 color;

varying vec3 vColor;
varying highp vec2 vTextureCoord;

void main(void) {
	gl_Position = n_matrix * vec4(position, 1.);
	vColor = color;
	vTextureCoord = textureCoord;
}

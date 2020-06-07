attribute vec3 position;
uniform mat4 Pmatrix;
uniform mat4 Vmatrix;
uniform mat4 Mmatrix;
attribute vec3 color;
varying vec3 vColor;
void main(void) {
	gl_Position = Pmatrix*Vmatrix*Mmatrix*vec4(position, 1.);
	vColor = color;
}

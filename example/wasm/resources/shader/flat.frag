precision mediump float;
varying vec3 vColor;

varying highp vec2 vTextureCoord;
uniform sampler2D uSampler;

void main(void) {
	gl_FragColor = vec4(vColor, 1.) * texture2D(uSampler, vTextureCoord);
}

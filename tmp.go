package noodle

const newvertShaderCode = `
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
`
const newfragShaderCode = `
precision mediump float;
varying vec3 vColor;

varying highp vec2 vTextureCoord;
uniform sampler2D uSampler;

uniform vec2 u_dimensions;
uniform vec2 u_border;


float map(float value, float originalMin, float originalMax, float newMin, float newMax) {
    return (value - originalMin) / (originalMax - originalMin) * (newMax - newMin) + newMin;
}

// Helper function, because WET code is bad code
// Takes in the coordinate on the current axis and the borders
float processAxis(float coord, float textureBorder, float windowBorder) {
    if (coord < windowBorder)
        return map(coord, 0.0, windowBorder, 0.0, textureBorder) ;
    if (coord < 1.0 - windowBorder)
        return map(coord,  windowBorder, 1.0 - windowBorder, textureBorder, 1.0- textureBorder);
    return map(coord, 1.0 - windowBorder, 1.0, 1.0 - textureBorder, 1.0);
}

void main(void) {
  vec2 newUV = vec2(
        processAxis(vTextureCoord.x, u_border.x, u_dimensions.x),
        processAxis(vTextureCoord.y, u_border.y, u_dimensions.y)
    );


    //newUV.xy += u_clip.xy / u_clip.wz;
    //newUV.xy *= u_clip.zw / u_texsize.xy;

    gl_FragColor =  texture2D(uSampler, newUV);
}
`

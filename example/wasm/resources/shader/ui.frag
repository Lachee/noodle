precision mediump float;

varying vec4 vTexCoords;
varying vec2 vSliceCoords;
varying vec2 vDimension;
varying vec4 vColor;

uniform sampler2D uSampler;
uniform vec2 uBorder;

float map(float value, float originalMin, float originalMax, float newMin, float newMax) {
    return (value - originalMin) / (originalMax - originalMin) * (newMax - newMin) + newMin;
}

float processAxis(float coord, float textureBorder, float windowBorder) {
    if (coord < windowBorder)
        return map(coord, 0.0, windowBorder, 0.0, textureBorder) ;
    if (coord < 1.0 - windowBorder)
        return map(coord,  windowBorder, 1.0 - windowBorder, textureBorder, 1.0 - textureBorder);
    return map(coord, 1.0 - windowBorder, 1.0, 1.0 - textureBorder, 1.0);
}

void main(void) {
    vec2 sliced = vec2(
        processAxis(vSliceCoords.x, uBorder.x, vDimension.x),
        processAxis(vSliceCoords.y, uBorder.y, vDimension.y)
    );

	vec2 size = vTexCoords.zw - vTexCoords.xy;
	vec2 nuv = vTexCoords.xy + size * sliced;
	gl_FragColor = vColor * texture2D(uSampler, nuv); 
}
precision mediump float;

varying vec2 vTexCoords;
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
    vec2 uv = vec2(
        processAxis(vTexCoords.x, uBorder.x, vDimension.x),
        processAxis(vTexCoords.y, uBorder.y, vDimension.y)
    );

    //newUV.x = vTextureCoord.x;


    //newUV.xy += u_clip.xy / u_clip.wz;
    //newUV.xy *= u_clip.zw / u_texsize.xy;

	gl_FragColor = vColor * texture2D(uSampler, uv); 
	//gl_FragColor = vec4(vTexCoords.x*vDimension.x/10.0, vTexCoords.y*vDimension.y/10.0, 0, 1);
}
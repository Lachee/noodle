precision mediump float;
varying vec4 var_Color;
varying vec2 var_TexCoords;
uniform sampler2D uf_Texture;
void main (void) { 
    gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords); 
}

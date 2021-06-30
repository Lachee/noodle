attribute vec2 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;
uniform vec2 uf_Projection;
uniform vec2 uf_Camera;
varying vec4 var_Color;
varying vec2 var_TexCoords;
const vec2 center = vec2(-1.0, 1.0);

void main() { 
    var_Color = in_Color; 
    var_TexCoords = in_TexCoords;	
    gl_Position = vec4(in_Position.x / uf_Projection.x + center.x, in_Position.y / -uf_Projection.y + center.y,  0.0, 1.0) + vec4(uf_Camera.x, uf_Camera.y, 0, 0); 
}
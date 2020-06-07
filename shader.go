package noodle

import (
	"io/ioutil"
	"net/http"
)

//Shader holds the shaders
type Shader struct {
	program WebGLShaderProgram
}

//LoadShaderFromURL loads a shader from a URL
func LoadShaderFromURL(vertURL, fragURL string) (*Shader, error) {

	//Load the vertext shader
	resp, err := http.Get(vertURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	vertCode := string(body)

	//Load the frag shader
	fragResp, err := http.Get(fragURL)
	if err != nil {
		return nil, err
	}
	defer fragResp.Body.Close()

	fragBody, err := ioutil.ReadAll(fragResp.Body)
	if err != nil {
		return nil, err
	}
	fragCode := string(fragBody)

	return LoadShader(vertCode, fragCode), nil
}

//LoadShader loads a shader from code
func LoadShader(vertCode, fragCode string) *Shader {
	vertex := GL.NewShader(GlVertexShader, vertCode)
	fragment := GL.NewShader(GlFragmentShader, fragCode)
	defer GL.DeleteShader(vertex)
	defer GL.DeleteShader(fragment)
	program := GL.NewProgram([]WebGLShader{vertex, fragment})
	return &Shader{program}
}

//GetProgram gets the shader program
func (shader *Shader) GetProgram() WebGLShaderProgram {
	return shader.program
}

//GetUniformLocation returns the location of a specific uniform variable which is part of a given WebGLProgram.
func (shader *Shader) GetUniformLocation(location string) WebGLUniformLocation {
	return GL.GetUniformLocation(shader.program, location)
}

//GetAttribLocation gets a location of an attribute
func (shader *Shader) GetAttribLocation(attribute string) int {
	return GL.GetAttribLocation(shader.program, attribute)
}

//BindVertexData binds a buffer of vertex data to an attribute
func (shader *Shader) BindVertexData(attribute string, bufferLocation GLEnum, buffer WebGLBuffer, bufferSize int, bufferType GLEnum, normalize bool, stride int, offset int) {
	pos := shader.GetAttribLocation(attribute)
	GL.BindBuffer(bufferLocation, buffer)
	GL.VertexAttribPointer(pos, bufferSize, bufferType, normalize, stride, offset)
	GL.EnableVertexAttribArray(pos)
}

//Use tells GL to use this shader
func (shader *Shader) Use() {
	GL.UseProgram(shader.program)
}

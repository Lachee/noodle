package noodle

//Shader holds the shaders
type Shader struct {
	program WebGLShaderProgram
}

//LoadShaderFromURL loads a shader from a URL
func LoadShaderFromURL(vertURL, fragURL string) (*Shader, error) {

	//Load the vertext shader
	vertCode, err := DownloadString(vertURL)
	if err != nil {
		return nil, err
	}

	//Load the frag shader
	fragCode, err := DownloadString(fragURL)

	return LoadShader(vertCode, fragCode)
}

//LoadShader loads a shader from code
func LoadShader(vertCode, fragCode string) (*Shader, error) {
	vertex, err := GL.NewShader(GlVertexShader, vertCode)
	defer GL.DeleteShader(vertex)
	if err != nil {
		return nil, err
	}

	fragment, err := GL.NewShader(GlFragmentShader, fragCode)
	defer GL.DeleteShader(fragment)
	if err != nil {
		return nil, err
	}

	program, err := GL.NewProgram([]WebGLShader{vertex, fragment})
	if err != nil {
		return nil, err
	}

	return &Shader{program}, nil
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
func (shader *Shader) GetAttribLocation(attribute string) WebGLAttributeLocation {
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

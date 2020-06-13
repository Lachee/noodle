package noodle

//Application handles the base framework
type Application interface {
	Start() bool
	Update(deltaTime float32)
	Render()
}

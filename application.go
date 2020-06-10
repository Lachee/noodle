package noodle

type Application interface {
	Start() bool
	Update(deltaTime float32)
	Render()
}

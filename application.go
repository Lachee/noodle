package noodle

type Application interface {
	Setup()
	Update(deltaTime float64)
	Render()
}

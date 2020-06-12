package main

import "github.com/lachee/noodle"

//Aliases are defined in aliases.go

func main() {
	var app noodle.Application

	//app = &RotatingCubeApp{}
	//app = &NineSliceApp{}
	app = &SpriteApp{}

	noodle.Initialize(app)
}

//BaseApplication handles the game. Put your variables in here
type BaseApplication struct {
}

//Start allows for setup
func (app *BaseApplication) Start() bool {
	return false
}

//Update runs once a frame
func (app *BaseApplication) Update(dt float32) {
}

//Render draws the frame
func (app *BaseApplication) Render() {
}

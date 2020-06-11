package main

import "github.com/lachee/noodle"

//Prepare some aliases, this will help us in the long run.
type Matrix = noodle.Matrix
type Vector2 = noodle.Vector2
type Vector3 = noodle.Vector3

func main() {
	app := &RotatingCubeApp{}
	//app := &NineSliceApp{}
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

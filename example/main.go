package main

import "github.com/lachee/noodle"

var app *MainApplication

func main() {
	app = &MainApplication{}
	noodle.Initialize(app)
}

type MainApplication struct {
}

func (app *MainApplication) Setup() {

}

func (app *MainApplication) Update(deltaTime float64) {

}

func (app *MainApplication) Render() {

}

package main

import (
	"log"
	"math"
	"math/rand"
	n "github.com/lachee/noodle"
)

//UIApp tests the sprite renderer
type UIApp struct {
	cursor   *n.Sprite
	sprite   *n.SliceSprite
	texture  *n.Texture
	batch    *n.UIRenderer
	boxs     []*Box
	boxCount int
}

//Box structure represents a sprite that will bounce around
type Box struct {
	sprite *n.SliceSprite
	rect   Rectangle
	speed   Vector2
}

func (box *Box) update(dt float32) {
	dt *= 0.10

	speedX := float64(0.0015 * box.speed.X)
	box.rect.Width = (float32(math.Sin(n.GetFrameTime()*speedX)/2) + 1) * 3
	
	speedY := float64(0.0015 * box.speed.Y)
	box.rect.Height = (float32(math.Cos(n.GetFrameTime()*speedY)/2) + 1) * 3
}

//Start is called by the noodle engine when ready
func (app *UIApp) Start() bool {
	//Setup the canvasss
	//n.SetCanvasSize(400, 300)

	//Prepare the image
	image, err := n.LoadImage("resources/slice.png") // The image URL
	if err != nil {
		log.Fatalln("Failed to spawn image", err)
		return false
	}

	//Setup the texture
	app.texture = image.CreateTexture()
	tileSize := float32(app.texture.Width()) / 6.0
	app.sprite = n.NewSliceSprite(app.texture, Rectangle{tileSize * 0, 0, tileSize, float32(app.texture.Height())}, Vector2{3, 4})
	app.batch = n.NewUIRenderer()

	cursor, _ := n.LoadImage("resources/cursors.svg")
	cursorTexture := cursor.CreateTexture()
	app.cursor = n.NewSprite(cursorTexture, Rectangle{0, 0, float32(cursorTexture.Width()) / 8.0, float32(cursorTexture.Height()) / 8.0})

	return true
}

//Update occurs once a frame
func (app *UIApp) Update(dt float32) {
	if n.Input().GetButtonDown(0) || app.boxCount == 0 {
		//TODO: Create a UI camera
		mouse := n.Input().GetMousePosition() //.Scale(1/10.0)
		for i := 0; i < 1; i++ {
			box := &Box{
				sprite: app.sprite,
				rect:   Rectangle{mouse.X, mouse.Y, 1, 1},
				speed:   Vector2{rand.Float32(), rand.Float32()},
			}
			app.boxs = append(app.boxs, box)
			app.boxCount++
		}

		log.Println("Boxs", app.boxCount)
	}

	//update the boxs
	for _, box := range app.boxs {
		box.update(dt)
	}
}

//Render occurs when the screen needs updating
func (app *UIApp) Render() {
	n.GL.ClearColor(1, 1, 1, 1)
	n.GL.Clear(n.GlColorBufferBit)
	
	app.batch.Begin()
	for _, box := range app.boxs {
		app.batch.SetSprite(box.sprite)
		app.batch.Draw(box.rect, n.White)
	}

	//mouse := n.Input().GetMousePosition()
	//t := n.NewTransform2D(mouse, 0, Vector2{1, 1})
	//app.batch.Draw(app.cursor, Vector2{0.5, 0.5}, t, 0xffffff, 1)

	//app.batch.Draw(app.sprite, Vector2{0, 0}, Vector2{0, 0}, Vector2{1, 1}, 0, 0xffffff, 1)
	app.batch.End()
}

package main

import (
	"log"
	"math/rand"

	n "github.com/lachee/noodle"
)

//Noodle Type Aliases are defined in aliases.go

//SpriteApp tests the sprite renderer
type SpriteApp struct {
	sprite    *n.Sprite
	texture   *n.Texture
	batch     *n.SpriteRenderer
	balls     []*Ball
	ballCount int
}

type Ball struct {
	sprite              *n.Sprite
	transform           Transform2D
	origin              Vector2
	velocity            Vector2
	angularVelocity     float32
	angularVelocitySign float32
}

func (ball *Ball) update(dt float32) {
	ball.transform.Position.X += ball.velocity.X * dt
	ball.transform.Position.Y += -ball.velocity.Y * dt

	//ball.transform.Rotation += ball.angularVelocity * dt * ball.angularVelocitySign

	if ball.transform.Position.X < 0 {
		ball.velocity.X *= -1
		ball.transform.Position.X = 0
		ball.angularVelocitySign = 1
	} else if ball.transform.Position.X > float32(n.Width()) {
		ball.velocity.X *= -1
		ball.transform.Position.X = float32(n.Width())
		ball.angularVelocitySign = -1
	}

	ball.velocity.Y += -0.0005 * dt
	if ball.transform.Position.Y > float32(n.Height()-ball.sprite.Height()) {
		ball.velocity.Y *= -.90
		ball.transform.Position.Y = float32(n.Height() - ball.sprite.Height())
		ball.angularVelocity *= 0.85
	}
}

func (app *SpriteApp) PrepareImage() (*n.Image, error) {
	return n.LoadImage("resources/tile.png") // The image URL
}

func (app *SpriteApp) Start() bool {

	//Setup the canvas
	n.SetCanvasSize(400, 300)

	//Prepare the image
	image, err := app.PrepareImage()
	if err != nil {
		log.Fatalln("Failed to spawn image", err)
		return false
	}

	//Setup the texture
	app.texture = image.CreateTexture()
	app.sprite = n.NewSprite(app.texture, Rectangle{0, 0, float32(app.texture.Width()), float32(app.texture.Height())})
	app.batch = n.NewSpriteRenderer()

	return true
}

//Update occurs once a frame
func (app *SpriteApp) Update(dt float32) {
	if n.Input().GetButton(0) {
		for i := 0; i < 100; i++ {
			ball := &Ball{
				sprite:              app.sprite,
				transform:           n.NewTransform2D(Vector2{rand.Float32() * 500, rand.Float32() * 500}, 0, Vector2{1, 1}),
				origin:              Vector2{0.5, 0.5},
				velocity:            Vector2{rand.Float32() * 0.5, 0},
				angularVelocity:     rand.Float32(),
				angularVelocitySign: 1,
			}
			app.balls = append(app.balls, ball)
			app.ballCount++
		}

		log.Println("Balls", app.ballCount)
	}

	//update the balls
	for _, ball := range app.balls {
		ball.update(dt)
	}
}

//Render occurs when the screen needs updating
func (app *SpriteApp) Render() {

	n.GL.Clear(n.GlColorBufferBit)

	app.batch.Begin()

	for _, ball := range app.balls {
		app.batch.Draw(ball.sprite, ball.origin, ball.transform, 0xffffff, 1)
	}
	//app.batch.Draw(app.sprite, Vector2{0, 0}, Vector2{0, 0}, Vector2{1, 1}, 0, 0xffffff, 1)
	app.batch.End()
}

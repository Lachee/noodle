package main

import (
	"log"
	"math/rand"

	n "github.com/lachee/noodle"
)

//Noodle Type Aliases are defined in aliases.go
const appSpritesAllowRotate = true

//SpriteApp tests the sprite renderer
type SpriteApp struct {
	cursor    *n.Sprite
	sprite    *n.Sprite
	texture   *n.Texture
	batch     *n.SpriteRenderer
	balls     []*Ball
	ballCount int
}

//Ball structure represents a sprite that will bounce around
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

	if appSpritesAllowRotate {
		ball.transform.Rotation += ball.angularVelocity * dt * ball.angularVelocitySign
	}

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

//Start is called by the noodle engine when ready
func (app *SpriteApp) Start() bool {
	//Setup the canvasss
	//n.SetCanvasSize(400, 300)

	//Prepare the image
	image, err := n.LoadImage("resources/snufkin.gif") // The image URL
	if err != nil {
		log.Fatalln("Failed to spawn image", err)
		return false
	}

	//Setup the texture
	app.texture = image.CreateTexture()
	app.sprite = n.NewSprite(app.texture, Rectangle{0, 0, float32(app.texture.Width()), float32(app.texture.Height())})
	app.batch = n.NewSpriteRenderer()

	cursor, _ := n.LoadImage("resources/cursors.svg")
	cursorTexture := cursor.CreateTexture()
	app.cursor = n.NewSprite(cursorTexture, Rectangle{0, 0, float32(cursorTexture.Width()) / 8.0, float32(cursorTexture.Height()) / 8.0})

	return true
}

//Update occurs once a frame
func (app *SpriteApp) Update(dt float32) {
	if n.Input().GetButton(0) {

		mouse := n.Input().GetMousePosition()
		t := n.NewTransform2D(mouse, 0, Vector2{1, 1})

		for i := 0; i < 1; i++ {
			ball := &Ball{
				sprite:              app.sprite,
				transform:           t,
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
		app.batch.Draw(ball.sprite, ball.origin, ball.transform, n.White)
	}

	mouse := n.Input().GetMousePosition()
	t := n.NewTransform2D(mouse, 0, Vector2{1, 1})
	app.batch.Draw(app.cursor, Vector2{0.5, 0.5}, t, n.White)

	//app.batch.Draw(app.sprite, Vector2{0, 0}, Vector2{0, 0}, Vector2{1, 1}, 0, 0xffffff, 1)
	app.batch.End()
}

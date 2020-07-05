package main

import (
	"log"
	"math"

	"github.com/golang/freetype/truetype"
	n "github.com/lachee/noodle"
)

//FontApp handles the game. Put your variables in here
type FontApp struct {
	font           *n.Font
	fontTTF        *n.Font
	fontBitmap     *n.Font
	cursor         *n.Sprite
	boxSprite      *n.SliceSprite
	spriteRenderer *n.SpriteRenderer
	uiRenderer     *n.UIRenderer
	previewRune    rune
}

//Start allows for setup
func (app *FontApp) Start() bool {

	//Load the TTF font
	fontData, err := n.DownloadFile("/resources/fonts/BalsamiqSans-Regular.ttf")
	//fontData, err := n.DownloadFile("/resources/fonts/ShareTechMono-Regular.ttf")
	//fontData, err := n.DownloadFile("/resources/fonts/LobsterTwo-Regular.ttf")
	//fontData, err := n.DownloadFile("/resources/fonts/Notable-Regular.ttf")
	//fontData, err := n.DownloadFile("/resources/fonts/luxirr.ttf")
	if err != nil {
		log.Fatalln("Failed to download font", err)
		return false
	}
	fontSrc, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatalln("Failed to parse the font", err)
		return false
	}
	options := &truetype.Options{Size: 100}
	fontFace := truetype.NewFace(fontSrc, options)
	app.fontTTF = n.LoadFont(fontFace, n.CharacterSetASCII)

	//Load the Bitmap Font
	fontImage, err := n.LoadImage("/resources/fonts/engi.png")
	if err != nil {
		log.Fatalln("Failed to load bitmap font", err)
		return false
	}
	app.fontBitmap = n.LoadFontBitmap(fontImage, n.CharacterSetCompleteASCII, 48, 6)
	app.font = app.fontBitmap

	//Load the renderers
	app.uiRenderer = n.NewUIRenderer()
	app.spriteRenderer = n.NewSpriteRenderer()
	cursor, _ := n.LoadImage("resources/cursors.svg")
	cursorTexture := cursor.CreateTexture()
	app.cursor = n.NewSprite(cursorTexture, Rectangle{0, 0, float32(cursorTexture.Width()) / 8.0, float32(cursorTexture.Height()) / 8.0})

	boxImage, err := n.LoadImage("resources/outline.png") // The image URL
	if err != nil {
		log.Fatalln("Failed to spawn image", err)
		return false
	}
	//Setup the texture
	boxTexture := boxImage.CreateTexture()
	tileSize := float32(boxTexture.Width()) / 1.0
	app.boxSprite = n.NewSliceSprite(boxTexture, Rectangle{tileSize * 0, 0, tileSize, float32(boxTexture.Height())}, Vector2{10, 10})
	app.uiRenderer = n.NewUIRenderer()

	return true
}

var index = 0

//Update runs once a frame
func (app *FontApp) Update(dt float32) {

	if n.Input().GetKeyDown(n.KeySpace) {
		n.DebugDraw = !n.DebugDraw
	}

	if n.Input().GetKeyDown(n.KeyOne) {
		app.font = app.fontTTF
	}
	if n.Input().GetKeyDown(n.KeyTwo) {
		app.font = app.fontBitmap
	}

	//Set a default font
	if app.font == nil {
		app.font = app.fontTTF
	}

	//Camera Control
	axis := n.Input().GetAxis2D(n.KeyArrowLeft, n.KeyArrowRight, n.KeyArrowDown, n.KeyArrowUp).Scale(1 / dt)
	app.spriteRenderer.Camera = app.spriteRenderer.Camera.Add(axis.Scale(-0.05 * app.spriteRenderer.Zoom))

	scroll := n.Input().GetMouseScroll()
	if scroll > 0 {
		app.spriteRenderer.Zoom += 0.01 * dt
	}
	if scroll < 0 {
		app.spriteRenderer.Zoom -= 0.01 * dt
	}

	if n.Input().GetButtonDown(1) {
		app.spriteRenderer.Zoom = 7.5
		app.spriteRenderer.Camera.X = -1.9
	}

	//index = int(math.Abs(math.Mod(float64(index-1), float64(len(app.font.Set)))))
	runeIndexerSpeed := float64(0.1)
	index = int(math.Mod(float64((n.GetFrameCount()+100))*runeIndexerSpeed, float64(len(app.font.GetCharset()))))
	app.previewRune = rune(app.font.GetCharset()[index])
}

//Render draws the frame
func (app *FontApp) Render() {

	//Set a default font
	if app.font == nil {
		app.font = app.fontTTF
	}

	//n.GL.ClearColor(1, 1, 1, 1)
	n.GL.Clear(n.GlColorBufferBit)

	app.renderFont()
	app.renderUI()
}

func (app *FontApp) renderFont() {

	app.spriteRenderer.Begin()
	//Draw the atlas
	//atlasTransform := n.NewTransform2D(Vector2{0, 0}, 0, Vector2{1, 1})

	//Draw the atlas
	fontTransform := n.NewTransform2D(n.NewVector2Zero(), 0, Vector2{1, 1})
	app.spriteRenderer.Draw(app.font.GetTexture(), Vector2{0, 0}, fontTransform, n.Black)

	//Draw the glyph demo
	glyphTransform := n.NewTransform2D(Vector2{0, 0}, 0, Vector2{1, 1})
	sprite := n.NewSprite(app.font.GetTexture(), app.font.GetGlyphs()[app.previewRune].Atlas)
	app.spriteRenderer.Draw(sprite, Vector2{0, 0}, glyphTransform, n.Black)

	//Draw the text
	app.font.GlyphString("\x02 0123465789 This is an example string! \x7F \xB0\xB1\xB2").RenderSprites(app.spriteRenderer, Vector2{0, float32(app.font.GetTexture().Height()) + 20}, 20.0/30.0, n.GopherBlue)

	app.spriteRenderer.End()
}

func (app *FontApp) renderUI() {
	app.uiRenderer.Begin()

	app.uiRenderer.SetSprite(app.boxSprite)

	mouse := n.Input().GetMousePosition()
	app.uiRenderer.Draw(n.NewRectangle(mouse.X, mouse.Y, 1, 1), n.Green)

	//Draw an outline against hte font texture
	w := float32(app.font.GetTexture().Width()) / app.uiRenderer.GetScale()
	h := float32(app.font.GetTexture().Height()) / app.uiRenderer.GetScale()
	app.uiRenderer.Draw(n.NewRectangle(0.0, 0.0, w, h), n.Red)

	b := app.font.GetGlyphs()[app.previewRune].Atlas
	pos := b.Position()
	size := b.Size().Scale(1.0 / app.uiRenderer.GetScale())
	app.uiRenderer.Draw(n.NewRectangleFromPositionSize(pos, size), n.GopherBlue)

	app.uiRenderer.End()
}

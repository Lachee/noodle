package main

import (
	"image"
	"log"
	"math"
	"strings"

	"github.com/galsondor/go-ascii"

	"github.com/golang/freetype/truetype"
	n "github.com/lachee/noodle"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//FontApp handles the game. Put your variables in here
type FontApp struct {
	font           *BitFont
	cursor         *n.Sprite
	boxSprite      *n.SliceSprite
	spriteRenderer *n.SpriteRenderer
	uiRenderer     *n.UIRenderer
	previewRune    rune
}

//Start allows for setup
func (app *FontApp) Start() bool {

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

	options := &truetype.Options{
		Size: 100,
	}
	fontFace := truetype.NewFace(fontSrc, options)

	kernA := rune('A')
	kernB := rune('V')
	fontKern := fontSrc.Kern(fixed.I(1), fontSrc.Index(kernA), fontSrc.Index(kernB))
	faceKern := fontFace.Kern(kernA, kernB)
	log.Println("Kern between", string(kernA), string(kernB), ". Face:", faceKern, ", Font:", fontKern)

	app.font = LoadBitFont(fontFace, CharacterSetASCII)
	//app.font.Spacing = -10
	/*
		// nGlyphs is the number of glyphs to generate: 95 characters in the range
		// [0x20, 0x7e], plus the replacement character.
		const nGlyphs = 95 + 1

		// The particular font (unicode.7x13.font) leaves the right-most column
		// empty in its ASCII glyphs. We don't have to include that column in the
		// generated glyphs, so we subtract one off the effective width.
		const width, height, ascent = 7 - 1, 13, 11

		dst := image.NewRGBA(image.Rect(0, 0, 500, 500))
		d := &font.Drawer{
			Dst:  dst,
			Src:  image.White,
			Face: fontFace,
			Dot:  fixed.P(0, ascent),
		}

		ascii := app.createASCII()
		log.Println("ASCII: ", ascii)
		d.DrawString(ascii)

		genImage, err := n.LoadImageRGBA(dst)
		if err != nil {
			log.Fatalln("Failed to convert the image", err)
			return false
		}
	*/

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
	//log.Println("scroll", scroll)

	if n.Input().GetButtonDown(1) {
		app.spriteRenderer.Zoom = 7.5
		app.spriteRenderer.Camera.X = -1.9
	}

	//index = int(math.Abs(math.Mod(float64(index-1), float64(len(app.font.Set)))))
	runeIndexerSpeed := float64(0.1)
	index = int(math.Mod(float64((n.GetFrameCount()+100))*runeIndexerSpeed, float64(len(app.font.Set))))
	app.previewRune = rune(app.font.Set[index])
}

//Render draws the frame
func (app *FontApp) Render() {

	n.GL.ClearColor(1, 1, 1, 1)
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
	app.spriteRenderer.Draw(app.font.Texture, Vector2{0, 0}, fontTransform, n.Black)

	//Draw the glyph demo
	glyphTransform := n.NewTransform2D(Vector2{0, 0}, 0, Vector2{1, 1})
	sprite := n.NewSprite(app.font.Texture, app.font.Glyphs[app.previewRune].AtlasBounding)
	app.spriteRenderer.Draw(sprite, Vector2{0, 0}, glyphTransform, n.Black)

	//Draw the text
	app.font.RenderSprites(app.spriteRenderer, "AV should be kerned.", Vector2{0, float32(app.font.Texture.Height())}, 30.0/20.0, n.Red)

	app.spriteRenderer.End()
}

func (app *FontApp) renderUI() {
	app.uiRenderer.Begin()

	app.uiRenderer.SetSprite(app.boxSprite)

	mouse := n.Input().GetMousePosition()
	app.uiRenderer.Draw(n.NewRectangle(mouse.X, mouse.Y, 1, 1), n.Green)

	//Draw an outline against hte font texture
	w := float32(app.font.Texture.Width()) / app.uiRenderer.GetScale()
	h := float32(app.font.Texture.Height()) / app.uiRenderer.GetScale()
	app.uiRenderer.Draw(n.NewRectangle(0.0, 0.0, w, h), n.Red)

	b := app.font.Glyphs[app.previewRune].AtlasBounding
	pos := b.Position() //.Scale(app.uiRenderer.GetScale())
	size := b.Size().Scale(1.0 / app.uiRenderer.GetScale())
	app.uiRenderer.Draw(n.NewRectangleFromPositionSize(pos, size), n.GopherBlue)

	//app.uiRenderer.Draw(Rectangle{10, 0, float32(app.font.Texture.Width()), float32(app.font.Texture.Height())}, n.GopherBlue)

	//bounds := app.font.Glyphs['a']
	//app.uiRenderer.Draw(bounds, n.White)

	app.uiRenderer.End()
}

//createASCII creates a set of ascii characters
func (app *FontApp) createASCII() string {
	var sb strings.Builder
	for i := 32; i < 255; i++ {
		if ascii.IsPrint(byte(i)) {
			sb.WriteRune(rune(i))
		}
	}

	return sb.String()
}

//CharacterSetASCII is a predefined set of visible ascii characters
const CharacterSetASCII = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

//BitFont respresents a font with a specified character set
type BitFont struct {
	//Set is the characters used in this font
	Set string

	//FontFace is the associated FontFace. If nil, the entire BitFont will be treated as a monospaced font.
	FontFace font.Face

	//Glyphs map each character in the set to a rectangle
	Glyphs map[rune]BitGlyph

	//Texture is the loaded texture that is in GPU memory
	Texture *n.Texture

	//Font Spacing
	Spacing float32
}

//BitGlyph is the metadata for a rune
type BitGlyph struct {
	AtlasBounding Rectangle //AtlasBounding is the position within the atlas the Glyph is in
	BoundMin      Vector2   //BoundMin is the minimum bounding box
	BoundMax      Vector2   //BoundMax is the maximum bounding box
	Advance       float32   //Advance is how far to move after this glpyh
}

//BitGlyphString is a representation of a glyph string, used for rendering
type BitGlyphString struct {
	Font       *BitFont    //Font is the font this string is using
	Positions  []Vector2   //Position is relative position of the sprite
	UV         []Rectangle //The UV of the sprite
	LineHeight float32     //LineHeight is the heighest character in the string. Use this for calculating multi-height strings.
}

//LoadBitFont loads a font
func LoadBitFont(fontFace font.Face, charset string) *BitFont {
	const CPL = 20    //CPL characters per line
	const padding = 5 // Padding between characters

	//Prepare hte font and metrics
	bf := &BitFont{
		Set:      charset,
		FontFace: fontFace,
		Glyphs:   make(map[rune]BitGlyph, len(charset)),
	}

	metrics := fontFace.Metrics()
	lines := int(math.Ceil(float64(len(charset)) / CPL))

	//Figure out the target image size
	width := font.MeasureString(fontFace, charset[:CPL]).Ceil() + (padding * CPL)
	height := lines*metrics.Height.Ceil() + metrics.Ascent.Round() + (padding * lines) + padding

	//Find the longest length. We discount the last line because its highly likely to be partial,
	// and it will introduce extra checks
	for i := 0; i < lines-1; i++ {
		lineLength := font.MeasureString(fontFace, charset[i*CPL:(i+1)*CPL]).Ceil()
		if lineLength > width {
			width = lineLength
		}
	}

	//Generate a new image and prepare the drawer
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	ascent := metrics.Ascent.Round()

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: fontFace,
		Dot:  fixed.P(0, ascent),
	}

	//Draw eeach character
	for _, r := range charset {
		drawer.Dot = drawer.Dot.Add(fixed.P(padding, 0))
		point, nextPos := drawer.BoundString(string(r))

		min := Vector2{float32(point.Min.X.Round()), float32(point.Min.Y.Round())}
		max := Vector2{float32(point.Max.X.Round()), float32(point.Max.Y.Round())}
		rect := n.NewRectangleFromMinMax(min, max)

		if drawer.Dot.X.Ceil()+nextPos.Ceil() > width {
			drawer.Dot = fixed.P(0, drawer.Dot.Y.Ceil()+metrics.Height.Ceil()+padding)
			rect.X = 0
			rect.Y += float32(metrics.Height.Ceil() + padding)
		}

		/*
			rect := Rectangle{
				X:      float32(point.Min.X.Floor()),
				Y:      float32(drawer.Dot.Y.Floor()),
				Width:  float32(charWidth.Floor()),
				Height: float32(metrics.Height.Ceil()),
			}
		*/
		bnd, adv, ok := fontFace.GlyphBounds(rune(r))
		if ok {
			bmin := Vector2{float32(bnd.Min.X.Round()), float32(bnd.Min.Y.Round())}
			bmax := Vector2{float32(bnd.Max.X.Round()), float32(bnd.Max.Y.Round())}
			bf.Glyphs[r] = BitGlyph{
				AtlasBounding: rect,
				BoundMin:      bmin,
				BoundMax:      bmax,
				Advance:       float32(adv.Round()),
			}
			drawer.DrawString(string(r))
		}
	}

	/*
		src := image.NewRGBA(image.Rect(0, 0, width, height))
		cyan := color.RGBA{100, 0, 0, 0xff}
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				switch {
				case x < width/2 && y < height/2: // upper left quadrant
					src.Set(x, y, cyan)
				case x >= width/2 && y >= height/2: // lower right quadrant
					src.Set(x, y, color.White)
				default:
					// Use zero value.
					src.Set(x, y, dst.RGBAAt(x, y))
				}
			}
		}
	*/

	image, _ := n.LoadImageRGBA(dst)
	bf.Texture = image.CreateTexture()

	//Log that atlas we just generated
	log.Println("Font Atlas: ", width, height, charset)
	return bf
}

//RenderSprites renders the message using the font as a series of sprites
func (f *BitFont) RenderSprites(renderer *n.SpriteRenderer, message string, position Vector2, scale float32, color n.Color) {
	bgs := f.GlyphString(message)
	bgs.RenderSprites(renderer, position, scale, color)
}

//GlyphString creates a list of every glyph required for a string, storing its UV in the atlas and its 2D position.
// This is used in rendering for drawing each character. A quad can be made at each position, with the UV supplied.
// All the origins of the sprites are assumed to be at 0,1 (bottom left). When drawing, draw from the bottom left corner. This will ensure the scaling is applied correctly.
// This function takes into account for their Kerning, but some TTF fonts may not support Kerning (Go Bug)
func (f *BitFont) GlyphString(message string) *BitGlyphString {

	//Prepare the BGS
	bgs := &BitGlyphString{}
	bgs.UV = make([]Rectangle, len(message))
	bgs.Positions = make([]Vector2, len(message))
	bgs.Font = f
	position := Vector2{0, 0}

	//Iterate over every character
	doKerning := true
	for i, c := range message {

		//Prepare the bounds and the sprite for the bounds
		r := rune(c)
		glyph := f.Glyphs[r]
		bgs.UV[i] = glyph.AtlasBounding

		//Get the sprite back to the baseline
		//baseline := -glyph.AtlasBounding.Height
		//offsetY := baseline + glyph.BoundMax.Y
		offsetY := glyph.BoundMax.Y
		offsetX := float32(0)

		//Update the position to account for the previous kerning, moving us backwards if required.
		if doKerning && i > 0 {
			pr := rune(message[i-1])
			kern := f.FontFace.Kern(pr, r)
			offsetX = float32(kern.Round())
		}

		//Store its position
		bgs.Positions[i] = position.Add(Vector2{offsetX, offsetY})

		//Store the heighest character
		if glyph.AtlasBounding.Height > bgs.LineHeight {
			bgs.LineHeight = glyph.AtlasBounding.Height
		}

		//Update the position progress
		position.X += glyph.Advance + offsetX + f.Spacing
		//position.X += glyph.BoundMax.X - glyph.BoundMin.X
		//position.Y += 5
	}

	return bgs
}

//RenderSprites renders the font using the sprite renderer
func (s *BitGlyphString) RenderSprites(renderer *n.SpriteRenderer, position Vector2, scale float32, color n.Color) {
	for i := range s.Positions {
		glyphTransform := n.NewTransform2D(position.Add(s.Positions[i]), 0, Vector2{1, 1})
		sprite := &n.Sprite{s.Font.Texture, s.UV[i]}
		renderer.Draw(sprite, Vector2{0, 1}, glyphTransform, color)
	}
}

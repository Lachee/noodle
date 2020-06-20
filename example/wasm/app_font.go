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
		Size: 30,
	}
	fontFace := truetype.NewFace(fontSrc, options)

	kernA := rune('A')
	kernB := rune('V')
	fontKern := fontSrc.Kern(fixed.I(1), fontSrc.Index(kernA), fontSrc.Index(kernB))
	faceKern := fontFace.Kern(kernA, kernB)
	log.Println("Kern between", string(kernA), string(kernB), ". Face:", faceKern, ", Font:", fontKern)

	app.font = LoadBitFont(fontFace, CharacterSetASCII)

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

//Update runs once a frame
func (app *FontApp) Update(dt float32) {

	runeIndexerSpeed := float64(0.1)
	runeIndex := int(math.Mod(float64((n.GetFrameCount()+100))*runeIndexerSpeed, float64(len(app.font.Set))))
	app.previewRune = rune(app.font.Set[runeIndex])
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

	//Draw the mouse
	//mouse := n.Input().GetMousePosition()
	glyphTransform := n.NewTransform2D(Vector2{350, 150}, 0, Vector2{1, 1})
	fontTransform := n.NewTransform2D(n.NewVector2Zero(), 0, Vector2{1, 1})
	//app.batch.Draw(app.cursor, Vector2{0.5, 0.5}, mouseTransform, 0xffffff, 1)

	app.spriteRenderer.Draw(app.font.Texture, Vector2{0, 0}, fontTransform, 0x0, 1)

	sprite := n.NewSprite(app.font.Texture, app.font.Glyphs[app.previewRune].AtlasBounding)
	app.spriteRenderer.Draw(sprite, Vector2{0, 0}, glyphTransform, 0x0, 1)

	app.font.RenderSprite(app.spriteRenderer, "AV should be kerned.", Vector2{400, 150})
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
}

//BitGlyph is the metadata for a rune
type BitGlyph struct {
	//AtlasBounding is the position within the atlas the Glyph is in
	AtlasBounding Rectangle

	//Width is the width of the character when drawing (this helps determine spacing)
	Width float32
	//Ascent is how far up from the base line it must travel
	Ascent float32
	//Descent is how far down from the base line it must travel
	Descent float32
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

	//Draw each line
	/*
		for i := 0; i < lines; i++ {
			from := i * CPL
			to := (i + 1) * CPL

			drawer.Dot = fixed.P(0, metrics.Ascent.Round()+(metrics.Height.Ceil()*i))
			if i == lines-1 {
				log.Println("FLB: ", charset[from:])
				drawer.DrawString(charset[from:])
			} else {
				log.Println("FLA: ", charset[from:to])
				drawer.DrawString(charset[from:to])
			}
		}
	*/
	for _, r := range charset {
		drawer.Dot = drawer.Dot.Add(fixed.P(padding, 0))
		point, nextPos := drawer.BoundString(string(r))

		min := Vector2{float32(point.Min.X.Round()), float32(point.Min.Y.Round())}
		max := Vector2{float32(point.Max.X.Round()), float32(point.Max.Y.Round())}
		rect := n.NewRectangleFromMinMax(min, max)

		if drawer.Dot.X.Ceil()+nextPos.Ceil() > width {
			drawer.Dot = fixed.P(0, drawer.Dot.Y.Ceil()+metrics.Height.Ceil()+padding)
			rect.X = 0
			rect.Y += float32(metrics.Height.Ceil())
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
			bf.Glyphs[r] = BitGlyph{
				AtlasBounding: rect,
				Width:         float32(adv.Round()),
				Ascent:        float32(bnd.Min.Y.Round()),
				Descent:       float32(bnd.Max.Y.Round()),
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

//RenderSprite renders the font as a sprite using the current SpriteRenderer
func (f *BitFont) RenderSprite(renderer *n.SpriteRenderer, message string, position Vector2) {

	for i, c := range message {

		//Prepare the bounds and the sprite for the bounds
		r := rune(c)
		glyph := f.Glyphs[r]
		sprite := n.NewSprite(f.Texture, glyph.AtlasBounding)

		//Update the position to account for the previous kerning, moving us backwards if required.
		if i > 0 {
			pr := rune(message[i-1])
			kern := f.FontFace.Kern(pr, r)
			//log.Println(string(pr), string(r), kern, kern.Round())
			position.X += float32(kern.Round())
		}

		//Get the sprite back to the baseline
		baseline := -glyph.AtlasBounding.Height
		offset := baseline + glyph.Descent

		//Draw the glyph
		glyphTransform := n.NewTransform2D(position.Add(Vector2{0, offset}), 0, Vector2{1, 1})
		renderer.Draw(sprite, Vector2{0, 0}, glyphTransform, 0x0, 1)

		//Update the position progress
		position.X += glyph.Width
	}
}

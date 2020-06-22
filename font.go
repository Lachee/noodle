package noodle

import (
	"image"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	//CharacterSetASCII provides standard printable ASCII characters
	CharacterSetASCII = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
)

//FontKerner handles kerning
type FontKerner interface {
	//Kern the font.
	Kern(a, b rune) fixed.Int26_6
}

//Font contains meta information about an available font atlas
type Font struct {
	charset string         //charset is the character set
	glyphs  map[rune]Glyph //glyphs is a map of glyphs
	texture *Texture       //texture is the atlas
	kerner  FontKerner     //kerner is what kerns the font

	Spacing float32 //Spacing is the offset between characters
}

//Glyph is a particular rune graphic in the atlas
type Glyph struct {
	Atlas   Rectangle //Atlas defines the UV of the glyph
	Ascent  float32   //Ascent is how far up it should be shifted
	Descent float32   //Descent is how far down it should be shifted
	Advance float32   //Advance is how far the next character should start
}

//GlyphString is a string made up out of glyphs.
type GlyphString struct {
	font       *Font
	Positions  []Vector2   //Positions is the relative position of each glyph
	UV         []Rectangle //UV is the UV of each glyph
	LineHeight float32     //LineHeight is the maximium height of each line
}

//LoadFont generates an atlas and prepares a font
func LoadFont(face font.Face, charset string) *Font {
	const CPL = 20
	const padding = 5

	f := &Font{
		charset: charset,
		glyphs:  make(map[rune]Glyph, len(charset)),
		Spacing: 0,
		kerner:  face,
	}

	metrics := face.Metrics()
	lines := int(math.Ceil(float64(len(charset)) / CPL))

	//Figure out the target image size
	width := font.MeasureString(face, charset[:CPL]).Ceil() + (padding * CPL)
	height := lines*metrics.Height.Ceil() + metrics.Ascent.Round() + (padding * lines) + padding

	//Find the longest length. We discount the last line because its highly likely to be partial,
	// and it will introduce extra checks
	for i := 0; i < lines-1; i++ {
		lineLength := font.MeasureString(face, charset[i*CPL:(i+1)*CPL]).Ceil()
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
		Face: face,
		Dot:  fixed.P(0, ascent),
	}

	//Draw eeach character
	for _, r := range charset {
		drawer.Dot = drawer.Dot.Add(fixed.P(padding, 0))
		point, nextPos := drawer.BoundString(string(r))

		min := Vector2{float32(point.Min.X.Round()), float32(point.Min.Y.Round())}
		max := Vector2{float32(point.Max.X.Round()), float32(point.Max.Y.Round())}
		rect := NewRectangleFromMinMax(min, max)

		if drawer.Dot.X.Ceil()+nextPos.Ceil() > width {
			drawer.Dot = fixed.P(0, drawer.Dot.Y.Ceil()+metrics.Height.Ceil()+padding)
			rect.X = 0
			rect.Y += float32(metrics.Height.Ceil() + padding)
		}

		bnd, adv, ok := face.GlyphBounds(rune(r))
		if ok {
			f.glyphs[r] = Glyph{
				Atlas:   rect,
				Ascent:  float32(bnd.Min.Y.Round()),
				Descent: float32(bnd.Max.Y.Round()),
				Advance: float32(adv.Round()),
			}
			drawer.DrawString(string(r))
		}
	}

	//Load the image and create the texture
	image, _ := LoadImageRGBA(dst)
	f.texture = image.CreateTexture()
	return f
}

//GetCharset gets the current character set
func (f *Font) GetCharset() string { return f.charset }

//GetGlyphs gets a map of glyphs
func (f *Font) GetGlyphs() map[rune]Glyph { return f.glyphs }

//GetTexture gets the current atlas texture
func (f *Font) GetTexture() *Texture { return f.texture }

//Kern gets the spacing between two runes
func (f *Font) Kern(left, right rune) float32 { return float32(f.kerner.Kern(left, right).Round()) }

/*
GlyphString turns the string into a series of Glyphs.
This is used in rendering for drawing each character. A quad can be made at each position, with the UV supplied.
All the origins of the sprites are assumed to be at 0,1 (bottom left). When drawing, draw from the bottom left corner. This will ensure the scaling is applied correctly.
This function takes into account for their Kerning, but some TTF fonts may not support Kerning (Go Bug)
*/
func (f *Font) GlyphString(str string) *GlyphString {
	gstr := &GlyphString{}
	gstr.UV = make([]Rectangle, len(str))
	gstr.Positions = make([]Vector2, len(str))
	gstr.font = f
	position := Vector2{0, 0}

	//Iterate over every character
	doKerning := true
	for i, c := range str {

		//Prepare the bounds and the sprite for the bounds
		r := rune(c)
		glyph := f.glyphs[r]
		gstr.UV[i] = glyph.Atlas

		//Get the sprite back to the baseline
		//baseline := -glyph.AtlasBounding.Height
		//offsetY := baseline + glyph.BoundMax.Y
		offsetY := glyph.Descent
		offsetX := float32(0)

		//Update the position to account for the previous kerning, moving us backwards if required.
		if doKerning && i > 0 {
			offsetX = f.Kern(rune(str[i-1]), r)
		}

		//Store its position
		gstr.Positions[i] = position.Add(Vector2{offsetX, offsetY})

		//Store the heighest character
		if glyph.Atlas.Height > gstr.LineHeight {
			gstr.LineHeight = glyph.Atlas.Height
		}

		//Update the position progress
		position.X += glyph.Advance + offsetX + f.Spacing
	}

	return gstr
}

//GetFont gets the current font in the glyph stirng
func (gstr *GlyphString) GetFont() *Font { return gstr.font }

//GetTexture gets the current atlas texture
func (gstr *GlyphString) GetTexture() *Texture { return gstr.font.GetTexture() }

//RenderSprites uses the SpriteRenderer to draw the glyphs. Its main purpose is to serve as an example on how a renderer could be writen for the fonts.
// see noodle/font.go for this function.
func (gstr *GlyphString) RenderSprites(renderer *SpriteRenderer, position Vector2, scale float32, color Color) {

	//Iterate over every position. This represents a new glyph
	for i := range gstr.Positions {

		//Prepare the position of the glyph, which is the current position, shifted
		pos := gstr.Positions[i].Scale(scale).Add(position)
		tex := gstr.font.texture //This is using internals, but you could use gstr.GetTexture() instead here.

		//Prepare the transform and sprite
		transform := NewTransform2D(pos, 0, Vector2{scale, scale})
		sprite := NewSprite(tex, gstr.UV[i])

		//Draw it, using the bottom left as the origin
		renderer.Draw(sprite, Vector2{0, 1}, transform, color)
	}
}

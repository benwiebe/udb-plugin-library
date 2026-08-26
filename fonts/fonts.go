package fonts

import (
	"embed"

	"github.com/golang/freetype/truetype"
)

// Embed the fonts exposed by the helper API. License text stays in the repo,
// but it is not bundled into the Go binary.
//
//go:embed assets/VT323/VT323-Regular.ttf assets/Pixelify_Sans/PixelifySans-Regular.ttf assets/Silkscreen/Silkscreen-Regular.ttf assets/Montserrat/Montserrat-Medium.ttf
var fontAssets embed.FS

/*
UdbFonts is a struct that holds references to various fonts and helpers
which are useful for rendering text on a pixel board display.
*/
type UdbFonts struct {
	vt323      *truetype.Font
	pixelify   *truetype.Font
	silkscreen *truetype.Font
	montserrat *truetype.Font
}

var udbFonts *UdbFonts

func GetUdbFonts() *UdbFonts {
	if udbFonts == nil {
		udbFonts = &UdbFonts{}
	}

	return udbFonts
}

/*
VT323 returns the VT323 font, which works well at very small sizes and has a
terminal-like pixel style.
*/
func (f *UdbFonts) VT323() *truetype.Font {
	if f.vt323 != nil {
		return f.vt323
	}

	fontBytes, err := fontAssets.ReadFile("assets/VT323/VT323-Regular.ttf")
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	f.vt323 = font
	return font
}

/*
PixelifySans returns the Pixelify Sans regular font. It is available for boards
that want the style explicitly, but it is not selected by FontForSize because it
gets muddy when dense text is squeezed onto small pixel boards.
*/
func (f *UdbFonts) PixelifySans() *truetype.Font {
	if f.pixelify != nil {
		return f.pixelify
	}

	fontBytes, err := fontAssets.ReadFile("assets/Pixelify_Sans/PixelifySans-Regular.ttf")
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	f.pixelify = font
	return font
}

/*
Silkscreen returns the Silkscreen regular font. It is available for boards that
want the style explicitly, but its stencil-like shapes are not a good default
for dense text on small pixel boards.
*/
func (f *UdbFonts) Silkscreen() *truetype.Font {
	if f.silkscreen != nil {
		return f.silkscreen
	}

	fontBytes, err := fontAssets.ReadFile("assets/Silkscreen/Silkscreen-Regular.ttf")
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	f.silkscreen = font
	return font
}

/*
Montserrat returns the Montserrat Medium font, which is a good general-purpose font
*/
func (f *UdbFonts) Montserrat() *truetype.Font {
	if f.montserrat != nil {
		return f.montserrat
	}

	fontBytes, err := fontAssets.ReadFile("assets/Montserrat/Montserrat-Medium.ttf")
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	f.montserrat = font
	return font
}

/*
FontForSize returns an appropriate font based on the given size. The font
should render nicely on a pixel board display.
*/
func (f *UdbFonts) FontForSize(size int) *truetype.Font {
	if size <= 18 {
		return f.VT323()
	}

	return f.Montserrat()
}

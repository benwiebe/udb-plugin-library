package fonts

import (
	"embed"

	"github.com/golang/freetype/truetype"
)

//go:embed assets
var fontAssets embed.FS

/*
UdbFonts is a struct that holds references to various fonts and helpers
which are useful for rendering text on a pixel board display.
*/
type UdbFonts struct {
	micro5     *truetype.Font
	jersey10   *truetype.Font
	jersey15   *truetype.Font
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
Micro5 returns the Micro 5 font, which is designed as 5px tall.
*/
func (f *UdbFonts) Micro5() *truetype.Font {
	if f.micro5 != nil {
		return f.micro5
	}

	fontBytes, err := fontAssets.ReadFile("assets/Micro_5/Micro5-Regular.ttf")
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	f.micro5 = font
	return font
}

/*
Jersey10 returns the Jersey 10 font, which is designed as 10px tall.
*/
func (f *UdbFonts) Jersey10() *truetype.Font {
	if f.jersey10 != nil {
		return f.jersey10
	}

	fontBytes, err := fontAssets.ReadFile("assets/Jersey_10/Jersey10-Regular.ttf")
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	f.jersey10 = font
	return font
}

/*
Jersey15 returns the Jersey 15 font, which is designed as 15px tall.
*/
func (f *UdbFonts) Jersey15() *truetype.Font {
	if f.jersey15 != nil {
		return f.jersey15
	}

	fontBytes, err := fontAssets.ReadFile("assets/Jersey_15/Jersey15-Regular.ttf")
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	f.jersey15 = font
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
	if size <= 8 {
		return f.Micro5()
	} else if size <= 12 {
		return f.Jersey10()
	} else if size <= 18 {
		return f.Jersey15()
	}

	return f.Montserrat()
}

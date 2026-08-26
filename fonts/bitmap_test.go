package fonts

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func TestFitPixelFaceSelectsScoreboardForNumericText(t *testing.T) {
	face := GetUdbFonts().FitPixelFace("12:34", 64, 32, PixelFontAuto)
	defer face.Close()

	if got := font.MeasureString(face, "12:34").Ceil(); got > 64 {
		t.Fatalf("expected fitted face to fit width, got %d", got)
	}
	if got := face.Metrics().Height.Ceil(); got > 32 {
		t.Fatalf("expected fitted face to fit height, got %d", got)
	}
}

func TestPixelFaceDrawsHardEdges(t *testing.T) {
	face := GetUdbFonts().PixelFace(PixelFaceOptions{
		Style:         PixelFontScoreboard,
		Scale:         3,
		LetterSpacing: -1,
	})
	defer face.Close()

	img := image.NewRGBA(image.Rect(0, 0, 40, 24))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	d := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(0, 21),
	}
	d.DrawString("2")

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0xffff {
				t.Fatalf("unexpected alpha at %d,%d: %d", x, y, a)
			}

			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			if c.R != 0 && c.R != 255 {
				t.Fatalf("expected hard red channel at %d,%d, got %d", x, y, c.R)
			}
			if c.G != 0 && c.G != 255 {
				t.Fatalf("expected hard green channel at %d,%d, got %d", x, y, c.G)
			}
			if c.B != 0 && c.B != 255 {
				t.Fatalf("expected hard blue channel at %d,%d, got %d", x, y, c.B)
			}
		}
	}
}

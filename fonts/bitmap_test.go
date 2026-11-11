package fonts

import (
	"fmt"
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

func TestAutoStyleFallsBackToCondensedForLongText(t *testing.T) {
	style := PixelFontStyleForTextBounds("1220d 21:25:39", 52, 26)
	if style != PixelFontCondensed {
		t.Fatalf("expected condensed style, got %v", style)
	}

	face := GetUdbFonts().FitPixelFace("1220d 21:25:39", 52, 26, PixelFontAuto)
	defer face.Close()
	if got := font.MeasureString(face, "1220d 21:25:39").Ceil(); got > 52 {
		t.Fatalf("expected condensed face to fit width, got %d", got)
	}
}

func TestColonAndSpaceKeepBlinkingLayoutStable(t *testing.T) {
	for _, style := range []PixelFontStyle{PixelFontDotMatrix, PixelFontScoreboard, PixelFontCondensed} {
		for _, scale := range []int{1, 3} {
			for _, spacing := range []int{-1, 0, 2} {
				t.Run(fmt.Sprintf("style=%d/scale=%d/spacing=%d", style, scale, spacing), func(t *testing.T) {
					face := GetUdbFonts().PixelFace(PixelFaceOptions{
						Style: style, Scale: scale, LetterSpacing: spacing,
					})
					defer face.Close()

					colon, ok := face.GlyphAdvance(':')
					if !ok {
						t.Fatal("expected colon glyph")
					}
					space, ok := face.GlyphAdvance(' ')
					if !ok || colon != space {
						t.Fatalf("colon and space advances differ: %v and %v (space supported: %v)", colon, space, ok)
					}
					if lit, dark := font.MeasureString(face, "12:34:56"), font.MeasureString(face, "12 34 56"); lit != dark {
						t.Fatalf("blinking changed the line width: %v and %v", lit, dark)
					}
				})
			}
		}
	}
}

func TestAutoStyleForTextBounds(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		width  int
		height int
		want   PixelFontStyle
	}{
		{"label at exact bounds", "AB", 12, 7, PixelFontDotMatrix},
		{"label too wide", "AB", 11, 7, PixelFontCondensed},
		{"label too tall", "AB", 12, 6, PixelFontCondensed},
		{"clock at exact bounds", "12:34", 19, 7, PixelFontScoreboard},
		{"clock too wide", "12:34", 18, 7, PixelFontCondensed},
		{"clock too tall", "12:34", 19, 5, PixelFontCondensed},
		{"long countdown", "1220d 21:25:39", 52, 26, PixelFontCondensed},
		{"empty text", "", 0, 7, PixelFontDotMatrix},
		{"impossible bounds", "AB", 0, 0, PixelFontCondensed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PixelFontStyleForTextBounds(tt.text, tt.width, tt.height); got != tt.want {
				t.Fatalf("style = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFitPixelFaceScaleAndFallback(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		width      int
		height     int
		style      PixelFontStyle
		wantWidth  int
		wantHeight int
	}{
		{"scoreboard auto", "12:34", 64, 32, PixelFontAuto, 57, 21},
		{"label auto", "AB", 24, 20, PixelFontAuto, 24, 14},
		{"width fallback", "AB", 8, 7, PixelFontAuto, 8, 5},
		{"height fallback", "12:34", 19, 5, PixelFontAuto, 18, 5},
		{"condensed explicit", "AB", 16, 10, PixelFontCondensed, 16, 10},
		{"explicit style retained", "AB", 8, 5, PixelFontDotMatrix, 12, 7},
		{"cannot fit", "AB", 1, 1, PixelFontAuto, 8, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			face := GetUdbFonts().FitPixelFace(tt.text, tt.width, tt.height, tt.style)
			defer face.Close()
			if got := font.MeasureString(face, tt.text).Ceil(); got != tt.wantWidth {
				t.Fatalf("width = %d, want %d", got, tt.wantWidth)
			}
			if got := face.Metrics().Height.Ceil(); got != tt.wantHeight {
				t.Fatalf("height = %d, want %d", got, tt.wantHeight)
			}
		})
	}
}

func TestCondensedGlyphRendering(t *testing.T) {
	face := GetUdbFonts().PixelFace(PixelFaceOptions{Style: PixelFontCondensed, Scale: 2, LetterSpacing: -1})
	defer face.Close()
	for _, r := range "aA?☃" {
		rect, mask, offset, advance, ok := face.Glyph(fixed.P(0, 10), r)
		if !ok || rect != image.Rect(0, 0, 6, 10) || advance != fixed.I(8) {
			t.Fatalf("unexpected condensed glyph %q: rect=%v advance=%v ok=%v", r, rect, advance, ok)
		}
		want := []string{"010", "101", "111", "101", "101"}
		if r == '?' || r == '☃' {
			want = []string{"110", "001", "010", "000", "010"}
		}
		for y := 0; y < rect.Dy(); y++ {
			for x := 0; x < rect.Dx(); x++ {
				_, _, _, alpha := mask.At(offset.X+x, offset.Y+y).RGBA()
				wantAlpha := uint32(0)
				if want[y/2][x/2] == '1' {
					wantAlpha = 0xffff
				}
				if alpha != wantAlpha {
					t.Fatalf("glyph %q at %d,%d: alpha=%d, want %d", r, x, y, alpha, wantAlpha)
				}
			}
		}
	}
}

func TestBoldPixelFacesRenderFullTextAtIntegerScales(t *testing.T) {
	for _, style := range []PixelFontStyle{PixelFontSilkscreenBold, PixelFontPixelifySansBold} {
		t.Run(fmt.Sprintf("style=%d", style), func(t *testing.T) {
			one := GetUdbFonts().PixelFace(PixelFaceOptions{Style: style, Scale: 1, LetterSpacing: -1})
			defer one.Close()
			two := GetUdbFonts().PixelFace(PixelFaceOptions{Style: style, Scale: 2, LetterSpacing: -1})
			defer two.Close()

			if two.Metrics().Height != one.Metrics().Height*2 {
				t.Fatalf("scale 2 height %v, want twice scale 1 height %v", two.Metrics().Height, one.Metrics().Height)
			}
			if font.MeasureString(one, "N N").Ceil() <= font.MeasureString(one, "NN").Ceil() {
				t.Fatal("space does not advance between glyphs")
			}

			for _, r := range []rune{'N', 'n', '8', ':'} {
				bounds, mask, _, _, ok := one.Glyph(fixed.P(0, one.Metrics().Ascent.Ceil()), r)
				if !ok || bounds.Empty() || mask == nil {
					t.Fatalf("glyph %q is missing: bounds=%v", r, bounds)
				}
				for y := 0; y < mask.Bounds().Dy(); y++ {
					for x := 0; x < mask.Bounds().Dx(); x++ {
						alpha := color.AlphaModel.Convert(mask.At(x, y)).(color.Alpha).A
						if alpha != 0 && alpha != 255 {
							t.Fatalf("glyph %q contains a non-binary pixel at (%d,%d): %d", r, x, y, alpha)
						}
					}
				}
				if r != 'N' {
					continue
				}

				doubledBounds, doubledMask, _, _, ok := two.Glyph(fixed.P(0, two.Metrics().Ascent.Ceil()), r)
				if !ok || doubledBounds.Dx() != bounds.Dx()*2 || doubledBounds.Dy() != bounds.Dy()*2 {
					t.Fatalf("scale 2 N bounds %v, want twice %v", doubledBounds, bounds)
				}
				for y := 0; y < mask.Bounds().Dy(); y++ {
					for x := 0; x < mask.Bounds().Dx(); x++ {
						want := color.AlphaModel.Convert(mask.At(x, y)).(color.Alpha).A
						for sy := 0; sy < 2; sy++ {
							for sx := 0; sx < 2; sx++ {
								got := color.AlphaModel.Convert(doubledMask.At(x*2+sx, y*2+sy)).(color.Alpha).A
								if got != want {
									t.Fatalf("scale 2 changed a pixel block at (%d,%d)", x, y)
								}
							}
						}
					}
				}
			}
		})
	}
}

func TestBoldPixelFaceSpacingAndFit(t *testing.T) {
	for _, style := range []PixelFontStyle{PixelFontSilkscreenBold, PixelFontPixelifySansBold} {
		t.Run(fmt.Sprintf("style=%d", style), func(t *testing.T) {
			defaultSpacing := GetUdbFonts().PixelFace(PixelFaceOptions{Style: style, Scale: 2, LetterSpacing: -1})
			defer defaultSpacing.Close()
			packed := GetUdbFonts().PixelFace(PixelFaceOptions{Style: style, Scale: 2, LetterSpacing: 0})
			defer packed.Close()
			spaced := GetUdbFonts().PixelFace(PixelFaceOptions{Style: style, Scale: 2, LetterSpacing: 1})
			defer spaced.Close()

			if got, want := font.MeasureString(defaultSpacing, "Nn8"), font.MeasureString(packed, "Nn8"); got != want {
				t.Fatalf("default and zero extra spacing differ: %v != %v", got, want)
			}
			if got, want := font.MeasureString(spaced, "Nn8")-font.MeasureString(packed, "Nn8"), fixed.I(6); got != want {
				t.Fatalf("positive spacing adds %v, want %v", got, want)
			}

			face := GetUdbFonts().FitPixelFace("Nn8:04", 128, 64, style)
			defer face.Close()
			boldFace, ok := face.(*boldPixelFontFace)
			if !ok || boldFace.scale < 1 {
				t.Fatalf("fit returned %T with invalid scale", face)
			}
			if got := font.MeasureString(face, "Nn8:04").Ceil(); got > 128 {
				t.Fatalf("fitted line width %d exceeds 128", got)
			}
			if got := face.Metrics().Height.Ceil(); got > 64 {
				t.Fatalf("fitted line height %d exceeds 64", got)
			}
			nextScale := GetUdbFonts().PixelFace(PixelFaceOptions{Style: style, Scale: boldFace.scale + 1})
			defer nextScale.Close()
			if font.MeasureString(nextScale, "Nn8:04").Ceil() <= 128 && nextScale.Metrics().Height.Ceil() <= 64 {
				t.Fatalf("fit stopped at scale %d although scale %d also fits", boldFace.scale, boldFace.scale+1)
			}

			fallback := GetUdbFonts().FitPixelFace("N", 1, 1, style)
			defer fallback.Close()
			if fallback.(*boldPixelFontFace).scale != 1 {
				t.Fatalf("impossible bounds returned scale %d, want scale 1", fallback.(*boldPixelFontFace).scale)
			}
		})
	}
}

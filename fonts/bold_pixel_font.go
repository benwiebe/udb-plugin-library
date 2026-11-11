package fonts

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	boldPixelFontsOnce sync.Once
	boldPixelFonts     map[PixelFontStyle]*opentype.Font
	boldPixelFontsErr  error
)

func isBoldPixelFontStyle(style PixelFontStyle) bool {
	return style == PixelFontSilkscreenBold || style == PixelFontPixelifySansBold
}

func boldPixelFont(style PixelFontStyle) (*opentype.Font, float64) {
	boldPixelFontsOnce.Do(func() {
		boldPixelFonts = make(map[PixelFontStyle]*opentype.Font)
		for style, path := range map[PixelFontStyle]string{
			PixelFontSilkscreenBold:   "assets/Silkscreen/Silkscreen-Bold.ttf",
			PixelFontPixelifySansBold: "assets/Pixelify_Sans/PixelifySans-Bold.ttf",
		} {
			fontBytes, err := fontAssets.ReadFile(path)
			if err != nil {
				boldPixelFontsErr = err
				return
			}
			parsed, err := opentype.Parse(fontBytes)
			if err != nil {
				boldPixelFontsErr = err
				return
			}
			boldPixelFonts[style] = parsed
		}
	})
	if boldPixelFontsErr != nil {
		panic(boldPixelFontsErr)
	}

	size := 8.0
	if style == PixelFontPixelifySansBold {
		size = 10
	}
	parsed, ok := boldPixelFonts[style]
	if !ok {
		panic(fmt.Errorf("unsupported bold pixel font style %d", style))
	}
	return parsed, size
}

func newBoldPixelFontFace(style PixelFontStyle, scale, letterSpacing int) font.Face {
	if scale < 1 {
		scale = 1
	}
	parsed, size := boldPixelFont(style)
	base, err := opentype.NewFace(parsed, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	if letterSpacing < 0 {
		letterSpacing = 0
	}
	return &boldPixelFontFace{base: base, scale: scale, letterSpacing: letterSpacing}
}

func fitBoldPixelFontFace(style PixelFontStyle, text string, maxWidth, maxHeight int) font.Face {
	var best font.Face
	for scale := 1; scale <= maxHeight; scale++ {
		face := newBoldPixelFontFace(style, scale, 0)
		if face.Metrics().Height.Ceil() > maxHeight || font.MeasureString(face, text).Ceil() > maxWidth {
			face.Close()
			break
		}
		if best != nil {
			best.Close()
		}
		best = face
	}
	if best != nil {
		return best
	}
	return newBoldPixelFontFace(style, 1, 0)
}

type boldPixelFontFace struct {
	base          font.Face
	scale         int
	letterSpacing int
}

func (f *boldPixelFontFace) Glyph(dot fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	bounds, mask, maskp, advance, ok := f.base.Glyph(fixed.Point26_6{}, r)
	if !ok {
		return image.Rectangle{}, nil, image.Point{}, 0, false
	}
	width, height := bounds.Dx(), bounds.Dy()
	scaledMask := image.NewAlpha(image.Rect(0, 0, width*f.scale, height*f.scale))
	if mask != nil {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				alpha := color.AlphaModel.Convert(mask.At(maskp.X+x, maskp.Y+y)).(color.Alpha).A
				if alpha < 128 {
					continue
				}
				for sy := y * f.scale; sy < (y+1)*f.scale; sy++ {
					for sx := x * f.scale; sx < (x+1)*f.scale; sx++ {
						scaledMask.SetAlpha(sx, sy, color.Alpha{A: 255})
					}
				}
			}
		}
	}
	x, y := dot.X.Ceil(), dot.Y.Ceil()
	scaledBounds := image.Rect(
		x+bounds.Min.X*f.scale,
		y+bounds.Min.Y*f.scale,
		x+bounds.Max.X*f.scale,
		y+bounds.Max.Y*f.scale,
	)
	return scaledBounds, scaledMask, image.Point{}, advance*f.fixedScale() + fixed.I(f.letterSpacing*f.scale), true
}

func (f *boldPixelFontFace) GlyphBounds(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
	bounds, advance, ok := f.base.GlyphBounds(r)
	if !ok {
		return fixed.Rectangle26_6{}, 0, false
	}
	scale := f.fixedScale()
	return fixed.Rectangle26_6{
		Min: fixed.Point26_6{X: bounds.Min.X * scale, Y: bounds.Min.Y * scale},
		Max: fixed.Point26_6{X: bounds.Max.X * scale, Y: bounds.Max.Y * scale},
	}, advance*scale + fixed.I(f.letterSpacing*f.scale), true
}

func (f *boldPixelFontFace) GlyphAdvance(r rune) (fixed.Int26_6, bool) {
	advance, ok := f.base.GlyphAdvance(r)
	return advance*f.fixedScale() + fixed.I(f.letterSpacing*f.scale), ok
}

func (f *boldPixelFontFace) Kern(r0, r1 rune) fixed.Int26_6 {
	return f.base.Kern(r0, r1) * f.fixedScale()
}

func (f *boldPixelFontFace) Metrics() font.Metrics {
	metrics := f.base.Metrics()
	scale := f.fixedScale()
	metrics.Height *= scale
	metrics.Ascent *= scale
	metrics.Descent *= scale
	metrics.XHeight *= scale
	metrics.CapHeight *= scale
	metrics.CaretSlope.X *= f.scale
	metrics.CaretSlope.Y *= f.scale
	return metrics
}

func (f *boldPixelFontFace) fixedScale() fixed.Int26_6 {
	return fixed.Int26_6(f.scale)
}

func (f *boldPixelFontFace) Close() error { return f.base.Close() }

package fonts

import (
	"image"
	"image/color"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// PixelFontStyle selects a hard-edged bitmap font.
type PixelFontStyle int

const (
	// PixelFontAuto lets FitPixelFace choose a style from the text and bounds.
	PixelFontAuto PixelFontStyle = iota
	// PixelFontDotMatrix is a 5×7 uppercase text font.
	PixelFontDotMatrix
	// PixelFontScoreboard is a 3×7 seven-segment numeric font.
	PixelFontScoreboard
	// PixelFontCondensed is a 3×5 uppercase text font for small rectangles.
	PixelFontCondensed
)

// PixelFaceOptions configures a bitmap face.
type PixelFaceOptions struct {
	// Style selects the font. PixelFontAuto uses PixelFontDotMatrix here;
	// use FitPixelFace to select a style from the text and available space.
	Style PixelFontStyle
	// Scale is an integer pixel multiplier, clamped to at least 1.
	Scale int
	// LetterSpacing is an extra gap in unscaled pixels after each glyph,
	// including the last. A negative value uses the font's default spacing;
	// zero packs glyphs together without an extra gap.
	LetterSpacing int
}

type bitmapGlyph struct {
	rows []string
}

type bitmapFont struct {
	height        int
	defaultWidth  int
	defaultGap    int
	glyphs        map[rune]bitmapGlyph
	uppercaseOnly bool
	scoreboard    bool
}

type bitmapFace struct {
	font              bitmapFont
	scale             int
	letterSpacing     int
	useDefaultSpacing bool
}

// PixelFace returns a bitmap face using the requested style, scale, and spacing.
// Text faces map ASCII lowercase letters to uppercase; unsupported runes use ?.
func (f *UdbFonts) PixelFace(opts PixelFaceOptions) font.Face {
	scale := opts.Scale
	if scale < 1 {
		scale = 1
	}

	bf := bitmapFontForStyle(opts.Style)

	letterSpacing := opts.LetterSpacing
	useDefaultSpacing := false
	if letterSpacing < 0 {
		letterSpacing = 0
		useDefaultSpacing = true
	}

	return &bitmapFace{
		font:              bf,
		scale:             scale,
		letterSpacing:     letterSpacing,
		useDefaultSpacing: useDefaultSpacing,
	}
}

// FitPixelFace returns the largest integer-scaled face that fits the text within
// maxWidth and maxHeight, using default letter spacing. PixelFontAuto selects a
// style from the text and bounds. Explicit styles never fall back. If the text
// cannot fit even at scale 1, the returned face remains at scale 1; callers can
// measure it to decide whether to wrap, clip, or shorten the text.
func (f *UdbFonts) FitPixelFace(text string, maxWidth, maxHeight int, style PixelFontStyle) font.Face {
	if style == PixelFontAuto {
		style = PixelFontStyleForTextBounds(text, maxWidth, maxHeight)
	}

	bf := bitmapFontForStyle(style)

	scale := 1
	for candidate := 1; candidate <= maxHeight; candidate++ {
		face := &bitmapFace{font: bf, scale: candidate, useDefaultSpacing: true}
		metrics := face.Metrics()
		if font.MeasureString(face, text).Ceil() > maxWidth || (metrics.Ascent+metrics.Descent).Ceil() > maxHeight {
			break
		}
		scale = candidate
	}

	return &bitmapFace{font: bf, scale: scale, useDefaultSpacing: true}
}

// PixelFontStyleForText selects scoreboard for digits with numeric separators,
// and dot-matrix for other text, without considering available space.
func PixelFontStyleForText(text string) PixelFontStyle {
	if isScoreboardText(text) {
		return PixelFontScoreboard
	}

	return PixelFontDotMatrix
}

// PixelFontStyleForTextBounds selects the preferred text style, falling back to
// condensed if that style exceeds either bound at scale 1 with default spacing.
// Condensed is returned even if the text cannot fit in the supplied bounds.
func PixelFontStyleForTextBounds(text string, maxWidth, maxHeight int) PixelFontStyle {
	style := PixelFontStyleForText(text)
	face := &bitmapFace{font: bitmapFontForStyle(style), scale: 1, useDefaultSpacing: true}
	metrics := face.Metrics()
	if font.MeasureString(face, text).Ceil() > maxWidth || (metrics.Ascent+metrics.Descent).Ceil() > maxHeight {
		return PixelFontCondensed
	}

	return style
}

func bitmapFontForStyle(style PixelFontStyle) bitmapFont {
	switch style {
	case PixelFontScoreboard:
		return scoreboard7SegmentFont()
	case PixelFontCondensed:
		return condensed3x5Font()
	default:
		return dotMatrix5x7Font()
	}
}

func (f *bitmapFace) Glyph(dot fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	glyph, ok := f.glyph(r)
	if !ok {
		return image.Rectangle{}, nil, image.Point{}, 0, false
	}

	width := glyphWidth(glyph)
	scaledWidth := width * f.scale
	scaledHeight := f.font.height * f.scale
	x := dot.X.Ceil()
	y := dot.Y.Ceil() - scaledHeight

	mask := image.NewAlpha(image.Rect(0, 0, scaledWidth, scaledHeight))
	for row, pattern := range glyph.rows {
		for col, bit := range pattern {
			if bit != '1' {
				continue
			}
			cell := image.Rect(col*f.scale, row*f.scale, (col+1)*f.scale, (row+1)*f.scale)
			fillAlpha(mask, cell)
		}
	}

	advance := fixed.I(f.advancePixels(r, glyph) * f.scale)
	return image.Rect(x, y, x+scaledWidth, y+scaledHeight), mask, image.Point{}, advance, true
}

func (f *bitmapFace) GlyphBounds(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
	glyph, ok := f.glyph(r)
	if !ok {
		return fixed.Rectangle26_6{}, 0, false
	}

	width := glyphWidth(glyph) * f.scale
	height := f.font.height * f.scale
	advance := fixed.I(f.advancePixels(r, glyph) * f.scale)
	return fixed.R(0, -height, width, 0), advance, true
}

func (f *bitmapFace) GlyphAdvance(r rune) (fixed.Int26_6, bool) {
	_, advance, ok := f.GlyphBounds(r)
	return advance, ok
}

func (f *bitmapFace) Kern(r0, r1 rune) fixed.Int26_6 {
	return 0
}

func (f *bitmapFace) Metrics() font.Metrics {
	ascent := f.font.height * f.scale
	return font.Metrics{
		Height:  fixed.I(ascent),
		Ascent:  fixed.I(ascent),
		Descent: 0,
	}
}

func (f *bitmapFace) Close() error {
	return nil
}

func (f *bitmapFace) glyph(r rune) (bitmapGlyph, bool) {
	if f.font.uppercaseOnly && r >= 'a' && r <= 'z' {
		r -= 'a' - 'A'
	}

	glyph, ok := f.font.glyphs[r]
	if ok {
		return glyph, true
	}

	glyph, ok = f.font.glyphs['?']
	return glyph, ok
}

func (f *bitmapFace) advancePixels(r rune, glyph bitmapGlyph) int {
	spacing := f.letterSpacing
	if f.useDefaultSpacing {
		spacing = f.font.defaultGap
		if f.font.scoreboard {
			switch r {
			case ':', '.', ' ', '/', '+', '-':
				spacing = 0
			}
		}
	}

	return glyphWidth(glyph) + spacing
}

func glyphWidth(g bitmapGlyph) int {
	width := 0
	for _, row := range g.rows {
		if len(row) > width {
			width = len(row)
		}
	}
	return width
}

func fillAlpha(img *image.Alpha, rect image.Rectangle) {
	drawRect := rect.Intersect(img.Bounds())
	for y := drawRect.Min.Y; y < drawRect.Max.Y; y++ {
		for x := drawRect.Min.X; x < drawRect.Max.X; x++ {
			img.SetAlpha(x, y, color.Alpha{A: 255})
		}
	}
}

func isScoreboardText(text string) bool {
	hasDigit := false
	for _, r := range text {
		if r >= '0' && r <= '9' {
			hasDigit = true
			continue
		}
		switch r {
		case ' ', ':', '.', '-', '/', '+':
			continue
		default:
			return false
		}
	}
	return hasDigit
}

func glyph(rows ...string) bitmapGlyph {
	return bitmapGlyph{rows: rows}
}

func dotMatrix5x7Font() bitmapFont {
	return bitmapFont{
		height:        7,
		defaultWidth:  5,
		defaultGap:    1,
		uppercaseOnly: true,
		glyphs: map[rune]bitmapGlyph{
			' ':  glyph("000", "000", "000", "000", "000", "000", "000"),
			'!':  glyph("1", "1", "1", "1", "1", "0", "1"),
			'"':  glyph("101", "101", "000", "000", "000", "000", "000"),
			'#':  glyph("01010", "11111", "01010", "01010", "11111", "01010", "01010"),
			'%':  glyph("11001", "11010", "00100", "01000", "10110", "00110", "00000"),
			'&':  glyph("01100", "10010", "10100", "01000", "10101", "10010", "01101"),
			'\'': glyph("1", "1", "0", "0", "0", "0", "0"),
			'(':  glyph("01", "10", "10", "10", "10", "10", "01"),
			')':  glyph("10", "01", "01", "01", "01", "01", "10"),
			'*':  glyph("00000", "10101", "01110", "11111", "01110", "10101", "00000"),
			'+':  glyph("00000", "00100", "00100", "11111", "00100", "00100", "00000"),
			',':  glyph("00", "00", "00", "00", "00", "10", "10"),
			'-':  glyph("00000", "00000", "00000", "11111", "00000", "00000", "00000"),
			'.':  glyph("0", "0", "0", "0", "0", "0", "1"),
			'/':  glyph("00001", "00010", "00100", "01000", "10000", "00000", "00000"),
			':':  glyph("000", "010", "000", "000", "000", "010", "000"),
			';':  glyph("0", "1", "0", "0", "0", "1", "1"),
			'?':  glyph("01110", "10001", "00001", "00010", "00100", "00000", "00100"),
			'0':  glyph("01110", "10001", "10011", "10101", "11001", "10001", "01110"),
			'1':  glyph("00100", "01100", "00100", "00100", "00100", "00100", "01110"),
			'2':  glyph("01110", "10001", "00001", "00010", "00100", "01000", "11111"),
			'3':  glyph("11110", "00001", "00001", "01110", "00001", "00001", "11110"),
			'4':  glyph("00010", "00110", "01010", "10010", "11111", "00010", "00010"),
			'5':  glyph("11111", "10000", "10000", "11110", "00001", "00001", "11110"),
			'6':  glyph("01110", "10000", "10000", "11110", "10001", "10001", "01110"),
			'7':  glyph("11111", "00001", "00010", "00100", "01000", "01000", "01000"),
			'8':  glyph("01110", "10001", "10001", "01110", "10001", "10001", "01110"),
			'9':  glyph("01110", "10001", "10001", "01111", "00001", "00001", "01110"),
			'A':  glyph("01110", "10001", "10001", "11111", "10001", "10001", "10001"),
			'B':  glyph("11110", "10001", "10001", "11110", "10001", "10001", "11110"),
			'C':  glyph("01111", "10000", "10000", "10000", "10000", "10000", "01111"),
			'D':  glyph("11110", "10001", "10001", "10001", "10001", "10001", "11110"),
			'E':  glyph("11111", "10000", "10000", "11110", "10000", "10000", "11111"),
			'F':  glyph("11111", "10000", "10000", "11110", "10000", "10000", "10000"),
			'G':  glyph("01111", "10000", "10000", "10011", "10001", "10001", "01111"),
			'H':  glyph("10001", "10001", "10001", "11111", "10001", "10001", "10001"),
			'I':  glyph("111", "010", "010", "010", "010", "010", "111"),
			'J':  glyph("00111", "00010", "00010", "00010", "00010", "10010", "01100"),
			'K':  glyph("10001", "10010", "10100", "11000", "10100", "10010", "10001"),
			'L':  glyph("10000", "10000", "10000", "10000", "10000", "10000", "11111"),
			'M':  glyph("10001", "11011", "10101", "10101", "10001", "10001", "10001"),
			'N':  glyph("10001", "11001", "10101", "10011", "10001", "10001", "10001"),
			'O':  glyph("01110", "10001", "10001", "10001", "10001", "10001", "01110"),
			'P':  glyph("11110", "10001", "10001", "11110", "10000", "10000", "10000"),
			'Q':  glyph("01110", "10001", "10001", "10001", "10101", "10010", "01101"),
			'R':  glyph("11110", "10001", "10001", "11110", "10100", "10010", "10001"),
			'S':  glyph("01111", "10000", "10000", "01110", "00001", "00001", "11110"),
			'T':  glyph("11111", "00100", "00100", "00100", "00100", "00100", "00100"),
			'U':  glyph("10001", "10001", "10001", "10001", "10001", "10001", "01110"),
			'V':  glyph("10001", "10001", "10001", "10001", "10001", "01010", "00100"),
			'W':  glyph("10001", "10001", "10001", "10101", "10101", "10101", "01010"),
			'X':  glyph("10001", "10001", "01010", "00100", "01010", "10001", "10001"),
			'Y':  glyph("10001", "10001", "01010", "00100", "00100", "00100", "00100"),
			'Z':  glyph("11111", "00001", "00010", "00100", "01000", "10000", "11111"),
		},
	}
}

func scoreboard7SegmentFont() bitmapFont {
	return bitmapFont{
		height:       7,
		defaultWidth: 3,
		defaultGap:   1,
		scoreboard:   true,
		glyphs: map[rune]bitmapGlyph{
			' ': glyph("000", "000", "000", "000", "000", "000", "000"),
			':': glyph("000", "010", "010", "000", "010", "010", "000"),
			'.': glyph("0", "0", "0", "0", "0", "0", "1"),
			'-': glyph("000", "000", "000", "111", "000", "000", "000"),
			'/': glyph("001", "001", "010", "010", "100", "100", "000"),
			'+': glyph("000", "010", "010", "111", "010", "010", "000"),
			'?': glyph("111", "001", "001", "111", "100", "000", "100"),
			'0': glyph("111", "101", "101", "101", "101", "101", "111"),
			'1': glyph("010", "110", "010", "010", "010", "010", "111"),
			'2': glyph("111", "001", "001", "111", "100", "100", "111"),
			'3': glyph("111", "001", "001", "111", "001", "001", "111"),
			'4': glyph("101", "101", "101", "111", "001", "001", "001"),
			'5': glyph("111", "100", "100", "111", "001", "001", "111"),
			'6': glyph("111", "100", "100", "111", "101", "101", "111"),
			'7': glyph("111", "001", "001", "010", "010", "100", "100"),
			'8': glyph("111", "101", "101", "111", "101", "101", "111"),
			'9': glyph("111", "101", "101", "111", "001", "001", "111"),
		},
	}
}

func condensed3x5Font() bitmapFont {
	return bitmapFont{
		height:        5,
		defaultWidth:  3,
		defaultGap:    1,
		uppercaseOnly: true,
		glyphs: map[rune]bitmapGlyph{
			' ':  glyph("0", "0", "0", "0", "0"),
			'!':  glyph("1", "1", "1", "0", "1"),
			'"':  glyph("101", "101", "000", "000", "000"),
			'#':  glyph("101", "111", "101", "111", "101"),
			'%':  glyph("101", "001", "010", "100", "101"),
			'&':  glyph("010", "101", "010", "101", "011"),
			'\'': glyph("1", "1", "0", "0", "0"),
			'(':  glyph("01", "10", "10", "10", "01"),
			')':  glyph("10", "01", "01", "01", "10"),
			'*':  glyph("101", "010", "111", "010", "101"),
			'+':  glyph("000", "010", "111", "010", "000"),
			',':  glyph("0", "0", "0", "1", "1"),
			'-':  glyph("000", "000", "111", "000", "000"),
			'.':  glyph("0", "0", "0", "0", "1"),
			'/':  glyph("001", "001", "010", "100", "100"),
			':':  glyph("0", "1", "0", "1", "0"),
			';':  glyph("0", "1", "0", "1", "1"),
			'?':  glyph("110", "001", "010", "000", "010"),
			'0':  glyph("111", "101", "101", "101", "111"),
			'1':  glyph("010", "110", "010", "010", "111"),
			'2':  glyph("111", "001", "111", "100", "111"),
			'3':  glyph("111", "001", "111", "001", "111"),
			'4':  glyph("101", "101", "111", "001", "001"),
			'5':  glyph("111", "100", "111", "001", "111"),
			'6':  glyph("111", "100", "111", "101", "111"),
			'7':  glyph("111", "001", "010", "100", "100"),
			'8':  glyph("111", "101", "111", "101", "111"),
			'9':  glyph("111", "101", "111", "001", "111"),
			'A':  glyph("010", "101", "111", "101", "101"),
			'B':  glyph("110", "101", "110", "101", "110"),
			'C':  glyph("011", "100", "100", "100", "011"),
			'D':  glyph("110", "101", "101", "101", "110"),
			'E':  glyph("111", "100", "110", "100", "111"),
			'F':  glyph("111", "100", "110", "100", "100"),
			'G':  glyph("011", "100", "101", "101", "011"),
			'H':  glyph("101", "101", "111", "101", "101"),
			'I':  glyph("111", "010", "010", "010", "111"),
			'J':  glyph("001", "001", "001", "101", "010"),
			'K':  glyph("101", "101", "110", "101", "101"),
			'L':  glyph("100", "100", "100", "100", "111"),
			'M':  glyph("101", "111", "111", "101", "101"),
			'N':  glyph("101", "111", "111", "111", "101"),
			'O':  glyph("111", "101", "101", "101", "111"),
			'P':  glyph("110", "101", "110", "100", "100"),
			'Q':  glyph("111", "101", "101", "111", "001"),
			'R':  glyph("110", "101", "110", "101", "101"),
			'S':  glyph("011", "100", "010", "001", "110"),
			'T':  glyph("111", "010", "010", "010", "010"),
			'U':  glyph("101", "101", "101", "101", "111"),
			'V':  glyph("101", "101", "101", "101", "010"),
			'W':  glyph("101", "101", "111", "111", "101"),
			'X':  glyph("101", "101", "010", "101", "101"),
			'Y':  glyph("101", "101", "010", "010", "010"),
			'Z':  glyph("111", "001", "010", "100", "111"),
		},
	}
}

func normalizeGlyphRows(rows []string, width, height int) []string {
	out := make([]string, height)
	for i := 0; i < height; i++ {
		if i < len(rows) {
			out[i] = rows[i]
		}
		if len(out[i]) < width {
			out[i] += strings.Repeat("0", width-len(out[i]))
		}
	}
	return out
}

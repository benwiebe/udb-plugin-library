# Changelog

## v0.7.0

- Release the bitmap font APIs: `UdbFonts.PixelFace`, `UdbFonts.FitPixelFace`,
  `PixelFaceOptions`, `PixelFontAuto`, `PixelFontDotMatrix`, `PixelFontScoreboard`,
  `PixelFontCondensed`, `PixelFontStyleForText`, and `PixelFontStyleForTextBounds`.
  The base bitmap implementation was committed in `866d0e9` after v0.6.0.
- Add a 3×5 condensed font. Automatic selection falls back when the preferred
  numeric or text face exceeds either bound at scale 1.
- Match colon and space advances in all bitmap styles for stable blinking.
- Document scaling, spacing, character coverage, and behavior when text cannot
  fit at scale 1.

Validation: library tests, race tests, and `go vet` passed. Race tests used Go
1.26.1 because the installed Go 1.27.1 failed to build `runtime/race`.

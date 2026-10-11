# udb-plugin-library

Interface library for UDB (Universal Display Board) plugins. Import this in your plugin to get the type definitions required by `udb-core`.

## Key Interfaces

### Board

Boards render content to the display. Every board implements the base `Board` interface plus one of three render interfaces:

| Type | Interface | Render returns | Use when |
|------|-----------|---------------|----------|
| `BoardTypeStatic` | `StaticBoard` | `image.Image` | Fixed image; no updates while displayed |
| `BoardTypeAnimated` | `AnimatedBoard` | `[]AnimationFrame` | Pre-baked frame sequence |
| `BoardTypeDynamic` | `DynamicBoard` | `AnimationFrame` | Live-updating data (clocks, scores) |

`Init(config, datasource, dimensions)` is called once at startup. Boards should pre-compute any layout values that depend on display dimensions here. `Render()` takes no parameters.

### Datasource

Datasources provide data to boards. `GetData() any` must always return immediately — it is called on the render path. Do the type assertion once in `Init()` and store the concrete type rather than asserting on every `Render()` call.

`Init(config)` is called after the datasource instance is constructed so it can parse its own JSON config. `Start(ctx)` is called at startup so datasources can launch background fetch goroutines. The goroutine should exit when `ctx` is cancelled. `DataChanged()` returns a channel the datasource can signal to trigger an immediate re-render; return `nil` for no push notifications.

## Fonts

The library also includes a `fonts` package for LED matrix and other low-resolution text rendering.

- `FitPixelFace(text, maxWidth, maxHeight, PixelFontAuto)` returns a hard-edged bitmap face scaled to fit the target rectangle.
- `PixelFontAuto` selects 3×7 scoreboard-style numerals for numeric text such as clocks and scores, and a 5×7 dot-matrix face for general labels. If the preferred face exceeds either bound at scale 1, it falls back to the 3×5 `PixelFontCondensed` face.
- `PixelFontStyleForTextBounds(text, maxWidth, maxHeight)` exposes that choice for custom layouts. `PixelFontStyleForText(text)` selects only by text, without considering bounds.
- `PixelFace(PixelFaceOptions{Style: PixelFontCondensed, Scale: 2, LetterSpacing: -1})` creates an explicitly styled face. Scale is an integer multiplier (minimum 1); negative letter spacing uses the font default, zero removes the gap, and positive values add that many unscaled pixels after every glyph, including the last. `PixelFace` treats `PixelFontAuto` as dot-matrix because it has no text or bounds to inspect.
- `FitPixelFace` uses default spacing and the largest integer scale that fits. Explicit styles never fall back. When no face fits at scale 1, it still returns a scale-1 face: measure the result before deciding to wrap, clip, or shorten the text.
- Colons and spaces have matching advances in all bitmap styles, so replacing `:` with a space for blinking keeps the text width and subsequent glyph positions stable.
- Text faces map ASCII lowercase letters to uppercase. Unsupported characters render as `?`.
- The bitmap faces implement `golang.org/x/image/font.Face`, so plugin authors can use `font.MeasureString` and `font.Drawer` normally.
- The older TrueType helpers remain available for boards that explicitly want outline fonts, but they are no longer the preferred path for pixel-board UI text.
- Third-party font license texts are included under `fonts/assets/*/OFL.txt`.

The bitmap APIs above are available starting in **v0.7.0**. Plugins can select that version with:

```sh
go get github.com/benwiebe/udb-plugin-library@v0.7.0
```

Only the font files currently used by the package are bundled, to keep module size down.

### Plugin Registration

Plugins call `Register()` from an `init()` function. The core picks them up automatically when the package is blank-imported.

```go
func init() {
    udb_plugin_library.Register(&MyPlugin{})
}
```

## See Also

- [udb-core](https://github.com/benwiebe/udb-core) — the runtime and plugin authoring guide (`docs/PLUGIN_AUTHORING.md`)
- [udb-plugin-samples](https://github.com/benwiebe/udb-plugin-samples) — reference implementations

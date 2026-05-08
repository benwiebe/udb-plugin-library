# udb-plugin-library

Interface library for UDB (Universal Display Board) plugins. Import this in your plugin to get the type definitions required by `udb-core`.

## Key Interfaces

### Board

Boards render content to the display. Every board implements the base `Board[T]` interface plus one of three render interfaces:

| Type | Interface | Render returns | Use when |
|------|-----------|---------------|----------|
| `BoardTypeStatic` | `StaticBoard[T]` | `image.Image` | Fixed image; no updates while displayed |
| `BoardTypeAnimated` | `AnimatedBoard[T]` | `[]AnimationFrame` | Pre-baked frame sequence |
| `BoardTypeDynamic` | `DynamicBoard[T]` | `AnimationFrame` | Live-updating data (clocks, scores) |

`Init(config, datasource, dimensions)` is called once at startup. Boards should pre-compute any layout values that depend on display dimensions here. `Render()` takes no parameters.

### Datasource

Datasources provide data to boards. `GetData()` must always return immediately — it is called on the render path.

`Start(ctx)` is called at startup so datasources can launch background fetch goroutines. The goroutine should exit when `ctx` is cancelled. `DataChanged()` returns a channel the datasource can signal to trigger an immediate re-render; return `nil` for no push notifications.

## See Also

- [udb-core](https://github.com/benwiebe/udb-core) — the runtime and plugin authoring guide (`docs/PLUGIN_AUTHORING.md`)
- [udb-plugin-samples](https://github.com/benwiebe/udb-plugin-samples) — reference implementations

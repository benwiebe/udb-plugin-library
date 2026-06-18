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

`Start(ctx)` is called at startup so datasources can launch background fetch goroutines. The goroutine should exit when `ctx` is cancelled. `DataChanged()` returns a channel the datasource can signal to trigger an immediate re-render; return `nil` for no push notifications.

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

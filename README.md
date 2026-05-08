# udb-plugin-library

Interface contracts for [UDB (Universal Display Board)](https://github.com/benwiebe/udb-core) plugins. This library defines what every plugin must implement — it contains no logic, only types and interfaces.

## Usage

Add it as a dependency in your plugin module:

```bash
go get github.com/benwiebe/udb-plugin-library
```

## What's In Here

| Package | Contents |
|---------|----------|
| Root package | `UdbPlugin`, `UdbBoardPlugin`, `UdbDatasourcePlugin`, `UdbCombinedPlugin` interfaces |
| `types` | `Board[T]`, `StaticBoard[T]`, `AnimatedBoard[T]`, `DynamicBoard[T]`, `Datasource[T]`, `BoardDimensions`, `AnimationFrame`, `PluginType`, `BoardType`, `PluginConfig` |

## Writing a Plugin

See the [Plugin Authoring Guide](https://github.com/benwiebe/udb-core/blob/main/docs/PLUGIN_AUTHORING.md) in `udb-core` for a complete walkthrough, including examples, build instructions, and config wiring.

## Interface Summary

```go
// Every plugin exports a Plugin variable of this type.
type UdbPlugin interface {
    GetId() string
    GetName() string
    GetPluginType() types.PluginType
    Configure(config types.PluginConfig) error
}

// Implement this to provide boards.
type UdbBoardPlugin interface {
    UdbPlugin
    GetBoardMap() map[string]types.Board[any]
    GetAllBoards() []types.Board[any]
}

// Implement this to provide datasources.
type UdbDatasourcePlugin interface {
    UdbPlugin
    GetDatasourceMap() map[string]types.Datasource[any]
    GetAllDatasources() []types.Datasource[any]
}

// Implement this to provide both.
type UdbCombinedPlugin interface {
    UdbBoardPlugin
    UdbDatasourcePlugin
}
```

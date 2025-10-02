package udb_plugin_library

import (
	"github.com/benwiebe/udb-plugin-library/types"
)

// UdbPlugin is an interface for all plugins' common methods.
type UdbPlugin interface {
	GetPluginType() types.PluginType
	Configure(config types.PluginConfig) error
}

// UdbBoardPlugin is an interface describing a plugin which only provides boards but
// does not provide any datasources.
type UdbBoardPlugin interface {
	UdbPlugin
	GetBoardMap() map[string]types.Board[any]
	GetAllBoards() []types.Board[any]
}

// UdbDatasourcePlugin is an interface describing a plugin which only provides datasources but
// does not provide any boards.
type UdbDatasourcePlugin interface {
	UdbPlugin
	GetDatasourceMap() map[string]types.Datasource[any]
	GetAllDatasources() []types.Datasource[any]
}

// UdbCombinedPlugin is an interface describing a plugin which provides both boards and datasources.
type UdbCombinedPlugin interface {
	UdbBoardPlugin
	UdbDatasourcePlugin
}

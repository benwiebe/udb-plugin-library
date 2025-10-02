package udb_plugin_library

import (
	"github.com/benwiebe/udb-plugin-library/types"
)

type UdbPlugin interface {
	GetPluginType() types.PluginType
	Configure(config types.PluginConfig) error
}

type UdbBoardPlugin interface {
	UdbPlugin
	GetBoardMap() map[string]types.Board
	GetAllBoards() []types.Board
}

type UdbDatasourcePlugin interface {
	UdbPlugin
	GetDatasourceMap() map[string]types.Datasource[any]
	GetAllDatasources() []types.Datasource[any]
}

type UdbCombinedPlugin interface {
	UdbBoardPlugin
	UdbDatasourcePlugin
}

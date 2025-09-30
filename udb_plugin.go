package udb_plugin_library

import (
	"github.com/benwiebe/udb_plugin_library/types"
)

type UdbPlugin interface {
	getPluginType() types.PluginType
	configure(config types.PluginConfig) error
}

type UdbBoardPlugin interface {
	UdbPlugin
	getBoardMap() map[string]types.Board
	getAllBoards() []types.Board
}

type UdbDatasourcePlugin interface {
	UdbPlugin
	getDatasourceMap() map[string]types.Datasource[any]
	getAllDatasources() []types.Datasource[any]
}

type UdbCombinedPlugin interface {
	UdbBoardPlugin
	UdbDatasourcePlugin
}

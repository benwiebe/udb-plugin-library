package types

type PluginType string

const (
	PluginTypeDatasource PluginType = "datasource"
	PluginTypeBoards     PluginType = "boards"
	PluginTypeCombined   PluginType = "combined"
)

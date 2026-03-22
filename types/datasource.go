package types

// Datasource is an interface for a datasource that will be provided by a plugin.
type Datasource[T any] interface {
	GetId() string
	GetName() string
	// GetType returns the type identifier for the data this datasource provides.
	// Type identifiers follow the convention "author/plugin-name/type-name",
	// e.g. "benwiebe/nhl-plugin/game-data". A board's GetDatasourceType() must
	// match the datasource's GetType() for them to be compatible.
	GetType() string
	GetData() T
}

package types

import "context"

// DatasourceFactory is a constructor that returns a fresh Datasource instance.
// Plugins return these from GetDatasourceMap() so the core can create one instance
// per datasource config entry.
type DatasourceFactory = func() Datasource

// Datasource is an interface for a datasource that will be provided by a plugin.
type Datasource interface {
	GetId() string
	GetName() string
	// GetType returns the type identifier for the data this datasource provides.
	// Type identifiers follow the convention "author/plugin-name/type-name",
	// e.g. "benwiebe/nhl-plugin/game-data". A board's GetDatasourceType() must
	// match the datasource's GetType() for them to be compatible.
	GetType() string
	GetData() any
	// Start is called once at startup. Datasources that poll external sources should
	// start their background goroutine here and stop it when ctx is cancelled.
	Start(ctx context.Context) error
	// DataChanged returns a channel that receives a value whenever the datasource's data
	// has changed and an immediate re-render is warranted. Return nil if the datasource
	// does not support push notifications; a nil channel blocks forever in a select,
	// which is the correct no-op behaviour.
	DataChanged() <-chan struct{}
}

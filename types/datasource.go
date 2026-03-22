package types

// Datasource is an interface for a datasource that will be provided by a plugin
type Datasource[T any] interface {
	GetId() string
	GetName() string
	GetData() T
}

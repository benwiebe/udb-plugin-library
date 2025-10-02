package types

// Datasource is a interface for a datasource that will be provided by a plugin
type Datasource[T any] interface {
	getId() string
	getName() string
	getData() T
}

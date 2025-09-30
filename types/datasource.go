package types

type Datasource[T any] interface {
	getId() string
	getName() string
	getData() T
}

package types

import (
	"image"
	"time"
)

// BoardType is a string to describe the type of a board to ensure that the correct type
// cast is used when creating a board. The supported types are: static, animated and dynamic
type BoardType string

const (
	BoardTypeStatic   BoardType = "static"
	BoardTypeAnimated BoardType = "animated"
	BoardTypeDynamic  BoardType = "dynamic"
)

// BoardDimensions is a struct to describe the dimensions of a board
type BoardDimensions struct {
	Width  int
	Height int
}

// Board is an interface for all boards' common methods. Board and its related interfaces
// are used to describe the different types of boards that are provided by a plugin and
// which can be displayed on the screen.
type Board[T any] interface {
	GetId() string
	GetName() string
	GetSupportedDimensions() []BoardDimensions
	GetType() BoardType
	Init(datasource Datasource[T]) error
}

// StaticBoard is an interface for a board which displays a static image
type StaticBoard[T any] interface {
	Board[T]
	Render() image.Image
}

// AnimatedBoard is an interface for a board which displays a pre-rendered animated image.
// If a board needs to update while it's being shown, then it should implement DynamicBoard instead.
type AnimatedBoard[T any] interface {
	Board[T]
	Render() Animation
}

// DynamicBoard is an interface for a board which displays content that may change while the board
// is being shown. If the board does not need to update while it's being shown, then it should implement
// AnimatedBoard instead.
type DynamicBoard[T any] interface {
	Board[T]
	Render() AnimationFrame
}

// Animation is an array of AnimationFrame that can be played on a board
type Animation = []AnimationFrame

// AnimationFrame is a struct that describes a single frame (image) of an animation and the duration
// it should be displayed for.
type AnimationFrame struct {
	Img      image.Image
	Duration time.Duration
}

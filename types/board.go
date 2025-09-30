package types

import "image"

type BoardDimensions struct {
	Width  int
	Height int
}

type Board interface {
	getId() string
	getName() string
	getSupportedDimensions() []BoardDimensions
	render() image.Image
}

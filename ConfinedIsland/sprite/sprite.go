package sprite

import "image/color"

type Sprite struct {
	X, Y          float64
	Width, Height float64
	Color         color.RGBA
}

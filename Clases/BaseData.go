package Clases

import (
	"../Vectores"
	"image/color"
)

type BaseObject struct {
	Id          int
	Material    int
	Translation Vectores.Vector
	Color       color.RGBA
}

package Clases

import (
	"../Vectores"
	"image/color"
)

type Objeto interface {
	Distancia(Vectores.Vector) float64
	GetColor() color.RGBA
	GetMaterial() int
}

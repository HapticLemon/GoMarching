package Clases

import (
	"../Vectores"
	"image/color"
)

// Definición de la clase esfera.
// OJO, para poder exportarlo fuera del paquete ha de comenzar por
// mayúscula.
type Esfera struct {
	BaseObject
	Radio float64
}

func (e Esfera) Distancia(punto Vectores.Vector) float64 {
	return punto.Sub(e.Translation).Length() - e.Radio
}

func (e Esfera) GetColor() color.RGBA {
	return e.Color
}

func (e Esfera) GetMaterial() int {
	return e.Material
}

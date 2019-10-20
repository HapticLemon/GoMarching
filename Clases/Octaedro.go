package Clases

import (
	"../Vectores"
	"image/color"
)

// Definición de la clase esfera.
// OJO, para poder exportarlo fuera del paquete ha de comenzar por
// mayúscula.
type Octaedro struct {
	BaseObject
	Radio float64
}

func (o Octaedro) Distancia(punto Vectores.Vector) float64 {
	var puntoTrasladado Vectores.Vector
	var puntoAbsoluto Vectores.Vector

	puntoTrasladado = punto.Sub(o.Translation)
	puntoAbsoluto = puntoTrasladado.Abs()

	// Hay que hacer que use la cte definida y no esté hardcodeado.
	//
	return (puntoAbsoluto.X + puntoAbsoluto.Y + puntoAbsoluto.Z - o.Radio) * 0.57735027
}

func (o Octaedro) GetColor() color.RGBA {
	return o.Color
}

func (o Octaedro) GetMaterial() int {
	return o.Material
}

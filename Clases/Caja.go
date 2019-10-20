package Clases

import (
	"../Vectores"
	"image/color"
	"math"
)

// Definición de la clase esfera.
// OJO, para poder exportarlo fuera del paquete ha de comenzar por
// mayúscula.
type Caja struct {
	BaseObject
	Posicion    Vectores.Vector
	Dimensiones Vectores.Vector
}

func (c Caja) Distancia(punto Vectores.Vector) float64 {
	// Por implementar
	var puntoT Vectores.Vector

	puntoT = punto.Sub(c.Translation)

	/*	return max(max(abs(puntoT[0] - self.posicion[0]) - self.dimensiones[0],
		abs(puntoT[1] - self.posicion[1]) - self.dimensiones[1]),
		abs(puntoT[2] - self.posicion[2]) - self.dimensiones[2]);*/

	return math.Max(math.Max(math.Abs(puntoT.X-c.Posicion.X)-c.Dimensiones.X,
		math.Abs(puntoT.Y-c.Posicion.Y)-c.Dimensiones.Y),
		math.Abs(puntoT.Z-c.Posicion.Z)-c.Dimensiones.Z)
}

func (c Caja) GetColor() color.RGBA {
	return c.Color
}

func (c Caja) GetMaterial() int {
	return c.Material
}

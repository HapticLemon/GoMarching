package Clases

import "../Vectores"

// Definición de la clase esfera.
// OJO, para poder exportarlo fuera del paquete ha de comenzar por
// mayúscula.
type Esfera struct {
	BaseObject
	Radio float64
}

func (e Esfera) Distancia(punto Vectores.Vector) float64 {
	return punto.Length() - 10
}
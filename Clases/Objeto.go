package Clases

import "../Vectores"

type Objeto interface {
	Distancia(Vectores.Vector) float64
}

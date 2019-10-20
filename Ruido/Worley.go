package Ruido

import (
	"../Vectores"
	"math"
	"math/rand"
	"sort"
)

// Distribución de Poisson
// No estaría mal encontrar la forma de generarla en automático.
//
var distP [10]int = [10]int{4, 4, 6, 5, 3, 4, 8, 8, 7, 5}

// Calculo una semilla diferente para cada uno de los cubos.
//
func calculateSeed(cube [3]int) int {
	var seed int

	seed = (541*cube[0] + 79*cube[1] + 31*cube[2]) % 4294967296
	return seed
}

// Devolvemos el número de puntos por cada cubo.
//
func pointNumber(seed int) int {
	var max int = 9
	var min int = 0
	var index int

	rand.Seed(int64(seed))

	index = rand.Intn(max-min) + min
	return distP[index]
}

// uniform returns a uniformly random float in [0,1).
// https://stackoverflow.com/users/5181219/ted
func uniform() float64 {
	var sig uint64

	sig = rand.Uint64() % (1 << 52)
	return (1 + float64(sig)/(1<<52)) / math.Pow(2, geometric())
}

// geometric returns a number picked from a geometric
// distribution of parameter 0.5.
// https://stackoverflow.com/users/5181219/ted
func geometric() float64 {
	var b float64 = 1
	for rand.Uint64()%2 == 0 {
		b++
	}
	return b
}

// Genero un punto 3d con coordenadas en el rango 0-1
// Les añado las del cubo para poder calcular la distancia al punto
// original, ya que las coordenadas de éste no están en el rango 0-1
//
func generatePoint(seed int, cube [3]int) [3]float64 {

	var point [3]float64

	point[0] = uniform() + float64(cube[0])
	point[1] = uniform() + float64(cube[1])
	point[2] = uniform() + float64(cube[2])

	return point
}

// Distancia euclídea entre dos puntos.
// Quizá exista dentro de alguna librería
//
func euclidean(punto Vectores.Vector, coord [3]float64) float64 {
	var distancia float64

	distancia = math.Sqrt(math.Pow((coord[0]-punto.X), 2) + math.Pow((coord[1]-punto.Y), 2) + math.Pow((coord[2]-punto.Z), 2))

	return distancia

}

func clip(valor float64, min float64, max float64) float64 {
	if valor < min {
		return min
	}
	if valor > max {
		return max
	}
	return valor
}

// Worley 3d con cubos.
// Referencias :
//   https://thebookofshaders.com/12/
//   http://www.rhythmiccanvas.com/research/papers/worley.pdf
//   https://github.com/bhickey/worley/blob/master/worley.c
//   https://www.kdnuggets.com/2017/08/comparing-distance-measurements-python-scipy.html
//
func Worley3D(punto Vectores.Vector) float64 {
	var minimo float64 = 1000
	var seed int
	var points int

	var cx int
	var cy int
	var cz int

	var cube [3]int
	var dummy [3]float64

	for cx = int(math.Floor(punto.X - 1)); cx <= int(math.Floor(punto.X+2)); cx++ {
		for cy = int(math.Floor(punto.Y - 1)); cy <= int(math.Floor(punto.Y+2)); cy++ {
			for cz = int(math.Floor(punto.Z - 1)); cz <= int(math.Floor(punto.Z+2)); cz++ {
				cube[0] = cx
				cube[1] = cy
				cube[2] = cz

				seed = calculateSeed(cube)
				points = pointNumber(seed)

				distancias := make([]float64, points)

				for cp := 0; cp < points; cp++ {
					dummy = generatePoint(seed, cube)
					distancias[cp] = euclidean(punto, dummy)
				}

				sort.Float64s(distancias)

				if distancias[1] < minimo {
					minimo = distancias[1]
				}
			}
		}
	}
	return clip(minimo, 0, 1)

}

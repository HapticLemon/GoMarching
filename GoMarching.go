package main

import (
	"./Clases"
	"./Vectores"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

// Para simplificar, considero una esfera de radio1 en el centro de coordenadas.
//
func distanciaEsfera(punto Vectores.Vector) float64 {
	// 10 es el radio de la esfera.
	return punto.Length() - 10
}

// Desde aquí llamaré a la función de distancia. De momento solemente será una esfera.
// Por implementar
//
func mapTheWorld(punto Vectores.Vector) float64 {
	// Distancia inicial arbitrariamente grande.
	//
	var distancia float64 = 1000
	var distanciaObjeto float64

	for _, elemento := range Objetos {
		distanciaObjeto = elemento.Distancia(punto)
		if distanciaObjeto < distancia {
			distancia = distanciaObjeto
		}
	}

	return distancia
}

// Cálculo de la normal (gradiente) en un punto.
//
func calculateNormal(punto Vectores.Vector) Vectores.Vector {
	var gradiente = Vectores.Vector{1, 0, 0}
	var EPSILON float64 = 0.01

	gradiente.X = mapTheWorld(Vectores.Vector{punto.X + EPSILON, punto.Y, punto.Z}) - mapTheWorld(Vectores.Vector{punto.X - EPSILON, punto.Y, punto.Z})
	gradiente.Y = mapTheWorld(Vectores.Vector{punto.X, punto.Y + EPSILON, punto.Z}) - mapTheWorld(Vectores.Vector{punto.X, punto.Y - EPSILON, punto.Z})
	gradiente.Z = mapTheWorld(Vectores.Vector{punto.X, punto.Y, punto.Z + EPSILON}) - mapTheWorld(Vectores.Vector{punto.X, punto.Y, punto.Z - EPSILON})

	gradiente.MultiplyByScalar(-1)
	return gradiente.Normalize()
}

// Cálculo de iluminación difusa
//
func ilumina(punto Vectores.Vector, diffuseIntensity float64, normal Vectores.Vector) color.RGBA {
	var luz Vectores.Vector
	var color = color.RGBA{0, 0, 0, 0}

	luz = COLOR.MultiplyByScalar(diffuseIntensity)

	color.R = uint8(luz.X)
	color.G = uint8(luz.Y)
	color.B = uint8(luz.Z)
	color.A = 255

	return color
}

func raymarch(ro Vectores.Vector, rd Vectores.Vector) color.RGBA {

	var punto Vectores.Vector
	var directionToLight Vectores.Vector
	var normal Vectores.Vector
	var t float64 = 0
	var diffuseIntensity float64 = 0
	var distancia float64 = 0
	var color = color.RGBA{0, 0, 0, 255}

	for x := 0; x < MAXSTEPS; x++ {
		punto = ro.Add(rd.MultiplyByScalar(t))
		distancia = mapTheWorld(punto)

		if distancia < MINIMUM_HIT_DISTANCE {
			directionToLight = punto.Sub(LIGHT).Normalize()
			normal = calculateNormal(punto)
			diffuseIntensity = math.Max(0.0, normal.Dot(directionToLight))
			color = ilumina(punto, diffuseIntensity, normal)
			return color
		}
		t += distancia
	}

	// Devuelvo el color negro de fondo.
	return color
}

func defineObjetos() {
	esfera_0 := Clases.Esfera{
		Clases.BaseObject{0, 0},
		2.0,
	}

	Objetos = append(Objetos, esfera_0)
}

func main() {
	//var V1 = Vectores.Vector{1,0,1}
	//var V2 = Vectores.Vector{1,2,0}

	//var v3 = V1.Add(V2).MultiplyByScalar(escalar)
	//fmt.Printf("Vectores.Vector %f, %f, %f", v3.X,v3.Y,v3.Z)

	var NDC_x float64
	var NDC_y float64
	var PixelScreen_x float64
	var PixelScreen_y float64
	var PixelCamera_x float64
	var PixelCamera_y float64

	var ro Vectores.Vector
	var rd Vectores.Vector
	var color color.RGBA

	defineObjetos()
	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	out, err := os.Create("./output.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for x := 0; x < WIDTH; x++ {
		for y := 0; y < WIDTH; y++ {
			// Hacemos las conversiones de espacios
			//
			NDC_x = (float64(x) + correccion) / float64(WIDTH)
			NDC_y = (float64(y) + correccion) / float64(HEIGHT)

			PixelScreen_x = 2*NDC_x - 1
			PixelScreen_y = 2*NDC_y - 1

			PixelCamera_x = PixelScreen_x * ImageAspectRatio
			PixelCamera_y = PixelScreen_y

			// Origen y dirección

			ro = EYE.Add(FORWARD.MultiplyByScalar(FL)).Add(RIGHT.MultiplyByScalar(PixelCamera_x)).Add(UP.MultiplyByScalar(PixelCamera_y))
			rd = ro.Sub(EYE).Normalize()

			if x == 139 && y == 97 {
				fmt.Println("GatetERlZ \n")
			}
			color = raymarch(ro, rd)

			img.Set(x, y, color)

		}
	}
	var opt jpeg.Options

	opt.Quality = 80
	// ok, write out the data into the new JPEG file

	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Generated image to output.jpg \n")
}

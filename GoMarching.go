package main

import (
	"./Clases"
	"./Vectores"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
)

type Escena struct {
	Escena []Objeto `json:"Escena"`
}

type Objeto struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	Material    string `json:"Material"`
	Radio       int    `json:"Radio"`
	Translation []int  `json:Translation`
	Color       []int  `json:Translation`
}

// En Go las enumeraciones se montan así.
const (
	NOMAT = iota
	WORLEY3D
	MARMOL
)

func cargaObjetos(escena Escena) {
	var id int = 0
	var material int = 0
	var translationRaw []int
	var translationVec = Vectores.Vector{0, 0, 0}
	var colorRaw []int
	var colorVec = color.RGBA{0, 0, 0, 0}

	for i := 0; i < len(escena.Escena); i++ {
		// Procesamos el tipo de material para que pase de string a enum.
		switch escena.Escena[i].Material {
		case "NOMAT":
			material = NOMAT
		case "WORLEY3D":
			material = WORLEY3D
		case "MARMOL":
			material = MARMOL
		}

		// Es una forma bastante ortopédica de hacerlo pero funciona.
		// ¿No puede hacerse de forma directa?
		// El campo es opcional.
		if escena.Escena[i].Translation != nil {
			translationRaw = escena.Escena[i].Translation
			translationVec.X = float64(translationRaw[0])
			translationVec.Y = float64(translationRaw[1])
			translationVec.Z = float64(translationRaw[2])
		}

		// Lo mismo con el caso del color.
		if escena.Escena[i].Color != nil {
			colorRaw = escena.Escena[i].Color
			colorVec.R = uint8(colorRaw[0])
			colorVec.G = uint8(colorRaw[1])
			colorVec.B = uint8(colorRaw[2])
		}

		switch escena.Escena[i].Type {
		case "Esfera":
			esfera := Clases.Esfera{
				Clases.BaseObject{
					id,
					material,
					translationVec,
					colorVec,
				},
				float64(escena.Escena[i].Radio),
			}
			Objetos = append(Objetos, esfera)

		}
	}
}

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
			currentColor = elemento.GetColor()
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
	var color = color.RGBA{0, 0, 0, 0}

	color.R = uint8(float64(currentColor.R) * diffuseIntensity)
	color.G = uint8(float64(currentColor.G) * diffuseIntensity)
	color.B = uint8(float64(currentColor.B) * diffuseIntensity)
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

func main() {
	var NDC_x float64
	var NDC_y float64
	var PixelScreen_x float64
	var PixelScreen_y float64
	var PixelCamera_x float64
	var PixelCamera_y float64

	var ro Vectores.Vector
	var rd Vectores.Vector
	var color color.RGBA

	// Open our jsonFile
	jsonFile, err := os.Open("escena.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var escena Escena

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &escena)

	//defineObjetos()
	cargaObjetos(escena)

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

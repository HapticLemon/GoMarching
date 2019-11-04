package main

import (
	"./Clases"
	"./Ruido"
	"./Vectores"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

// Carga de objetos desde archivo Json
//
func cargaObjetos(escena Escena) {
	var id int = 0
	var material int = 0
	var translationRaw []int
	var positionRaw []int
	var dimensionRaw []int
	var positionVec = Vectores.Vector{0, 0, 0}
	var translationVec = Vectores.Vector{0, 0, 0}
	var dimensionVec = Vectores.Vector{0, 0, 0}
	var colorRaw []int
	var colorVec = color.RGBA{0, 0, 0, 0}

	for i := 0; i < len(escena.Escena); i++ {
		// Procesamos el tipo de material para que pase de string a enum.
		switch escena.Escena[i].Material {
		case "NOMAT":
			material = NOMAT
		case "WORLEY3D":
			material = WORLEY3D
		case "SIMPLEX":
			material = SIMPLEX
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

		if escena.Escena[i].Position != nil {
			positionRaw = escena.Escena[i].Translation
			positionVec.X = float64(positionRaw[0])
			positionVec.Y = float64(positionRaw[1])
			positionVec.Z = float64(positionRaw[2])
		}

		if escena.Escena[i].Dimensions != nil {
			dimensionRaw = escena.Escena[i].Dimensions
			dimensionVec.X = float64(dimensionRaw[0])
			dimensionVec.Y = float64(dimensionRaw[1])
			dimensionVec.Z = float64(dimensionRaw[2])
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
		case "Octaedro":
			octaedro := Clases.Octaedro{
				Clases.BaseObject{
					id,
					material,
					translationVec,
					colorVec,
				},
				float64(escena.Escena[i].Radio),
			}
			Objetos = append(Objetos, octaedro)
		case "Caja":
			caja := Clases.Caja{
				Clases.BaseObject{
					id,
					material,
					translationVec,
					colorVec,
				},
				positionVec,
				dimensionVec,
			}
			Objetos = append(Objetos, caja)
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
/*func mapTheWorld(punto Vectores.Vector) (float64, uint8) {
	// Distancia inicial arbitrariamente grande.
	//
	var distancia float64 = 1000
	var distanciaObjeto float64
	var material int
	var cont uint8 = 0
	var indiceObjeto uint8 = 0

	for _, elemento := range Objetos {
		distanciaObjeto = elemento.Distancia(punto)
		if distanciaObjeto < distancia {
			distancia = distanciaObjeto
			currentColor = elemento.GetColor()
			material = elemento.GetMaterial()
			indiceObjeto = cont
		}
		cont += 1
	}
	CurrentMaterial = material
	return distancia, indiceObjeto
}*/

// Hay que hacer que mapTheWorld funcione pasando la escena como parámetro.
// Supongo que la forma correcta de hacerlo sería con el vector de objetos en local
// y no en global.
//
func mapTheWorld(punto Vectores.Vector, Escena []Clases.Objeto) (float64, uint8) {
	// Distancia inicial arbitrariamente grande.
	//
	var distancia float64 = 1000
	var distanciaObjeto float64
	var material int
	var cont uint8 = 0
	var indiceObjeto uint8 = 0

	for _, elemento := range Escena {
		distanciaObjeto = elemento.Distancia(punto)
		if distanciaObjeto < distancia {
			distancia = distanciaObjeto
			currentColor = elemento.GetColor()
			material = elemento.GetMaterial()
			indiceObjeto = cont
		}
		cont += 1
	}
	CurrentMaterial = material
	return distancia, indiceObjeto
}

// Cálculo de la normal (gradiente) en un punto.
//
func calculateNormal(punto Vectores.Vector, posObjeto uint8) Vectores.Vector {
	var gradiente = Vectores.Vector{1, 0, 0}

	gradiente.X = Objetos[posObjeto].Distancia(Vectores.Vector{punto.X + EPSILON, punto.Y, punto.Z}) - Objetos[posObjeto].Distancia(Vectores.Vector{punto.X - EPSILON, punto.Y, punto.Z})
	gradiente.Y = Objetos[posObjeto].Distancia(Vectores.Vector{punto.X, punto.Y + EPSILON, punto.Z}) - Objetos[posObjeto].Distancia(Vectores.Vector{punto.X, punto.Y - EPSILON, punto.Z})
	gradiente.Z = Objetos[posObjeto].Distancia(Vectores.Vector{punto.X, punto.Y, punto.Z + EPSILON}) - Objetos[posObjeto].Distancia(Vectores.Vector{punto.X, punto.Y, punto.Z - EPSILON})

	gradiente.MultiplyByScalar(-1)
	return gradiente.Normalize()
}

// Cálculo de iluminación difusa
// Hay que pasar punto, material, difusa y normal
//
func ilumina(punto Vectores.Vector, diffuseIntensity float64, normal Vectores.Vector) color.RGBA {
	var color = color.RGBA{0, 0, 0, 0}

	if CurrentMaterial == NOMAT {
		color.R = uint8(float64(currentColor.R) * diffuseIntensity)
		color.G = uint8(float64(currentColor.G) * diffuseIntensity)
		color.B = uint8(float64(currentColor.B) * diffuseIntensity)
		color.A = 255
	} else if CurrentMaterial == WORLEY3D {
		var worley3dValue = Ruido.Worley3D(punto)

		color.R = uint8(float64(currentColor.R) * worley3dValue)
		color.G = uint8(float64(currentColor.G) * worley3dValue)
		color.B = uint8(float64(currentColor.B) * worley3dValue)
		color.A = 255
	} else if CurrentMaterial == SIMPLEX {
		// Lo dejo con las coordenadas por respeto a la implementación original.
		// Podría cambiarse por Vectores.Vector
		//
		var SimplexValue = Ruido.Noise3(punto.X, punto.Y, punto.Z)

		SimplexValue = Ruido.Clip(SimplexValue, 0, 1)

		color.R = uint8(float64(currentColor.R) * SimplexValue)
		color.G = uint8(float64(currentColor.G) * SimplexValue)
		color.B = uint8(float64(currentColor.B) * SimplexValue)
		color.A = 255
	}

	return color
}

// Implementación de niebla según idea de Íñigo Quílez.
// https://iquilezles.org/www/articles/fog/fog.htm
func applyFog(color color.RGBA, distancia float64) color.RGBA {
	var fogAmount float32 = 0.0

	fogAmount = float32(1.0 - math.Pow(math.E, -distancia*DENSIDAD))

	return mixColor(color, FOGCOLOR, fogAmount)
}

// Interpolación entre dos colores.
//
func mixColor(x color.RGBA, y color.RGBA, a float32) color.RGBA {
	var resultado color.RGBA

	resultado.R = uint8(float32(x.R)*(1-a) + float32(y.R)*a)
	resultado.G = uint8(float32(x.G)*(1-a) + float32(y.G)*a)
	resultado.B = uint8(float32(x.B)*(1-a) + float32(y.B)*a)

	return resultado
}

func softShadow(ro Vectores.Vector, posObjeto uint8, normal Vectores.Vector) float32 {
	var restoObjetos []Clases.Objeto
	var cont uint8 = 0
	var rd Vectores.Vector
	var punto Vectores.Vector
	var angulo float32 = 0
	var t float64 = 0

	//var posObjeto uint8 = 0
	var shadow float32 = 0

	// Monto un sclice con los todos los elementos excepto el procesado ya que no
	// puede darse sombra a sí mismo.
	//
	for _, elemento := range Objetos {
		if posObjeto != cont {
			restoObjetos = append(restoObjetos, elemento)
		}
		cont += 1
	}

	rd = LIGHT.Sub(ro).Normalize()

	// Compruebo el ángulo entre el vector hacia la luz y la normal del objeto.
	// Si dicho ángulo pasa de 90 grados (en radianes), no hay sombra. Nos sirve para
	// evitar "dobles sombras", por ejemplo en la parte superior e inferior de una esfera.
	//
	angulo = float32(math.Acos(normal.Dot(rd) / normal.Length() * rd.Length()))

	if angulo > NOVENTAGRADOSRAD {
		return 1.0
	}

	for x := 1; x < MAXSTEPS; x++ {
		punto = rd.MultiplyByScalar(t).Add(ro)
		distancia, _ := mapTheWorld(punto, restoObjetos)
		if distancia < MINIMUM_HIT_DISTANCE {
			return 0.0
		}

		shadow = min(float32((8.0*distancia)/float64(x)), shadow)
		t += distancia
	}

	return clip(shadow, 0, 1)

	/*	for i in 1..MAXSTEPS{
			punto = Add(ro,MultiplyByScalar(rd,t));
			let (distancia, idObjeto, colorObjeto, materialObjeto)  = mapTheWorld(punto, &restoEscena);
			if distancia < MINIMUM_HIT_DISTANCE {
				return 0.0;
			}
			shadow = ((8.0 * distancia) / i as f32).min(shadow);
			t += distancia
		}

		return clip(shadow,0.0,1.0)*/
}

func min(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func raymarch(ro Vectores.Vector, rd Vectores.Vector) color.RGBA {

	var punto Vectores.Vector
	var directionToLight Vectores.Vector
	var normal Vectores.Vector
	var t float64 = 0
	var diffuseIntensity float64 = 0
	var valorSombra float32 = 0

	//var distancia float64 = 0
	var color = color.RGBA{30, 30, 150, 255}
	//var posObjeto uint = 0

	for x := 0; x < MAXSTEPS; x++ {
		punto = ro.Add(rd.MultiplyByScalar(t))
		distancia, posObjeto := mapTheWorld(punto, Objetos)

		if distancia < MINIMUM_HIT_DISTANCE {
			//directionToLight = punto.Sub(LIGHT).Normalize()
			directionToLight = LIGHT.Sub(punto).Normalize()
			normal = calculateNormal(punto, posObjeto)

			// Si sólo tenemos un objeto no calcularemos sombras.
			// TODO : Hay que revisar la sombra para que se muestre bien; parece que "atraviesa".
			// TODO : También debe de verse afectada por la niebla.
			if posObjeto == 0 {
				//print("GateteeRLz")
			}
			if CASTSHADOWS == true {
				if len(Objetos) > 1 {
					valorSombra = softShadow(punto, posObjeto, normal)
					if valorSombra == 0.0 {
						if FOG == true {
							return applyFog(SHADOWCOLOR, t)
						}
						return SHADOWCOLOR
					}
				}
			}

			diffuseIntensity = math.Max(0.0, normal.Dot(directionToLight))
			color = ilumina(punto, diffuseIntensity, normal)
			if FOG == true {
				color = applyFog(color, t)
			}
			return color
		}
		t += distancia
	}

	// Devuelvo el color de fondo.
	return color
}

func clip(valor float32, max float32, min float32) float32 {
	if valor > max {
		return max
	}
	if valor < min {
		return min
	}
	return valor
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
	var nuevo Vectores.Vector
	var color color.RGBA

	var fileIn string
	var fileOut string

	start := time.Now()

	//argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	fileIn = argsWithoutProg[0]
	fileOut = argsWithoutProg[1]

	fmt.Printf("Files In : %s, Out %s\n", fileIn, fileOut)

	// Open our jsonFile
	jsonFile, err := os.Open(fileIn)
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

	cargaObjetos(escena)

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	out, err := os.Create(fileOut)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Calculo el Field of View. El ángulo es de 45 grados.
	//
	var FOV float64 = float64(math.Tan(float64(ALPHA / 2.0 * math.Pi / 180.0)))

	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			// Hacemos las conversiones de espacios
			//
			NDC_x = (float64(x) + correccion) / float64(WIDTH)
			NDC_y = (float64(y) + correccion) / float64(HEIGHT)

			PixelScreen_x = 2*NDC_x - 1
			PixelScreen_y = 2*NDC_y - 1

			PixelCamera_x = PixelScreen_x * ImageAspectRatio * FOV
			PixelCamera_y = PixelScreen_y * FOV

			// Origen y dirección

			//ro = EYE.Add(FORWARD.MultiplyByScalar(FL)).Add(RIGHT.MultiplyByScalar(PixelCamera_x)).Add(UP.MultiplyByScalar(PixelCamera_y))
			//rd = ro.Sub(EYE).Normalize()

			ro = EYE
			nuevo.X = PixelCamera_x
			nuevo.Y = PixelCamera_y
			nuevo.Z = -1

			rd = nuevo.Sub(ro).Normalize()
			//rd = Normalize(Sub(Point3{x : PixelCamera_X, y: PixelCamera_Y, z : -1.0}, ro));

			color = raymarch(ro, rd)

			img.Set(x, y, color)

		}
	}
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	var opt jpeg.Options

	opt.Quality = 80
	// ok, write out the data into the new JPEG file

	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Generated image to %s \n", fileOut)

}

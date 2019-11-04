package main

import (
	"./Clases"
	"./Vectores"
	"image/color"
)

var FL float64 = 0.5

const CTEOCTAEDRO = 0.57735027

var EYE = Vectores.Vector{0, 0, 0}
var UP = Vectores.Vector{0, 1, 0}
var RIGHT = Vectores.Vector{1, 0, 0}
var FORWARD = Vectores.Vector{0, 0, 1}
var LIGHT = Vectores.Vector{0, -200, -150.0}
var COLOR = Vectores.Vector{0, 0, 200}
var FOGCOLOR = color.RGBA{R: 128, G: 0, B: 128, A: 255}

var WIDTH int = 800
var HEIGHT int = 600

// Ángulo para el FOV. Actúa como una especie de zoom.
var ALPHA float32 = 55.0

var correccion float64 = 0.5
var ImageAspectRatio float64 = float64(WIDTH) / float64(HEIGHT)
var MAXSTEPS = 32
var MINIMUM_HIT_DISTANCE = 0.35
var EPSILON float64 = 0.01
var FOG bool = true
var DENSIDAD float64 = 0.005

// Slice genérica en la que almacenaremos todos los objetos
//
var Objetos []Clases.Objeto
var currentColor color.RGBA

type Escena struct {
	Escena []Objeto `json:"Escena"`
}

type Objeto struct {
	Name        string `json:"Name"`
	Type        string `json:"Type"`
	Material    string `json:"Material"`
	Radio       int    `json:"Radio"`
	Position    []int  `json:Position`
	Translation []int  `json:Translation`
	Dimensions  []int  `json:Dimensions`
	Color       []int  `json:Color`
}

// En Go las enumeraciones se montan así.
const (
	NOMAT = iota
	WORLEY3D
	SIMPLEX
)

var CurrentMaterial int

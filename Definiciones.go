package main

import (
	"./Clases"
	"./Vectores"
)

var FL float64 = 1.0

const CTEOCTAEDRO = 0.57735027

var EYE = Vectores.Vector{0, 0, -15}
var UP = Vectores.Vector{0, 1, 0}
var RIGHT = Vectores.Vector{1, 0, 0}
var FORWARD = Vectores.Vector{0, 0, 1}
var LIGHT = Vectores.Vector{0, 30, 0.0}
var COLOR = Vectores.Vector{0, 0, 200}

var WIDTH = 640
var HEIGHT = 480

var correccion float64 = 0.5
var ImageAspectRatio float64 = 1.6
var MAXSTEPS = 32
var MINIMUM_HIT_DISTANCE = 0.05

// Slice genérica en la que almacenaremos todos los objetos
//
var Objetos []Clases.Objeto

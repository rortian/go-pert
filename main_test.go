package main

import (
	"image"
	"image/color"
	"sync"
	"testing"

	"github.com/rortian/go-pert/fractal"
	"github.com/rortian/go-pert/paint"
)

var m = 2.0
var n = 2.0

var lambda_x = 1e-6
var lambda_y = 0.0

var width = 500
var height = 500

var x_min = -1.5
var x_max = 1.5

var y_min = -1.5
var y_max = 1.5

var finished = make(chan int, height)

var pert = fractal.SingPert{complex(m, 0), complex(n, 0), complex(lambda_x, lambda_y)}
var grid2 = fractal.Grid{width, height, x_max, y_max, x_min, y_min, &pert, &sync.WaitGroup{}}

var img2 = image.NewNRGBA(image.Rect(0, 0, width, height))

var red = color.NRGBA{0xFF, 0, 0, 0xFF}
var blue = color.NRGBA{0, 0xFF, 0, 0xFF}
var green = color.NRGBA{0, 0, 0xFF, 0xFF}

var simple = []color.Color{red, blue, green}

var first = &paint.SimplePaint{&paint.SimpColors{simple}, img2}

func TestGrid(t *testing.T) {
	first.PaintFrac(grid2.Solve())
}

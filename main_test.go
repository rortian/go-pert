package main

import (
	"image"
	"image/color"
	"sync"
	"testing"
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

var pert = SingPert{complex(m, 0), complex(n, 0), complex(lambda_x, lambda_y)}
var grid = Grid{width, height, x_max, y_max, x_min, y_min, &pert, finished}
var grid2 = GridWG{width, height, x_max, y_max, x_min, y_min, &pert, &sync.WaitGroup{}}

var img = image.NewNRGBA(image.Rect(0, 0, width, height))
var img2 = image.NewNRGBA(image.Rect(0, 0, width, height))

var red = color.NRGBA{0xFF, 0, 0, 0xFF}
var blue = color.NRGBA{0, 0xFF, 0, 0xFF}
var green = color.NRGBA{0, 0, 0xFF, 0xFF}

var simple = []color.Color{red, blue, green}

var first = &SimplePaint{&SimpColors{simple}, img2}

func TestGrid(t *testing.T) {

	rows := grid.Solve()
	for needs := height; needs > 0; needs-- {
		select {
		case x := <-finished:
			for y, speed := range rows[x] {
				img.Set(x, y, simple[speed%3])
			}
		}
	}
}

func TestGridWG(t *testing.T) {
	first.PaintFrac(grid2.Solve())
}

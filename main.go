package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/rortian/go-pert/fractal"
	"github.com/rortian/go-pert/paint"
)

var _ = image.Rect(0, 0, 1, 1)
var _ = fractal.SingPert{}

func init() {
	http.HandleFunc("/fractals", errorHandler(fractalHandler))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}

func parseFloat(s string, d float64) float64 {
	p, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return d
	}
	return p
}

func parseInt(s string, d int64) int64 {
	p, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return d
	}
	return p
}

func parseSingPert(r *http.Request) *fractal.SingPert {
	var m = parseFloat(r.FormValue("m"), 2.0)
	var n = parseFloat(r.FormValue("n"), 2.0)

	var lambda_x = parseFloat(r.FormValue("lambda_x"), 1e-6)
	var lambda_y = parseFloat(r.FormValue("lambda_y"), 0.0)

	return &fractal.SingPert{complex(m, 0), complex(n, 0), complex(lambda_x, lambda_y)}
}

func parseGrid(r *http.Request, sp *fractal.SingPert) *fractal.Grid {
	var width = parseInt(r.FormValue("width"), 500)
	var height = parseInt(r.FormValue("height"), 500)

	var x_min = parseFloat(r.FormValue("x_min"), -1.25)
	var x_max = parseFloat(r.FormValue("x_max"), 1.25)

	var y_min = parseFloat(r.FormValue("y_min"), -1.25)
	var y_max = parseFloat(r.FormValue("y_max"), 1.25)

	return &fractal.Grid{int(width), int(height), x_max, y_max, x_min, y_min, sp, &sync.WaitGroup{}}
}

func fractalHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	singpert := parseSingPert(r)

	grid := parseGrid(r, singpert)

	img := image.NewNRGBA(image.Rect(0, 0, grid.X, grid.Y))

	var red = color.NRGBA{0xFF, 0, 0, 0xFF}
	var blue = color.NRGBA{0, 0xFF, 0, 0xFF}
	var green = color.NRGBA{0, 0, 0xFF, 0xFF}

	var simple = []color.Color{red, blue, green}

	var first = &paint.SimplePaint{&paint.SimpColors{simple}, img}

	first.PaintFrac(grid.Solve())

	png.Encode(w, img)
	return nil

}

func main() {

	log.Fatal(http.ListenAndServe(":8899", nil))
}

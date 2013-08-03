package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"os"
	"sync"
)

type SingPert struct {
	m, n, lambda complex128
}

func (p *SingPert) Step(z complex128) complex128 {
	zm := cmplx.Pow(z, p.m)
	zn := cmplx.Pow(z, p.n)
	return zm + p.lambda/zn
}

func (p *SingPert) Path(z complex128) chan complex128 {
	c := make(chan complex128)
	go func() {
		current := z
		for {
			c <- current
			current = p.Step(current)
		}
	}()
	return c
}

func (p *SingPert) Escape(z complex128) uint16 {
	c := p.Path(z)
	i := uint16(0)
	for current := <-c; cmplx.Abs(current) < 3; current = <-c {
		i++
	}
	return i
}

func (p *SingPert) EscapeS(z complex128) uint16 {
	i := uint16(0)
	for current := z; cmplx.Abs(current) < 3; i++ {
		current = p.Step(current)
	}
	return i
}

type Grid struct {
	x, y                       int
	x_max, y_max, x_min, y_min float64
	*SingPert
	finished chan int
}

type GridWG struct {
	x, y                       int
	x_max, y_max, x_min, y_min float64
	*SingPert
	*sync.WaitGroup
}

func (g *Grid) Solve() [][]uint16 {
	ret := make([][]uint16, g.x)
	x_delta := (g.x_max - g.x_min) / float64(g.x-1)
	y_delta := (g.y_max - g.y_min) / float64(g.y-1)
	for i, _ := range ret {
		x_current := g.x_min + float64(i)*x_delta
		ret[i] = make([]uint16, g.y)
		go g.CalcRow(ret[i], complex(x_current, 0), y_delta, i)
	}
	return ret
}

func (g *GridWG) Solve() [][]uint16 {
	ret := make([][]uint16, g.x)
	g.Add(g.x)
	x_delta := (g.x_max - g.x_min) / float64(g.x-1)
	y_delta := (g.y_max - g.y_min) / float64(g.y-1)
	for i, _ := range ret {
		x_current := g.x_min + float64(i)*x_delta
		ret[i] = make([]uint16, g.y)
		go g.CalcRow(ret[i], complex(x_current, 0), y_delta, i)
	}
	g.Wait()
	return ret
}

func (g *Grid) CalcRow(row []uint16, x complex128, y_delta float64, row_num int) {
	for i := range row {
		pos := x + complex(0, g.y_max-y_delta*float64(i))
		row[i] = g.EscapeS(pos)
	}
	g.finished <- row_num
}

func init() {
	http.HandleFunc("/heyfractals", errorHandler(fractalHandler))
}

func (g *GridWG) CalcRow(row []uint16, x complex128, y_delta float64, row_num int) {
	for i := range row {
		pos := x + complex(0, g.y_max-y_delta*float64(i))
		row[i] = g.EscapeS(pos)
	}
	g.Done()
}

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

type Colorer interface {
	Colorize(uint16) color.Color
}

type Paintable interface {
	Set(int, int, color.Color)
}

type ColorPaint interface {
	Colorer
	Paintable
}

type Painter interface {
	PaintFrac([][]uint16)
}

type SimpColors struct {
	Colors []color.Color
}

func (s *SimpColors) Colorize(n uint16) color.Color {
	return s.Colors[int(n)%len(s.Colors)]
}

type SimplePaint struct {
	Colorer
	Paintable
}

func (s *SimplePaint) PaintFrac(vs [][]uint16) {
	for x, _ := range vs {
		for y, speed := range vs[x] {
			s.Set(x, y, s.Colorize(speed))
		}
	}
}

func fractalHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func main() {
	var width, height int
	var m, n, lambda_x, lambda_y, x_min, x_max, y_min, y_max float64

	flag.Float64Var(&m, "m", 2, "the m in z^m + lambda / z^n")
	flag.Float64Var(&n, "n", 2, "the n in z^m + lambda / z^n")

	flag.Float64Var(&lambda_x, "lx", 1e-6, "the real part of lambda in z^m + lambda / z^n")
	flag.Float64Var(&lambda_y, "ly", 0, "the imaginary part of lambda in z^m + lambda / z^n")

	flag.IntVar(&width, "width", 500, "the width of the image")
	flag.IntVar(&height, "height", 500, "the height of the image")

	flag.Float64Var(&x_min, "x_min", -1, "the minimum x value of the image")
	flag.Float64Var(&x_max, "x_max", 1, "the maximum x value of the image")

	flag.Float64Var(&y_min, "y_min", -1, "the minimum y value of the image")
	flag.Float64Var(&y_max, "y_max", 1, "the maximum y value of the image")

	flag.Parse()

	finished := make(chan int, height)

	pert := SingPert{complex(m, 0), complex(n, 0), complex(lambda_x, lambda_y)}
	grid := Grid{width, height, x_max, y_max, x_min, y_min, &pert, finished}
	grid2 := GridWG{width, height, x_max, y_max, x_min, y_min, &pert, &sync.WaitGroup{}}
	rows := grid.Solve()

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	img2 := image.NewNRGBA(image.Rect(0, 0, width, height))

	red := color.NRGBA{0xFF, 0, 0, 0xFF}
	blue := color.NRGBA{0, 0xFF, 0, 0xFF}
	green := color.NRGBA{0, 0, 0xFF, 0xFF}

	simple := []color.Color{red, blue, green}

	first := &SimplePaint{&SimpColors{simple}, img2}
	first.PaintFrac(grid2.Solve())

	fmt.Printf("%v", first)

	out, _ := os.Create("fun.png")
	defer out.Close()

	fmt.Printf("%v is simple\n", simple)

	for needs := height; needs > 0; needs-- {
		select {
		case x := <-finished:
			for y, speed := range rows[x] {
				img.Set(x, y, simple[speed%3])
			}
		}
	}

	var b bytes.Buffer
	png.Encode(&b, img2)

	http.HandleFunc("/fractal", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		rdbuf := bytes.NewReader(b.Bytes())
		io.Copy(w, rdbuf)
	})

	png.Encode(out, img)
	log.Fatal(http.ListenAndServe(":8899", nil))
}

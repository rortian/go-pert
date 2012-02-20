
package main

import(
  "math/cmplx"
  "fmt"
  "time"
  "runtime"
)

type SingPert struct {
  m,n,lambda complex128
}

func (p *SingPert) Step(z complex128) complex128 {
  zm := cmplx.Pow(z,p.m)
  zn := cmplx.Pow(z,p.n)
  return zm + p.lambda / zn
}

func (p *SingPert) Path(z complex128) chan complex128 {
  c := make(chan complex128)
  go func(){
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
  for current := <- c ; cmplx.Abs(current) < 3 ; current = <- c {
    i++
  }
  return i
}

type Grid struct {
  x,y int
  x_max,y_max,x_min,y_min float64
  *SingPert
}

func (g *Grid) Solve() [][]uint16 {
  ret := make([][]uint16,g.x)
  x_delta := (g.x_max - g.x_min)/ float64(g.x - 1)
  y_delta := (g.y_max - g.y_min)/ float64(g.y - 1)
  for i,_ := range ret {
    x_current := g.x_min + float64(i)*x_delta
    ret[i] = make([]uint16,g.y)
    go g.CalcRow(ret[i],complex(x_current,0),y_delta)
  }
  return ret
}

func (g *Grid) CalcRow(row []uint16,x complex128,y_delta float64){
  for i := range row {
    func(y int){
      pos := x + complex(0,g.y_max-y_delta*float64(y))
      row[y] = g.Escape(pos)
    }(i)
  }
}


func main(){
  fmt.Printf("There are %v goroutines now",runtime.Goroutines())
  pert := SingPert{ 2,2,0.001i }
  grid := Grid { 100, 100, 1, 1, -1, -1, &pert }
  hi := grid.Solve()
  time.Sleep(200*time.Millisecond)
  fmt.Printf("%v\n",hi[3])
}

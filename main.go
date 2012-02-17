
package main

import(
  "fmt"
  "math/cmplx"
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
  c := make(chan complex128,5)
  go func(){
    current := z
    for {
      select {
        case c <- current:
      }
      current = p.Step(current)
    }
  }()
  return c
}

func (p *SingPert) Escape(z complex128) int {
  c := p.Path(z)
  i := 0
  for current := <- c ; cmplx.Abs(current) < 3 ; current = <- c {
    i++
    fmt.Printf("%v\t%v\n",current,cmplx.Abs(current))
  }
  return i
}

type Grid struct {
  x,y int
  x_max,y_max,x_min,y_min float64
}

func (g *Grid) Escape() [][]uint16 {
  ret := make([][]uint16,g.x)
  for i,_ := range ret {
    ret[i] = make([]uint16,g.y)
  }
  return ret
}

func main(){
  pert := SingPert{ 2,2,0.001i }
  fmt.Printf("%v\n",pert)
  fmt.Printf("%v\n",pert.Step(0.5+0.5i))
  pert.Escape(0.5+0.5i)
  grid := Grid { 100,100,1,1,-1,-1 }
  hmm := grid.Escape()
  fmt.Printf("%v\n",hmm)
  fmt.Printf("%v\n",cap(hmm[0]))
}

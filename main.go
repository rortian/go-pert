
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

func main(){
  pert := SingPert{ 2,2,0.001i }
  fmt.Printf("%v\n",pert)
  fmt.Printf("%v\n",pert.Step(0.5+0.5i))
  pert.Escape(0.5+0.5i)
}

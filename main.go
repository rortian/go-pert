
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

func main(){
  pert := SingPert{ 2,2,0.001i }
  fmt.Printf("%v\n",pert)
  fmt.Printf("%v\n",pert.Step(0.5+0.5i))
}

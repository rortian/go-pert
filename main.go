
package main

import(
  "fmt"
  "math/cmplx"
)

type SingPert struct {
  m,n,lambda complex128
}

func (p *SingPert) step(z complex128) complex128 {
  zm := cmplx.Pow(z,p.m)
  zn := cmplx.Pow(z,p.n)
  return zm + p.lambda / zn
}

func main(){
  fmt.Println("hi")
  pert := SingPert{ 2,2,0.001i }
  fmt.Printf("%v\n",pert)
}

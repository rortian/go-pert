package fractal

import (
	"math/cmplx"
	"sync"
)

type SingPert struct {
	M, N, Lambda complex128
}

func Power(z complex128,n int) complex128 {
  if(n < 0) {
    return 1 / Power(z,-n)
  }
  switch n {
    case 0: return complex(1.0,0)
    case 1: return z
    case 2: return z*z;
    case 3: return z*z*z;
    case 4: return z*z*z*z;
    default: return z*z*z*z*z*Power(z,n - 5)
  }
}

func (p *SingPert) Step(z complex128) complex128 {
	zm := Power(z, p.M)
	zn := Power(z, p.N)
	return zm + p.Lambda/zn
}

func (p *SingPert) Escape(z complex128) uint16 {
	i := uint16(0)
	for current := z; cmplx.Abs(current) < 3; i++ {
		current = p.Step(current)
	}
	return i
}

type Grid struct {
	X, Y                       int
	X_max, Y_max, X_min, Y_min float64
	*SingPert
	*sync.WaitGroup
}

func (g *Grid) Solve() [][]uint16 {
	ret := make([][]uint16, g.X)
	g.Add(g.X)
	x_delta := (g.X_max - g.X_min) / float64(g.X-1)
	y_delta := (g.Y_max - g.Y_min) / float64(g.Y-1)
	for i, _ := range ret {
		x_current := g.X_min + float64(i)*x_delta
		ret[i] = make([]uint16, g.Y)
		go g.CalcRow(ret[i], complex(x_current, 0), y_delta, i)
	}
	g.Wait()
	return ret
}

func (g *Grid) CalcRow(row []uint16, x complex128, y_delta float64, row_num int) {
	for i := range row {
		pos := x + complex(0, g.Y_max-y_delta*float64(i))
		row[i] = g.Escape(pos)
	}
	g.Done()
}

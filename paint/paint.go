package paint

import (
  "image"
  "image/color"
  "image/png"
)

type Colorer interface {
	Colorize(uint16) color.Color
}

type Paintable interface {
	Set(int, int, color.Color)
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


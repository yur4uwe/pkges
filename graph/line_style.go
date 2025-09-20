package graph

import (
	"math"

	"github.com/fogleman/gg"
)

type LineStyle struct {
	solid        bool
	dots         bool
	pillars      bool
	solidWidth   float64
	dotsRadius   float64
	pillarsWidth float64
}

func NewLS() *LineStyle {
	return &LineStyle{
		solid:        false,
		dots:         false,
		pillars:      false,
		solidWidth:   2,
		dotsRadius:   3,
		pillarsWidth: 2,
	}
}

func (ls *LineStyle) IsSolid() bool {
	return ls.solid
}

func (ls *LineStyle) IsDots() bool {
	return ls.dots
}

func (ls *LineStyle) IsPillars() bool {
	return ls.pillars
}

func (ls *LineStyle) Solid(width ...float64) {
	ls.solid = true
	if len(width) > 0 {
		ls.solidWidth = width[0]
	}
}

func (ls *LineStyle) Dots(radius ...float64) {
	ls.dots = true
	if len(radius) > 0 {
		ls.dotsRadius = radius[0]
	}
}

func (ls *LineStyle) Pillars(width ...float64) {
	ls.pillars = true
	if len(width) > 0 {
		ls.pillarsWidth = width[0]
	}
}

func (ls *LineStyle) SetLineParams(dc *gg.Context) {
	if ls.solid {
		dc.SetLineWidth(ls.solidWidth)
	}

	if ls.dots {
		dc.SetLineWidth(ls.dotsRadius)
	}

	if ls.pillars {
		dc.SetLineWidth(ls.pillarsWidth)
		dc.SetLineCap(gg.LineCapSquare)
	}
}

func (ls *LineStyle) DrawLine(dc *gg.Context, x, y []float64, originY float64) {
	if ls.solid && len(x) > 0 {
		dc.NewSubPath()
		dc.MoveTo(x[0], y[0])
		for i := 1; i < len(x); i++ {
			dc.LineTo(x[i], y[i])
		}
		dc.Stroke()
	}

	if ls.dots {
		for i := range x {
			dc.DrawCircle(x[i], y[i], 5)
			dc.Fill()
		}
	}

	if ls.pillars {
		for i := 0; i < int(math.Min(float64(len(x)), float64(len(y)))); i++ {
			dc.NewSubPath()
			dc.MoveTo(x[i], originY)
			dc.LineTo(x[i], y[i])
			dc.Stroke()
		}
	}
}

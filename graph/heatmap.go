package graph

import (
	"fmt"
	"math"
)

func (g *Graph) Heatmap(x, y []float64, values [][]float64) {
	if g.gtype != -1 {
		panic("Graph type already set. Cannot add another heatmap.")
	}

	if len(x) == 0 || len(y) == 0 || len(values) == 0 {
		return
	}
	if len(values) != len(y) {
		return
	}
	for _, row := range values {
		if len(row) != len(x) {
			return
		}
	}

	g.values = values
	g.plots = append(g.plots, Plot{
		x: x,
		y: y,
	})
	g.gtype = HeatmapType
}

func valueBounds(values [][]float64) (float64, float64) {
	if len(values) == 0 || len(values[0]) == 0 {
		return 0, 0
	}

	minVal := values[0][0]
	maxVal := values[0][0]

	for _, row := range values {
		for _, v := range row {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
	}

	return minVal, maxVal
}

func colorForValue(v, vmin, vmax float64) (r, g, b float64) {
	t := (v - vmin) / (vmax - vmin)
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// blue -> cyan -> green -> yellow -> red
	switch {
	case t < 0.25:
		t2 := t / 0.25
		r = 0
		g = t2
		b = 1
	case t < 0.5:
		t2 := (t - 0.25) / 0.25
		r = 0
		g = 1
		b = 1 - t2
	case t < 0.75:
		t2 := (t - 0.5) / 0.25
		r = t2
		g = 1
		b = 0
	default:
		t2 := (t - 0.75) / 0.25
		r = 1
		g = 1 - t2
		b = 0
	}
	return
}

func (g *Graph) drawHeatmapScale(x, y, height float64) {
	const width = 20.0

	// Draw the heatmap scale rectangle
	g.dc.SetRGB(0, 0, 0)
	g.dc.DrawRectangle(x, y, width, height)
	g.dc.Stroke()

	vmin, vmax := valueBounds(g.values)

	steps := math.Ceil(height)
	if steps < 2 {
		steps = 2
	}
	for s := 0.0; s < steps; s++ {
		t := s / (steps - 1)
		val := vmax - t*(vmax-vmin)
		R, G, B := colorForValue(val, vmin, vmax)
		y0 := y + s*(height/steps)
		g.dc.SetRGB(R, G, B)
		g.dc.DrawRectangle(x, y0, width, height/steps+0.5)
		g.dc.Fill()
	}

	for _, v := range computeTicks(vmin, vmax) {
		scaledY := y + height*(1-v/(vmax-vmin))

		g.dc.SetRGB(0, 0, 0)
		g.dc.DrawLine(x+width, scaledY, x+width+6, scaledY)
		g.dc.Stroke()

		g.dc.DrawStringAnchored(fmt.Sprintf("%.3f", v), x+width+10, scaledY, 0, 0.5)
	}
}

func (g *Graph) drawHeatmap(scaleX, scaleY, plotHeight, plotWidth, offsetX, offsetY float64) {
	g.drawHeatmapScale(offsetX+plotWidth+20, offsetY, plotHeight)

	x := g.plots[0].x
	y := g.plots[0].y

	vmin, vmax := valueBounds(g.values)
	if vmax == vmin {
		vmax = vmin + 1e-9
	}

	scalerX := g.xScaler(scaleX, offsetX)
	scalerY := g.yScaler(scaleY, offsetY)

	yTop := make([]float64, len(y))
	yBottom := make([]float64, len(y))
	if len(y) == 1 {
		yTop[0] = g.bounds.minY
		yBottom[0] = g.bounds.maxY
	} else {
		dy := (g.bounds.maxY - g.bounds.minY) / float64(len(y))
		for j := range y {
			yTop[j] = g.bounds.minY + float64(j)*dy
			yBottom[j] = yTop[j] + dy
		}
	}

	xLeft := make([]float64, len(x))
	xRight := make([]float64, len(x))
	if len(x) == 1 {
		xLeft[0] = g.bounds.minX
		xRight[0] = g.bounds.maxX
	} else {
		dx := (g.bounds.maxX - g.bounds.minX) / float64(len(x))
		for i := range x {
			xLeft[i] = g.bounds.minX + float64(i)*dx
			xRight[i] = xLeft[i] + dx
		}
	}

	for j := range y {
		for i := range x {
			val := g.values[j][i]
			R, G, B := colorForValue(val, vmin, vmax)

			x0 := scalerX(xLeft[i])
			x1 := scalerX(xRight[i])
			y0 := scalerY(yTop[j])
			y1 := scalerY(yBottom[j])

			xmin := math.Min(x0, x1)
			ymin := math.Min(y0, y1)
			w := math.Abs(x1 - x0)
			h := math.Abs(y1 - y0)

			if xmin < offsetX {
				w -= (offsetX - xmin)
				xmin = offsetX
			}
			if ymin < offsetY {
				h -= (offsetY - ymin)
				ymin = offsetY
			}
			if xmin+w > offsetX+plotWidth {
				w = offsetX + plotWidth - xmin
			}
			if ymin+h > offsetY+plotHeight {
				h = offsetY + plotHeight - ymin
			}

			if w <= 0 || h <= 0 {
				continue
			}

			g.dc.SetRGB(R, G, B)
			g.dc.DrawRectangle(xmin, ymin, w, h)
			g.dc.Fill()
		}
	}
}

package graph

import "fmt"

type Plot struct {
	x      []float64
	y      []float64
	labels []string
	ls     *LineStyle
}

func (g *Graph) Plot(x, y []float64, ls *LineStyle, labels ...[]string) {
	if g.gtype == HeatmapType {
		panic("Heatmap type already set. Cannot add plot.")
	}

	if len(x) == 0 || len(y) == 0 || (len(labels) > 0 && len(labels[0]) != len(x)) || ls == nil {
		return
	}

	if len(x) != len(y) {
		panic(fmt.Sprintf("x and y arrays must have the same length: %d != %d", len(x), len(y)))
	}

	lbs := make([]string, 0)
	if len(labels) > 0 {
		lbs = labels[0]
	}

	plot := Plot{
		x:      x,
		y:      y,
		labels: lbs,
		ls:     ls,
	}

	g.plots = append(g.plots, plot)
	g.gtype = GraphType
}

func (g *Graph) drawPlots(scaleX, scaleY, offsetX, offsetY, originY float64) {
	xScale := g.xScaler(scaleX, offsetX)
	yScale := g.yScaler(scaleY, offsetY)

	for i := range g.plots {
		color := plotColors[i%len(plotColors)]
		g.dc.SetRGB(color[0], color[1], color[2])

		currPlot := g.plots[i]

		currPlot.ls.SetLineParams(g.dc)

		x := ScaleArray(currPlot.x, xScale)
		y := ScaleArray(currPlot.y, yScale)

		if currPlot.ls.IsSolid() {
			g.dc.MoveTo(x[0], y[0])
		}

		currPlot.ls.DrawLine(g.dc, x, y, originY)

		g.dc.SetRGB(0, 0, 0)

		if len(currPlot.labels) == 0 {
			continue
		}

		for j := range x {
			g.dc.DrawStringAnchored(currPlot.labels[j], x[j]+10, y[j]+10, 0, 0)
		}
	}

	g.dc.Stroke()
}

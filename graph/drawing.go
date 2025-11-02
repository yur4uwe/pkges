package graph

import (
	"fmt"
	"math"
)

var plotColors = [][3]float64{
	{1, 0, 0},     // red
	{0, 0, 1},     // blue
	{0, 1, 0},     // green
	{1, 0.5, 0},   // orange
	{0.5, 0, 0.5}, // purple
	{0, 0.7, 0.7}, // teal
}

func (g *Graph) drawAxes(scaleX, scaleY, plotHeight, plotWidth, xOffset, yOffset float64) (originX float64, originY float64) {
	originX = xOffset + (-g.bounds.minX * scaleX)
	originY = yOffset + (g.bounds.maxY * scaleY)

	insideYAxis := isInsideAxis(originX, xOffset, xOffset+plotWidth)
	insideXAxis := isInsideAxis(originY, yOffset, yOffset+plotHeight)

	if insideXAxis {
		g.dc.SetRGB(0.3, 0.3, 0.3)
		g.dc.SetLineWidth(2)
		g.dc.DrawLine(xOffset, originY, xOffset+plotWidth, originY)
		g.dc.Stroke()
	}

	if insideYAxis {
		g.dc.SetRGB(0.3, 0.3, 0.3)
		g.dc.SetLineWidth(2)
		g.dc.DrawLine(originX, yOffset, originX, yOffset+plotHeight)
		g.dc.Stroke()
	}

	g.dc.SetLineWidth(1)

	xTicks := computeTicks(g.bounds.minX, g.bounds.maxX)
	yTicks := computeTicks(g.bounds.minY, g.bounds.maxY)

	if g.gtype == GraphType {
		g.dc.SetRGB(0.85, 0.85, 0.85)
		for _, i := range xTicks {
			xTick := xOffset + (float64(i)-g.bounds.minX)*scaleX
			g.dc.DrawLine(xTick, yOffset, xTick, yOffset+plotHeight)
			g.dc.Stroke()
		}

		for _, i := range yTicks {
			yTick := yOffset + (g.bounds.maxY-float64(i))*scaleY
			g.dc.DrawLine(xOffset, yTick, xOffset+plotWidth, yTick)
			g.dc.Stroke()
		}
	}

	g.dc.SetRGB(0, 0, 0)
	const minLabelDist = 20.0
	lastXLabel := -1000.0
	lastYLabel := -1000.0

	xTickBaseY := yOffset + plotHeight
	if insideXAxis {
		xTickBaseY = originY
	}

	for _, i := range xTicks {
		xTick := xOffset + (float64(i)-g.bounds.minX)*scaleX
		if math.Abs(xTick-lastXLabel) <= minLabelDist {
			continue
		}

		g.dc.DrawLine(xTick, xTickBaseY-5, xTick, xTickBaseY+5)
		g.dc.Stroke()

		label := fmt.Sprintf("%.2f", i)

		labelY := 0.0
		if insideXAxis {
			labelY = xTickBaseY + 14
		} else {
			labelY = yOffset + plotHeight + 14
		}

		g.dc.DrawStringAnchored(label, xTick, labelY, 0.5, 0.5)

		lastXLabel = xTick
	}

	yTickBaseX := xOffset
	if insideYAxis {
		yTickBaseX = originX
	}

	for _, i := range yTicks {
		yTick := yOffset + (g.bounds.maxY-float64(i))*scaleY
		if math.Abs(yTick-lastYLabel) <= minLabelDist {
			continue
		}

		g.dc.DrawLine(yTickBaseX-5, yTick, yTickBaseX+5, yTick)
		g.dc.Stroke()

		label := fmt.Sprintf("%.2f", i)

		labelX := 0.0
		if insideYAxis {
			labelX = yTickBaseX - 12
		} else {
			labelX = xOffset - 12
		}

		g.dc.DrawStringAnchored(label, labelX, yTick, 1, 0.5)

		lastYLabel = yTick
	}

	return originX, originY
}

func (g *Graph) Draw() error {
	if len(g.plots) == 0 {
		return fmt.Errorf("no data to plot")
	}

	font, err := GetFontFace(14)
	if err != nil {
		return err
	}

	padding := 40.0
	heatmapTempScalePadding := 100.0

	plotHeight := float64(g.height) - padding*2
	plotWidth := float64(g.width) - padding*2

	if g.gtype == HeatmapType {
		plotWidth -= heatmapTempScalePadding
	}

	offsetX := padding * 1.5
	offsetY := padding

	g.dc.SetFontFace(font)

	g.dc.SetRGB(0.98, 0.98, 0.98)
	g.dc.DrawRectangle(offsetX, offsetY, plotWidth, plotHeight)
	g.dc.Fill()

	g.dc.SetRGB(0, 0, 0)
	g.dc.SetLineWidth(1)
	g.dc.DrawRectangle(offsetX, offsetY, plotWidth, plotHeight)
	g.dc.Stroke()

	g.computeBounds()
	scaleX, scaleY := g.getScaleFactors(plotHeight, plotWidth)

	_, originY := g.drawAxes(scaleX, scaleY, plotHeight, plotWidth, offsetX, offsetY)

	switch g.gtype {
	case GraphType:
		g.drawPlots(scaleX, scaleY, offsetX, offsetY, originY)
	case HeatmapType:
		g.drawHeatmap(scaleX, scaleY, plotHeight, plotWidth, offsetX, offsetY)
	}

	return nil
}

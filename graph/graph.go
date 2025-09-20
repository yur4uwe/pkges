package graph

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

const (
	wwidth  = 800.0
	wheight = 400.0
)

var plotColors = [][3]float64{
	{1, 0, 0},     // red
	{0, 1, 0},     // green
	{0, 0, 1},     // blue
	{1, 0.5, 0},   // orange
	{0.5, 0, 0.5}, // purple
	{0, 0.7, 0.7}, // teal
}

type Graph struct {
	dc    *gg.Context
	ls    *LineStyle
	xArgs [][]float64
	yArgs [][]float64
}

func NewGraph() *Graph {
	dc := gg.NewContext(wwidth, wheight)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	return &Graph{dc: dc, xArgs: make([][]float64, 0), yArgs: make([][]float64, 0)}
}

func (g *Graph) SetLineStyle(ls *LineStyle) {
	g.ls = ls
}

func (g *Graph) SavePNG(filename string, replace ...bool) error {
	name := strings.TrimSuffix(filename, ".png")
	ext := ".png"
	if len(replace) > 0 && !replace[0] {
		if _, err := os.Stat(name + ext); err == nil {
			for i := 1; ; i++ {
				if _, err := os.Stat(name + "_" + strconv.Itoa(i) + ext); err != nil && os.IsNotExist(err) {
					name = name + "_" + strconv.Itoa(i)
					break
				}
			}
		}
	}

	return g.dc.SavePNG(name + ext)
}

func (g *Graph) Plot(x, y []float64) {
	if len(x) == 0 || len(y) == 0 {
		return
	}

	if len(x) != len(y) {
		panic("x and y arrays must have the same length")
	}

	g.xArgs = append(g.xArgs, x)
	g.yArgs = append(g.yArgs, y)
}

func (g *Graph) Borders() (float64, float64, float64, float64) {
	if len(g.xArgs) == 0 || len(g.yArgs) == 0 {
		return 0, 0, 0, 0
	}

	x := flatten(g.xArgs)
	y := flatten(g.yArgs)

	if len(x) == 0 || len(y) == 0 {
		return 0, 0, 0, 0
	}

	minX, maxX := x[0], x[0]
	minY, maxY := y[0], y[0]

	for i := 1; i < len(x); i++ {
		if x[i] < minX {
			minX = x[i]
		}
		if x[i] > maxX {
			maxX = x[i]
		}
		if y[i] < minY {
			minY = y[i]
		}
		if y[i] > maxY {
			maxY = y[i]
		}
	}

	deltaX := maxX - minX
	deltaY := maxY - minY

	minX = math.Min(minX, -maxX/10) - math.Max(deltaX/10, 0.1)
	maxX = math.Max(maxX, -minX/10) + math.Max(deltaX/10, 0.1)
	minY = math.Min(minY, -maxY/10) - math.Max(deltaY/10, 0.1)
	maxY = math.Max(maxY, -minY/10) + math.Max(deltaY/10, 0.1)

	return minX, maxX, minY, maxY
}

func (g *Graph) Draw() error {
	if len(g.xArgs) == 0 || len(g.yArgs) == 0 {
		return fmt.Errorf("no data to plot")
	}

	minX, maxX, minY, maxY := g.Borders()
	scaleX := wwidth / (maxX - minX)
	scaleY := wheight / (maxY - minY)

	originX := -minX * scaleX
	originY := wheight - (-minY * scaleY)

	g.dc.SetRGB(0.3, 0.3, 0.3)
	g.dc.SetLineWidth(2)
	g.dc.DrawLine(0, originY, wwidth, originY)
	g.dc.Stroke()
	g.dc.DrawLine(originX, 0, originX, wheight)
	g.dc.Stroke()

	g.dc.SetLineWidth(1)

	if err := drawAxesLabels(g.dc, minX, maxX, minY, maxY, originX, originY, scaleX, scaleY); err != nil {
		return err
	}

	xScale := scaledX(minX, scaleX)
	yScale := scaledY(minY, scaleY)
	g.ls.SetLineParams(g.dc)

	for i := range g.xArgs {
		color := plotColors[i%len(plotColors)]
		g.dc.SetRGB(color[0], color[1], color[2])

		x := ScaleArray(g.xArgs[i], xScale)
		y := ScaleArray(g.yArgs[i], yScale)

		if g.ls.IsSolid() {
			g.dc.MoveTo(x[0], y[0])
		}

		g.ls.DrawLine(g.dc, x, y, originY)
	}

	g.dc.Stroke()

	return nil
}

func (g *Graph) Clear() {
	g.dc.SetRGB(1, 1, 1)
	g.dc.Clear()
	g.xArgs = make([][]float64, 0)
	g.yArgs = make([][]float64, 0)
}

func flatten(arrays [][]float64) []float64 {
	var result []float64
	for _, array := range arrays {
		result = append(result, array...)
	}
	return result
}

func IntLinearArray(min, max int) []float64 {
	if max == min {
		return nil
	}

	result := make([]float64, max-min+1)

	for i := range result {
		result[i] = float64(min + i)
	}

	return result
}

func LinearArray(min, max float64, length int) []float64 {
	if max == min || length <= 0 {
		return nil
	}

	if length == 1 {
		return []float64{min}
	}

	norm := make([]float64, length)
	step := (max - min) / float64(length-1)
	for i := 0; i < length; i++ {
		norm[i] = min + step*float64(i)
	}

	return norm
}

func scaledX(minX, scaleX float64) func(x float64) float64 {
	return func(x float64) float64 {
		return (x - minX) * scaleX
	}
}
func scaledY(minY, scaleY float64) func(y float64) float64 {
	return func(y float64) float64 {
		return wheight - (y-minY)*scaleY
	}
}

func ScaleArray(arr []float64, scale func(float64) float64) []float64 {
	scaled := make([]float64, len(arr))
	for i, v := range arr {
		scaled[i] = scale(v)
	}
	return scaled
}

// Plot draws axes and the plot for given x and y arrays, and saves the PNG file.
func Plot(x, y []float64, filename string, line_style *LineStyle, replace ...bool) error {
	if len(x) == 0 || len(y) == 0 {
		return fmt.Errorf("x and y arrays must not be empty")
	}

	dc := gg.NewContext(wwidth, wheight)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// minX, maxX, minY, maxY := Borders(x, y)
	// scaleX := wwidth / (maxX - minX)
	// scaleY := wheight / (maxY - minY)

	// originX := -minX * scaleX
	// originY := wheight - (-minY * scaleY)

	// dc.SetRGB(0.3, 0.3, 0.3)
	// dc.SetLineWidth(2)
	// dc.DrawLine(0, originY, wwidth, originY)
	// dc.Stroke()
	// dc.DrawLine(originX, 0, originX, wheight)
	// dc.Stroke()

	// dc.SetLineWidth(1)

	// if err := drawAxesLabels(dc, minX, maxX, minY, maxY, originX, originY, scaleX, scaleY); err != nil {
	// 	return err
	// }

	// dc.SetRGB(0, 0, 1)
	// line_style.SetLineParams(dc)

	// xScale := scaledX(minX, scaleX)
	// yScale := scaledY(minY, scaleY)

	// if line_style.IsSolid() {
	// 	dc.MoveTo(xScale(x[0]), yScale(y[0]))
	// }

	// for i := range x {
	// 	line_style.DrawLine(dc, xScale(x[i]), yScale(y[i]), originY)
	// }
	// dc.Stroke()

	name := strings.TrimSuffix(filename, ".png")
	ext := ".png"
	if len(replace) > 0 && !replace[0] {
		for i := 1; ; i++ {
			if _, err := os.Stat(name + "_" + strconv.Itoa(i) + ext); err != nil && os.IsNotExist(err) {
				name = name + "_" + strconv.Itoa(i)
				break
			}
		}
	}

	return dc.SavePNG(name + ext)
}

// IntegersInRange returns a slice of all integers between min and max (inclusive).
func IntegersInRange(min, max float64) []int {
	start := int(math.Ceil(min))
	end := int(math.Floor(max))
	if end < start {
		return []int{}
	}
	result := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func toFloatSlice(ints []int) []float64 {
	floats := make([]float64, len(ints))
	for i, v := range ints {
		floats[i] = float64(v)
	}
	return floats
}

func computeTicks(ints []int, min, max float64) []float64 {
	if len(ints) <= 10 {
		new_tick_num := 12

		ticks := make([]float64, new_tick_num)
		step := (max - min) / float64(new_tick_num-1)
		for i := 0; i < new_tick_num; i++ {
			ticks[i] = min + step*float64(i)
		}
		return ticks
	}

	if len(ints) <= 20 {
		return toFloatSlice(ints)
	}

	step := int(math.Ceil(float64(len(ints)) / 20.0))
	ticks := make([]int, 0, (len(ints)+step-1)/step)

	for i := 0; i < len(ints); i += step {
		ticks = append(ticks, ints[i])
	}

	return toFloatSlice(ticks)
}

func drawAxesLabels(dc *gg.Context, minX, maxX, minY, maxY, originX, originY, scaleX, scaleY float64) error {
	err := dc.LoadFontFace("ArialMT.ttf", 14)
	if err != nil {
		return err
	}

	xInts := IntegersInRange(minX, maxX)
	yInts := IntegersInRange(minY, maxY)

	isIntegerTicksX := len(xInts) >= 10
	isIntegerTicksY := len(yInts) >= 10

	xTicks := computeTicks(xInts, minX, maxX)
	yTicks := computeTicks(yInts, minY, maxY)

	dc.SetRGB(0.5, 0.5, 0.5)
	for _, i := range xTicks {
		xTick := (float64(i) - minX) * scaleX
		dc.DrawLine(xTick, 0, xTick, wheight)
		dc.Stroke()
	}

	for _, i := range yTicks {
		yTick := wheight - (float64(i)-minY)*scaleY
		dc.DrawLine(0, yTick, wwidth, yTick)
		dc.Stroke()
	}

	dc.SetRGB(0, 0, 0)
	const minLabelDist = 20.0 // minimum pixel distance between labels
	lastXLabel := -1000.0
	lastYLabel := -1000.0

	xScaled := scaledX(minX, scaleX)
	yScaled := scaledY(minY, scaleY)

	for _, i := range xTicks {
		xTick := xScaled(i)
		if math.Abs(xTick-lastXLabel) <= minLabelDist {
			continue
		}

		dc.DrawLine(xTick, originY-5, xTick, originY+5)
		dc.Stroke()

		label := ""
		if isIntegerTicksX {
			label = fmt.Sprintf("%.0f", i)
		} else {
			label = fmt.Sprintf("%.2f", i)
		}

		dc.DrawStringAnchored(label, xTick, originY+18, 0.5, 0.5)

		lastXLabel = xTick

	}

	for _, i := range yTicks {
		yTick := yScaled(i)
		if math.Abs(yTick-lastYLabel) <= minLabelDist {
			continue
		}

		dc.DrawLine(originX-5, yTick, originX+5, yTick)
		dc.Stroke()
		label := ""

		if isIntegerTicksY {
			label = fmt.Sprintf("%.0f", i)
		} else {
			label = fmt.Sprintf("%.2f", i)
		}

		dc.DrawStringAnchored(label, originX-18, yTick, 0.5, 0.5)

		lastYLabel = yTick

	}

	return nil
}

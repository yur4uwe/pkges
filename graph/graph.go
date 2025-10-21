package graph

import (
	_ "embed"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var plotColors = [][3]float64{
	{1, 0, 0},     // red
	{0, 1, 0},     // green
	{0, 0, 1},     // blue
	{1, 0.5, 0},   // orange
	{0.5, 0, 0.5}, // purple
	{0, 0.7, 0.7}, // teal
}

//go:embed fonts/ArialMT.ttf
var fontBytes []byte

var (
	parsedFont *truetype.Font
	fontOnce   sync.Once
	fontErr    error
)

func loadFont() (*truetype.Font, error) {
	fontOnce.Do(func() {
		parsedFont, fontErr = truetype.Parse(fontBytes)
	})
	return parsedFont, fontErr
}

func GetFontFace(size float64) (font.Face, error) {
	f, err := loadFont()
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		Size: size,
	}), nil
}

type Plot struct {
	x      []float64
	y      []float64
	labels []string
	ls     *LineStyle
}

type Graph struct {
	dc     *gg.Context
	width  int
	height int
	plots  []Plot
}

func NewGraph(w, h int) *Graph {
	dc := gg.NewContext(w, h)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	return &Graph{dc: dc, plots: make([]Plot, 0), width: w, height: h}
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

func (g *Graph) Plot(x, y []float64, ls *LineStyle, labels ...[]string) {
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
}

func (g *Graph) getFlattenedData() ([]float64, []float64) {
	var xData []float64
	var yData []float64

	for _, plot := range g.plots {
		xData = append(xData, plot.x...)
		yData = append(yData, plot.y...)
	}

	return xData, yData
}

func (g *Graph) Borders() (float64, float64, float64, float64) {
	if len(g.plots) == 0 {
		return 0, 0, 0, 0
	}

	x, y := g.getFlattenedData()

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
	if len(g.plots) == 0 {
		return fmt.Errorf("no data to plot")
	}

	font, err := GetFontFace(14)
	if err != nil {
		return err
	}

	g.dc.SetFontFace(font)

	minX, maxX, minY, maxY := g.Borders()
	scaleX := float64(g.width) / (maxX - minX)
	scaleY := float64(g.height) / (maxY - minY)

	originX := -minX * scaleX
	originY := float64(g.height) - (-minY * scaleY)

	g.dc.SetRGB(0.3, 0.3, 0.3)
	g.dc.SetLineWidth(2)
	g.dc.DrawLine(0, originY, float64(g.width), originY)
	g.dc.Stroke()
	g.dc.DrawLine(originX, 0, originX, float64(g.height))
	g.dc.Stroke()

	g.dc.SetLineWidth(1)

	if err := g.drawAxesLabels(g.dc, minX, maxX, minY, maxY, originX, originY, scaleX, scaleY); err != nil {
		return err
	}

	xScale := g.scaledX(minX, scaleX)
	yScale := g.scaledY(minY, scaleY)

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

	return nil
}

func (g *Graph) Clear() {
	g.dc.SetRGB(1, 1, 1)
	g.dc.Clear()
	g.plots = make([]Plot, 0)
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

func (g *Graph) scaledX(minX, scaleX float64) func(x float64) float64 {
	return func(x float64) float64 {
		return (x - minX) * scaleX
	}
}
func (g *Graph) scaledY(minY, scaleY float64) func(y float64) float64 {
	return func(y float64) float64 {
		return float64(g.height) - (y-minY)*scaleY
	}
}

func ScaleArray(arr []float64, scale func(float64) float64) []float64 {
	scaled := make([]float64, len(arr))
	for i, v := range arr {
		scaled[i] = scale(v)
	}
	return scaled
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

func ToFloatSlice(ints []int) []float64 {
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
		return ToFloatSlice(ints)
	}

	step := int(math.Ceil(float64(len(ints)) / 20.0))
	ticks := make([]int, 0, (len(ints)+step-1)/step)

	for i := 0; i < len(ints); i += step {
		ticks = append(ticks, ints[i])
	}

	return ToFloatSlice(ticks)
}

func (g *Graph) drawAxesLabels(dc *gg.Context, minX, maxX, minY, maxY, originX, originY, scaleX, scaleY float64) error {
	xInts := IntegersInRange(minX, maxX)
	yInts := IntegersInRange(minY, maxY)

	isIntegerTicksX := len(xInts) >= 10
	isIntegerTicksY := len(yInts) >= 10

	xTicks := computeTicks(xInts, minX, maxX)
	yTicks := computeTicks(yInts, minY, maxY)

	dc.SetRGB(0.5, 0.5, 0.5)
	for _, i := range xTicks {
		xTick := (float64(i) - minX) * scaleX
		dc.DrawLine(xTick, 0, xTick, float64(g.height))
		dc.Stroke()
	}

	for _, i := range yTicks {
		yTick := float64(g.height) - (float64(i)-minY)*scaleY
		dc.DrawLine(0, yTick, float64(g.width), yTick)
		dc.Stroke()
	}

	dc.SetRGB(0, 0, 0)
	const minLabelDist = 20.0
	lastXLabel := -1000.0
	lastYLabel := -1000.0

	xScaled := g.scaledX(minX, scaleX)
	yScaled := g.scaledY(minY, scaleY)

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

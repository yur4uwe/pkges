package graph

import (
	_ "embed"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	GraphType = iota
	HeatmapType
)

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

type bounds struct {
	minX float64
	maxX float64
	minY float64
	maxY float64
}

type Graph struct {
	dc     *gg.Context
	width  int
	height int
	plots  []Plot
	values [][]float64
	gtype  int
	bounds bounds
}

func NewGraph(w, h int) *Graph {
	dc := gg.NewContext(w, h)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	return &Graph{dc: dc, plots: make([]Plot, 0), width: w, height: h, gtype: -1}
}

func (g *Graph) ClearWithNewSize(w, h int) {
	g.dc = gg.NewContext(w, h)
	g.width = w
	g.height = h
	g.Clear()
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

func (g *Graph) getFlattenedData() ([]float64, []float64) {
	var xData []float64
	var yData []float64

	for _, plot := range g.plots {
		xData = append(xData, plot.x...)
		yData = append(yData, plot.y...)
	}

	return xData, yData
}

func (g *Graph) computeBounds() {
	if len(g.plots) == 0 {
		g.bounds = bounds{0, 0, 0, 0}
		return
	}

	x, y := g.getFlattenedData()

	if len(x) == 0 || len(y) == 0 {
		g.bounds = bounds{0, 0, 0, 0}
		return
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
	}

	for i := 1; i < len(y); i++ {
		if y[i] < minY {
			minY = y[i]
		}
		if y[i] > maxY {
			maxY = y[i]
		}
	}

	if minX == maxX {
		minX -= 1
		maxX += 1
	}

	if minY == maxY {
		minY -= 1
		maxY += 1
	}

	g.bounds = bounds{minX, maxX, minY, maxY}
}

func (g *Graph) getScaleFactors(plotHeight, plotWidth float64) (float64, float64) {
	scaleX := plotWidth / (g.bounds.maxX - g.bounds.minX)
	scaleY := plotHeight / (g.bounds.maxY - g.bounds.minY)
	return scaleX, scaleY
}

func isInsideAxis(origin, min, max float64) bool {
	return origin >= min && origin <= max
}

func (g *Graph) Clear() {
	g.dc.SetRGB(1, 1, 1)
	g.dc.Clear()
	g.plots = make([]Plot, 0)
	g.gtype = -1
}

func (g *Graph) xScaler(scaleX, offsetX float64) func(x float64) float64 {
	return func(x float64) float64 {
		return offsetX + (x-g.bounds.minX)*scaleX
	}
}
func (g *Graph) yScaler(scaleY, offsetY float64) func(y float64) float64 {
	return func(y float64) float64 {
		return offsetY + (g.bounds.maxY-y)*scaleY
	}
}

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

func computeTicks(min, max float64) []float64 {
	if max <= min {
		return []float64{min}
	}

	ints := math.Ceil(max - min)

	desired := 6
	if ints > 12 {
		desired = 10
	} else if ints > 6 {
		desired = 8
	}

	ticks := make([]float64, desired)
	span := max - min
	step := span / float64(desired-1)
	for i := 0; i < desired; i++ {
		ticks[i] = min + float64(i)*step
	}

	return ticks
}

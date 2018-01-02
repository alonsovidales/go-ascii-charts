package asciicharts

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type BarsData struct {
	title     string
	keys      []string
	vals      []float64
	useColors bool
}

var ErrMissmatchKeysVals = errors.New("the number of keys doesn't match with the number of vals")
var ErrKeyWiderThanCanvas = errors.New("the width of a key is wider than the specified canvas size, please increase the canvas width")
var ErrNoEnoughCanvasHeight = errors.New("there is no room left for the metric in the canvas, please increase the canvas heght")

func InitBars(title string, keys []string, vals []float64, useColors bool) (*BarsData, error) {
	if len(keys) != len(vals) {
		return nil, ErrMissmatchKeysVals
	}

	return &BarsData{
		title:     title,
		keys:      keys,
		vals:      vals,
		useColors: useColors,
	}, nil
}

func (bd *BarsData) AddData(keys []string, vals []float64) error {
	if len(keys) != len(vals) {
		return ErrMissmatchKeysVals
	}

	bd.keys = append(bd.keys, keys...)
	bd.vals = append(bd.vals, vals...)

	return nil
}

func (bd *BarsData) RenderSingleBar(width, barHeight int, legend bool) (string, error) {
	totalValues := 0.0

	for _, v := range bd.vals {
		totalValues += v
	}

	widthPerKey := make([]int, len(bd.vals))
	percPerKey := make([]float64, len(bd.vals))
	for i, v := range bd.vals {
		percPerKey[i] = v / totalValues
		widthPerKey[i] = int(percPerKey[i] * float64(width))
	}

	result := ""
	for h := 0; h < barHeight; h++ {
		for i, w := range widthPerKey {
			result += bd.addColor(strings.Repeat(bd.getSymbol(i), w), i)
		}
		result += "\n"
	}

	if legend {
		for i, k := range bd.keys {
			result += bd.addColor(fmt.Sprintf("%s: %f %%%%\n", k, percPerKey[i]), i)
		}
	}

	return result, nil
}

func (bd *BarsData) RenderMultiBar(width, heght int) ([]string, error) {
	// Calculate the width for the vals and the max with for every data point
	// name
	maxNameWidth := 0
	maxNameLines := 0
	maxValueWidth := 0
	maxValue := 0.0
	minValue := math.Inf(1)

	for i, v := range bd.vals {
		if v > maxValue {
			maxValue = v
		}
		if v < minValue {
			minValue = v
		}
		keyParts := strings.Split(bd.keys[i], " ")
		if len(keyParts) > maxNameLines {
			maxNameLines = len(keyParts)
		}
		for _, keyPart := range keyParts {
			if len(keyPart) > maxNameWidth {
				maxNameWidth = len(keyPart)
			}
		}

		valstr := strconv.FormatFloat(v, 'f', -1, 64)
		if len(valstr) > maxValueWidth {
			maxValueWidth = len(valstr)
		}
	}

	// The -1 is to leave some room for the separation lines
	canvasW := width - maxValueWidth - 1
	canvasH := heght - maxNameLines - 1

	if bd.title != "" {
		canvasH--
	}

	if canvasH <= 1 {
		return nil, ErrNoEnoughCanvasHeight
	}

	unitsPerLine := (maxValue - minValue) / float64(canvasH-1)

	namesPerChunk := int(math.Floor(float64(canvasW) / float64(maxNameWidth+1)))

	if namesPerChunk <= 0 {
		return nil, ErrKeyWiderThanCanvas
	}

	chunks := int(math.Ceil(float64(len(bd.keys)) / float64(namesPerChunk)))

	result := make([]string, chunks)
	for c := 0; c < chunks; c++ {
		result[c] += bd.addColor("- %s:\n", 0, bd.title)
		keysSlice := bd.keys[c*namesPerChunk : int(math.Min(float64(len(bd.keys)), float64(namesPerChunk+c*namesPerChunk)))]
		valsSlice := bd.vals[c*namesPerChunk : int(math.Min(float64(len(bd.keys)), float64(namesPerChunk+c*namesPerChunk)))]
		for i := 0; i < canvasH; i++ {
			lineVal := maxValue - float64(i)*unitsPerLine
			result[c] += fmt.Sprintf("%s |", fmt.Sprintf("%f", lineVal)[:maxValueWidth])
			for ki, _ := range keysSlice {
				if valsSlice[ki] >= lineVal-unitsPerLine/2 {
					result[c] += bd.addColor(" %s", ki, strings.Repeat(bd.getSymbol(i), maxNameWidth))
				} else {
					result[c] += fmt.Sprintf(" %s", strings.Repeat(" ", maxNameWidth))
				}
			}

			result[c] += "\n"
		}
		result[c] += fmt.Sprintf("%s + %s\n", strings.Repeat("-", maxValueWidth), strings.Repeat("-", canvasW))
		for kl := 0; kl < maxNameLines; kl++ {
			result[c] += fmt.Sprintf("%s |", strings.Repeat(" ", maxValueWidth))
			for ki, k := range keysSlice {
				keyParts := strings.Split(k, " ")

				if len(keyParts) > kl {
					result[c] += bd.addColor(" %s", ki, strPad(keyParts[kl], " ", maxNameWidth, false))
				} else {
					result[c] += fmt.Sprintf(" %s", strings.Repeat(" ", maxNameWidth))
				}
			}
			result[c] += "\n"
		}
	}

	return result, nil
}

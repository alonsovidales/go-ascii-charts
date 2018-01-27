package asciichart_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alonsovidales/go-ascii-charts"
)

func TestRenderMultiBaringAndPagination(t *testing.T) {
	br, err := asciichart.InitBars("This is a test", []string{
		"test1 mult line",
		"test2 line2 line3",
		"test3 line2 line3",
	}, []float64{
		13.2,
		10.2,
		11.0,
	})

	if err != nil {
		t.Errorf("The bars chart can't be initialized: %s", err)
	}

	result, err := br.RenderMultiBar(30, 10, false)
	if err != nil {
		t.Errorf("the chart can't be plotted: %s", err)
	}

	// This should fit in a single page, we have 30 width and we should use up
	// to 5 per key plus separators and so
	if len(result) != 1 {
		t.Errorf("the amounth of pages expected was 1, but the number of returned pages was: %d", len(result))
	}

	result, err = br.RenderMultiBar(15, 10, false)
	if err != nil {
		t.Errorf("the chart can't be plotted: %s", err)
	}

	// This should fit in three pages, one per metric since we have about 7 for
	// the left side of the chart + 5 per metric
	if len(result) != 3 {
		t.Errorf("the amounth of pages expected was 3, but the number of returned pages was: %d", len(result))
	}

	err = br.AddData([]string{"aa", "bb"}, []float64{10, 12})

	if err != nil {
		t.Errorf("the data couldn't be added to the chart : %s", err)
	}

	result, err = br.RenderMultiBar(15, 10, false)
	if err != nil {
		t.Errorf("the chart can't be plotted: %s", err)
	}

	// Now we have 3 keys that we used to have before + the 2 new that we just
	// added and one per page
	if len(result) != 5 {
		t.Errorf("the amounth of pages expected was 5, but the number of returned pages was: %d", len(result))
	}

}

func TestErrors(t *testing.T) {
	_, err := asciichart.InitBars("This is a test", []string{
		"test1 mult line",
		"EEOOO new?",
	}, []float64{
		13.2,
	})

	if err == nil {
		t.Log("the number of keys doesn't match with the number of values, but no error was returned")
	}

	br, _ := asciichart.InitBars("This is a test", []string{
		"test1 mult line",
		"EEOOO new?",
	}, []float64{
		11,
		13.2,
	})

	br.AddData([]string{"AddedData"}, []float64{123})

	_, err = br.RenderMultiBar(90, 4, false)
	if err != asciichart.ErrNoEnoughCanvasHeight {
		t.Log("the table height is not sufficient, but no error was returned")
	}

	_, err = br.RenderMultiBar(2, 100, false)
	if err != asciichart.ErrKeyWiderThanCanvas {
		t.Log("the table width is not sufficient to plot any key, but no error was returned")
	}

	err = br.AddData([]string{"aaa"}, []float64{1.2, 1.1})
	if err == nil {
		t.Log("the number of keys doesn't match with the number of values while adding new data, but no error was returned")
	}
}

func TestRenderSingleBar(t *testing.T) {
	keys := []string{
		"key1",
		"key2",
		"key3",
		"key4",
		"key5",
		"key6",
	}
	br, err := asciichart.InitBars("", keys, []float64{
		13.2,
		22.5,
		12.3,
		33.2,
		15.3,
		16.3,
	})

	if err != nil {
		t.Fatalf("not possible to initialize, error: %s", err)
	}

	for i := 1; i < 10; i++ {
		t.Run(fmt.Sprintf("Single bar %d", i), func(t *testing.T) {
			result, err := br.RenderSingleBar(i, i, false, false)
			if err != nil {
				t.Fatalf("error while rendering %d x %d single bar: %s", i, i, err)
			}
			resultLines := strings.Split(result, "\n")
			if len(resultLines) != i+1 {
				t.Fatalf("the number of returned lines doesn't match with the requested one, returned: %d requested: %d", len(resultLines), i+1)
			}
			// We can have a metric smaller than the requested width if we
			// can't plot some of the data, for instance 4 keys in a width of
			// one will not return nothing since no key can be rendered
			if len(resultLines[0]) > i {
				t.Fatalf("the width of the plotted graph doesn't match with the requested one, returned: %d requested: %d", len(resultLines[0]), i, resultLines[0][50], "-")
			}

			// Checking that the number of lines is as before but adding one
			// line per key for the legend
			result, err = br.RenderSingleBar(i, i, true, false)
			if err != nil {
				t.Fatalf("error while rendering %d x %d single bar: %s", i, i, err)
			}
			resultLines = strings.Split(result, "\n")
			if len(resultLines) != i+1+len(keys) {
				t.Fatalf("the number of returned lines doesn't match with the requested one, returned: %d requested: %d", len(resultLines), i+1)
			}
		})
	}
}

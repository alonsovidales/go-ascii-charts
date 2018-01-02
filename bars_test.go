package asciicharts

import (
	"fmt"
	"testing"
)

func TestRenderMultiBaringAndPagination(t *testing.T) {
	br, err := InitBars("This is a test", []string{
		"test1 mult line",
		"test2 line2 line3",
		"test3 line2 line3",
	}, []float64{
		13.2,
		10.2,
		11.0,
	}, true)

	if err != nil {
		t.Errorf("The bars chart can't be initialized: %s", err)
	}

	result, err := br.RenderMultiBar(30, 10)
	if err != nil {
		t.Errorf("The chart can't be plotted: %s", err)
	}

	// This should fit in a single page, we have 30 width and we should use up
	// to 5 per key plus separators and so
	if len(result) != 1 {
		t.Errorf("The amounth of pages expected was 1, but the number of returned pages was: %d", len(result))
	}

	result, err = br.RenderMultiBar(15, 10)
	if err != nil {
		t.Errorf("The chart can't be plotted: %s", err)
	}

	// This should fit in three pages, one per metric since we have about 7 for
	// the left side of the chart + 5 per metric
	if len(result) != 3 {
		t.Errorf("The amounth of pages expected was 3, but the number of returned pages was: %d", len(result))
	}

	err = br.AddData([]string{"aa", "bb"}, []float64{10, 12})

	if err != nil {
		t.Errorf("The data couldn't be added to the chart : %s", err)
	}

	result, err = br.RenderMultiBar(15, 10)
	if err != nil {
		t.Errorf("The chart can't be plotted: %s", err)
	}

	// Now we have 3 keys that we used to have before + the 2 new that we just
	// added and one per page
	if len(result) != 5 {
		t.Errorf("The amounth of pages expected was 5, but the number of returned pages was: %d", len(result))
	}

}

func TestErrors(t *testing.T) {
	_, err := InitBars("This is a test", []string{
		"test1 mult line",
		"EEOOO new?",
	}, []float64{
		13.2,
	}, true)

	if err == nil {
		t.Log("The number of keys doesn't match with the number of values, but no error was returned")
	}

	br, _ := InitBars("This is a test", []string{
		"test1 mult line",
		"EEOOO new?",
	}, []float64{
		11,
		13.2,
	}, true)

	br.AddData([]string{"AddedData"}, []float64{123})

	_, err = br.RenderMultiBar(90, 4)
	if err != ErrNoEnoughCanvasHeight {
		t.Log("The table height is not sufficient, but no error was returned")
	}

	_, err = br.RenderMultiBar(2, 100)
	if err != ErrKeyWiderThanCanvas {
		t.Log("The table width is not sufficient to plot any key, but no error was returned")
	}

	err = br.AddData([]string{"aaa"}, []float64{1.2, 1.1})
	if err == nil {
		t.Log("The number of keys doesn't match with the number of values while adding new data, but no error was returned")
	}
}

func TestRenderSingleBar(t *testing.T) {
	br, err := InitBars("This is a test", []string{
		"key1",
		"key2",
		"key3",
		"key4",
		"key5",
		"key6",
	}, []float64{
		13.2,
		22.5,
		12.3,
		33.2,
		15.3,
		16.3,
	}, true)

	if err != nil {
		t.Fatalf("Not possible to initialize, error: %s", err)
	}

	result, _ := br.RenderSingleBar(100, 2, true)

	fmt.Println(result)
}

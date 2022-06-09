package app

import (
	"bufio"
	"encoding/csv"
	"os"
	"testing"
)

var rawList = ``

func TestParser(t *testing.T) {
	fr, err := os.Open("D:\\Code\\poketto\\data\\test_data.txt")
	if err != nil {
		panic(err)
	}
	defer fr.Close()
	sc := bufio.NewScanner(fr)

	var eps []*Episode
	for sc.Scan() {
		ep := NewEpisode(sc.Text())
		ep.TryParse()
		eps = append(eps, ep)
	}

	fw, err := os.Create("D:\\Code\\poketto\\data\\test_out.txt")
	if err != nil {
		panic(err)
	}
	defer fw.Close()
	w := csv.NewWriter(fw)
	for _, ep := range eps {
		if err := w.Write(ep.ToFields()); err != nil {
			panic(err)
		}
	}
	w.Flush()
}

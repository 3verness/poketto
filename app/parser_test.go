package app

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

var rawList = `[NC-Raws] 小書痴的下剋上：為了成為圖書管理員不擇手段！第三季 / Honzuki no Gekokujou S3 - 33 (Baha 1920x1080 AVC AAC MP4)
[NC-Raws] 小书痴的下克上：为了成为图书管理员不择手段！第三季 / Honzuki no Gekokujou S3 - 33 (B-Global 1920x1080 HEVC AAC MKV)`

func TestParser(t *testing.T) {
	var eps []*Episode
	sc := bufio.NewScanner(strings.NewReader(rawList))
	for sc.Scan() {
		ep := NewEpisode(sc.Text())
		ep.TryParse()
		eps = append(eps, ep)
		fmt.Printf("%+v\n", ep)
	}
}

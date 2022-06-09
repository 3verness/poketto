package app

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Episode struct {
	TitleRaw string

	Name   string
	Season int
	Ep     int
	Group  string
	Dpi    string
	Sub    string
	Source string

	ParseErr error
}

func NewEpisode(raw string) *Episode {
	return &Episode{TitleRaw: raw}
}

func (ep *Episode) TryParse() {
	ep.ParseErr = ep.parse()
}

func (ep *Episode) ToFields() []string {
	return []string{ep.TitleRaw, ep.Name, fmt.Sprint(ep.Season), fmt.Sprint(ep.Ep), ep.Group, ep.Dpi, ep.Sub, ep.Source}
}

func (ep *Episode) parse() error {
	raw := ep.TitleRaw
	if raw == "" {
		return errors.New("原始标题为空，无法解析。")
	}
	raw = strings.NewReplacer("【", "[", "】", "]").Replace(raw)
	var group string
	if regexp.MustCompile(`[\[\]]`).MatchString(raw) {
		group = regexp.MustCompile(`[\[\]]`).Split(raw, -1)[1]
	}
	matcher := regexp.MustCompile(`(.*|\[.*])( -? \d{1,3} |\[\d{1,3}]|\[\d{1,3}.?[vV]\d{1}]|[第第]\d{1,3}[话話集集]|\[\d{1,3}.?END])(.*)`).FindStringSubmatch(raw)
	if matcher == nil {
		return errors.New("无法解析")
	}
	name, season, err := getSeason(matcher[1])
	if err != nil {
		return errors.New("无法解析季数")
	}
	name, err = getName(name)
	if err != nil {
		return errors.New("无法解析标题")
	}
	epNum, err := getEp(matcher[2])
	if err != nil {
		return errors.New("无法解析集数")
	}
	dpi, sub, source, err := getTag(matcher[3])
	if err != nil {
		return errors.New("无法解析标签")
	}

	ep.Name = name
	ep.Season = season
	ep.Ep = epNum
	ep.Group = group
	ep.Dpi = dpi
	ep.Sub = sub
	ep.Source = source
	return nil
}

func getSeason(raw string) (name string, season int, err error) {
	if regexp.MustCompile(`新番|月?番`).MatchString(raw) {
		raw = regexp.MustCompile(`.*新番.`).ReplaceAllString(raw, "")
	} else {
		raw = regexp.MustCompile(`^[^]】]*[]】]`).ReplaceAllString(raw, "")
		raw = strings.TrimSpace(raw)
	}
	raw = regexp.MustCompile(`[\[\]]`).ReplaceAllString(raw, "")
	seasonRe := regexp.MustCompile(`S\d{1,2}|Season \d{1,2}|[第].[季期]`)
	seasonMatcher := seasonRe.FindAllString(raw, -1)
	if seasonMatcher == nil {
		return raw, 1, nil
	} else {
		name = seasonRe.ReplaceAllString(raw, "")
		for _, s := range seasonMatcher {
			if regexp.MustCompile(`S|Season`).MatchString(s) {
				season, err = strconv.Atoi(regexp.MustCompile(`S|Season`).ReplaceAllString(s, ""))
				if err == nil {
					return
				}
			} else if regexp.MustCompile(`[第 ].*[季期]`).MatchString(s) {
				seasonBuf := regexp.MustCompile(`[第季期 ]`).ReplaceAllString(s, "")
				if season, err = strconv.Atoi(seasonBuf); err == nil {
					return
				}
				if season, err = getNum(seasonBuf); err == nil {
					return
				}
			}
		}
	}
	return "", 0, errors.New("无法识别季数")
}

var numDict = map[rune]int{
	'一': 1, '二': 2, '三': 3, '四': 4, '五': 5,
	'六': 6, '七': 7, '八': 8, '九': 9, '十': 10,
}

func getNum(raw string) (int, error) {
	for _, r := range []rune(raw) {
		if n, ok := numDict[r]; ok {
			return n, nil
		}
	}
	return 0, errors.New("无法转换为数字")
}

func getName(raw string) (name string, err error) {
	raw = strings.TrimSpace(raw)
	raw = strings.ReplaceAll(raw, "（仅限港澳台地区）", "")
	slicesRaw := regexp.MustCompile(`/|  |-  `).Split(raw, -1)
	var slices []string
	for _, s := range slicesRaw {
		if s != "" {
			slices = append(slices, s)
		}
	}
	if len(slices) == 1 {
		if strings.Contains(raw, "_") {
			slices = strings.Split(raw, "_")
		} else if strings.Contains(raw, " - ") {
			slices = strings.Split(raw, "-")
		}
	}
	if len(slices) == 1 {
		matcher := regexp.MustCompile(`([^\x00-\xff]{1,})(\s)([\x00-\xff]{4,})`).FindStringSubmatch(raw)
		if matcher != nil && matcher[3] != "" {
			return matcher[3], nil
		}
	}
	maxLen := 0
	for _, s := range slices {
		if l := len(regexp.MustCompile(`[aA-zZ]`).FindAllString(s, -1)); l > maxLen {
			maxLen = l
			name = s
		}
	}
	name = strings.TrimSpace(name)
	return name, nil
}

func getEp(raw string) (epNum int, err error) {
	if epRaw := regexp.MustCompile(`\d{1,3}`).FindString(raw); epRaw != "" {
		epNum, err = strconv.Atoi(epRaw)
		return
	}
	return 0, errors.New("无法解析集数")
}

func getTag(raw string) (dpi, sub, source string, err error) {
	raw = regexp.MustCompile(`[\[\]()（）]`).ReplaceAllString(raw, " ")
	tagsRaw := strings.Split(raw, " ")
	var tags []string
	for _, t := range tagsRaw {
		if t != "" {
			tags = append(tags, t)
		}
	}
	for _, t := range tags {
		if regexp.MustCompile(`[简繁日字幕]|CH|BIG5|GB`).MatchString(t) {
			sub = strings.ReplaceAll(t, "_MP4", "")
		} else if regexp.MustCompile(`1080|720|2160|4K`).MatchString(t) {
			dpi = t
		} else if regexp.MustCompile(`B-Global|[Bb]aha|[Bb]ilibili|AT-X|Web`).MatchString(t) {
			source = t
		}
	}
	err = nil
	return
}

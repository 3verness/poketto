package app

import (
	"crypto/tls"
	"github.com/mmcdole/gofeed"
	"net/http"
)

type RSSParser struct {
	feedParser *gofeed.Parser
	rssURI     string
}

func NewRSSParser(rssURI string) *RSSParser {
	p := &RSSParser{}
	p.feedParser = gofeed.NewParser()
	p.feedParser.Client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	p.rssURI = rssURI
	return p
}

func (p *RSSParser) Grab() (*gofeed.Feed, error) {
	feed, err := p.feedParser.ParseURL(p.rssURI)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

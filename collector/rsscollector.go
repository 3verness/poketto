package collector

import (
	"crypto/tls"
	"github.com/mmcdole/gofeed"
	"net/http"
)

type RSSCollector struct {
	feedParser *gofeed.Parser
	rssURI     string
}

func NewRSSCollector(rssURI string) *RSSCollector {
	c := &RSSCollector{}
	c.feedParser = gofeed.NewParser()
	c.feedParser.Client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	c.rssURI = rssURI
	return c
}

func (p *RSSCollector) Collect() ([]string, error) {
	feed, err := p.feedParser.ParseURL(p.rssURI)
	if err != nil {
		return nil, err
	}

	var torrents []string
	for _, item := range feed.Items {
		for _, enclosure := range item.Enclosures {
			torrents = append(torrents, enclosure.URL)
		}
	}
	return torrents, nil
}

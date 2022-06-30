package app

import "testing"

func TestRSS(t *testing.T) {
	p := NewRSSParser("https://mikanani.me/RSS/MyBangumi?token=4nyg8SRIggQOpbUdpYXSJZSmuHImtJBv8VUDHrikwoM%3d")
	f, err := p.Grab()
	if err != nil {
		t.Error(err)
	}
	for _, item := range f.Items {
		for _, enclosure := range item.Enclosures {
			t.Log(enclosure.URL)
		}
	}
}

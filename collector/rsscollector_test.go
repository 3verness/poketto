package collector

import "testing"

func TestRSSCollect(t *testing.T) {
	c := NewRSSCollector("https://mikanani.me/RSS/MyBangumi?token=4nyg8SRIggQOpbUdpYXSJZSmuHImtJBv8VUDHrikwoM%3d")
	ts, err := c.Collect()
	if err != nil {
		t.Error(err)
	}
	for _, s := range ts {
		t.Log(s)
	}
}

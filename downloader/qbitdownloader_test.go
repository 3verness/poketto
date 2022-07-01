package downloader

import "testing"

func TestQbitDownloader(t *testing.T) {
	torrents := []string{"https://mikanani.me/Download/20220626/8bcf415e471ecbeeb919b25821d8de69b84146ee.torrent",
		"https://mikanani.me/Download/20220626/616de67ecf2889b0bb2b5e60eb2b6c07620e07dc.torrent",
		"https://mikanani.me/Download/20220626/eb79dcd941845eaf3430a7d129fed5fd9b1e8dd7.torrent"}

	d := NewQbitDownloader("http://localhost:8080/api/v2")
	a, _ := d.Login("admin", "adminadmin")
	t.Log(a)

	for _, torrent := range torrents {
		hash, err := d.AddNewTorrent(torrent, "D:\\dev")
		t.Log(hash, err)
		status, err := d.GetTorrentStatus(hash)
		t.Log(status, err)
		files, err := d.GetTorrentFiles(hash)
		t.Log(files, err)
	}
}

package downloader

type TorrentStatus int

const (
	NotFound TorrentStatus = iota
	Downloading
	Complete
	Errored
	Unknown
)

type Downloader interface {
	AddNewTorrent(torrent string, savePath string) (string, error)
	GetTorrentStatus(hash string) (TorrentStatus, error)
	GetTorrentFiles(hash string) ([]string, error)
	RenameFiles(hash string, oldName string, newName string) error
}

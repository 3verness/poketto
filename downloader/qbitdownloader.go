package downloader

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type QbitDownloader struct {
	client  *http.Client
	baseUrl string
	authed  bool
	jar     http.CookieJar
}

type torrentInfo struct {
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	SavePath string `json:"save_path"`
	State    string `json:"state"`
}

type fileInfo struct {
	Index    int    `json:"index"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

func NewQbitDownloader(url string) *QbitDownloader {
	c := &QbitDownloader{}
	c.jar, _ = cookiejar.New(&cookiejar.Options{})
	c.client = &http.Client{Jar: c.jar}
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	c.baseUrl = url
	return c
}

func (d *QbitDownloader) post(endpoint string, params map[string]string) (*http.Response, error) {
	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", d.baseUrl+endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "Failed build request.")
	}

	req.Header.Set("User-Agent", "poketto v0.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get response.")
	}

	return resp, nil
}

func (d *QbitDownloader) Login(username, password string) (bool, error) {
	params := make(map[string]string)
	params["username"] = username
	params["password"] = password

	resp, err := d.post("auth/login", params)
	if err != nil {
		return false, errors.Wrap(err, "Login failed.")
	}

	if resp.Status != "200 OK" {
		return false, errors.Errorf("Login failed with status code %s", resp.Status)
	}
	d.authed = true
	return true, nil
}

func (d QbitDownloader) AddNewTorrent(torrent string, savePath string) (string, error) {
	params := make(map[string]string)
	params["urls"] = torrent
	params["savepath"] = savePath
	params["category"] = "Poketto Managed"

	resp, err := d.post("torrents/add", params)
	if err != nil {
		return "", errors.Wrap(err, "Failed to add "+torrent)
	}

	if resp.Status != "200 OK" {
		return "", errors.Errorf("Failed to add %s with status code %s", torrent, resp.Status)
	}

	time.Sleep(1 * time.Second)

	resp, err = d.post("torrents/info", map[string]string{"category": "Poketto Managed", "sort": "added_on", "reverse": "true"})
	if err != nil {
		return "", errors.Wrap(err, "Failed to check "+torrent+" info")
	}

	if resp.Status != "200 OK" {
		return "", errors.Errorf("Failed to check info with status code %s", resp.Status)
	}

	var info []torrentInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decode info")
	}

	for _, i := range info {
		if i.SavePath == savePath {
			return i.Hash, nil
		}
	}
	return "", errors.New("Added but can not find info")
}

func (d QbitDownloader) GetTorrentStatus(hash string) (TorrentStatus, error) {
	resp, err := d.post("torrents/info", map[string]string{"hash": hash})
	if err != nil {
		return NotFound, errors.Wrap(err, "Failed to check "+hash+" info")
	}

	if resp.Status != "200 OK" {
		return NotFound, errors.Errorf("Failed to check %s info with status code %s", hash, resp.Status)
	}

	var info []torrentInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return NotFound, errors.Wrap(err, "Failed to decode info")
	}

	for _, i := range info {
		if i.State == "error" || i.State == "missingFiles" {
			return Errored, nil
		} else if i.State == "uploading" || strings.Contains(i.State, "UP") {
			return Complete, nil
		} else if i.State == "allocating" || i.State == "downloading" || i.State == "checkingResumeData" || strings.Contains(i.State, "DL") {
			return Downloading, nil
		} else {
			return Unknown, nil
		}
	}

	return NotFound, nil
}

func (d QbitDownloader) GetTorrentFiles(hash string) ([]string, error) {
	resp, err := d.post("torrents/files", map[string]string{"hash": hash})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to check "+hash)
	}

	if resp.Status != "200 OK" {
		return nil, errors.Errorf("Failed to check %s with status code %s", hash, resp.Status)
	}

	var info []fileInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode info")
	}

	var files []string
	for _, i := range info {
		files = append(files, i.Name)
	}
	return files, nil
}

func (d QbitDownloader) RenameFiles(hash string, oldName string, newName string) error {
	resp, err := d.post("torrents/renameFile", map[string]string{"hash": hash, "oldPath": oldName, "newPath": newName})
	if err != nil {
		return errors.Wrap(err, "Failed to rename "+hash)
	}

	if resp.Status != "200 OK" {
		return errors.Errorf("Failed to rename %s with status code %s", hash, resp.Status)
	}

	return nil
}

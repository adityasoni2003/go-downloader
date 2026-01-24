package downloader

import (
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func downloadSingle(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(100),
		progressbar.OptionClearOnFinish(),
	)

	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	return err
}

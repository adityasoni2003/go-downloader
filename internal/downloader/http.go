package downloader

import "net/http"

func supportsRange(url string) (bool, int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	size := resp.ContentLength
	acceptRanges := resp.Header.Get("Accept-Ranges") == "bytes"

	return acceptRanges && size > 0, size, nil
}

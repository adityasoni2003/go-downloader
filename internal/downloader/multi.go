package downloader

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func downloadPart(
	url string,
	file io.WriterAt,
	start, end int64,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range",
		fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 32*1024)
	offset := start

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			file.WriteAt(buf[:n], offset)
			offset += int64(n)
		}
		if err != nil {
			break
		}
	}
}

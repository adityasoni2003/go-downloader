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
	bar interface{ Add(int) error },
) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Range",
		fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 206 = Partial Content (important!)
	if resp.StatusCode != http.StatusPartialContent {
		return
	}

	buf := make([]byte, 32*1024)
	offset := start

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, werr := file.WriteAt(buf[:n], offset)
			if werr != nil {
				return
			}
			offset += int64(n)

			// 🔥 progress update (thread-safe)
			bar.Add(n)
		}
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
	}
}

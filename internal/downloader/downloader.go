package downloader

import (
	"fmt"
	"sync"

	"github.com/adityasoni2003/go-downloader/internal/utils"
	"github.com/schollz/progressbar/v3"
)

func Download(url, path string, threads int) error {
	supported, size, err := supportsRange(url)
	if err != nil || !supported || threads <= 1 {
		fmt.Print("The server does not support multi-threaded downloads or an error occurred. Falling back to single-threaded download.\n")
		return downloadSingle(url, path)
	}

	file, err := utils.CreateFile(path, size)
	if err != nil {
		return err
	}
	defer file.Close()

	// 🔥 Shared progress bar
	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionClearOnFinish(),
	)

	partSize := size / int64(threads)
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		start := int64(i) * partSize
		end := start + partSize - 1

		if i == threads-1 {
			end = size - 1
		}

		wg.Add(1)
		go downloadPart(url, file, start, end, &wg, bar)
	}

	wg.Wait()
	return nil
}

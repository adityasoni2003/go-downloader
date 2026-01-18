package downloader

import (
	"sync"

	"github.com/adityasoni2003/go-downloader/internal/utils"
)

func Download(url, path string, threads int) error {
	supported, size, err := supportsRange(url)
	if err != nil || !supported || threads <= 1 {
		return downloadSingle(url, path)
	}

	file, err := utils.CreateFile(path, size)
	if err != nil {
		return err
	}
	defer file.Close()

	partSize := size / int64(threads)
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		start := int64(i) * partSize
		end := start + partSize - 1
		if i == threads-1 {
			end = size - 1
		}

		wg.Add(1)
		go downloadPart(url, file, start, end, &wg)
	}

	wg.Wait()
	return nil
}

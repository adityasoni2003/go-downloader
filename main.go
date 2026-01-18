package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

func downloadRange(url string, filePath *os.File, start int64, end int64, wg *sync.WaitGroup) {
	defer wg.Done()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()
	buf := make([]byte, 32*1024) // 32KB buffer
	offset := start

	for {
		n, err := res.Body.Read(buf)
		if n > 0 {
			_, wErr := filePath.WriteAt(buf[:n], offset)
			if wErr != nil {
				return
			}
			offset += int64(n)
		}

		if err != nil {
			break
		}
	}
}

func getFileSize(fileUrl string) (int64, error) {
	response, err := http.Head(fileUrl)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Got bad statuscode while getting file size: %s", response.StatusCode)
	}

	return response.ContentLength, nil
}

func downloadFileToPathMultithreads(fileUrl, filePath string, threads int) error {

	size, err := getFileSize(fileUrl)

	if err != nil {
		return err
	}

	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("Got error creating output file %w", err)
	}
	defer file.Close()

	chunkSize := size / int64(threads)
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1

		if i == threads-1 {
			end = size - 1
		}

		wg.Add(1)
		go downloadRange(fileUrl, file, start, end, &wg)
	}

	wg.Wait()

	return nil

}

func safeFilenameFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "download.bin"
	}

	name := filepath.Base(u.Path)
	if name == "" || name == "/" {
		return "download.bin"
	}

	return name
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : go run main.go <url>")
		return
	}

	url := os.Args[1]
	fmt.Println("Fetching file...", url)
	output := safeFilenameFromURL(url)

	// calling the main function to download the file
	err := downloadFileToPathMultithreads(url, output, 4)
	if err != nil {
		fmt.Println("Got error while downloading file %w", err)
		return
	}

	fmt.Println("File is downloaded to ", output)

}

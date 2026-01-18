package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadFileToPath(fileUrl, filePath string) error {

	response, err := http.Get(fileUrl)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Got bad status after requesting %s", response.Status)
	}

	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("Got error creating output file %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)

	if err != nil {
		return fmt.Errorf("Got error while copying downloaded file to loacation %w", err)

	}

	return nil

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : go run main.go <url>")
		return
	}

	url := os.Args[1]
	fmt.Println("Fetching file...", url)
	output := filepath.Base(url)

	// calling the main function to download the file
	err := downloadFileToPath(url, output)
	if err != nil {
		fmt.Println("Got error while downloading file %w", err)
		return
	}

	fmt.Println("File is downloaded to ", output)

}

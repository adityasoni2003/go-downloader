package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adityasoni2003/go-downloader/internal/downloader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: downloader <url>")
		os.Exit(1)
	}

	url := os.Args[1]
	output := filepath.Base(url)

	fmt.Println("Downloading:", url)

	err := downloader.Download(url, output, 4)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Downloaded to:", output)
}

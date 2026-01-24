package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adityasoni2003/go-downloader/internal/downloader"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter file URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	fmt.Print("Enter output directory (press Enter for current dir): ")
	dir, _ := reader.ReadString('\n')
	dir = strings.TrimSpace(dir)

	if dir == "" {
		dir = "."
	}

	output := filepath.Join(dir, filepath.Base(url))

	fmt.Print("Enter number of threads: ")
	tStr, _ := reader.ReadString('\n')
	tStr = strings.TrimSpace(tStr)

	threads, err := strconv.Atoi(tStr)
	if err != nil || threads < 1 {
		fmt.Println("Invalid thread count, using 1")
		threads = 1
	}

	fmt.Println("\nStarting download...")
	err = downloader.Download(url, output, threads)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("\nDownload completed 🎉")
}

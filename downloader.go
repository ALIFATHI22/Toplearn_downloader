package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func downloadFile(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(resp.Body)
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: topcdw <file_name>")
		return
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	time.Sleep(2 * time.Second)

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	var totalSize int64
	for _, url := range urls {
		size, err := getFileSize(url)
		if err != nil {
			fmt.Println("Error getting file size for", url, ":", err)
			return
		}
		totalSize += size
	}

	fmt.Printf("Total size of all files: %d bytes\n", totalSize)
	fmt.Print("Do you want to download the files? (Y/n): ")

	reader := bufio.NewReader(os.Stdin)
	for {
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)

		if response == "" {
			fmt.Print("Please enter 'Y' or 'N': ")
			continue
		}

		if response == "Y" || response == "y" {
			for i, url := range urls {
				fileName := fmt.Sprintf("file%d", i+1)
				fmt.Println("Downloading", url, "to", fileName)
				err := downloadFile(url, fileName)
				if err != nil {
					fmt.Println("Error downloading", url, ":", err)
				} else {
					fmt.Println("Downloaded", url, "to", fileName)
				}
			}
			break
		} else if response == "N" || response == "n" {
			fmt.Println("Download canceled.")
			break
		} else {
			fmt.Println("Invalid response. Please enter 'Y' or 'N'.")
		}
	}
}

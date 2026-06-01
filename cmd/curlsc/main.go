package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ashish-dhakane/curlsc/internal/checker"
	"github.com/ashish-dhakane/curlsc/internal/models"
	"github.com/ashish-dhakane/curlsc/internal/output"
)

func main() {
	// Flags
	filePath := flag.String("file", "", "Path to file with URLs (one per line)")
	shortFilePath := flag.String("f", "", "Path to file with URLs (alias)")
	concurrency := flag.Int("concurrency", 20, "Max parallel requests")
	shortConcurrency := flag.Int("c", 0, "Max parallel requests (alias)")
	timeoutStr := flag.String("timeout", "10s", "Per-request timeout")
	shortTimeoutStr := flag.String("t", "", "Per-request timeout (alias)")
	method := flag.String("method", "GET", "HTTP method: GET or HEAD")
	shortMethod := flag.String("m", "", "HTTP method (alias)")
	outFormat := flag.String("format", "table", "Output format: table, json, csv")
	outPath := flag.String("output", "", "Write results to file path")
	shortOutPath := flag.String("o", "", "Write results to file path (alias)")
	verbose := flag.Bool("verbose", false, "Show headers and full error details")
	shortVerbose := flag.Bool("v", false, "Show headers and full error details (alias)")
	noColor := flag.Bool("no-color", false, "Disable ANSI color codes")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification")
	shortInsecure := flag.Bool("k", false, "Skip TLS certificate verification (alias)")
	retries := flag.Int("retries", 0, "Number of retries on failure")
	follow := flag.Int("follow", 5, "Max redirects to follow")
	userAgent := flag.String("user-agent", "curlsc/1.0", "Custom User-Agent header")
	version := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *version {
		fmt.Println("curlsc version 1.0.0")
		return
	}

	// Alias resolution
	if *shortFilePath != "" { *filePath = *shortFilePath }
	if *shortConcurrency != 0 { *concurrency = *shortConcurrency }
	if *shortTimeoutStr != "" { *timeoutStr = *shortTimeoutStr }
	if *shortMethod != "" { *method = *shortMethod }
	if *shortOutPath != "" { *outPath = *shortOutPath }
	if *shortVerbose { *verbose = *shortVerbose }
	if *shortInsecure { *insecure = *shortInsecure }

	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing timeout: %v\n", err)
		os.Exit(1)
	}

	var urls []string

	// 1. Positional arguments
	urls = append(urls, flag.Args()...)

	// 2. File input
	if *filePath != "" {
		fileUrls, err := readLines(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		urls = append(urls, fileUrls...)
	}

	// 3. Stdin (only if no args and no file)
	if len(urls) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					urls = append(urls, line)
				}
			}
		}
	}

	if len(urls) == 0 {
		fmt.Println("Usage: curlsc [flags] [url1] [url2] ...")
		flag.PrintDefaults()
		return
	}

	c := &checker.Checker{
		Concurrency:  *concurrency,
		Timeout:      timeout,
		Method:       strings.ToUpper(*method),
		Insecure:     *insecure,
		UserAgent:    *userAgent,
		MaxRedirects: *follow,
		Retries:      *retries,
	}

	results, summary := c.CheckAll(urls)

	// Output
	var out io.Writer = os.Stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	renderer := output.GetRenderer(*outFormat, out, *noColor)
	if err := renderer.Render(results, summary); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering output: %v\n", err)
		os.Exit(1)
	}

	// Exit code logic
	if summary.Failure > 0 {
		os.Exit(1)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

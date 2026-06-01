package checker

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ashish-dhakane/curlsc/internal/models"
)

type Checker struct {
	Concurrency int
	Timeout     time.Duration
	Method      string
	Insecure    bool
	UserAgent   string
	MaxRedirects int
	Retries      int
}

func (c *Checker) CheckAll(urls []string) ([]models.CheckResult, models.Summary) {
	start := time.Now()
	results := make([]models.CheckResult, len(urls))
	
	urlChan := make(chan struct {
		index int
		url   string
	}, len(urls))
	
	resChan := make(chan struct {
		index int
		res   models.CheckResult
	}, len(urls))

	var wg sync.WaitGroup
	// Semaphore to limit concurrency
	sem := make(chan struct{}, c.Concurrency)

	// Dispatcher
	go func() {
		for i, u := range urls {
			urlChan <- struct {
				index int
				url   string
			}{i, u}
		}
		close(urlChan)
	}()

	// Workers
	for i := 0; i < c.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range urlChan {
				sem <- struct{}{}
				res := c.Check(task.url)
				resChan <- struct {
					index int
					res   models.CheckResult
				}{task.index, res}
				<-sem
			}
		}()
	}

	// Closer
	go func() {
		wg.Wait()
		close(resChan)
	}()

	var successCount, failureCount int
	for r := range resChan {
		results[r.index] = r.res
		if r.res.Status == models.StatusUp {
			successCount++
		} else {
			failureCount++
		}
	}

	summary := models.Summary{
		Total:       len(urls),
		Success:     successCount,
		Failure:     failureCount,
		ElapsedTime: time.Since(start),
	}

	return results, summary
}

func (c *Checker) Check(targetURL string) models.CheckResult {
	// Validate URL format
	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Scheme == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return models.CheckResult{
			URL:       targetURL,
			Status:    models.StatusError,
			ErrorType: "INVALID_URL",
		}
	}

	client := &http.Client{
		Timeout: c.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= c.MaxRedirects {
				return fmt.Errorf("stopped after %d redirects", c.MaxRedirects)
			}
			return nil
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Insecure},
		},
	}

	var result models.CheckResult
	result.URL = targetURL

	// Retry loop
	for i := 0; i <= c.Retries; i++ {
		start := time.Now()
		req, _ := http.NewRequest(c.Method, targetURL, nil)
		req.Header.Set("User-Agent", c.UserAgent)

		resp, err := client.Do(req)
		result.ResponseTime = time.Since(start)
		result.ResponseTimeMs = result.ResponseTime.Milliseconds()

		if err != nil {
			result.Status = models.StatusError
			result.ErrorType = classifyError(err)
			result.ErrorMessage = err.Error()
			// If it's a transient error, maybe retry, but for now we'll just continue or finish
			continue
		}

		result.StatusCode = resp.StatusCode
		result.FinalURL = resp.Request.URL.String()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result.Status = models.StatusUp
		} else {
			result.Status = models.StatusDown
		}
		
		resp.Body.Close()
		break // Success (or expected failure code), no need to retry
	}

	return result
}

func classifyError(err error) string {
	if err == nil {
		return ""
	}
	
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return "TIMEOUT"
	}

	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "no such host"):
		return "DNS_ERROR"
	case strings.Contains(errStr, "connection refused"):
		return "CONNECTION_REFUSED"
	case strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") || strings.Contains(errStr, "ssl"):
		return "TLS_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

package progress

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// WebProgressChecker checks progress from an HTML website
type WebProgressChecker struct {
	URL string
}

// GetProgress gets latest works in progress from brandonsanderson.com
func (wpc WebProgressChecker) GetProgress() ([]WorkInProgress, error) {
	req, err := http.NewRequest(http.MethodGet, wpc.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return parseProgressFromHTML(string(responseBody))
}

func parseProgressFromHTML(html string) ([]WorkInProgress, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("get document from HTML: %w", err)
	}

	progressEntrySelectors := doc.Find(".progress-item-uniq p")

	textEntries := progressEntrySelectors.Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})

	wips := []WorkInProgress{}
	for i := range textEntries {
		if i%2 == 0 {
			continue
		}

		progressStr := strings.TrimRight(strings.TrimSpace(textEntries[i-1]), "%")
		if progressStr == "" {
			return nil, fmt.Errorf("failed to parse progress from progress entry[%d]", i)
		}
		progress, err := strconv.Atoi(progressStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse progress from progress entry[%d] %q: %w", i, progressStr, err)
		}

		title := strings.TrimSpace(textEntries[i])
		if title == "" {
			return nil, fmt.Errorf("failed to parse title from progress entry[%d]", i)
		}

		wips = append(wips, WorkInProgress{Title: title, Progress: progress})
	}

	if len(wips) == 0 {
		fmt.Println("No progress entries found from HTML:\n", html)
		return nil, errors.New("no progress entries found")
	}

	return wips, nil
}

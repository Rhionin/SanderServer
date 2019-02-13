package progress

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

// WorkInProgress represents each work and its progress
type WorkInProgress struct {
	Title        string `json:"title"`
	Progress     int    `json:"progress"`
	PrevProgress int    `json:"prevProgress"`
}

// ToString prints to string
func (wip *WorkInProgress) ToString() string {
	return fmt.Sprintf("%s (%d%%)", wip.Title, wip.Progress)
}

// CheckProgress gets latest works in progress from brandonsanderson.com
func CheckProgress() []WorkInProgress {
	doc, err := goquery.NewDocument("http://brandonsanderson.com")
	if err != nil {
		log.Fatal("error!", err)
	}

	progressDiv := doc.Find("#pagewrap .progress-titles").First()
	titleSelectors := progressDiv.Find(".book-title")
	progressSelectors := progressDiv.Find(".progress")

	if titleSelectors.Length() != progressSelectors.Length() {
		log.Fatal("Mismatched book progress count!", progressDiv.Text())
	}

	titles := titleSelectors.Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	progresses := progressSelectors.Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})

	wips := make([]WorkInProgress, len(titles))
	for i := range titles {
		title := strings.TrimSpace(titles[i])
		progress, err := strconv.Atoi(strings.TrimSpace(strings.Trim(strings.TrimSpace(progresses[i]), "%")))

		if err != nil {
			log.Fatal("Error parsing progress for ", title)
		} else {
			wips[i] = WorkInProgress{Title: title, Progress: progress}
		}
	}

	return wips
}

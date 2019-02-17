package progress

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WebProgressChecker struct{}

// GetProgress gets latest works in progress from brandonsanderson.com
func (WebProgressChecker) GetProgress() []WorkInProgress {
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

package progress

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var entryRegex = regexp.MustCompile("^(.*)+ ([\\d]+)%$")

// WebProgressChecker checks progress from an HTML website
type WebProgressChecker struct {
	URL string
}

// GetProgress gets latest works in progress from brandonsanderson.com
func (wpc WebProgressChecker) GetProgress() []WorkInProgress {
	doc, err := goquery.NewDocument(wpc.URL)
	if err != nil {
		log.Fatal("error!", err)
	}

	progressDiv := doc.Find(".wpb_wrapper .vc_progress_bar").First()
	progressEntrySelectors := progressDiv.Find(".vc_label")

	progressEntries := progressEntrySelectors.Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})

	wips := make([]WorkInProgress, len(progressEntries))
	for i := range progressEntries {
		entry := strings.TrimSpace(progressEntries[i])
		submatchResult := entryRegex.FindStringSubmatch(entry)
		if len(submatchResult) != 3 {
			log.Fatalf("Could not parse title and progress from progress entry: \"%s\"", entry)
		}

		title := strings.TrimSpace(submatchResult[1])
		progress, err := strconv.Atoi(strings.TrimSpace(submatchResult[2]))

		if err != nil {
			log.Fatal("Error parsing progress for ", entry)
		} else {
			wips[i] = WorkInProgress{Title: title, Progress: progress}
		}
	}

	return wips
}

package progress

import (
	"reflect"
	"testing"
)

const htmlScrape = `<div
  class="section-template--23098442055954__ss_progress_circles_87N98z progress-template--23098442055954__ss_progress_circles_87N98z"
  style="background-color:#ffffff; background-image: ;">
  <div class="section-template--23098442055954__ss_progress_circles_87N98z-settings">
    <div class="progress-heading-template--23098442055954__ss_progress_circles_87N98z">
      <p>BRANDON'S PROGRESS</p>
    </div>
    <div class="progress-items-template--23098442055954__ss_progress_circles_87N98z">
      <div class="progress-item-template--23098442055954__ss_progress_circles_87N98z progress-item-uniq in-view">
        <div class="progress-bar-template--23098442055954__ss_progress_circles_87N98z">
          <div class="progress-circle-progress_circle_T3mtH3"></div>
          <p class="progress-percent-template--23098442055954__ss_progress_circles_87N98z">100%</p>
        </div>
        <p class="progress-title-template--23098442055954__ss_progress_circles_87N98z">Moment Zero 2.0</p>
      </div>
      <div class="progress-item-template--23098442055954__ss_progress_circles_87N98z progress-item-uniq in-view">
        <div class="progress-bar-template--23098442055954__ss_progress_circles_87N98z">
          <div class="progress-circle-progress_circle_knrNWB"></div>
          <p class="progress-percent-template--23098442055954__ss_progress_circles_87N98z">81%</p>
        </div>
        <p class="progress-title-template--23098442055954__ss_progress_circles_87N98z">Blightfall 3.0 (Skyward Legacy)
        </p>
      </div>
      <div class="progress-item-template--23098442055954__ss_progress_circles_87N98z progress-item-uniq in-view">
        <div class="progress-bar-template--23098442055954__ss_progress_circles_87N98z">
          <div class="progress-circle-progress_circle_RUtA3D"></div>
          <p class="progress-percent-template--23098442055954__ss_progress_circles_87N98z">0%</p>
        </div>
        <p class="progress-title-template--23098442055954__ss_progress_circles_87N98z">White Sand (Prose Version)</p>
      </div>
      <div class="progress-item-template--23098442055954__ss_progress_circles_87N98z progress-item-uniq in-view">
        <div class="progress-bar-template--23098442055954__ss_progress_circles_87N98z">
          <div class="progress-circle-progress_circle_chtWGi"></div>
          <p class="progress-percent-template--23098442055954__ss_progress_circles_87N98z">34%</p>
        </div>
        <p class="progress-title-template--23098442055954__ss_progress_circles_87N98z">Words of Radiance $650 Signed Edition tier</p>
      </div>
    </div>
  </div>
</div>`

func TestParseProgressFromHTML(t *testing.T) {
	wips, err := parseProgressFromHTML(htmlScrape)
	if err != nil {
		t.Fatal(err)
	}

	expectedWips := []WorkInProgress{
		{Title: "Moment Zero 2.0", Progress: 100},
		{Title: "Blightfall 3.0 (Skyward Legacy)", Progress: 81},
		{Title: "White Sand (Prose Version)", Progress: 0},
		{Title: "Words of Radiance $650 Signed Edition tier", Progress: 34},
	}
	if !reflect.DeepEqual(wips, expectedWips) {
		t.Fatalf("mismatch\nexpected\t%v\ngot\t\t\t%v", expectedWips, wips)
	}
}

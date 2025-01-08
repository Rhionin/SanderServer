package progress

import (
	"reflect"
	"testing"
)

func TestProgressUpdate_String(t *testing.T) {
	testCases := []struct {
		name     string
		pu       ProgressUpdate
		expected string
	}{
		{
			name: "No previous progress",
			pu: ProgressUpdate{
				Title:    "Task A",
				Progress: 50,
			},
			expected: "Task A (50%)",
		},
		{
			name: "With previous progress",
			pu: ProgressUpdate{
				Title:        "Task B",
				Progress:     100,
				PrevProgress: 75,
			},
			expected: "Task B (75% => 100%)",
		},
		{
			name: "Previous progress matches current progress",
			pu: ProgressUpdate{
				Title:        "Task B",
				Progress:     75,
				PrevProgress: 75,
			},
			expected: "Task B (75%)",
		},
		{
			name: "Zero progress",
			pu: ProgressUpdate{
				Title:    "Task C",
				Progress: 0,
			},
			expected: "Task C (0%)",
		},
		{
			name: "Zero previous progress",
			pu: ProgressUpdate{
				Title:        "Task D",
				Progress:     75,
				PrevProgress: 0,
			},
			expected: "Task D (75%)",
		},
		{
			name: "Empty Title",
			pu: ProgressUpdate{
				Title:    "",
				Progress: 25,
			},
			expected: " (25%)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.pu.String()
			if actual != tc.expected {
				t.Errorf("Expected: %q, Actual: %q", tc.expected, actual) // Use %q for quoted strings
			}
		})
	}
}

func TestGetProgressUpdate(t *testing.T) {
	testCases := []struct {
		name           string
		latestProgress []WorkInProgress
		prevProgress   []WorkInProgress
		expected       []ProgressUpdate
	}{
		{
			name: "No previous progress",
			latestProgress: []WorkInProgress{
				{Title: "Task A", Progress: 50},
				{Title: "Task B", Progress: 100},
			},
			prevProgress: []WorkInProgress{
				{Title: "Task C", Progress: 20}, // No matching titles
				{Title: "Task D", Progress: 0},
			},
			expected: []ProgressUpdate{
				{Title: "Task A", Progress: 50, PrevProgress: 0},
				{Title: "Task B", Progress: 100, PrevProgress: 0},
			},
		},
		{
			name: "Some previous progress",
			latestProgress: []WorkInProgress{
				{Title: "Task A", Progress: 75},
				{Title: "Task B", Progress: 100},
			},
			prevProgress: []WorkInProgress{
				{Title: "Task A", Progress: 50},
				{Title: "Task B", Progress: 75},
			},
			expected: []ProgressUpdate{
				{Title: "Task A", Progress: 75, PrevProgress: 50},
				{Title: "Task B", Progress: 100, PrevProgress: 75},
			},
		},
		{
			name:           "Empty input",
			latestProgress: []WorkInProgress{},
			prevProgress:   []WorkInProgress{},
			expected:       []ProgressUpdate{},
		},
		{
			name: "Different order",
			latestProgress: []WorkInProgress{
				{Title: "Task B", Progress: 100},
				{Title: "Task A", Progress: 75},
			},
			prevProgress: []WorkInProgress{
				{Title: "Task A", Progress: 50},
				{Title: "Task B", Progress: 75},
			},
			expected: []ProgressUpdate{
				{Title: "Task B", Progress: 100, PrevProgress: 75},
				{Title: "Task A", Progress: 75, PrevProgress: 50},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetProgressUpdate(tc.latestProgress, tc.prevProgress)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected: %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}

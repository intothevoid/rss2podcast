package rss

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestParser_ParseURL(t *testing.T) {
	parser := &Parser{
		fp: &gofeed.Parser{},
	}

	tests := []struct {
		name     string
		feedURL  string
		expected []*gofeed.Item
		err      error
	}{
		{
			name:     "TestParseURL",
			feedURL:  "https://www.flipboard.com/topic/world.rss",
			expected: []*gofeed.Item{},
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := parser.ParseURL(tt.feedURL)
			if err != tt.err {
				t.Errorf("Unexpected error: got %v, want %v", err, tt.err)
			}
			if len(items) <= 0 {
				t.Errorf("Unexpected number of items: got %d, want > 0", len(items))
			}
		})
	}
}

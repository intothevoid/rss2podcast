package rss

import "strings"

type RSSItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	HtmlContent string `json:"html_content"`
}

func (rssi *RSSItem) IsFiltered(filterList []string) bool {
	// if rssi.Title contains any of the filterList, return true
	// Convert the title to lowercase to make the comparison case-insensitive
	for _, filter := range filterList {
		if filter != "" && strings.Contains(strings.ToLower(rssi.Title), strings.ToLower(filter)) {
			return true
		}
	}
	return false
}

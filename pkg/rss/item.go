package rss

type RSSItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	HtmlContent string `json:"html_content"`
}

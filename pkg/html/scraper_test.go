package html

import "testing"

func TestScrape(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestScrape",
			args: args{url: "https://news.google.com/rss/articles/CBMitAFodHRwczovL3d3dy5ldXJvc3BvcnQuY29tL21vdG8vcG9ydHVnYWwtZ3JhbmQtcHJpeC8yMDI0L21hcmMtbWFycXVlei1hdC1wb3J0dWd1ZXNlLWdyYW5kLXByaXgtaS13aWxsLW5vdC1iZS1mYXN0ZXItdGhhbi1pbi10aGUtcGFzdC1idXQtaS1oYXZlLWV4cGVyaWVuY2UtYW5fc3RvMTAwNzM2Njgvc3Rvcnkuc2h0bWzSAQA?oc=5"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Scrape(tt.args.url); got != tt.want {
				t.Errorf("Scrape() = %v, want %v", got, tt.want)
			}
		})
	}
}

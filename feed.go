package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Podcast stores data used in the <channel> section of RSS feed.
type Podcast struct {
	rssURL, atomURL string
	// Title: <title>
	Title string
	// Summary: <itunes:summary>
	Summary string
	// LogoURL: <image><url>
	LogoURL string
	// PageName: after last / in <link>
	PageName string
	Episodes []Episode
}

func (p *Podcast) FeedURL() string {
	n := strings.LastIndex(p.PageName, "/")
	if n == -1 {
		log.Fatalf("Bad PageURL %q", p.PageName)
	}
	name := p.PageName[n+1:]
	base := "https://localhost:8080/rss/"
	return base + name
}

// Episode stores unique attributes for each <item> in the RSS feed.
// The fields are derived from an <entry> in the atom feed.
type Episode struct {
	ArticleID   string
	Title       string
	Subtitle    string
	Description string
	Duration    time.Duration
	Published   time.Time
	Mp3ID       string
	Mp3Size     uint
}

func (e *Episode) PublishedTimestamp() string {
	return e.Published.Local().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

func (e *Episode) DurationFmt() string {
	d := e.Duration
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

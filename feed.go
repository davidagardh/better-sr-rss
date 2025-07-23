package main

import (
	"fmt"
	"time"
)

type Podcast struct {
	ID       string
	Title    string
	Subtitle string
	LogoURL  string
	Episodes []Episode
}

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

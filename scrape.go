package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func New(atomUrl, rssURL string) Podcast {
	p := Podcast{atomURL: atomUrl, rssURL: rssURL}
	p.UpdatePodcastData()

	return p
}

func (p *Podcast) UpdatePodcastData() {
	resp, err := http.Get(p.rssURL)
	if err != nil {
		log.Fatalln("Failed getting RSS feed\n" + err.Error())
	}
	if resp.StatusCode != 200 {
		log.Fatalln("Failed getting RSS feed\n" + resp.Status)
	}
	feed, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Failed reading RSS feed\n" + err.Error())
	}
	p.parseRSS(string(feed))
}

func (p *Podcast) parseRSS(feed string) {
	p.Title = getTagContents(feed, "title")
	p.Summary = getTagContents(feed, "itunes:summary")
	p.LogoURL = getTagContents(feed, "url")
	link := getTagContents(feed, "link")
	n := strings.LastIndex(link, "/")
	p.PageName = link[n+1:]
}

func (p *Podcast) UpdateEpisodesData() {
	resp, err := http.Get(p.atomURL)
	if err != nil {
		log.Fatalln("Failed getting atom feed\n" + err.Error())
	}
	if resp.StatusCode != 200 {
		log.Fatalln("Failed getting atom feed\n" + resp.Status)
	}
	feed, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Failed reading atom feed\n" + err.Error())
	}
	p.parseAtom(string(feed))
}

func (p *Podcast) parseAtom(feed string) {

}

func getTagContents(source, tag string) string {
	_, remain, found := strings.Cut(source, "<"+tag+">")
	if !found {
		return ""
	}
	contents, _, found := strings.Cut(remain, "</"+tag+">")
	if !found {
		return ""
	}
	return contents
}

package main

import (
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// UpdateEpisodesData parses the feed at RssURL and updates all fields except Episodes.
func (p *Podcast) UpdatePodcastData() {
	resp, err := Get(p.RssURL)
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
	p.atomURL = "https://api.sr.se/api/rss/program/" + getTagContents(feed, "sr:programid")
}

// UpdateEpisodesData parses the feed at atomURL and updates Episodes.
func (p *Podcast) UpdateEpisodesData() {
	resp, err := Get(p.atomURL)
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

	articles := getLinkHrefs(string(feed))

	p.Episodes = p.Episodes[:0]
	for _, v := range articles {
		r, err := Get(v)
		if err != nil {
			log.Println("While parsing " + p.Title)
			log.Fatalln("Failed to get article " + v + ": " + err.Error())
		}
		if r.StatusCode != 200 {
			log.Println("While parsing " + p.Title)
			log.Println("Failed to get article " + v + ": " + r.Status)
			continue
		}
		text, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln("Failed to read article " + err.Error())
		}
		episode := parseArticleHtml(string(text))

		episode.ArticleID = getAfterLast(v, "/")
		episode.Mp3Size = getContetLenght("https://sverigesradio.se/topsy/ljudfil/srse/" + episode.Mp3ID + ".mp3")
		p.Episodes = append(p.Episodes, episode)
	}

	log.Println("Done parsing " + p.Title)
}

func getLinkHrefs(feed string) []string {
	r := regexp.MustCompile(`\s*<link href="(.*?)"\s?/>`)
	matches := r.FindAllStringSubmatch(feed, -1)
	hrefs := make([]string, len(matches))
	for i, v := range matches {
		hrefs[i] = v[1]
	}
	return hrefs
}

func parseArticleHtml(text string) (e Episode) {
	// TOOD fix for nextJS rewrite
	relevantJson := getFirstCaptureMatch(text, `({\\"contentTitle\\":.+?)</script>`)
	if relevantJson == "" {
		//log.Println(text)
		log.Fatalln("Failed to find relevant json in article")
	}

	e.Mp3ID = getFirstCaptureMatch(relevantJson, `https://www.sverigesradio.se/topsy/ljudfil/(\d+)`)

	e.Title = getFirstCaptureMatch(relevantJson, `\\"contentTitle\\":\\"\\\\\\"(.+?)\\\\\\"\\",`)

	datePublished := getFirstCaptureMatch(relevantJson, `\\"publishUtc\\":\\"([0-9\-:TZ]+?)\\",`)
	published, err := time.Parse("2006-01-02T15:04:05Z", datePublished)
	if err != nil {
		log.Printf("Failed to parse date published. Article %v\n", e.Title)
		e.Published = time.Now()
	} else {
		e.Published = published
	}

	e.EpisodeImageURL = getFirstCaptureMatch(relevantJson, `\\"imageUrl\\":\\"(.+?)\\",`)

	e.Subtitle = getFirstCaptureMatch(text, `\\"subtitle\\":\\"(.+?)\\",`)

	description := `"` + getFirstCaptureMatch(text, `\\"fullDescription\\":\\"(.+?)\\",`) + `"`
	// Decode unicode characters
	err = json.Unmarshal([]byte(description), &e.Description)
	if err != nil {
		log.Println("Failed to decode description")
		e.Description = description
	}

	seconds := getFirstCaptureMatch(text, `\\"duration\\":(\d+),`)
	duration, _ := strconv.Atoi(seconds)
	e.Duration = time.Second * time.Duration(duration)

	return e
}

func getContetLenght(url string) uint {
	resp, err := Head(url)
	if err != nil {
		log.Println("Could not get podcast " + url)
		return 0
	}
	size, _ := strconv.Atoi(resp.Header.Get("content-length"))
	if size == 0 {
		log.Fatalln(url)
	}
	return uint(size)
}

// getTagContents searches source for the first occurance of <tag> and returns all the text from there to </tag> not including the tags.
// If starting or ending tags are not found, returns empty string.
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

func getTagWithAttrubutesContents(source, tag, attritbutes string) string {
	_, remain, found := strings.Cut(source, "<"+tag+" "+attritbutes+">")
	if !found {
		return ""
	}
	contents, _, found := strings.Cut(remain, "</"+tag+">")
	if !found {
		return ""
	}
	return contents
}

func getFirstCaptureMatch(text, pattern string) string {
	r := regexp.MustCompile(pattern)
	match := r.FindStringSubmatch(text)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func getAfterLast(text, substring string) string {
	parts := strings.Split(text, substring)
	if len(parts) == 1 {
		return ""
	}
	return parts[len(parts)-1]
}

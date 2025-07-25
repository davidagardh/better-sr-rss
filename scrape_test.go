package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetTagContents(t *testing.T) {
	source := "test<search>data</search>test"
	got := getTagContents(source, "search")
	require.Equal(t, "data", got)
}

func TestGetTagWithAttributesContents(t *testing.T) {
	source := "test<search type=\"text\">data</search>test"
	got := getTagWithAttrubutesContents(source, "search", "type=\"text\"")
	require.Equal(t, "data", got)
}

func TestParseRSS(t *testing.T) {
	source := `<rss version="2.0">
	<title>A Title</title>
	<itunes:summary>A long summary</itunes:summary>
	<image>
		<link>link.page/api/name</link>
		<url>imageurl</url>
	</image>
	<sr:programid>5419</sr:programid>
	`
	got := Podcast{}
	got.parseRSS(source)
	require.Equal(t, Podcast{Title: "A Title", Summary: "A long summary", LogoURL: "imageurl", PageName: "name", atomURL: "https://api.sr.se/api/rss/program/5419"}, got)
}

func TestGetLinkHrefs(t *testing.T) {
	source := `<id>rss:sr.se/article/8928686</id>
	<link href="link1" />
	<category term="" />
	<link href="link2"/>
	<link href="link3" />
	`
	got := getLinkHrefs(source)
	require.Equal(t, []string{"link1", "link2", "link3"}, got)
}

func TestParseArticleHtml(t *testing.T) {
	source := `
	<a href="//sverigesradio.se/topsy/ljudfil/srse/id.mp3"
	<meta name="description" content="my subtitle" />
	<script type="application/ld+json">
		{"@type":"NewsArticle","headline":"A Title","logo":{"@type":"ImageObject","url":"wrongurl"}},"image":{"@type":"ImageObject","url":"image.url/"},"datePublished":"2025-06-06 14:00:00Z"}
	</script>
	<div>
	<div class="publication-text text-editor-content" >Episode description.</div>            </div>
	<abbr title="29 minuter">29 min</abbr>
	`
	expected := Episode{
		Mp3ID:           "id",
		Title:           "A Title",
		Subtitle:        "my subtitle",
		Description:     "Episode description.",
		Duration:        time.Minute * 29,
		Published:       time.Date(2025, 6, 6, 14, 0, 0, 0, time.UTC),
		EpisodeImageURL: "image.url/",
	}
	got := parseArticleHtml(source)
	require.Equal(t, expected, got)
}

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTagContents(t *testing.T) {
	source := "test<search>data</search>test"
	got := getTagContents(source, "search")
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
	`
	got := Podcast{}
	got.parseRSS(source)
	require.Equal(t, Podcast{Title: "A Title", Summary: "A long summary", LogoURL: "imageurl", PageName: "name"}, got)
}

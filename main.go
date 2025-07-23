package main

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

func main() {
	testPodcast := New("https://api.sr.se/api/rss/program/4923", "https://api.sr.se/api/rss/pod/22712")

	tmpl := template.Must(template.New("sample.rss").
		Funcs(template.FuncMap{
			"now": Timestamp,
		}).ParseFiles("sample.rss"))

	http.HandleFunc("/rss/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		tmpl.Execute(w, testPodcast)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func Timestamp() string {
	return time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

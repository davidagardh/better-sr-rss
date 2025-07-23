package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Podcasts map[string]Podcast

func (ps *Podcasts) Add(rssUrl string) {
	p := Podcast{RssURL: rssUrl}
	p.UpdatePodcastData()
	(*ps)[p.PageName] = p
}

func main() {
	podcasts := make(Podcasts)
	podcasts.Add("https://api.sr.se/api/rss/pod/22712")
	podcasts.Add("https://api.sr.se/api/rss/pod/34530")
	podcasts.Add("https://api.sr.se/api/rss/pod/4023")
	index := make([]string, len(podcasts))
	for k := range podcasts {
		index = append(index, k)
	}
	fmt.Println(index)

	tmpl := template.Must(template.New("sample.rss").
		Funcs(template.FuncMap{
			"now": Timestamp,
		}).ParseFiles("sample.rss"))

	http.HandleFunc("/rss/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		p, ok := podcasts[name]
		if !ok {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(name + " not found"))
			w.WriteHeader(404)
			return
		}
		w.Header().Set("Content-Type", "text/xml")
		tmpl.Execute(w, p)
	})

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func Timestamp() string {
	return time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

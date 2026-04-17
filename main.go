package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

type Podcasts map[string]*Podcast

func (ps *Podcasts) Add(rssUrl string) {
	p := Podcast{RssURL: rssUrl}
	p.UpdatePodcastData()
	(*ps)[p.PageName] = &p
}

func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed creating request url %q\n%v\n", url, err.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	return http.DefaultClient.Do(req)
}

func Head(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed creating request url %q\n%v\n", url, err.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Method = "HEAD"

	return http.DefaultClient.Do(req)
}

func main() {
	podcasts := make(Podcasts)
	podcasts.Add("https://api.sr.se/api/rss/pod/22712")
	podcasts.Add("https://api.sr.se/api/rss/pod/34530")
	podcasts.Add("https://api.sr.se/api/rss/pod/44054")
	podcasts.Add("https://api.sr.se/api/rss/pod/21597")
	index := make([]string, len(podcasts))
	for k := range podcasts {
		index = append(index, k)
		go podcasts[k].UpdateEpisodesData()
	}
	fmt.Println(index)

	tmpl := template.Must(template.New("podcast.rss").
		Funcs(template.FuncMap{
			"now": Timestamp,
		}).ParseFiles("templates/podcast.rss"))

	http.HandleFunc("/rss/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		for _, v := range index {
			w.Write([]byte(v + "\n"))
		}
	})

	go refreshEpisodesLoop(podcasts)

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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func Timestamp() string {
	return time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

func refreshEpisodesLoop(ps Podcasts) {
	env := os.Getenv("REFRESH_INTERVAL")
	waitTime, err := strconv.Atoi(env)
	if err != nil {
		waitTime = 12
	}
	log.Printf("Refresh interval set to %d hours\n", waitTime)
	for {
		time.Sleep(time.Duration(waitTime) * time.Hour)
		log.Println("Updating episodes...")
		for k := range ps {
			ps[k].UpdateEpisodesData()
		}
		log.Println("Done updating episodes!")
	}
}

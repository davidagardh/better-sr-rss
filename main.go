package main

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

func main() {
	testPodcast := Podcast{
		ID:       "uuid:53e5ae1b-39ab-467e-a2b9-aaa027e8d087;id=19535",
		Title:    "USApodden (better)",
		Subtitle: "Sveriges största och skarpaste podd om amerikansk politik.",
		LogoURL:  "https://static-cdn.sr.se/images/4923/0ff122cc-8c46-4f0e-9142-df6e0d51d912.jpg?preset=api-itunes-presentation-image",
		Episodes: []Episode{
			{
				ArticleID:   "8974370",
				Title:       "Trump och Musk: episk skilsmässa inför öppen ridå",
				Duration:    time.Duration(time.Minute * 29),
				Subtitle:    "Radarparet&amp;nbsp;ryker ihop.",
				Mp3ID:       "9792414",
				Mp3Size:     27855651,
				Published:   time.Now().Add(-time.Hour),
				Description: "<p>Radarparet ryker ihop.</p><p> <a href=\"https://sverigesradio.se/play/program/4923?utm_source=thirdparty&utm_medium=rss&utm_campaign=episode_usapodden\">Lyssna på alla avsnitt i Sveriges Radio Play.</a></p> <p>Världens rikaste man frontalkraschade med världens mäktigaste man framför ögonen på en chockad omvärld. Donald Trumps och Elon Musks vänskap är definitivt över och påhoppen haglar på respektive sociala medier-plattform. </p><p>Vad var det egentligen som hände – ett underhållande gräl på internet, eller en konflikt som kommer få politiska konsekvenser? Det och mycket mer behandlar vi i det här avsnittet av USA-podden.</p><p>Medverkande: Ginna Lindberg och Roger Wilson, Sveriges Radios USA-korrespondenter</p><p>Programledare: Sara Stenholm</p><p>Producent: Anna Roxvall</p>",
			},
			{
				ArticleID:   "8981044",
				Title:       "Så ska Trump ducka storkriget",
				Duration:    time.Duration(time.Minute * 32),
				Subtitle:    "Eskaleringen mellan Israel och Iran sätter press på USA.",
				Mp3ID:       "9802166",
				Mp3Size:     54300000,
				Published:   time.Now().Add(-30 * time.Minute),
				Description: "&lt;p&gt;&lt;img src=\"https://static-cdn.sr.se/images/4923/0ff122cc-8c46-4f0e-9142-df6e0d51d912.jpg?preset=api-default-rectangle\" /&gt;&lt;/p&gt;Eskaleringen mellan Israel och Iran sätter press på USA.&lt;p&gt;Israel har attackerat kärnvapenanläggningar och dödat flera höga militära ledare i Iran – mot USA:s vilja. Vad innebär det för relationen mellan länderna och kommer amerikanerna kunna hålla sig utanför ett annalkande storkrig?&lt;/p&gt;&lt;p&gt;Hör också om protesterna i Los Angeles, den juridiska kampen om nationalgardet, en nedbrottad senator och Trumps födelsedagspresent – den stundande militärparaden som lamslår Washington.&lt;br&gt;&lt;br&gt;Medverkande: Ginna Lindberg och Roger Wilson, Sveriges Radios USA-korrespondenter&lt;/p&gt;&lt;p&gt;Programledare: Sara Stenholm&lt;/p&gt;&lt;p&gt;Producent: Anna Roxvall&lt;/p&gt;",
			},
		},
	}

	tmpl := template.Must(template.New("sample.atom").
		Funcs(template.FuncMap{
			"now": Timestamp,
		}).ParseFiles("sample.atom"))

	http.HandleFunc("/api/rss/", func(w http.ResponseWriter, r *http.Request) {
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

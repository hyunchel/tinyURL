package tinyUrl

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"regexp"
)

func shortenUrl(originalUrl string) string {
	url := HashIn(originalUrl)
	log.Printf("Shortened %q to %q", originalUrl, url)
	if len(url) == 0 {
		// Throw an Error here.
		log.Printf("Unable to hash %q", originalUrl)
		return ""
	}
	return url
}

func lookupForUrl(shortenedUrl string) string {
	PrintHashMapper()
	url := HashOut(shortenedUrl)
	if len(url) == 0 {
		log.Printf("Tried to HashOut %q, but got %q", shortenedUrl, url)
		log.Println("Should throw 404.")
		return ""
	}
	return url
}

func CreateAndRunServer() {

	var shortenRegex = regexp.MustCompile(`\/shorten`)
	var redirectRegex = regexp.MustCompile(`\/redirect`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case shortenRegex.MatchString(r.URL.Path):
			r.ParseForm()
			var message string
			if longUrl := r.Form["url"]; longUrl != nil {
				if shortenedUrl := shortenUrl(longUrl[0]); shortenedUrl != "" {
					message = `
						Your shortened URL:
						%v
					`
					fmt.Fprintf(w, message, shortenedUrl)
				} else {
					message = `
						Sorry, failed to shorten %q
					`
					http.Error(w, fmt.Sprintf(message, longUrl), 400)
				}
			} else {
				message = `
					Sorry, "url" parameter is missing.
				`
				http.Error(w, message, 400)
			}
		case redirectRegex.MatchString(r.URL.Path):
			urlPath := redirectRegex.Split(r.URL.Path, 2)
			originalUrl := lookupForUrl(urlPath[1])
			code := 301
			if originalUrl == "" {
				http.NotFound(w, r)
			}
			http.Redirect(w, r, originalUrl, code)
		default:
			// FIXME: Use links.
			welcomeMessage := `
				Welcome to tinyUrl.
				In order to shorten, please head over to %q.
				To use the shortened URL, head over to %q.
			`
			fmt.Fprintf(w, welcomeMessage, "/shorten/", "/redirect/")
		}
	})

	s := &http.Server{
		Addr:	":8080",
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
		MaxHeaderBytes:	1 << 20,
	}
	log.Println("Listening at port 8080.")
	log.Fatal(s.ListenAndServe())
}
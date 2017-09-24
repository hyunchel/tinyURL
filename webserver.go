package tinyURL

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
	"regexp"
)

func shortenURL(originalURL string) string {
	url := HashIn(originalURL)
	if len(url) == 0 {
		log.Printf("Unable to hash %q", originalURL)
		return ""
	}
	return url
}

func lookupForURL(shortenedURL string) string {
	url := HashOut(shortenedURL)
	if len(url) == 0 {
		log.Printf("Tried to HashOut %q, but got %q", shortenedURL, url)
		log.Println("Should throw 404.")
		return ""
	}
	return url
}

func CreateAndRunServer() {

	var shortenRegex = regexp.MustCompile(`\/shorten`)
	var redirectRegex = regexp.MustCompile(`\/redirect\/`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case shortenRegex.MatchString(r.URL.Path):
			r.ParseForm()
			var message string
			if longURL := r.Form["url"]; longURL != nil {
				if shortenedURL := shortenURL(longURL[0]); shortenedURL != "" {
					resultURL := fmt.Sprintf("%v/redirect/%v", r.Host, shortenedURL)
					message := fmt.Sprintf("<p>Your shortened URL:</p><a href=%q>%v</a>", resultURL, resultURL)
					fmt.Fprintf(w, html.UnescapeString(message))
				} else {
					message = `
						Sorry, failed to shorten %q
					`
					http.Error(w, fmt.Sprintf(message, longURL), 400)
				}
			} else {
				message = `
					Sorry, "url" parameter is missing.
				`
				http.Error(w, message, 400)
			}
		case redirectRegex.MatchString(r.URL.Path):
			urlPath := redirectRegex.Split(r.URL.Path, 2)
			originalURL := lookupForURL(urlPath[1])
			code := 301
			if originalURL == "" {
				http.NotFound(w, r)
			}
			http.Redirect(w, r, originalURL, code)
		default:
			// FIXME: Use links.
			welcomeMessage := `
				Welcome to tinyURL.
				In order to shorten, please head over to %q.
				To use the shortened URL, head over to %q.
			`
			fmt.Fprintf(w, welcomeMessage, "/shorten/", "/redirect/")
		}
	})

	s := &http.Server{
		Addr:			":8080",
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
		MaxHeaderBytes:	1 << 20,
	}
	log.Println("Listening at port 8080...")
	log.Fatal(s.ListenAndServe())
}
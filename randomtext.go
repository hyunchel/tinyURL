package tinyURL

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"math/rand"
)

type Post struct {
	UserId int
	Id int
	Title string
	Body string
}

func getPostsJSON() []byte {
	res, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		log.Fatal(err)
	}
	posts, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return posts
}

func decodeJSON(postsJSON []byte) []Post {
	// Decode into an array of Post.
	var collection []Post
	err := json.Unmarshal(postsJSON, &collection)
	if err != nil {
		log.Fatal(err)
	}
	return collection
}

func GetRandomText(count int) []string {
	collection := decodeJSON(getPostsJSON())
	var randomTexts []string
	for i := 0; i < count; i++ {
		var text string
		for j := 0; j < 5; j++ {
			randomIndex := rand.Intn(len(collection))
			text += collection[randomIndex].Title
		}
		randomTexts = append(randomTexts, text)
	}
	return randomTexts
}
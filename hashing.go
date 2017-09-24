package tinyURL

import (
	"log"
	"hash/fnv"
	"encoding/hex"
)

// Global?
var hashMapper map[string]string = make(map[string]string)

func HashIn(text string) string {
	h := fnv.New32()
	_, err := h.Write([]byte(text))
	if err != nil {
		log.Printf("err: %q\n", err)
	}
	encodedString := hex.EncodeToString(h.Sum(nil))
	
	if hashMapper[encodedString] != "" && hashMapper[encodedString] != text {
		// Hash collision.
		log.Fatal("Hash collision occured!")
	}

	if err := SaveURL(encodedString, text); err != nil {
		hashMapper[encodedString] = text
	}
	return encodedString
}

func HashOut(hashedText string) string {
	url, err := GetURL(hashedText)
	if err != nil {
		log.Printf("Failed to get the original value from %q\n", hashedText)
		url = hashMapper[hashedText]
	}
	return url
}

func PrintHashMapper() {
	log.Println(hashMapper)
}
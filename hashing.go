package tinyUrl

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
	hashMapper[encodedString] = text
	return encodedString
}

func HashOut(hashedText string) string {
	url := hashMapper[hashedText]
	return url
}

func PrintHashMapper() {
	log.Println(hashMapper)
}
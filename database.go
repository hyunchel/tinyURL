package tinyURL

import (
	"fmt"
	"errors"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

func connect() (redis.Conn, error) {
	var url string
	url = os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://"
	}
	conn, err := redis.DialURL(url)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to establish Redis connection with %q", url)
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}
	return conn, nil
}

func SaveURL(shortURL string, longURL string) error {
	conn, err := connect()
	defer conn.Close()
	if err != nil {
		return err
	}
	value, err := redis.Bool(conn.Do("SET", shortURL, longURL))
	if err != nil || value != true{
		errorMessage := fmt.Sprintf("Failed to  set %q as %q to the database.", shortURL, longURL)
		log.Print(errorMessage)
		return errors.New(errorMessage)
	}
	return nil
}

func GetURL(shortURL string) (string, error) {
	conn, err := connect()
	defer conn.Close()
	if err != nil {
		return "", err
	}
	value, err := redis.String(conn.Do("GET", shortURL))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get %q from the database.", shortURL)
		log.Printf(errorMessage)
		return "", errors.New(errorMessage)
	}
	return value, nil
}
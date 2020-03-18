package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandStrWithCharset returns a random string of length as mentioned in the argument, created
// out of characters in the charset string
func RandStrWithCharset(length int, charset string) string {
	rand := make([]byte, length)
	for i := range rand {
		rand[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(rand)
}

func RandStr(length int) string {
	return RandStrWithCharset(length, charset)
}

type urlShortener struct {
	urlMap   map[string]string
	checkMap map[string]string
	mutex    sync.Mutex
}

type Request struct {
	Url string `json:"url"`
}

func NewURLShortener() *urlShortener {
	return &urlShortener{
		urlMap:   make(map[string]string),
		checkMap: make(map[string]string),
	}
}

func MapUrl(originalURL string, u *urlShortener) error {
	return nil
}

func GetShortenedURL(nus *urlShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---Running in GetShortenedURL---")
		var reqBody Request
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			fmt.Println("Error in parsing request body: ", err)
			return
		}
		fmt.Println("url input: ", reqBody.Url)
		if _, ok := nus.checkMap[reqBody.Url]; ok {
			w.Write([]byte("Provided url exists. Shorted URL is: " + nus.checkMap[reqBody.Url]))
			return
		}
		short := RandStr(5)

		nus.mutex.Lock()
		nus.urlMap[short] = reqBody.Url
		nus.checkMap[reqBody.Url] = short
		nus.mutex.Unlock()
		w.Write([]byte("The shortened URL is: " + short))
	}
}

func getOriginalURL(nus *urlShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---Running in getOriginalURL---")
		vars := mux.Vars(r)
		short, ok := vars["short"]
		if !ok {
			fmt.Println("Could not retrive the short URL: ", short)
			return
		}
		w.Write([]byte(nus.urlMap[short]))
	}
}
func main() {
	nus := NewURLShortener()
	router := mux.NewRouter()
	router.HandleFunc("/short", GetShortenedURL(nus)).Methods("POST")
	router.HandleFunc("/original/{short}", getOriginalURL(nus))
	log.Fatal(http.ListenAndServe(":8192", router))
}

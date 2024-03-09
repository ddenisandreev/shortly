package main

import (
	"fmt"
	"strings"
)

const ALPHABET string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var urls = map[int]string{}
var lastId = 0

func main() {
	url := "https://google.com"
	shortUrl := getShortUrl(url)
	fmt.Println(shortUrl)
	fullUrl := getFullUrl(shortUrl)
	fmt.Println(fullUrl)
}

func getShortUrl(url string) string {
	lastId++
	urls[lastId] = url
	return idToShortUrl(lastId)
}

func getFullUrl(shortUrl string) string {
	var id = shortUrlToId(shortUrl)
	return urls[id]
}

func idToShortUrl(id int) string {
	var shortUrl = ""
	for id > 0 {
		shortUrl = string(ALPHABET[id%62]) + shortUrl
		id = id / 62
	}
	return shortUrl
}

func shortUrlToId(shortUrl string) int {
	var id = 0
	for i := 0; i < len(shortUrl); i++ {
		id = id*len(ALPHABET) + strings.Index(ALPHABET, string(shortUrl[i]))
	}
	return id
}

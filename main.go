package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const ALPHABET string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var db *sql.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	log.Println(psqlConn)

	var err error
	db, err = sql.Open("postgres", psqlConn)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/shortly", func(w http.ResponseWriter, r *http.Request) {
		originalUrl := r.URL.Query().Get("url")
		if originalUrl != "" {
			fmt.Fprintf(w, "original : %q short : %q", originalUrl, getShortUrl(originalUrl))
		} else {
			fmt.Fprintf(w, "set req param like ?url=<your_url> to get short url")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortUrl := html.EscapeString(r.URL.Path)[1:]
		fullUrl := getFullUrl(shortUrl)
		http.Redirect(w, r, fullUrl, http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getShortUrl(url string) string {
	lastInsertId := 0
	err := db.QueryRow("INSERT INTO shortly.urls (url_) VALUES($1) RETURNING id", url).Scan(&lastInsertId)

	if err != nil {
		log.Fatal(err)
	}

	return idToShortUrl(lastInsertId)
}

func getFullUrl(shortUrl string) string {
	var id = shortUrlToId(shortUrl)

	rowsRs, err := db.Query("SELECT url_ FROM shortly.urls where id = $1", id)
	if err != nil {
		log.Fatal(rowsRs)
	}
	var url_ string
	for rowsRs.Next() {
		err = rowsRs.Scan(&url_)
		if err != nil {
			log.Fatal(err)
		}
	}

	return url_
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

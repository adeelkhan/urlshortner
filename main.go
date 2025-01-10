package main

import (
	"fmt"
	"net/http"

	"github.com/adeelkhan/urlshortner/shortener"
)

func main() {
	fmt.Println("URL Shortner Service Starting at 8090")
	http.HandleFunc("/s", shortener.GetShortUrl)
	http.HandleFunc("/l", shortener.GetLongUrl)
	http.ListenAndServe(":8090", nil)
}

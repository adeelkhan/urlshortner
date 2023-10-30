package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const URL_LENGTH = 10

type UrlShortnerRequest struct {
	LongUrl string `json:"long_url"`
}

type UrlShortnerResponse struct {
	ShortUrl string `json:"short_url"`
}

type UrlLongerRequest struct {
	ShortUrl string `json:"short_url"`
}

type UrlLongerResponse struct {
	LongUrl string `json:"long_url"`
}

// store for urls
var UrlMap = map[string]string{}

func GetShortUrl(w http.ResponseWriter, req *http.Request) {
	var request = UrlShortnerRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := url.Parse(request.LongUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	query := make([]string, 0)
	query = append(query, url.Path, url.RawQuery)

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(query, "")))
	sumHex := hex.EncodeToString(hasher.Sum(nil))
	shortUrl := strings.TrimSpace(sumHex)[0:URL_LENGTH]

	urlResponse := &UrlShortnerResponse{
		ShortUrl: fmt.Sprintf(`http://localhost:8090/%s`, shortUrl),
	}

	UrlMap[urlResponse.ShortUrl] = url.String()
	responseByte, _ := json.Marshal(urlResponse)

	fmt.Fprintf(w, string(responseByte))
}

func GetLongUrl(w http.ResponseWriter, req *http.Request) {
	var request = UrlLongerRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := url.Parse(request.ShortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	var longUrl string
	for k,v := range UrlMap {
		if k == url.String() {
			longUrl = v 
			break
		}
	}

	urlResponse := &UrlLongerResponse{
		LongUrl: longUrl,
	}

	responseByte, _ := json.Marshal(urlResponse)
	fmt.Fprintf(w, string(responseByte))
}

func GetAllUrls(w http.ResponseWriter, req *http.Request) {
	for k,v := range UrlMap {
		fmt.Println(k,":",v)
	}
	fmt.Fprintf(w, "longurl")
}

func main() {
	fmt.Println("URL Shortner Service Starting at 8090")

	http.HandleFunc("/shorturl", GetShortUrl)
	http.HandleFunc("/longurl", GetLongUrl)
	http.HandleFunc("/allurls", GetAllUrls)
	http.ListenAndServe(":8090", nil)

}
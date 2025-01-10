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
	query = append(query, url.String())

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(query, "")))
	sumHex := hex.EncodeToString(hasher.Sum(nil))
	shortUrl := strings.TrimSpace(sumHex)[0:URL_LENGTH]

	urlResponse := &UrlShortnerResponse{
		ShortUrl: fmt.Sprintf(`http://localhost:8090/l?q=%s`, shortUrl),
	}
	UrlMap[shortUrl] = url.String()
	responseByte, _ := json.Marshal(urlResponse)

	fmt.Fprintf(w, string(responseByte))
}

func GetLongUrl(w http.ResponseWriter, req *http.Request) {
	requestURL := req.URL.RequestURI()
	_, err := url.Parse(requestURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url_code := req.URL.Query()
	if _, ok := url_code["q"]; !ok {
		http.Error(w, "Invalid query", http.StatusBadRequest)
		return
	}
	if len(url_code["q"]) == 0 {
		http.Error(w, "Invalid query", http.StatusBadRequest)
		return
	}
	hash_code := url_code["q"][0]
	var longUrl string
	for k, v := range UrlMap {
		if k == hash_code {
			longUrl = v
			break
		}
	}
	http.Redirect(w, req, longUrl, http.StatusPermanentRedirect)
}

func main() {
	fmt.Println("URL Shortner Service Starting at 8090")
	http.HandleFunc("/s", GetShortUrl)
	http.HandleFunc("/l", GetLongUrl)
	http.ListenAndServe(":8090", nil)
}

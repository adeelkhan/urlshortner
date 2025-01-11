package shortener

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

	urlWithoutScheme := strings.Split(url.String(), "://")
	if len(urlWithoutScheme) == 1 {
		http.Error(w, "Invalid Url format", http.StatusBadRequest)
		return
	}
	query := make([]string, 0)
	query = append(query, urlWithoutScheme[1])

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(query, "")))
	sumHex := hex.EncodeToString(hasher.Sum(nil))
	shortUrl := strings.TrimSpace(sumHex)[0:URL_LENGTH]

	urlResponse := &UrlShortnerResponse{
		ShortUrl: fmt.Sprintf(`http://localhost:8090/l?q=%s`, shortUrl),
	}
	UrlMap[shortUrl] = urlWithoutScheme[1]
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
	finalUrl := fmt.Sprintf("https://%s", longUrl)
	http.Redirect(w, req, finalUrl, http.StatusSeeOther)
}

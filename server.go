package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/jenovs/api-image-search/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

type ResponseData struct {
	Items []struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Pagemap struct {
			CseImage []struct {
				Src string `json:"src,omitempty"`
			} `json:"cse_image,omitempty"`
		} `json:"pagemap"`
	} `json:"items"`
}

type Image struct {
	Url   string `json:"page_url"`
	Src   string `json:"img_src"`
	Title string `json:"alt_text"`
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	return ":" + port
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query()

	if len(query) == 0 {
		return
	}

	p := map[string]string{
		"key":   os.Getenv("API_KEY"),
		"cx":    os.Getenv("SE_ID"),
		"q":     strings.Join(query["q"], ""),
		"start": strings.Join(query["offset"], ""),
	}

	if len(p["q"]) == 0 {
		return
	}

	if _, err := strconv.Atoi(p["start"]); len(p["start"]) == 0 || err != nil {
		p["start"] = "1"
	}

	rootUrl := "https://www.googleapis.com/customsearch/v1"
	u, _ := url.Parse(rootUrl)
	q := u.Query()
	for key, val := range p {
		q.Set(key, val)
	}
	u.RawQuery = q.Encode()

	resp, err := client.Get(u.String())

	if err != nil {
		fmt.Println("Error getting request", err)
	}


	var data ResponseData
	var images []interface{}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&data)

	for _, item := range data.Items {
		if len(item.Pagemap.CseImage) > 0 {
			img := Image{item.Link, item.Pagemap.CseImage[0].Src, item.Title}
			images = append(images, img)
		}
	}

	res, err := json.Marshal(images)

	if err != nil {
		fmt.Println("JSON marshal error: ", err)
		return
	}

	updateLatest(p["q"])

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func updateLatest(s string) {
	latestData, err := ioutil.ReadFile("latest")
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	latest := s + "\n\n" + string(latestData)

	ioutil.WriteFile("latest", []byte(latest), 0644)
}

func latestHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./latest")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/latest/", latestHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(getPort(), nil))
}

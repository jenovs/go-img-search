package main

import (
	"fmt"
	_ "github.com/jenovs/api-image-search/config"
	"log"
	"net/http"
	"os"
)

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	return ":" + port
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Env variables:")
	fmt.Println(os.Getenv("API_KEY"))
	fmt.Println(os.Getenv("SE_ID"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(getPort(), nil))
}

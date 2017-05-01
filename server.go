package main

import (
	"fmt"
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
	os.Setenv("TEST", "this is a test")
	fmt.Println(os.Getenv("TEST"))
	fmt.Println(os.Getenv("TEST2"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(getPort(), nil))
}

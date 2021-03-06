package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Millisecond * 50)
	fmt.Fprintf(w, strings.Repeat("haha", 1024))
}

func main() {
	http.HandleFunc("/", myHandler)
	log.Fatal(http.ListenAndServe(":3001", nil))
}

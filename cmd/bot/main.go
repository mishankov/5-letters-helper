package main

import (
	"fmt"
	"net/http"
	"time"
)

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Time: %v", time.Now())
}

func main() {
	http.HandleFunc("/healthcheck", healthcheck)
	http.ListenAndServe(":4444", nil)
}

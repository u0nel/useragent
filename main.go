package main

import (
	"log"
	"net/http"

	"github.com/u0nel/accept"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		types := []string{"text/plain"}
		switch accept.ServeType(types, r.Header.Get("Accept")) {
		case "text/plain":
			w.Write([]byte(r.Header.Get("User-Agent")))
		default:
			http.Error(w, "Could not serve requested Type", http.StatusNotAcceptable)
		}
	})
	log.Fatal(http.ListenAndServe(":8090", nil))
}

package main

import (
	"log"
	"net/http"

	"github.com/u0nel/accept"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		useragent := r.Header.Get("User-Agent")
		types := []string{"text/plain", "text/html", "application/json"}
		switch accept.ServeType(types, r.Header.Get("Accept")) {
		case "text/plain":
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte(r.Header.Get("User-Agent")))
		case "text/html":
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(`
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/water.css@2/out/water.css">
	<h1>What's your user Agent?</h1>
	<pre>` + useragent))
		default:
			http.Error(w, "Could not serve requested Type", http.StatusNotAcceptable)
		}
	})
	log.Fatal(http.ListenAndServe(":8090", nil))
}

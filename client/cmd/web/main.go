package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	const dir = "./client/wasm"
	port := ":8080"

	fs := http.FileServer(http.Dir(dir))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		}
		fs.ServeHTTP(w, r)
	}))

	log.Printf("Open: http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

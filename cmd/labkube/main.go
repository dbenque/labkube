package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request on URI /ready")
		w.WriteHeader(200)
	})
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request on URI /hello")
		fmt.Fprintf(w, "Hello I am pod %s\nWelcome to kubernetes lab.\n", os.Getenv("HOSTNAME"))
	})
	r.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", os.Environ())
	})
	r.HandleFunc("/book/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	log.Printf("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

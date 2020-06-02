package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Handler).Methods("GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

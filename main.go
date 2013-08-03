package main

import (
	"fractal"
	"log"
	"net/http"
	"image"
)

var _ = image.Rect(0,0,1,1)
var _ = fractal.SingPert{}

func init() {
	http.HandleFunc("/fractals", errorHandler(fractalHandler))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}

func fractalHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func main() {

	log.Fatal(http.ListenAndServe(":8899", nil))
}

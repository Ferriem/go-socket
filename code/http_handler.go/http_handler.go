package main

import (
	"fmt"
	"net/http"
	"os"
)

func CheckStatusOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `Fine!`)
}

func StatusNotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served: %s\n", r.Host)
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served: %s\n", r.Host)
}

func main() {
	PORT := ":8081"
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Using default port:", PORT)
	} else {
		PORT = ":" + arguments[1]
	}
	http.HandleFunc("/CheckStatusOK", CheckStatusOK)
	http.HandleFunc("/StatusNotFound", MyHandler)
	http.HandleFunc("/", StatusNotFound)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

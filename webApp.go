package main

import (
	"fmt"
	"net/http"
)
//main handler that takes users input and returns the chat bots response
func userInputHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "you said , %s", r.URL.Query().Get("value"))
}

func main() {

	
	fs := http.FileServer(http.Dir("Html"))
	http.Handle("/", fs)

	http.HandleFunc("/user-input", userInputHandler)
	http.ListenAndServe(":8080", nil)
}
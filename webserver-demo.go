package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", top)
	http.HandleFunc("/simple", simple)
	http.ListenAndServe(":8080", nil)
}

func top(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world\n")
}

func simple(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "simple-template.html", UserInfo{123, "mockers"})
}

// UserInfo is a bodged type for use investigating how templates work
type UserInfo struct {
	// Id is some sort of ID or reference number.
	// It will probably be replaced by a string or UUID of some sort
	ID       int
	Username string
}

func renderTemplate(w io.Writer, templateFilename string, user UserInfo) {
	t, err := template.ParseFiles(templateFilename)
	if err != nil {
		log.Fatal("Could not parse template", templateFilename)
	}
	t.Execute(w, user)
}

package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

const addr = ":8080"

func main() {
	printStartUpMessage(addr)

	http.HandleFunc("/", top)
	http.HandleFunc("/simple", simple)
	http.ListenAndServe(addr, nil)
}

func printStartUpMessage(addr string) {
	fmt.Println("Demo webserver")
	fmt.Printf("Visit http://localhost%v\n", addr)
	fmt.Println("Press ^C (Ctrl-C) to finish")
}

func top(w http.ResponseWriter, r *http.Request) {
	homePageText := `
<html>
<body>
<p>Hello world</p>
</body>
</html>
`
	fmt.Fprintf(w, homePageText)
}

func simple(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/simple-template.html", UserInfo{123, "mockers"})
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

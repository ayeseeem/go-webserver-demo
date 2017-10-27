package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const addr = ":8080"

// Page represents a page in a wiki
type Page struct {
	Title string
	Body  string
}

func (p Page) save() error {
	log.Println("Saving page", p.Title)
	filename := p.Title + ".txt"
	err := ioutil.WriteFile(filename, []byte(p.Body), os.ModePerm)
	return err
}

func loadPage(title string) (Page, error) {
	log.Println("Loading page", title)
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return Page{Title: title, Body: "Could not retrieve file - do not save this!!!"}, err
	}
	return Page{Title: title, Body: string(body)}, nil
}

func wikiViewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/wiki/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/wiki/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "wikiView", p)
}

func wikiEditHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/wiki/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = Page{Title: title}
	}
	renderTemplate(w, "wikiEdit", p)
}

func renderTemplate(w http.ResponseWriter, templateName string, p Page) {
	t, err := template.ParseFiles("./templates/" + templateName + ".html")
	if err != nil {
		log.Println("Problem with template", templateName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, p)
}
func main() {

	p1 := Page{Title: "TestPage", Body: "This is a sample Page."}
	p1.save()

	p2, _ := loadPage("TestPage")
	fmt.Println(p2.Title)
	fmt.Println(p2.Body)

	printStartUpMessage(addr)

	http.HandleFunc("/", top)
	http.HandleFunc("/simple", simple)
	http.HandleFunc("/wiki/view/", wikiViewHandler)
	http.HandleFunc("/wiki/edit/", wikiEditHandler)
	//	http.HandleFunc("/wiki/sace/", wikiSaveHandler)
	http.ListenAndServe(addr, nil)
}

func printStartUpMessage(addr string) {
	fmt.Println("Demo webserver")
	fmt.Printf("Visit http://localhost%v\n", addr)
	fmt.Println("Press ^C (Ctrl-C) to finish")
}

func top(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./home.html")
	if err != nil {
		log.Println("Problem with template", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer f.Close()
	io.Copy(w, f)
}

func simple(w http.ResponseWriter, r *http.Request) {
	renderUserTemplate(w, "./templates/simple-template.html", UserInfo{123, "mockers"})
}

// UserInfo is a bodged type for use investigating how templates work
type UserInfo struct {
	// Id is some sort of ID or reference number.
	// It will probably be replaced by a string or UUID of some sort
	ID       int
	Username string
}

func renderUserTemplate(w http.ResponseWriter, templateFilename string, user UserInfo) {
	t, err := template.ParseFiles(templateFilename)
	if err != nil {
		log.Println("Problem with template", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, user)
}

package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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

var validPath = regexp.MustCompile("^/wiki/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		m := validPath.FindStringSubmatch(path)
		if m == nil {
			log.Println("Could not extract title from path", path)
			http.NotFound(w, r)
			return
		}
		viewType := m[1]
		log.Println("View type:", viewType)
		pageTitle := m[2]
		log.Println("Page title:", pageTitle)
		fn(w, r, pageTitle)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	path := r.URL.Path
	m := validPath.FindStringSubmatch(path)
	if m == nil {
		log.Println("Could not extract title from path", path)
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	viewType := m[1]
	log.Println("View type:", viewType)
	pageTitle := m[2]
	log.Println("Page title:", pageTitle)
	return pageTitle, nil
}

func wikiViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/wiki/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "wikiView", p)
}

func wikiEditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = Page{Title: title}
	}
	renderTemplate(w, "wikiEdit", p)
}

func wikiSaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := Page{Title: title, Body: body}
	p.save()
	http.Redirect(w, r, "/wiki/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("./templates/wikiEdit.html", "./templates/wikiView.html"))

func renderTemplate(w http.ResponseWriter, templateName string, p Page) {
	err := templates.ExecuteTemplate(w, templateName+".html", p)
	if err != nil {
		log.Println("Problem with template", templateName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	sandbox()

	printStartUpMessage(addr)

	http.HandleFunc("/", top)
	http.HandleFunc("/simple", simple)
	http.HandleFunc("/wiki/view/", makeHandler(wikiViewHandler))
	http.HandleFunc("/wiki/edit/", makeHandler(wikiEditHandler))
	http.HandleFunc("/wiki/save/", makeHandler(wikiSaveHandler))
	http.ListenAndServe(addr, nil)
}

func sandbox() {
	sandbox, err := loadPage("SandBox")
	if err != nil {
		log.Println("Creating SandBox wiki page")
		sandbox = Page{Title: "SandBox", Body: "This is the SandBox Page. Play around. Do not rely on this page surviving forever"}
		sandbox.save()
	}
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

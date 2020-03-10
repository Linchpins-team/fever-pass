package main

import (
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"

	"net/http"
)

func (h Handler) doc(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	path := "doc/" + title + ".md"
	input, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	output := blackfriday.Run(input)

	page := struct {
		Title   string
		Content template.HTML
	}{
		title,
		template.HTML(string(output)),
	}
	h.HTML(w, r, "doc.htm", page)
}

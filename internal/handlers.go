package internal

import (
	"fmt"
	"net/http"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "HOME PAGE")

}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.createPostPOST(w, r)
		return
	} else if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	ts, err := template.ParseFiles("./web/html/postCreate.html")
	if err != nil {
		app.logger.ErrorLog.Println(err)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	ts.Execute(w, nil)

}

func (app *application) createPostPOST(w http.ResponseWriter, r *http.Request) {

}

package delivery

import (
	"fmt"
	"net/http"
	"strings"

	"forum/internal/service"
	"forum/internal/utils"
	tmpl2 "forum/pkg/tmpl"
)

func (a *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.log.Debug(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodGet {
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	data, err := service.GetHomePageData(a.storage, r.Context())
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/home.html", data)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (a *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:

		categories, err := service.GetCategories(a.storage)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/post_create.html", categories)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		return

	case http.MethodPost:

		post, err := service.DecodePost(r)
		if err != nil {
			a.log.Error(err.Error())
			data := struct {
				Error string
			}{
				Error: err.Error(),
			}
			tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/home.html", data)
			return
		}

		lastID, err := a.storage.CreatePost(post)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		fmt.Println(lastID)

		http.Redirect(w, r, "/createdposts", http.StatusSeeOther)

	default:
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return

	}
}

func (a *application) ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	id, err := utils.Etoi(idStr)
	if err != nil {
		a.log.Warn(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if id > service.LengthOfPosts {
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	postData, err := service.GetPostRelatedData(r.Context(), a.storage, id)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/post_view.html", postData)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (a *application) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	comment, err := service.DecodeComment(r)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, err.Error())
		return
	}

	err = a.storage.CreateComment(*comment)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

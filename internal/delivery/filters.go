package delivery

import (
	"fmt"
	"net/http"
	"strings"

	"forum/internal/service"
	"forum/internal/utils"
	tmpl2 "forum/pkg/tmpl"
)

func (a *application) SortedByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/category/")
	id, err := utils.Etoi(idStr)
	if err != nil {
		a.log.Warn(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	posts, err := service.GetAllPostRelatedDataByCategory(a.storage, id)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	a.log.Debug("Sorted by category id: %v", id)
	a.log.Debug("Posts: %v", posts)

	err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/sorted.html", posts)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (a *application) MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userID := r.Context().Value("userID").(int)

	posts, err := service.GetAllPostRelatedDataByUserID(a.storage, userID)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/sorted.html", posts)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (a *application) MyLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userID := r.Context().Value("userID").(int)

	posts, err := service.GetAllLikedPostsById(a.storage, userID)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/sorted.html", posts)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

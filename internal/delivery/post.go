package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/service"
	"forum/internal/utils"
	tmpl2 "forum/pkg/tmpl"
	"net/http"
	"strings"
)

func (a *application) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := strings.TrimPrefix(r.URL.Path, "/edit/")
		id, err := utils.Etoi(idStr)
		if err != nil {
			a.log.Warn(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		postData, err := service.GetPostRelatedData(r.Context(), a.storage, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				a.log.Error(fmt.Sprintf("Post with ID %d not found", id))
				tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		isLogin, ok := r.Context().Value("IsLogin").(bool)
		if !ok {
			isLogin = false
		}
		postData.IsLogin = isLogin

		currentUserID, ok := r.Context().Value("userID").(int)
		if !ok {
			a.log.Warn("User ID not found in context")
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusUnauthorized, "User not authorized")
			return
		}

		if postData.Post.AuthorID != currentUserID {
			a.log.Warn("User is not the author of the post")
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusForbidden, "You are not the author of this post")
			return
		}

		err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/post_edit.html", postData)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

	case http.MethodPost:
		idStr := strings.TrimPrefix(r.URL.Path, "/edit/")
		id, err := utils.Etoi(idStr)
		if err != nil {
			a.log.Warn(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		post, err := service.DecodePost(r, a.storage)
		if err != nil {
			a.log.Error(err.Error())
			return
		}

		postData, err := service.GetPostRelatedData(r.Context(), a.storage, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				a.log.Error(fmt.Sprintf("Post with ID %d not found", id))
				tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		currentUserID, ok := r.Context().Value("userID").(int)
		if !ok {
			a.log.Warn("User ID not found in context")
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusUnauthorized, "User not authorized")
			return
		}

		if postData.Post.AuthorID != currentUserID {
			a.log.Warn("User is not the author of the post")
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusForbidden, "You are not the author of this post")
			return
		}

		err = a.storage.UpdatePost(id, post)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (a *application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting post...")
}

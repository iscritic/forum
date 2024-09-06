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

func (a *application) EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := strings.TrimPrefix(r.URL.Path, "/comment/edit/")
		id, err := utils.Etoi(idStr)
		if err != nil {
			a.log.Warn(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		commentData, err := a.storage.GetCommentByID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				a.log.Error(fmt.Sprintf("Comment with ID %d not found", id))
				tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		authorID, ok := r.Context().Value("userID").(int)
		if !ok || authorID != commentData.AuthorID {
			a.log.Warn("User not authorized to edit this comment")
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusForbidden, "You are not authorized to edit this comment")
			return
		}

		err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/comment_edit.html", commentData)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

	case http.MethodPost:
		idStr := strings.TrimPrefix(r.URL.Path, "/comment/edit/")
		id, err := utils.Etoi(idStr)
		if err != nil {
			a.log.Warn(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		updatedComment, err := service.DecodeComment(r, a.storage)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, err.Error())
			return
		}

		authorID, ok := r.Context().Value("userID").(int)
		if !ok || authorID != updatedComment.AuthorID {
			a.log.Warn("User not authorized to edit this comment")
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusForbidden, "You are not authorized to edit this comment")
			return
		}

		updatedComment.ID = id

		err = a.storage.UpdateComment(updatedComment)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, err.Error())
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%d", updatedComment.PostID), http.StatusSeeOther)

	default:
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (a *application) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	a.log.Warn("WE ARE HERE")
	a.log.Warn("WE ARE HERE")
	a.log.Warn("WE ARE HERE")

	idStr := strings.TrimPrefix(r.URL.Path, "/comment/delete/")
	id, err := utils.Etoi(idStr)
	if err != nil {
		a.log.Warn("Invalid comment ID: " + err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	commentData, err := a.storage.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.log.Error(fmt.Sprintf("Comment with ID %d not found", id))
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		a.log.Error("Failed to get comment: " + err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	authorID, ok := r.Context().Value("userID").(int)
	if !ok || authorID != commentData.AuthorID {
		a.log.Warn("User not authorized to delete this comment")
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusForbidden, "You are not authorized to delete this comment")
		return
	}

	err = a.storage.DeleteComment(id)
	if err != nil {
		a.log.Error("Failed to delete comment: " + err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", commentData.PostID), http.StatusSeeOther)
}

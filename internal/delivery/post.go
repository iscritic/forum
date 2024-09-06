package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/entity"
	"forum/internal/service"
	"forum/internal/utils"
	tmpl2 "forum/pkg/tmpl"
	"net/http"
	"strings"
)

func (a *application) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		// Получаем ID поста из URL
		idStr := strings.TrimPrefix(r.URL.Path, "/post/edit/")
		id, err := utils.Etoi(idStr)
		if err != nil {
			a.log.Warn(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		// Получаем данные поста для редактирования
		postData, err := service.GetPostRelatedData(r.Context(), a.storage, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				a.log.Error(fmt.Sprintf("Post with ID %d not found", id))
				tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		// Получаем все категории
		categories, err := service.GetCategories(a.storage)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		// Создаем структуру для передачи в шаблон
		pageData := struct {
			Post       *entity.PostRelatedData
			Categories []entity.Category
		}{
			Post:       postData,
			Categories: categories,
		}

		// Отображаем форму редактирования поста
		err = tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/post_edit.html", pageData)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

	case http.MethodPost:
		// Получаем ID поста из URL
		idStr := strings.TrimPrefix(r.URL.Path, "/post/edit/")
		id, err := utils.Etoi(idStr)
		if err != nil {
			a.log.Warn(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		// Декодируем данные поста из запроса
		post, err := service.DecodePost(r, a.storage)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusBadRequest, err.Error())
			return
		}

		// Устанавливаем ID поста для обновления
		post.ID = id

		// Обновляем пост в базе данных
		err = a.storage.UpdatePost(post)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		// Перенаправляем пользователя на страницу просмотра поста
		http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusSeeOther)

	default:
		a.log.Debug(fmt.Sprintf("Method Not Allowed %s %s", r.Method, r.URL.Path))
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

package delivery

import (
	"forum/internal/helpers/template"
	"forum/internal/service"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) SortedByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	posts, err := service.GetAllPostRelatedDataByCategory(app.storage, id)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.logger.InfoLog.Println("Sorted by category id:", id)
	app.logger.InfoLog.Println("Posts:", posts)

	err = template.RenderTemplate(w, app.templateCache, "./web/html/sorted.html", posts)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(int)

	posts, err := service.GetAllPostRelatedDataByUserID(app.storage, userID)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		return
	}

	template.RenderTemplate(w, app.templateCache, "./web/html/sorted.html", posts)
}

func (app *application) MyLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(int)

	posts, err := service.GetAllLikedPostsById(app.storage, userID)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		return
	}

	template.RenderTemplate(w, app.templateCache, "./web/html/sorted.html", posts)
}

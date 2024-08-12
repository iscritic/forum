package delivery

import (
	"net/http"
	"strconv"
)

// //// likes and dislikes
func (app *application) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	hasLiked, err := app.storage.HasLikedPost(userID, postID)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasLiked {
		err = app.storage.RemoveLike(userID, postID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasDisliked, err := app.storage.HasDislikedPost(userID, postID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasDisliked {
			err = app.storage.RemoveDislike(userID, postID)
			if err != nil {
				app.logger.ErrorLog.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = app.storage.LikePost(userID, postID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

func (app *application) DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	hasDisliked, err := app.storage.HasDislikedPost(userID, postID)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasDisliked {
		err = app.storage.RemoveDislike(userID, postID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasLiked, err := app.storage.HasLikedPost(userID, postID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasLiked {
			err = app.storage.RemoveLike(userID, postID)
			if err != nil {
				app.logger.ErrorLog.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = app.storage.DislikePost(userID, postID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

func (app *application) LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	hasLiked, err := app.storage.HasLikedComment(userID, commentID)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasLiked {
		err = app.storage.UnlikeComment(userID, commentID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasDisliked, err := app.storage.HasDislikedComment(userID, commentID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasDisliked {
			err = app.storage.UndislikeComment(userID, commentID)
			if err != nil {
				app.logger.ErrorLog.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = app.storage.LikeComment(userID, commentID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (app *application) DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	hasDisliked, err := app.storage.HasDislikedComment(userID, commentID)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasDisliked {
		err = app.storage.UndislikeComment(userID, commentID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasLiked, err := app.storage.HasLikedComment(userID, commentID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasLiked {
			err = app.storage.UnlikeComment(userID, commentID)
			if err != nil {
				app.logger.ErrorLog.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = app.storage.DislikeComment(userID, commentID)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

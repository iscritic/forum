package delivery

import (
	"net/http"
	"strconv"
)

// //// likes and dislikes
func (a *application) LikePostHandler(w http.ResponseWriter, r *http.Request) {
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

	hasLiked, err := a.storage.HasLikedPost(userID, postID)
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasLiked {
		err = a.storage.RemoveLike(userID, postID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasDisliked, err := a.storage.HasDislikedPost(userID, postID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasDisliked {
			err = a.storage.RemoveDislike(userID, postID)
			if err != nil {
				a.log.Error(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = a.storage.LikePost(userID, postID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	referer := r.Referer()

	if referer == "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, referer, http.StatusSeeOther)
	}

}

func (a *application) DislikePostHandler(w http.ResponseWriter, r *http.Request) {
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

	hasDisliked, err := a.storage.HasDislikedPost(userID, postID)
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasDisliked {
		err = a.storage.RemoveDislike(userID, postID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasLiked, err := a.storage.HasLikedPost(userID, postID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasLiked {
			err = a.storage.RemoveLike(userID, postID)
			if err != nil {
				a.log.Error(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = a.storage.DislikePost(userID, postID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	referer := r.Referer()

	if referer == "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, referer, http.StatusSeeOther)
	}
}

func (a *application) LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	hasLiked, err := a.storage.HasLikedComment(userID, commentID)
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasLiked {
		err = a.storage.UnlikeComment(userID, commentID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasDisliked, err := a.storage.HasDislikedComment(userID, commentID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasDisliked {
			err = a.storage.UndislikeComment(userID, commentID)
			if err != nil {
				a.log.Error(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = a.storage.LikeComment(userID, commentID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (a *application) DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	hasDisliked, err := a.storage.HasDislikedComment(userID, commentID)
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasDisliked {
		err = a.storage.UndislikeComment(userID, commentID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		hasLiked, err := a.storage.HasLikedComment(userID, commentID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if hasLiked {
			err = a.storage.UnlikeComment(userID, commentID)
			if err != nil {
				a.log.Error(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		err = a.storage.DislikeComment(userID, commentID)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

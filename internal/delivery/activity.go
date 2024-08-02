package delivery

import (
	"fmt"
	"net/http"
)

func (app *application) MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "My created posts page...")
}

func (app *application) LikedPostsHanlers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I liked those posts...")
}

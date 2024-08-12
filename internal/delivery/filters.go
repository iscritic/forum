package delivery

import (
	"fmt"
	"net/http"
)

func (app *application) SortedByCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CategoryID page")
}

func (app *application) CreatedPostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "My created posts page...")
}

func (app *application) LikedPostsHanlers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I liked those posts...")
}

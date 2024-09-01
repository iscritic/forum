package delivery

import (
	"fmt"
	"net/http"
)

func (a *application) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Editing post...")
}

func (a *application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting post...")
}

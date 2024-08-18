package delivery

import (
	"fmt"
	"net/http"
)

func (a *application) adminPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the admin page!")

}

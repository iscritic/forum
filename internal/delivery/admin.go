package delivery

import (
	"fmt"
	"net/http"
)

func (app *application) adminPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the admin page!")

}

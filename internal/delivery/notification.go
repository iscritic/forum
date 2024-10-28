package delivery

import (
	tmpl2 "forum/pkg/tmpl"
	"net/http"
)

func (a *application) ViewNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	notifications, err := a.storage.GetUnreadNotifications(userID)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, err.Error())
		return
	}

	err = tmpl2.RenderTemplate(w, a.tmplcache, "notifications.html", notifications)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, err.Error())
		return
	}
}

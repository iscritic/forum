package internal

import (
	"forum/pkg/logger"
	"net/http"
)

type application struct {
	logger *logger.Logger
}

func Routes(l *logger.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	app := &application{
		logger: l,
	}

	mux.HandleFunc("/", app.home)
	return mux
}

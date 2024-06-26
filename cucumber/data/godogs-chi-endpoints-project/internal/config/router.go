// Code generated by kk; DO NOT EDIT.

package config

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/waler4ik/godogs-rest-project/internal/api"
	"github.com/waler4ik/godogs-rest-project/internal/endpoints/machines/data"
)

func ConfigureRouter(a *api.API) http.Handler {
	r := chi.NewRouter()
	ConfigureMiddleware(r)
	data.ConfigureRouter(a, r)
	r.HandleFunc("/ws", a.Ws.Handler)
	return r
}

// Code generated by kk; DO NOT EDIT.

package config

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/waler4ik/godogs-rest-project/internal/api"
)

func ConfigureRouter(a *api.API) http.Handler {
	r := chi.NewRouter()
	ConfigureMiddleware(r)
	return r
}

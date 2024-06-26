// Code generated by kk; DO NOT EDIT.

package data

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"

	"github.com/waler4ik/godogs-rest-project/internal/api"
	"github.com/waler4ik/godogs-rest-project/internal/errors"
)

type Controller struct {
	p  Provider
	sd *schema.Decoder
}

func NewController(p Provider) Controller {
	return Controller{
		p:  p,
		sd: schema.NewDecoder(),
	}
}

const RoutePath = "/machines/data"

func ConfigureRouter(a *api.API, r *chi.Mux) {
	dataProvider := NewProvider(a)
	dataController := NewController(dataProvider)
	r.Route(RoutePath, func(r chi.Router) {
		r.With(Paginate).Get("/", dataController.ListData)
		r.Post("/", dataController.CreateDatum)

		r.Route("/{datumID}", func(r chi.Router) {
			r.Use(dataController.DatumCtx)
			r.Get("/", dataController.GetDatum)
			r.Put("/", dataController.UpdateDatum)
			r.Delete("/", dataController.DeleteDatum)
		})
	})
}

func (c Controller) ListData(w http.ResponseWriter, r *http.Request) {
	var qp QueryParameter
	if err := c.sd.Decode(&qp, r.URL.Query()); err != nil {
		render.Render(w, r, errors.ErrBadRequest(err))
		return
	}

	if data, err := c.p.GetData(r.Context(), qp); err != nil {
		render.Render(w, r, errors.ErrInternalServerError(err))
		return
	} else {
		if err := render.RenderList(w, r, data); err != nil {
			render.Render(w, r, errors.ErrRender(err))
			return
		}
	}
}

func (c Controller) CreateDatum(w http.ResponseWriter, r *http.Request) {
	datum := &Datum{}
	if err := render.Bind(r, datum); err != nil {
		render.Render(w, r, errors.ErrBadRequest(err))
		return
	}

	if datum, err := c.p.CreateDatum(r.Context(), datum); err != nil {
		render.Render(w, r, errors.ErrInternalServerError(err))
		return
	} else {
		render.Status(r, http.StatusCreated)
		render.Render(w, r, datum)
	}
}

func (c Controller) GetDatum(w http.ResponseWriter, r *http.Request) {
	datum := r.Context().Value(DatumCtxKey).(*Datum)
	if err := render.Render(w, r, datum); err != nil {
		render.Render(w, r, errors.ErrRender(err))
		return
	}
}

func (c Controller) UpdateDatum(w http.ResponseWriter, r *http.Request) {
	datum := r.Context().Value(DatumCtxKey).(*Datum)

	if err := render.Bind(r, datum); err != nil {
		render.Render(w, r, errors.ErrBadRequest(err))
		return
	}

	if err := c.p.UpdateDatum(r.Context(), datum); err != nil {
		render.Render(w, r, errors.ErrInternalServerError(err))
		return
	}
}

func (c Controller) DeleteDatum(w http.ResponseWriter, r *http.Request) {
	datum := r.Context().Value(DatumCtxKey).(*Datum)
	if err := c.p.DeleteDatum(r.Context(), datum); err != nil {
		render.Render(w, r, errors.ErrInternalServerError(err))
		return
	}
}

type ContextKeyType string

const DatumCtxKey = ContextKeyType("DatumCtxKey")

func (c Controller) DatumCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var datum *Datum
		var err error

		if datumID := chi.URLParam(r, "datumID"); datumID != "" {
			datum, err = c.p.GetDatum(r.Context(), datumID)
			if err != nil {
				render.Render(w, r, errors.ErrNotFound)
				return
			}
		} else {
			render.Render(w, r, errors.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), DatumCtxKey, datum)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

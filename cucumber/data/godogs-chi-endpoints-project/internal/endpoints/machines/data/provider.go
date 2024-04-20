// Code generated by kk; BUT FEEL FREE TO EDIT.

package data

import (
	"context"
	"fmt"

	"github.com/go-chi/render"
	"github.com/waler4ik/godogs-rest-project/internal/api"
)

type Provider struct {
}

func NewProvider(a *api.API) Provider {
	return Provider{}
}

// QueryParameter implements github.com/gorilla/schema
type QueryParameter struct {
	Name  string `schema:"name,required"` // custom name, must be supplied
	Phone string `schema:"phone"`         // custom name
	Admin bool   `schema:"-"`             // this field is never set
}

// GetData godoc
//
//	@Summary		List data
//	@Description	get data
//	@Tags			data
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	true	"name attribute in query"
//	@Param			phone	query		string	false	"phone attribute in query"
//	@Success		200	{array}		data.Datum
//	@Failure		400	{object}	errors.ErrResponse
//	@Failure		404	{object}	errors.ErrResponse
//	@Failure		500	{object}	errors.ErrResponse
//	@Router			/machines/data [get]
func (p Provider) GetData(ctx context.Context, qp QueryParameter) ([]render.Renderer, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetDatum godoc
//
//	@Summary		Show datum
//	@Description	get Datum by ID
//	@Tags			data
//	@Accept			json
//	@Produce		json
//	@Param			datumID	path		string	true	"Datum ID"
//	@Success		200	{object}	data.Datum
//	@Failure		400	{object}	errors.ErrResponse
//	@Failure		404	{object}	errors.ErrResponse
//	@Failure		500	{object}	errors.ErrResponse
//	@Router			/machines/data/{datumID} [get]
func (p Provider) GetDatum(ctx context.Context, id string) (*Datum, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateDatum godoc
//
//	@Summary		Add datum
//	@Description	add by json datum
//	@Tags			data
//	@Accept			json
//	@Produce		json
//	@Param			datum	body		data.Datum	true	"Add datum"
//	@Success		201		{object}	data.Datum
//	@Failure		400		{object}	errors.ErrResponse
//	@Failure		404		{object}	errors.ErrResponse
//	@Failure		500		{object}	errors.ErrResponse
//	@Router			/machines/data [post]
func (p Provider) CreateDatum(ctx context.Context, c *Datum) (*Datum, error) {
	return nil, fmt.Errorf("not implemented")
}

// UpdateDatum godoc
//
//	@Summary		Update datum
//	@Description	Update by json datum
//	@Tags			data
//	@Accept			json
//	@Produce		json
//	@Param			datumID	path		string	true	"Datum ID"
//	@Param			datum	body		data.Datum	true	"Update datum"
//	@Success		200
//	@Failure		400		{object}	errors.ErrResponse
//	@Failure		404		{object}	errors.ErrResponse
//	@Failure		500		{object}	errors.ErrResponse
//	@Router			/machines/data/{datumID} [put]
func (p Provider) UpdateDatum(ctx context.Context, c *Datum) error {
	return fmt.Errorf("not implemented")
}

// DeleteDatum godoc
//
//	@Summary		Delete datum
//	@Description	Delete by datum ID
//	@Tags			data
//	@Accept			json
//	@Produce		json
//	@Param			datumID	path		string	true	"Datum ID"
//	@Success		200
//	@Failure		400	{object}	errors.ErrResponse
//	@Failure		404	{object}	errors.ErrResponse
//	@Failure		500	{object}	errors.ErrResponse
//	@Router			/machines/data/{datumID} [delete]
func (p Provider) DeleteDatum(ctx context.Context, c *Datum) error {
	return fmt.Errorf("not implemented")
}

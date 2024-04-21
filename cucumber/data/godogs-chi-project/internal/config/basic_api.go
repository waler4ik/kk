// Code generated by kk; BUT FEEL FREE TO EDIT.

package config

import (
	"log"

	"github.com/waler4ik/godogs-rest-project/internal/api"
)

// configureBasicAPI is an entry point for the basic api configuration part.
func configureBasicAPI(a *api.Basic) {
	a.Logger = log.Printf
	a.PreServerShutdown = func() {}
	a.ServerShutdown = func() {}
}
package httpRouter

import (
	"github.com/edgexfoundry/powerflow/traderPowerflow/routing/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Start",
		"GET",
		"/start",
		handlers.Start,
	},
}

package routing

import (
	"github.com/edgexfoundry/powerflow/registerPowerflow/routing/handlers"
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
	Route{
		"Peers",
		"GET",
		"/peers",
		handlers.Peers,
	},
	Route{
		"RegisterSelfTo",
		"POST",
		"/registerSelfTo",
		handlers.RegisterSelfTo,
	},
	Route{
		"Register",
		"POST",
		"/register",
		handlers.Register,
	},
	Route{
		"RegisterForBlockchain",
		"POST",
		"/register/blockchain",
		handlers.RegisterForBlockchain,
	},
}

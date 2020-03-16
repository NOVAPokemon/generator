package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

const StatusName = "STATUS"
const RegisterName = "REGISTER"
const LoginName = "LOGIN"
const RefreshName = "REFRESH"

const GET = "GET"
const POST = "POST"

var routes = Routes{
	Route{
		StatusName,
		GET,
		"/",
		Status,
	},
	Route{
		RegisterName,
		POST,
		"/register",
		Register,
	},
	Route{
		LoginName,
		POST,
		"/login",
		Login,
	},
	Route{
		RefreshName,
		POST,
		"/refresh",
		Refresh,
	},
}

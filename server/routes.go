package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type requestContext struct {
	w http.ResponseWriter
	r *http.Request
	p httprouter.Params
}

func (c *requestContext) InternalServerError(err error) {
	// TODO: report it somewhere
	log.Println(err.Error())
	c.WriteJSON(errorResponse{"error", "internal_server_error"})
}

func (c *requestContext) WriteJSON(thing interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(c.w).Encode(thing)
	if err != nil {
		panic(err)
	}
}

type handleFunc func(c *requestContext)

func route(h handleFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := requestContext{
			w: w,
			r: r,
			p: p,
		}
		h(&c)
	}
}

func initRoutes() {
	router := httprouter.New()

	router.GET("/", route(routeMain))

	router.GET("/artifacts/get", route(routeMain))

	router.GET("/jobs/get", route(routeJobsGet))
	router.POST("/jobs/start", route(routeMain))

	http.Handle("/", router)
}

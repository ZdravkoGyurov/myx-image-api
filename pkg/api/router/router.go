package router

import (
	"net/http"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/api/handlers"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/api/middlewares"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/controller"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/log"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
	controller controller.Controller
}

func New(ctrl controller.Controller) Router {
	r := mux.NewRouter()
	router := Router{
		Router:     r,
		controller: ctrl,
	}

	router.Use(middlewares.PanicRecovery)
	router.Use(middlewares.LoggerMiddleware)
	router.Use(middlewares.CorrelationIDMiddleware)
	router.mountImageRoutes()

	logRoutes(r)
	return router
}

func logRoutes(router *mux.Router) {
	logger := log.DefaultLogger()
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		if path != "" || len(methods) > 0 {
			logger.Info().Msgf("%v %s", methods, path)
		}
		return nil
	})
	if err != nil {
		logger.Err(err).Msg("could not list routes")
	}
}

func (r Router) mountImageRoutes() {
	imageHandler := handlers.Image{Controller: r.controller}
	r.Methods(http.MethodPost).Path(imagePath).HandlerFunc(imageHandler.Upload)
	r.Methods(http.MethodGet).Path(imagePath).HandlerFunc(imageHandler.GetAll)
	r.Methods(http.MethodDelete).Path(imagePathWithName).HandlerFunc(imageHandler.Delete)
}

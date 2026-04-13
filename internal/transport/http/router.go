package transporthttp

import (
	"net/http"
	"subscriptions/internal/logging"

	"github.com/gorilla/mux"

	swaggerdocs "subscriptions/internal/transport/http/docs"
	"subscriptions/internal/transport/http/handlers"
)

func NewRouter(subHandler *handlers.SubscriptionHandler, docsHandler *swaggerdocs.Handler) *mux.Router {
	logging.Logger.Info("initializing router")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
	router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
	router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/subscriptions", subHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/subscriptions", subHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/subscriptions/{id:[0-9]+}", subHandler.GetByID).Methods(http.MethodGet)
	api.HandleFunc("/subscriptions/{id:[0-9]+}", subHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/subscriptions/{id:[0-9]+}", subHandler.Delete).Methods(http.MethodDelete)
	api.HandleFunc("/subscriptions/total_price", subHandler.TotalPrice).Methods(http.MethodPost)

	logging.Logger.Info("router setup complete")

	return router
}

package application

import (
	"net/http"

	"github.com/RohanKhatua/go-chi/handler"
	"github.com/RohanKhatua/go-chi/repository/order"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (a *App) loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteByID)
}

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", a.loadOrderRoutes)

	a.router = router
}

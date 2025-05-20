package api

import (
	"github.com/dsrosen6/cw-ticket-bot/internal/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func pingRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", handleGetPing)
	return r
}

func handleGetPing(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, util.ResultBody("ping successful"))
}

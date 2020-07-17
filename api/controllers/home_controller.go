package controllers

import (
	"net/http"
	"os"

	"github.com/norfabagas/auth/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), "Index")
}

func (server *Server) ApiSecret(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), os.Getenv("API_SECRET"))
}

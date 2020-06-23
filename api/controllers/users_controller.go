package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/norfabagas/auth/api/auth"
	"github.com/norfabagas/auth/api/models"
	"github.com/norfabagas/auth/api/responses"
	"github.com/norfabagas/auth/api/utils/crypto"
	"github.com/norfabagas/auth/api/utils/formaterror"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("register")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, true, http.StatusText(http.StatusCreated), struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		PublicID  string    `json:"public_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		Name:      userCreated.Name,
		Email:     userCreated.Email,
		PublicID:  userCreated.PublicID,
		CreatedAt: userCreated.CreatedAt,
	})
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	user := models.User{}
	userData, err := user.FindUserByID(server.DB, uint32(tokenID))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userData.Name, err = crypto.Decrypt(userData.Name, os.Getenv("APP_KEY"))
	if err != nil {
		fmt.Println(userData.Name)
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"last_updated"`
	}{
		Name:      userData.Name,
		Email:     userData.Email,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	})
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateUser(server.DB, uint32(tokenID))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"last_updated"`
	}{
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	})
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(tokenID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	_, err = user.DeleteUser(server.DB, uint32(tokenID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, true, "Deleted", "")
}

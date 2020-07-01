package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/norfabagas/auth/api/auth"
	"github.com/norfabagas/auth/api/models"
	"github.com/norfabagas/auth/api/responses"
	"github.com/norfabagas/auth/api/utils/crypto"
	"github.com/norfabagas/auth/api/utils/formaterror"
)

func (server *Server) SignIn(email, password string) (string, error) {
	var err error

	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.PublicID)
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
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

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", user.Email).Take(&user).Error
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	user.Name, err = crypto.Decrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), struct {
		Token string `json:"token"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		Token: token,
		Email: user.Email,
		Name:  user.Name,
	})
}

package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MuriloAbranches/Go-Expert-API/internal/dto"
	"github.com/MuriloAbranches/Go-Expert-API/internal/entity"
	"github.com/MuriloAbranches/Go-Expert-API/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(userDB database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDB:       userDB,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var userDto dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByEmail(userDto.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !user.ValidatePassword(userDto.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := h.Jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDto dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := entity.NewUser(userDto.Name, userDto.Email, userDto.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

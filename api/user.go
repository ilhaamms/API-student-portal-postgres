package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	var creds model.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	if creds.Username == "" || creds.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "username atau password wajib diisi"})
		return
	}

	if api.userService.CheckPassLength(creds.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "password wajib 5 atau lebih dari 5 characters"})
		return
	}

	err = api.userService.Register(creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: fmt.Sprintf("Internal Server Error: %v", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: creds.Username,
		Message:  "User Registered",
	})
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var creds model.User

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	if creds.Username == "" || creds.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "username atau password wajib diisi"})
		return
	}

	if api.userService.CheckPassLength(creds.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "password wajib 5 atau lebih dari 5 characters"})
		return
	}

	err = api.userService.Login(creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username atau password salah"})
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(5 * time.Hour)
	session := model.Session{Token: sessionToken, Username: creds.Username, Expiry: expiresAt}

	err = api.sessionService.SessionAvailName(session.Username)
	if err != nil {
		err = api.sessionService.AddSession(session)
	} else {
		err = api.sessionService.UpdateSession(session)
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: creds.Username,
		Message:  "Login Success",
	})
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value("username").(string)
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	sessionToken := c.Value

	api.sessionService.DeleteSession(sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: username,
		Message:  "Logout Success",
	})
}

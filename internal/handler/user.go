package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/service"
	response "github.com/sprint-id/eniqilo-server/pkg/resp"
)

type userHandler struct {
	userSvc *service.UserService
}

func newUserHandler(userSvc *service.UserService) *userHandler {
	return &userHandler{userSvc}
}

func (h *userHandler) RegisterIT(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqRegister

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.RegisterIT(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "User registered successfully"
	successRes.Data = res

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) RegisterNurse(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqRegisterNurse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.RegisterNurse(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "User registered successfully"
	successRes.Data = res

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) LoginIT(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqLogin

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.LoginIT(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "User logged successfully"
	successRes.Data = res

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) LoginNurse(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqLogin

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.LoginNurse(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "User logged successfully"
	successRes.Data = res

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetUser

	param.UserID = queryParams.Get("userId")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.Name = queryParams.Get("name")
	param.NIP = queryParams.Get("nip")
	param.Role = queryParams.Get("role")
	param.CreatedAt = queryParams.Get("createdAt")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.GetUser(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "User data retrieved successfully"
	successRes.Data = res

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *userHandler) UpdateNurse(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[4]
	fmt.Printf("id: %s\n", id)
	var req dto.ReqUpdateNurse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	err = h.userSvc.UpdateNurse(r.Context(), req, id)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) DeleteNurse(w http.ResponseWriter, r *http.Request) {
	// id := r.PathValue("userId")
	id := strings.Split(r.URL.Path, "/")[4]
	fmt.Printf("id: %s\n", id)

	err := h.userSvc.DeleteNurse(r.Context(), id)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) AccessNurse(w http.ResponseWriter, r *http.Request) {
	// id := r.PathValue("userId")
	id := strings.Split(r.URL.Path, "/")[4]
	fmt.Printf("id: %s\n", id)
	var req dto.ReqAccessNurse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	err = h.userSvc.AccessNurse(r.Context(), req, id)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

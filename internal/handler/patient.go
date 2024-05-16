package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/service"
	response "github.com/sprint-id/eniqilo-server/pkg/resp"
)

type patientHandler struct {
	patientSvc *service.PatientService
}

func newPatientHandler(patientSvc *service.PatientService) *patientHandler {
	return &patientHandler{patientSvc}
}

// IdentityNumber      string `json:"identityNumber" validate:"required,len=16"`
// PhoneNumber         string `json:"phoneNumber" validate:"required,min=10,max=15"`
// Name                string `json:"name" validate:"required,min=3,max=30"`
// BirthDate           string `json:"birthDate" validate:"required"`
// Gender              string `json:"gender" validate:"required,oneof=male female"`
// IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,url"`

func (h *patientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreatePatient
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"identityNumber", "phoneNumber", "name", "birthDate", "gender", "identityCardScanImg"}
	for key := range jsonData {
		if !contains(expectedFields, key) {
			http.Error(w, "unexpected field in request body: "+key, http.StatusBadRequest)
			return
		}
	}

	// Convert the jsonData map into the req struct
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// show request
	fmt.Printf("RegisterCustomer request: %+v\n", req)

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.patientSvc.CreatePatient(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	// should return 201 if success
	w.WriteHeader(http.StatusCreated)
}

func (h *patientHandler) GetPatient(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetPatient

	param.IdentityNumber = queryParams.Get("identityNumber")
	param.PhoneNumber = queryParams.Get("phoneNumber")
	param.Name = queryParams.Get("name")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.CreatedAt = queryParams.Get("createdAt")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	customers, err := h.patientSvc.GetPatient(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	// show response
	// fmt.Printf("GetMatch response: %+v\n", customers)

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = customers

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

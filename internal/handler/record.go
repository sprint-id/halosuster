package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/service"
	response "github.com/sprint-id/eniqilo-server/pkg/resp"
)

type recordHandler struct {
	recordSvc *service.RecordService
}

func newRecordHandler(recordSvc *service.RecordService) *recordHandler {
	return &recordHandler{recordSvc}
}

// {
// 	"identityNumber": 123123, // not null, should be 16 digit
// 	"symptoms": "", // not null, minLength 1, maxLength 2000,
// 	"medications" : "" // not null, minLength 1, maxLength 2000
// }

func (h *recordHandler) AddRecord(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqAddRecord
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if the payload is empty
	if len(jsonData) == 0 {
		http.Error(w, "empty payload", http.StatusUnauthorized)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"identityNumber", "symptoms", "medications"}
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

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.recordSvc.AddRecord(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusCreated) // Set HTTP status code to 201
}

func (h *recordHandler) GetRecord(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetRecord

	param.ID = queryParams.Get("id")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.IdentityNumber = queryParams.Get("identityDetail.identityNumber")
	param.UserId = queryParams.Get("createdBy.userId")
	param.NIP = queryParams.Get("createdBy.nip")
	param.CreatedAt = queryParams.Get("createdAt")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	records, err := h.recordSvc.GetRecord(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = records

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

// The contains function checks if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

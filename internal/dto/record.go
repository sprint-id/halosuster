package dto

import (
	"github.com/sprint-id/eniqilo-server/internal/entity"
)

// {
// 	"identityNumber": 123123, // not null, should be 16 digit
// 	"symptoms": "", // not null, minLength 1, maxLength 2000,
// 	"medications" : "" // not null, minLength 1, maxLength 2000
// }

type (
	ReqAddRecord struct {
		IdentityNumber string `json:"identityNumber" validate:"required,len=16"`
		Symptoms       string `json:"symptoms" validate:"required,min=1,max=2000"`
		Medications    string `json:"medications" validate:"required,min=1,max=2000"`
	}

	ParamGetRecord struct {
		ID             string `json:"id"`
		Limit          int    `json:"limit"`
		Offset         int    `json:"offset"`
		IdentityNumber string `json:"identityDetail.identityNumber"`
		UserId         string `json:"createdBy.userId"`
		NIP            string `json:"createdBy.nip"`
		CreatedAt      string `json:"createdAt"`
	}

	ResGetRecord struct {
		ID             string            `json:"id"`
		IdentityDetail ResIdentityDetail `json:"identityDetail"`
		Symptoms       string            `json:"symptoms"`
		Medications    string            `json:"medications"`
		CreatedAt      string            `json:"createdAt"`
		CreatedBy      ResCreatedBy      `json:"createdBy"`
	}

	ResIdentityDetail struct {
		IdentityNumber      string `json:"identityNumber"`
		PhoneNumber         string `json:"phoneNumber"`
		Name                string `json:"name"`
		BirthDate           string `json:"birthDate"`
		Gender              string `json:"gender"`
		IdentityCardScanImg string `json:"identityCardScanImg"`
	}

	ResCreatedBy struct {
		NIP    string `json:"nip"`
		Name   string `json:"name"`
		UserID string `json:"userId"`
	}
)

func (d *ReqAddRecord) ToRecordEntity(userId string) entity.Record {
	return entity.Record{
		IdentityNumber: d.IdentityNumber,
		Symptoms:       d.Symptoms,
		Medications:    d.Medications,
		UserID:         userId,
	}
}

package dto

import "github.com/sprint-id/eniqilo-server/internal/entity"

// {
// 	"identityNumber": 123123, // not null, should be 16 digit
// 	"phoneNumber": "+62123123123", /** not null,
// 	 - should be string,
// 	 - starts with `+62`,
// 	 - minLength 10,
// 	 - maxLength 15
// 	*/
// 	"name": "", // not null, should be string, minLength 3, maxLength 30
// 	"birthDate" : "", // not null, should be string with ISO 8601 format
// 	"gender" : "", // not null, should be enum of 'male'|'female',
// 	"identityCardScanImg": "" // not null, should be an image url
// }

type (
	ReqCreatePatient struct {
		IdentityNumber      string `json:"identityNumber" validate:"required,len=16"`
		PhoneNumber         string `json:"phoneNumber" validate:"required,min=10,max=15"`
		Name                string `json:"name" validate:"required,min=3,max=30"`
		BirthDate           string `json:"birthDate" validate:"required"`
		Gender              string `json:"gender" validate:"required,oneof=male female"`
		IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,url"`
	}

	ParamGetPatient struct {
		IdentityNumber string `json:"identityNumber"`
		Limit          int    `json:"limit"`
		Offset         int    `json:"offset"`
		Name           string `json:"name"`
		PhoneNumber    string `json:"phoneNumber"`
		CreatedAt      string `json:"createdAt"`
	}

	ResGetPatient struct {
		IdentityNumber string `json:"identityNumber"`
		PhoneNumber    string `json:"phoneNumber"`
		Name           string `json:"name"`
		BirthDate      string `json:"birthDate"`
		Gender         string `json:"gender"`
		CreatedAt      string `json:"createdAt"`
	}
)

// ToEntity to convert dto to entity
func (d *ReqCreatePatient) ToPatientEntity() entity.Patient {
	return entity.Patient{
		IdentityNumber:      d.IdentityNumber,
		PhoneNumber:         d.PhoneNumber,
		Name:                d.Name,
		BirthDate:           d.BirthDate,
		Gender:              d.Gender,
		IdentityCardScanImg: d.IdentityCardScanImg,
	}
}

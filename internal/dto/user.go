package dto

import (
	"strconv"

	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/pkg/auth"
)

type (
	ReqRegister struct {
		NIP      int    `json:"nip" validate:"required"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
		Password string `json:"password" validate:"required,min=5,max=33"`
	}

	ReqRegisterNurse struct {
		NIP                 int    `json:"nip" validate:"required"`
		Name                string `json:"name" validate:"required,min=5,max=50"`
		IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,url"`
	}

	ReqLogin struct {
		NIP      int    `json:"nip" validate:"required"`
		Password string `json:"password" validate:"required,min=5,max=33"`
	}

	ParamGetUser struct {
		UserID    string `json:"userId"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		Name      string `json:"name"`
		NIP       string `json:"nip"`
		Role      string `json:"role"`
		CreatedAt string `json:"createdAt"`
		Search    string `json:"search"`
	}

	ReqUpdateNurse struct {
		NIP  int    `json:"nip" validate:"required"`
		Name string `json:"name" validate:"required,min=5,max=50"`
	}

	ReqAccessNurse struct {
		Password string `json:"password" validate:"required,min=5,max=33"`
	}

	ResRegister struct {
		UserID      string `json:"userId"`
		NIP         int    `json:"nip,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}

	ResRegisterNurse struct {
		UserID string `json:"userId"`
		NIP    int    `json:"nip,omitempty"`
		Name   string `json:"name"`
	}

	ResLogin struct {
		UserID      string `json:"userId"`
		NIP         int    `json:"nip,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}

	ResGetUser struct {
		UserID    string `json:"userId"`
		NIP       int    `json:"nip,omitempty"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
	}

	ReqUpdateAccount struct {
		ImageURL string `json:"imageUrl" validate:"required,url"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
	}
)

func (d *ReqRegister) ToUserEntity(cryptCost int) entity.User {
	return entity.User{Name: d.Name, Password: auth.HashPassword(d.Password, cryptCost), NIP: strconv.Itoa(d.NIP)}
}

func (d *ReqRegisterNurse) ToNurseEntity(cryptCost int) entity.User {
	return entity.User{Name: d.Name, Password: auth.HashPassword("password", cryptCost), NIP: strconv.Itoa(d.NIP)}
}

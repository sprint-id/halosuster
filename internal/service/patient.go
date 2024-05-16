package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type PatientService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newPatientService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *PatientService {
	return &PatientService{repo, validator, cfg}
}

func (u *PatientService) CreatePatient(ctx context.Context, body dto.ReqCreatePatient, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	patient := body.ToPatientEntity()
	err = u.repo.Patient.CreatePatient(ctx, sub, patient)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrDuplicate
		}
		return err
	}

	return nil
}

func (u *PatientService) GetPatient(ctx context.Context, param dto.ParamGetPatient, sub string) ([]dto.ResGetPatient, error) {
	res, err := u.repo.Patient.GetPatient(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

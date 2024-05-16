package service

import (
	"github.com/go-playground/validator/v10"

	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type Service struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg

	User    *UserService
	Record  *RecordService
	Patient *PatientService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.User = newUserService(repo, validator, cfg)
	service.Record = newRecordService(repo, validator, cfg)
	service.Patient = newPatientService(repo, validator, cfg)

	return &service
}

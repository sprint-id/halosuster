package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type RecordService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newRecordService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *RecordService {
	return &RecordService{repo, validator, cfg}
}

// {
// 	"identityNumber": 123123, // not null, should be 16 digit
// 	"symptoms": "", // not null, minLength 1, maxLength 2000,
// 	"medications" : "" // not null, minLength 1, maxLength 2000
// }

func (u *RecordService) AddRecord(ctx context.Context, body dto.ReqAddRecord, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return ierr.ErrBadRequest
	}

	record := body.ToRecordEntity(sub)
	err = u.repo.Record.AddRecord(ctx, sub, record)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}

func (u *RecordService) GetRecord(ctx context.Context, param dto.ParamGetRecord, sub string) ([]dto.ResGetRecord, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Record.GetRecord(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

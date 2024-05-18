package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
	"github.com/sprint-id/eniqilo-server/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newUserService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *UserService {
	return &UserService{repo, validator, cfg}
}

func (u *UserService) RegisterIT(ctx context.Context, body dto.ReqRegister) (dto.ResRegister, error) {
	res := dto.ResRegister{}
	var nip string

	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return res, ierr.ErrBadRequest
	}

	// convert int to string
	nip = strconv.Itoa(body.NIP)
	// validate nip must between 13-15 characters
	if len(nip) < 13 && len(nip) > 15 {
		return res, ierr.ErrBadRequest
	}
	// - first until third digit, should start with `615`
	if nip[:3] != "615" {
		return res, ierr.ErrBadRequest
	}
	// - the fourth digit, if it's male, fill it with `1`, else `2`
	if nip[3] != '1' && nip[3] != '2' {
		return res, ierr.ErrBadRequest
	}
	// - the fifth and ninth digit, fill it with a year, starts from `2000` till current year
	currentYear := time.Now().Year()
	year, err := strconv.Atoi(nip[4:8])
	if err != nil || year < 2000 || year > currentYear {
		return res, ierr.ErrBadRequest
	}
	// - the tenth and twelfth, fill it with month, starts from `01` till `12`
	month, err := strconv.Atoi(nip[8:10])
	if err != nil || month < 1 || month > 12 {
		return res, ierr.ErrBadRequest
	}

	user := body.ToUserEntity(u.cfg.BCryptSalt)
	userID, err := u.repo.User.Insert(ctx, user)
	if err != nil {
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: userID})
	if err != nil {
		return res, err
	}

	res.UserID = userID
	res.NIP = body.NIP
	res.Name = body.Name
	res.AccessToken = token

	return res, nil
}

func (u *UserService) RegisterNurse(ctx context.Context, body dto.ReqRegisterNurse) (dto.ResRegisterNurse, error) {
	res := dto.ResRegisterNurse{}
	var nip string

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	// check Image URL if invalid or not complete URL
	if !isValidURL(body.IdentityCardScanImg) {
		return res, ierr.ErrBadRequest
	}

	// convert int to string
	nip = strconv.Itoa(body.NIP)
	// validate nip must between 13-15 characters
	if len(nip) < 13 && len(nip) > 15 {
		return res, ierr.ErrBadRequest
	}
	// - first until third digit, should start with `303`
	if nip[:3] != "303" {
		return res, ierr.ErrBadRequest
	}
	// - the fourth digit, if it's male, fill it with `1`, else `2`
	if nip[3] != '1' && nip[3] != '2' {
		return res, ierr.ErrBadRequest
	}
	// - the fifth and ninth digit, fill it with a year, starts from `2000` till current year
	currentYear := time.Now().Year()
	year, err := strconv.Atoi(nip[4:8])
	if err != nil || year < 2000 || year > currentYear {
		return res, ierr.ErrBadRequest
	}
	// - the tenth and twelfth, fill it with month, starts from `01` till `12`
	month, err := strconv.Atoi(nip[8:10])
	if err != nil || month < 1 || month > 12 {
		return res, ierr.ErrBadRequest
	}

	user := body.ToNurseEntity(u.cfg.BCryptSalt)
	userID, err := u.repo.User.Insert(ctx, user)
	if err != nil {
		return res, err
	}

	res.UserID = userID
	res.NIP = body.NIP
	res.Name = body.Name

	return res, nil
}

func (u *UserService) LoginIT(ctx context.Context, body dto.ReqLogin) (dto.ResLogin, error) {
	res := dto.ResLogin{}
	var nip string

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	// convert int to string
	nip = strconv.Itoa(body.NIP)

	user, err := u.repo.User.GetByNIP(ctx, nip)
	if err != nil {
		return res, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.UserID = user.ID
	res.NIP, err = strconv.Atoi(user.NIP)
	if err != nil {
		return res, err
	}
	res.Name = user.Name
	res.AccessToken = token

	return res, nil
}

func (u *UserService) LoginNurse(ctx context.Context, body dto.ReqLogin) (dto.ResLogin, error) {
	res := dto.ResLogin{}
	var nip string

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	// convert int to string
	nip = strconv.Itoa(body.NIP)

	user, err := u.repo.User.GetByNIP(ctx, nip)
	if err != nil {
		return res, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.UserID = user.ID
	res.NIP, err = strconv.Atoi(user.NIP)
	if err != nil {
		return res, err
	}
	res.Name = user.Name
	res.AccessToken = token

	return res, nil
}

func (u *UserService) GetUser(ctx context.Context, param dto.ParamGetUser, sub string) ([]dto.ResGetUser, error) {
	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.User.GetUser(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserService) UpdateNurse(ctx context.Context, body dto.ReqUpdateNurse, id string) error {
	var nip string
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return ierr.ErrBadRequest
	}

	// convert int to string
	nip = strconv.Itoa(body.NIP)
	// validate nip must between 13-15 characters
	if len(nip) < 13 || len(nip) > 15 {
		return ierr.ErrBadRequest
	}
	// check if nip is exist
	var count int
	count, err = u.repo.User.CountNIP(ctx, nip)
	if err != nil {
		return err
	}

	if count > 0 {
		return ierr.ErrDuplicate
	}

	// - first until third digit, should start with `303`
	if nip[:3] != "303" {
		return ierr.ErrNotFound
	}
	// - the fourth digit, if it's male, fill it with `1`, else `2`
	if nip[3] != '1' && nip[3] != '2' {
		return ierr.ErrBadRequest
	}
	// - the fifth and ninth digit, fill it with a year, starts from `2000` till current year
	currentYear := time.Now().Year()
	year, err := strconv.Atoi(nip[4:8])
	if err != nil || year < 2000 || year > currentYear {
		return ierr.ErrBadRequest
	}
	// - the tenth and twelfth, fill it with month, starts from `01` till `12`
	month, err := strconv.Atoi(nip[8:10])
	if err != nil || month < 1 || month > 12 {
		return ierr.ErrBadRequest
	}

	err = u.repo.User.UpdateNurse(ctx, body, id)
	return err
}

func (u *UserService) DeleteNurse(ctx context.Context, id string) error {
	err := u.repo.User.DeleteNurse(ctx, id)
	return err
}

func (u *UserService) AccessNurse(ctx context.Context, body dto.ReqAccessNurse, id string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	password := auth.HashPassword(body.Password, u.cfg.BCryptSalt)

	err = u.repo.User.AccessNurse(ctx, password, id)
	return err
}

func isValidURL(urlString string) bool {
	// url validation using regex
	fmt.Printf("urlString: %s\n", urlString)
	regex := regexp.MustCompile(`^(https?|ftp)://[^/\s]+\.[^/\s]+(?:/.*)?(?:\.[^/\s]+)?$`)
	return regex.MatchString(urlString)
}

package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
)

type patientRepo struct {
	conn *pgxpool.Pool
}

func newPatientRepo(conn *pgxpool.Pool) *patientRepo {
	return &patientRepo{conn}
}

func (mr *patientRepo) CreatePatient(ctx context.Context, sub string, patient entity.Patient) error {
	// Start a transaction with serializable isolation level
	tx, err := mr.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO patients (user_id, identity_number, phone_number, name, birth_date, gender, identity_card_scan_img, created_at)
	VALUES ( $1, $2, $3, $4, $5, $6, $7, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var id string
	err = tx.QueryRow(ctx, q, sub,
		patient.IdentityNumber,
		patient.PhoneNumber,
		patient.Name,
		patient.BirthDate,
		patient.Gender,
		patient.IdentityCardScanImg,
	).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// func (mr *patientRepo) RegisterCustomer(ctx context.Context, sub string, customer entity.Customer) (dto.ResGetPatient, error) {
// 	q := `INSERT INTO patients (user_id, phone_number, name, created_at)
// 	VALUES ( $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

// 	var id string
// 	err := mr.conn.QueryRow(ctx, q, sub, customer.PhoneNumber, customer.Name).Scan(&id)
// 	if err != nil {
// 		if pgErr, ok := err.(*pgconn.PgError); ok {
// 			if pgErr.Code == "23505" {
// 				return dto.ResGetPatient{}, ierr.ErrDuplicate
// 			}
// 		}
// 		return dto.ResGetPatient{}, err
// 	}

// 	return dto.ResGetPatient{UserID: id, PhoneNumber: customer.PhoneNumber, Name: customer.Name}, nil
// }

func (mr *patientRepo) GetPatient(ctx context.Context, param dto.ParamGetPatient, sub string) ([]dto.ResGetPatient, error) {
	var query strings.Builder

	query.WriteString("SELECT id, identity_number, phone_number, name, birth_date, gender, created_at FROM patients WHERE 1=1 ")

	// param identityNumber limit the output based on the user id
	if param.IdentityNumber != "" {
		query.WriteString(fmt.Sprintf("AND identity_number = '%s' ", param.IdentityNumber))
	}

	// param phone number: it should search by wildcard (ex: if search by phoneNumber=+62 then customer with phone number +628123... should appear, but phoneNumber=123 will not show that)
	if param.PhoneNumber != "" {
		query.WriteString(fmt.Sprintf("AND phone_number LIKE '%s' ", fmt.Sprintf("%%%s%%", param.PhoneNumber)))
	}

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	// param createdAt sort by created time asc or desc, if value is wrong, just ignore the param
	if param.CreatedAt == "asc" && param.Offset == 0 {
		query.WriteString("ORDER BY created_at ASC ")
	} else if param.CreatedAt == "desc" && param.Offset == 0 {
		query.WriteString("ORDER BY created_at DESC ")
	} else if param.Offset == 0 {
		query.WriteString("ORDER BY created_at DESC ")
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	rows, err := mr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []dto.ResGetPatient
	for rows.Next() {
		var patient dto.ResGetPatient
		err = rows.Scan(&patient.IdentityNumber, &patient.PhoneNumber, &patient.Name, &patient.BirthDate, &patient.Gender, &patient.CreatedAt)
		if err != nil {
			return nil, err
		}

		patients = append(patients, patient)
	}

	return patients, nil
}

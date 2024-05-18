package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	timepkg "github.com/sprint-id/eniqilo-server/pkg/time"
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
	fmt.Printf("sub: %s\n", sub)

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

func (mr *patientRepo) GetPatient(ctx context.Context, param dto.ParamGetPatient, sub string) ([]dto.ResGetPatient, error) {
	var query strings.Builder

	query.WriteString("SELECT identity_number, phone_number, name, birth_date, gender, created_at FROM patients WHERE 1=1 ")

	// param identityNumber if 615 then 3452615789123456 should appear
	if param.IdentityNumber != "" {
		query.WriteString(fmt.Sprintf("AND identity_number LIKE '%s' ", fmt.Sprintf("%%%s%%", param.IdentityNumber)))
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

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	// show query
	fmt.Println(query.String())

	rows, err := mr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []dto.ResGetPatient
	for rows.Next() {
		var patient dto.ResGetPatient
		var identityNumber string
		var createdAt int64
		err = rows.Scan(&identityNumber, &patient.PhoneNumber, &patient.Name, &patient.BirthDate, &patient.Gender, &createdAt)
		if err != nil {
			return nil, err
		}

		patient.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		patient.IdentityNumber, _ = strconv.Atoi(identityNumber)
		patients = append(patients, patient)
	}

	return patients, nil
}

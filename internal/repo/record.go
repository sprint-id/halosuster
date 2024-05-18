package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	timepkg "github.com/sprint-id/eniqilo-server/pkg/time"
)

type recordRepo struct {
	conn *pgxpool.Pool
}

func newRecordRepo(conn *pgxpool.Pool) *recordRepo {
	return &recordRepo{conn}
}

// {
// 	"identityNumber": 123123, // not null, should be 16 digit
// 	"symptoms": "", // not null, minLength 1, maxLength 2000,
// 	"medications" : "" // not null, minLength 1, maxLength 2000
// }

func (cr *recordRepo) AddRecord(ctx context.Context, sub string, record entity.Record) error {
	// check if identity number is exist
	qCheck := `SELECT COUNT(*) FROM patients WHERE identity_number=$1`
	var count int

	err := cr.conn.QueryRow(ctx, qCheck, record.IdentityNumber).Scan(&count)
	if err != nil {
		return err
	}

	// check if identity number is exist, if not exist, return error bad request
	if count == 0 {
		return ierr.ErrBadRequest
	}

	// add record
	q := `INSERT INTO medical_records (user_id, patient_identifier, symptoms, medications, created_at)
	VALUES ( $1, $2, $3, $4, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var id string
	err = cr.conn.QueryRow(ctx, q, sub, record.IdentityNumber, record.Symptoms, record.Medications).Scan(&id)
	if err != nil {
		fmt.Printf("error query: %v\n", err)
		return err
	}

	return nil
}

func (cr *recordRepo) GetRecord(ctx context.Context, param dto.ParamGetRecord, sub string) ([]dto.ResGetRecord, error) {
	var query strings.Builder

	query.WriteString("SELECT patient_identifier, user_id, symptoms, medications, created_at FROM medical_records WHERE 1=1 ")

	// param id
	if param.ID != "" {
		id, err := strconv.Atoi(param.ID)
		if err != nil {
			return nil, err
		}
		query.WriteString(fmt.Sprintf("AND id = %d ", id))
	}

	// param identityNumber
	if param.IdentityNumber != "" {
		query.WriteString(fmt.Sprintf("AND identity_number = '%s' ", param.IdentityNumber))
	}

	// param userId
	if param.UserId != "" {
		query.WriteString(fmt.Sprintf("AND user_id = '%s' ", param.UserId))
	}

	// param NIP, value must be number. If nip=615 user with nip 61512321321 should be included and user_id should be creator of the record
	if param.NIP != "" {
		query.WriteString(fmt.Sprintf("AND nip LIKE '%s' ", fmt.Sprintf("%%%s%%", param.NIP)))
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

	rows, err := cr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetRecord, 0, 10)
	for rows.Next() {
		var identityCardScanImg sql.NullString
		var createdAt int64
		var patientIdentifierInRecord, userId string

		result := dto.ResGetRecord{}
		err := rows.Scan(
			&patientIdentifierInRecord,
			&userId,
			&result.Symptoms,
			&result.Medications,
			&createdAt)
		if err != nil {
			return nil, err
		}

		patient := dto.ResIdentityDetail{}
		var patientIdentifierInPatient string
		// get patient detail
		q := `SELECT identity_number, phone_number, name, birth_date, gender, identity_card_scan_img FROM patients WHERE identity_number=$1`

		err = cr.conn.QueryRow(ctx, q, patientIdentifierInRecord).Scan(&patientIdentifierInPatient, &patient.PhoneNumber, &patient.Name, &patient.BirthDate, &patient.Gender, &patient.IdentityCardScanImg)
		if err != nil {
			return nil, err
		}

		// get user detail
		createdBy := dto.ResCreatedBy{}
		var nip string
		q = `SELECT nip, name, id FROM users WHERE id=$1`

		err = cr.conn.QueryRow(ctx, q, userId).Scan(&nip, &createdBy.Name, &createdBy.UserID)
		if err != nil {
			return nil, err
		}

		patient.IdentityNumber, _ = strconv.Atoi(patientIdentifierInPatient)
		patient.IdentityCardScanImg = identityCardScanImg.String
		createdBy.NIP, _ = strconv.Atoi(nip)
		result.IdentityDetail = patient
		result.CreatedBy = createdBy
		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		results = append(results, result)
	}

	return results, nil
}

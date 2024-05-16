package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	timepkg "github.com/sprint-id/eniqilo-server/pkg/time"
)

type userRepo struct {
	conn *pgxpool.Pool
}

func newUserRepo(conn *pgxpool.Pool) *userRepo {
	return &userRepo{conn}
}

func (u *userRepo) Insert(ctx context.Context, user entity.User) (string, error) {
	credVal := user.NIP
	q := `INSERT INTO users (id, name, nip, password, created_at)
	VALUES (gen_random_uuid(), $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var userID string
	err := u.conn.QueryRow(ctx, q,
		user.Name, credVal, user.Password).Scan(&userID)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", ierr.ErrDuplicate
			}
		}
		return "", err
	}

	return userID, nil
}

func (u *userRepo) GetByNIP(ctx context.Context, cred string) (entity.User, error) {
	user := entity.User{}
	q := `SELECT id, name, nip, password FROM users
	WHERE nip = $1`

	var nip sql.NullString

	err := u.conn.QueryRow(ctx,
		q, cred).Scan(&user.ID, &user.Name, &nip, &user.Password)

	user.NIP = nip.String

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	user := entity.User{}
	q := `SELECT nip, name, password FROM users
	WHERE id = $1`

	var nip sql.NullString

	err := u.conn.QueryRow(ctx,
		q, id).Scan(&nip, &user.Name, &user.Password)

	user.NIP = nip.String

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) GetUser(ctx context.Context, param dto.ParamGetUser, sub string) ([]dto.ResGetUser, error) {
	var query strings.Builder

	query.WriteString("SELECT id, nip, name, created_at FROM users WHERE 1=1 ")

	// param user id
	if param.UserID != "" {
		query.WriteString(fmt.Sprintf("AND id = '%s' ", param.UserID))
	}
	// param name: case insensitive, if een in between the name it will be included
	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}
	// param nip, value must be number. If nip=615 user with nip 61512321321 should be included
	if param.NIP != "" {
		query.WriteString(fmt.Sprintf("AND nip LIKE '%s' ", fmt.Sprintf("%%%s%%", param.NIP)))
	}
	// param role, role 'it' starts with '615' and role 'nurse' starts with '303'
	if param.Role != "" {
		if param.Role == "it" {
			query.WriteString("AND nip LIKE '615%' ")
		} else if param.Role == "nurse" {
			query.WriteString("AND nip LIKE '303%' ")
		}
	}
	// param created at, value must be ascending or descending
	if param.CreatedAt != "" {
		query.WriteString(fmt.Sprintf("ORDER BY created_at %s ", param.CreatedAt))
	}
	// param limit, value must be number. default limit is 5
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	rows, err := u.conn.Query(ctx, query.String()) // Replace $1 with sub
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetUser, 0, 10)
	for rows.Next() {
		var createdAt int64

		result := dto.ResGetUser{}
		err := rows.Scan(
			&result.UserID,
			&result.NIP,
			&result.Name,
			&createdAt)
		if err != nil {
			return nil, err
		}

		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		results = append(results, result)
	}

	return results, nil
}

func (u *userRepo) UpdateNurse(ctx context.Context, body dto.ReqUpdateNurse, sub string) error {
	q := `UPDATE users SET name = $1, password = $2 WHERE id = $3`
	_, err := u.conn.Exec(ctx, q,
		body.Name, body.Password, sub)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (u *userRepo) DeleteNurse(ctx context.Context, id string) error {
	q := `DELETE FROM users WHERE id = $1`
	_, err := u.conn.Exec(ctx, q, id)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) AccessNurse(ctx context.Context, password, id string) error {
	q := `UPDATE users SET password = $1 WHERE id = $2`

	_, err := u.conn.Exec(ctx, q, password, id)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) LookUp(ctx context.Context, id string) error {
	q := `SELECT 1 FROM users WHERE id = $1`

	v := 0
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

func (u *userRepo) UpdateAccount(ctx context.Context, id, name, url string) error {
	q := `UPDATE users SET image_url = $1, name = $2 WHERE id = $3`
	_, err := u.conn.Exec(ctx, q,
		url, name, id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (u *userRepo) GetNameBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT name FROM users WHERE id = $1`

	v := ""
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		return "", err
	}

	return v, nil
}

func (u *userRepo) GetEmailBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT email FROM users WHERE id = $1`

	v := ""
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		return "", err
	}

	return v, nil
}

// func (u *userRepo) GetNameByID(ctx context.Context, id string) (string, error) {
// 	name := ""
// 	err := u.conn.QueryRow(ctx,
// 		`SELECT name FROM users
// 		WHERE id = $1`,
// 		id).Scan(&name)
// 	if err != nil {
// 		if err.Error() == "no rows in result set" {
// 			return "", ierr.ErrNotFound
// 		}
// 		if pgErr, ok := err.(*pgconn.PgError); ok {
// 			if pgErr.Code == "22P02" {
// 				return "", ierr.ErrNotFound
// 			}
// 		}
// 		return "", err
// 	}

// 	return name, nil
// }

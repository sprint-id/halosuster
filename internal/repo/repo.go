package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User    *userRepo
	Record  *recordRepo
	Patient *patientRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Record = newRecordRepo(conn)
	repo.Patient = newPatientRepo(conn)

	return &repo
}

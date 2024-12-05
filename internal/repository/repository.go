package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repository interface {
	SetUserInfo(userID uuid.UUID, userIP string, token []byte) error
}

type pgRepo struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) SetUserInfo(userID uuid.UUID, userIP string, token []byte) error {
	result, err := r.db.Exec(querySetUserInfo, userID, userIP, token)
	if err != nil {
		err := errors.Errorf("pgRepo.SetUserInfo %v", err)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.SetUserInfo %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.SetUserInfo: no rows affected")
		return err
	}

	return nil
}



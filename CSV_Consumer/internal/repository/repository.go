package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
)

type Store interface {
	Save(user *model.User) error
	CheckForAccepted(userId string) (bool, error)
}

type Repository struct {
	Store
}

func NewStorage(db *sqlx.DB) *Repository {
	return &Repository{
		Store: NewStoreRepository(db),
	}
}

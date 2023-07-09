package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
)

type StoreRepository struct {
	db *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

func (sr *StoreRepository) Save(user *model.User) error {
	const query = `insert into users (id, full_name, username, email, phone_number) VALUES ($1, $2, $3, $4, $5)`

	_, err := sr.db.Exec(query, user.Id, user.FullName, user.Username, user.Email, user.PhoneNumber)

	return err
}

func (sr *StoreRepository) CheckForAccepted(userId string) (bool, error) {
	var isAccepted bool

	const query = `select u.accepted from users u where u.id = $1 `

	row := sr.db.QueryRow(query, userId)

	err := row.Scan(&isAccepted)

	return isAccepted, err
}

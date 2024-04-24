package repository

import "github.com/jmoiron/sqlx"

type Repostitory struct {
	repo *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repostitory {
	return &Repostitory{
		repo: db,
	}
}

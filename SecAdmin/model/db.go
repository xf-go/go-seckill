package model

import "github.com/jmoiron/sqlx"

var DB *sqlx.DB

func Init(db *sqlx.DB) {
	DB = db
}

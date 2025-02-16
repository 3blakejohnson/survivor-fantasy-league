package dao

import (
	"database/sql"
	"log"

	"github.com/3blakejohnson/survivor-fantasy-league/backend/models"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) Create(userInfo models.User) error {
	_, err := dao.db.Exec(`INSERT INTO users (username, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4)`, userInfo.Username, userInfo.PasswordHash, userInfo.FirstName, userInfo.LastName)
	if err != nil {
		return err
	}
	log.Println("User successfully created")
	return nil
}

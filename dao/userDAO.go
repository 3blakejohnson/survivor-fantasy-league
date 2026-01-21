package dao

import (
	"database/sql"
	"log"
	"time"

	"github.com/3blakejohnson/survivor-fantasy-league/models"
)

type UserDAO interface {
	Create(models.User) error
	Get(int64) (*models.User, error)
	GetByUsername(string) (*models.User, error)
}

type userDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) UserDAO {
	return &userDAO{db: db}
}

func (dao *userDAO) Create(userInfo models.User) error {
	if userInfo.CreatedAt == nil {
		now := time.Now()
		userInfo.CreatedAt = &now
	}
	_, err := dao.db.Exec(`INSERT INTO users (username, password_hash, first_name, last_name, created_at)
		VALUES ($1, $2, $3, $4, $5)`, userInfo.Username, userInfo.PasswordHash, userInfo.FirstName, userInfo.LastName, userInfo.CreatedAt)
	if err != nil {
		return err
	}
	log.Println("User successfully created")
	return nil
}

func (dao *userDAO) Get(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT
		id,
		created_at,
		username,
		password_hash,
		first_name,
		last_name
	FROM users
	WHERE id = $1;`

	result := dao.db.QueryRow(query, id)
	err := result.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (dao *userDAO) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT
		id,
		created_at,
		username,
		password_hash,
		first_name,
		last_name
	FROM users
	WHERE username = $1;`

	result := dao.db.QueryRow(query, username)
	err := result.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

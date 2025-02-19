package dao

import "database/sql"

type DAOManager interface {
	User() UserDAO
}

type daoManager struct {
	UserDAO UserDAO
	db      *sql.DB
}

func NewDAOManager(database *sql.DB) DAOManager {
	return &daoManager{db: database}
}

func (m *daoManager) User() UserDAO {
	if m.UserDAO == nil {
		m.UserDAO = NewUserDAO(m.db)
	}
	return m.UserDAO
}

package dao

import "database/sql"

type DAOManager interface {
	User() UserDAO
	Episode() EpisodeDAO
}

type daoManager struct {
	UserDAO    UserDAO
	EpisodeDAO EpisodeDAO
	db         *sql.DB
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

func (m *daoManager) Episode() EpisodeDAO {
	if m.EpisodeDAO == nil {
		m.EpisodeDAO = NewEpisodeDAO(m.db)
	}
	return m.EpisodeDAO
}

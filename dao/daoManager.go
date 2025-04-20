package dao

import "database/sql"

type DAOManager interface {
	User() UserDAO
	Episode() EpisodeDAO
	League() LeagueDAO
}

type daoManager struct {
	UserDAO    UserDAO
	EpisodeDAO EpisodeDAO
	LeagueDAO  LeagueDAO
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

func (m *daoManager) League() LeagueDAO {
	if m.LeagueDAO == nil {
		m.LeagueDAO = NewLeagueDAO(m.db)
	}
	return m.LeagueDAO
}

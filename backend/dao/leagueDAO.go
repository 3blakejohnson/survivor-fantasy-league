package dao

import (
	"database/sql"
	"log"

	"github.com/3blakejohnson/survivor-fantasy-league/backend/models"
)

type LeagueDAO interface {
	Create(models.League) error
	GetByInviteCode(string) (*models.League, error)
}

type leagueDAO struct {
	db *sql.DB
}

func NewLeagueDAO(db *sql.DB) LeagueDAO {
	return &leagueDAO{db: db}
}

func (l *leagueDAO) Create(league models.League) error {
	_, err := l.db.Exec(`INSERT INTO leagues (created_at, name, owner_id, invite_code, season)
		VALUES (now(), $1, $2, $3, $4)`, league.Name, league.OwnerID, league.InviteCode, league.Season)
	if err != nil {
		return err
	}
	log.Println("League successfully created")
	return nil
}

func (l *leagueDAO) GetByInviteCode(code string) (*models.League, error) {
	league := &models.League{}
	query := `
		SELECT
			id,
			created_at,
			name,
			owner_id,
			invite_code,
			season
		FROM leagues
		WHERE invite_code = $1
	`

	result := l.db.QueryRow(query, code)
	err := result.Scan(
		&league.ID,
		&league.CreatedAt,
		&league.Name,
		&league.OwnerID,
		&league.InviteCode,
		&league.Season,
	)
	if err != nil {
		return nil, err
	}

	return league, nil
}

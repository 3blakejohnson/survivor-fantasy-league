package dao

import (
	"database/sql"
	"log"

	"github.com/3blakejohnson/survivor-fantasy-league/models"
)

type EpisodeDAO interface {
	Create(models.Episode) error
	Get(int64, int64) (*models.Episode, error)
}

type episodeDAO struct {
	db *sql.DB
}

func NewEpisodeDAO(db *sql.DB) EpisodeDAO {
	return &episodeDAO{db: db}
}

func (e *episodeDAO) Create(episode models.Episode) error {
	_, err := e.db.Exec(`INSERT INTO episodes (season, episode, title, air_date)
		VALUES ($1, $2, $3, $4)`, episode.Season, episode.Episode, episode.Title, episode.AirDate)
	if err != nil {
		return err
	}
	log.Println("Episode successfully created")
	return nil
}

func (e *episodeDAO) Get(seasonNum int64, episodeNum int64) (*models.Episode, error) {
	episode := &models.Episode{}
	query := `
		SELECT
			id,
			season,
			episode,
			title,
			air_date
		FROM episodes
		WHERE season = $1
		AND episode = $2;
	`
	result := e.db.QueryRow(query, seasonNum, episodeNum)
	err := result.Scan(
		&episode.ID,
		&episode.Season,
		&episode.Episode,
		&episode.Title,
		&episode.AirDate,
	)
	if err != nil {
		return nil, err
	}

	return episode, nil
}

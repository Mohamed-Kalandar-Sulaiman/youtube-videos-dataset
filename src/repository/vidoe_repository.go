package repository

import (
	"database/sql"
	"log"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/models"
)

type VideoRepository struct {
	DB *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{DB: db}
}



func (r *VideoRepository) GetVideos() ([]models.Video, error) {
	var videos []models.Video
	rows, err := r.DB.Query("SELECT * FROM videos")
	if err != nil {
		log.Printf("Error fetching videos: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var video models.Video
		if err := rows.Scan(&video.ID, &video.Title); err != nil {
			log.Printf("Error scanning video: %v", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

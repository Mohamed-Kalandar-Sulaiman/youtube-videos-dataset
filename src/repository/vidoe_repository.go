package repository

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/models"
)

type VideoRepository struct {
	DB *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{DB: db}
}



	
func (r *VideoRepository) GetFilteredVideos(options models.VideoFilterOptions) ([]models.Video, string, error) {
	query := "SELECT * FROM videos WHERE 1=1 "
	args := []interface{}{}
	placeholderIndex := 1

	if options.PublishedAfter != "" {
		query += fmt.Sprintf(" AND published_date >= $%d", placeholderIndex)
		args = append(args, options.PublishedAfter)
		placeholderIndex++
	}

	if options.PublishedBefore != "" {
		query += fmt.Sprintf(" AND published_date <= $%d", placeholderIndex)
		args = append(args, options.PublishedBefore)
		placeholderIndex++
	}

	if options.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", placeholderIndex)
		args = append(args, "%"+options.Title+"%")
		placeholderIndex++
	}

	if options.ChannelName != "" {
		query += fmt.Sprintf(" AND channel_name ILIKE $%d", placeholderIndex)
		args = append(args, "%"+options.ChannelName+"%")
		placeholderIndex++
	}
	

	if options.PageKey != "" {
		query += fmt.Sprintf(" AND published_date < $%d", placeholderIndex)
		args = append(args, options.PageKey)
		placeholderIndex++
	}

	// ! By default we will add PageSize
	if options.PageSize != -1 {
		query += fmt.Sprintf(" ORDER BY published_date DESC LIMIT $%d", placeholderIndex)
		args = append(args, options.PageSize)
		placeholderIndex++
	}


	// fmt.Println("Queryy : ", query)
	// fmt.Println("Argss : ", args)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var videos []models.Video
	var nextPageKey string

	for rows.Next() {
		var video models.Video
		if err := rows.Scan(&video.ID,
			 				&video.Title,
							&video.Description,
							&video.PublishedDate, 
							&video.ThumbnailURL, 
							&video.ChannelName, 
							&video.ChannelId); 
							err != nil {
									return nil, "", err
								}
		videos = append(videos, video)
		nextPageKey = video.PublishedDate 
	}

	if len(videos) == 0 {
		nextPageKey = ""
	}
	return videos, nextPageKey, nil
}







func (r *VideoRepository) AddVideo(video models.Video) error {
	sqlStatement := `
		INSERT INTO videos (id, title, description, published_date, thumbnail_url, channel_name, channel_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING;` 

	_, err := r.DB.Exec(sqlStatement,
		video.ID,
		video.Title,
		video.Description,
		video.PublishedDate,
		video.ThumbnailURL,
		video.ChannelName,
		video.ChannelId,
	)

	if err != nil {
		log.Printf("Error inserting video: %v", err)
		return fmt.Errorf("failed to insert video: %w", err)
	}
	return nil
}


func (r *VideoRepository) RunMigrations(filePath string) error {
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	sqlScript := string(sqlBytes)

	sqlStatements := strings.Split(sqlScript, ";")

	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue 
		}
		_, err := r.DB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %v. Statement: %s", err, stmt)
		}
	}
	log.Println("Migrations executed successfully")
	return nil
}
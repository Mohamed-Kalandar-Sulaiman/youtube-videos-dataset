package services

import (
	"log"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/models"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/repository"
)

type VideoService struct {
	VideoRepo *repository.VideoRepository
}

func NewVideoService(videoRepo *repository.VideoRepository) *VideoService {
	return &VideoService{
		VideoRepo: videoRepo,
	}
}

func (s *VideoService) GetAllVideos() ([]models.Video, error) {
	videos, err := s.VideoRepo.GetVideos()
	if err != nil {
		log.Printf("Error in GetAllVideos: %v", err)
		return nil, err
	}
	return videos, nil
}

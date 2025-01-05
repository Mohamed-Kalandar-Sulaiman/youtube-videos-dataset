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

func (s *VideoService) GetAllVideos(filters models.VideoFilterOptions) ([]models.Video,string,  error) {
	videos, nextPageKey,  err := s.VideoRepo.GetFilteredVideos( filters)
	

	if err != nil {
		log.Printf("Error in GetAllVideos: %v", err)
		return nil, "", err
	}

	if len(videos) == 0 {
		videos = []models.Video{}
	}
	return videos,nextPageKey,  nil
}

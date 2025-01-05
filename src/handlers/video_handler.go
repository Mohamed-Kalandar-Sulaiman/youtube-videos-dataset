package handlers

import (
	"log"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/services"
	"github.com/gofiber/fiber/v2"
)

type VideoHandler struct {
	VideoService *services.VideoService
}

func NewVideoHandler(videoService *services.VideoService) *VideoHandler {
	return &VideoHandler{
		VideoService: videoService,
	}
}

func (h *VideoHandler) GetVideos(c *fiber.Ctx) error {
	videos, err := h.VideoService.GetAllVideos()
	if err != nil {
		log.Printf("Error fetching videos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching videos"})
	}
	return c.JSON(videos)
}

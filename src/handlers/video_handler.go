package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/models"
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
	// Query parameters
	pageSizeStr := c.Query("pageSize", "10")
	pageKey     := c.Query("pageKey", "")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid pageSize. Must be a positive integer."})
	}

	publishedAfter  := c.Query("PublishedAfter","")
	publishedBefore := c.Query("PublishedBefore","")
	title           := c.Query("title","")
	channelName     := c.Query("channelName","")


	if publishedAfter != "" {
		_ , err := time.Parse(time.RFC3339, publishedAfter)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid PublishedAfter date format. Use RFC3339 format."})
		}
	}

	if publishedBefore != "" {
		_, err := time.Parse(time.RFC3339, publishedBefore)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid PublishedBefore date format. Use RFC3339 format."})
		}
	}

	// Build filter options
	filterOptions := models.VideoFilterOptions{
		PageSize       : pageSize,
		PageKey        : pageKey,
		PublishedAfter : publishedAfter,
		PublishedBefore: publishedBefore,
		Title          : title,
		ChannelName    : channelName,
	}

	videos, nextPageKey, err := h.VideoService.GetAllVideos(filterOptions)
	if err != nil {
		log.Printf("Error fetching videos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching videos"})
	}

	return c.JSON(fiber.Map{
		"data"       : videos,
		"nextPageKey": nextPageKey,
	})
}

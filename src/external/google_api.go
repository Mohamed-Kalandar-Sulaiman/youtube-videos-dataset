package external

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/models"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/repository"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)


type BackGroundProcess struct{
	service *youtube.Service
	videoRepo *repository.VideoRepository

}



func NewBackGroundProcess(apiKey string, videoRepo *repository.VideoRepository) *BackGroundProcess {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
	instance := &BackGroundProcess{service: service, videoRepo: videoRepo}
	return instance
}


func (b *BackGroundProcess) Start(ctx context.Context) error {
	/* 
	! On each day the latest data from youtube is fetched
	! Assumed that data up until yesterday is loaded already
	! On each day the delta data is loaded
	! Theres no way a vidoe can be uploaded with past date so If historical data is loaded already Its enough to handle new data alone everyday 
	! But for now I am going backwards to colect 500k data Starting from Jan 3 2025...
	*/
	lastProcessedDate, err := b.loadLastProcessedDate()
	if err != nil {
		log.Printf("Error loading last processed date: %v", err)
		return err
	}

	for {
		if err := b.fetchDataForDay(ctx, lastProcessedDate); err != nil {
			log.Printf("Error fetching data: %v", err)
			return err
		}

		lastProcessedDate = lastProcessedDate.AddDate(0, 0, -1) // Decrement and store the currentdate

		if err := b.saveLastProcessedDate(lastProcessedDate); err != nil {
			log.Printf("Error saving last processed date: %v", err)
			return err
		}

		time.Sleep(24 * time.Hour)
	}
}



func (b *BackGroundProcess) fetchDataForDay(ctx context.Context, date time.Time) error {
	defer ctx.Done()
	// if true{
	// 	return nil //! quota exceeded

	// }
	publishedAfter := date.Format(time.RFC3339) // Todays date
	publishedBefore := date.AddDate(0, 0, 1).Format(time.RFC3339) // Tomorrow date

	query := "leetcode" // The search query
	maxResults := int64(50) 

	var nextPageToken string
	totalVideos := 0

	for {
		call := b.service.Search.List([]string{"snippet"}).
			Q(query).
			Type("video").
			PublishedAfter(publishedAfter).
			PublishedBefore(publishedBefore).
			MaxResults(maxResults)

		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		videos, err := call.Do()
		if err != nil {
			if apiErr, ok := err.(*googleapi.Error); ok && apiErr.Code == 403 && apiErr.Message == "Quota exceeded" {
				backoffDuration := time.Second *5
				fmt.Printf("Rate limit exceeded. Retrying in %v...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			}
			return fmt.Errorf("error searching videos: %v", err)
		}
		log.Printf("Retrieved %d videos for date: %s (Page Token: %s)", len(videos.Items), date.Format("2006-01-02"), nextPageToken)

		for _, item := range videos.Items {
			video := models.Video{
				ID:            item.Id.VideoId,
				Title:         item.Snippet.Title,
				Description:   item.Snippet.Description,
				PublishedDate: item.Snippet.PublishedAt,
				ThumbnailURL:  item.Snippet.Thumbnails.High.Url,
				ChannelName:   item.Snippet.ChannelTitle,
				ChannelId:     item.Snippet.ChannelId,
			}

			if err := b.videoRepo.AddVideo(video); err != nil {
				log.Printf("Error adding video to database (ID: %s): %v", video.ID, err)
			}
		}

		totalVideos += len(videos.Items)

		if videos.NextPageToken == "" {
			break
		}
		nextPageToken = videos.NextPageToken
		time.Sleep(time.Second*5) // ! Sleep for 5 secs to avoid bursty rate limit
	}

	log.Printf("Completed fetching %d videos for date: %v", totalVideos, publishedAfter)
	return nil
}






func (b *BackGroundProcess) loadLastProcessedDate() (time.Time, error) {
	filePath := "last_processed_date.txt"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to gte date")
		if os.IsNotExist(err) {
			defaultDate := "2025-01-03T00:00:00Z"
			return time.Parse(time.RFC3339, defaultDate)
		}
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, string(data))
}

func (b *BackGroundProcess) saveLastProcessedDate(date time.Time) error {
	filePath := "last_processed_date.txt"
	return os.WriteFile(filePath, []byte(date.Format(time.RFC3339)), 0644)
}
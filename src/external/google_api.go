package external

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)


type BackGroundProcess struct{
	service *youtube.Service

}



func NewBackGroundProcess(apiKey string) *BackGroundProcess {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
	instance := &BackGroundProcess{service: service}
	return instance
}





func (b *BackGroundProcess)Fetchdata( nextPageToken string) (string, error){

	query := "leetcode"
	publishedAfter := "2025-01-01T00:00:00Z" 
	maxResults := int64(1)

	call := b.service.Search.List([]string{"snippet"}).
		Q(query).
		Type("video").
		PublishedAfter(publishedAfter).
		MaxResults(maxResults).
		PageToken(nextPageToken)

	videos, err := call.Do()

	if err != nil {
		log.Fatalf("Error searching videos: %v", err)
	}

	resultJSON, _ := json.MarshalIndent(videos.Items, "", "  ")
	fmt.Println(string(resultJSON))
	if videos.NextPageToken == "" {
		nextPageToken=""
	}
	nextPageToken = videos.NextPageToken
	return nextPageToken, nil
}


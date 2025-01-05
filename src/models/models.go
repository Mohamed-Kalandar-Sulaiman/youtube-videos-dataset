package models

// Model to represent the video
type Video struct {
	ID              string    `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	PublishedDate   string `json:"published_date"`
	ThumbnailURL    string `json:"thumbnail_url"`
	ChannelName     string `json:"channel_name"`
	ChannelId 		string `json:"channel_id"`
}


type VideoFilterOptions struct {
	PageSize        int
	PageKey         string
	PublishedAfter  string
	PublishedBefore string
	Title           string
	ChannelName     string
}

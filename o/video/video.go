package video

import (
	"LongTM/basic/x/db/mongodb"
)

var VideoTable = mongodb.NewTable("tk_booking", "tkbk", 18)

type Video struct {
	mongodb.BaseModel `bson:",inline"`
	Status            VideoStatus `json:"status"  bson:"status"`
}

type VideoStatus string

const ()

type SearchListReq struct {
	Etag  string `json:"etag"`
	Items []struct {
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Kind    string `json:"kind"`
		Snippet struct {
			ChannelID            string     `json:"channelId"`
			ChannelTitle         string     `json:"channelTitle"`
			Description          string     `json:"description"`
			LiveBroadcastContent string     `json:"liveBroadcastContent"`
			PublishedAt          string     `json:"publishedAt"`
			Thumbnails           Thumbnails `json:"thumbnails"`
			Title                string     `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
	Kind          string `json:"kind"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		ResultsPerPage int `json:"resultsPerPage"`
		TotalResults   int `json:"totalResults"`
	} `json:"pageInfo"`
	RegionCode string `json:"regionCode"`
}

type Thumbnails struct {
	Default DesVideo `json:"default"`
	High    DesVideo `json:"high"`
	Medium  DesVideo `json:"medium"`
}
type DesVideo struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

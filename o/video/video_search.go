package video

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
)

func printChannelsListResults(response *youtube.ChannelListResponse) {
	for _, item := range response.Items {
		fmt.Println(item.Id, ": ", item.Snippet.Title)
	}
}

func channelsListById(service *youtube.Service, part string, id string) {
	call := service.Channels.List(part)
	if id != "" {
		call = call.Id(id)
	}
	response, err := call.Do()
	fmt.Printf("Req", err)
	printChannelsListResults(response)
}

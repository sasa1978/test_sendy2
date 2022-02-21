package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "test1/proto"
)

var (
	port = flag.Int("port", 5555, "The server port")
)

type server struct {
	pb.UnimplementedPlaylistServer
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// Retrieve playlistItems in the specified playlist
func playlistItemsList(service *youtube.Service, part string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	arrpart := []string{part}
	call := service.PlaylistItems.List(arrpart)
	call = call.PlaylistId(playlistId)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	handleError(err, "")
	return response
}

// Retrieve resource for the authenticated user's channel
func channelsListMine(service *youtube.Service, part string) *youtube.ChannelListResponse {
	arrpart := []string{part}
	call := service.Channels.List(arrpart)
	call = call.Mine(true)
	response, err := call.Do()
	handleError(err, "")
	return response
}

func (s *server) Playlist(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	client := getClient(youtube.YoutubeReadonlyScope)
	service, err := youtube.New(client)

	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	response := channelsListMine(service, "contentDetails")
	result := ""

	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads

		fmt.Printf("List %s\r\n", playlistId)

		if in.GetMessage() != playlistId {
			continue
		}

		nextPageToken := ""
		for {
			// Retrieve next set of items in the playlist.
			playlistResponse := playlistItemsList(service, "snippet", playlistId, nextPageToken)

			for _, playlistItem := range playlistResponse.Items {
				title := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId
				result += fmt.Sprintf("%v, (%v)\r\n", title, videoId)
			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
			//fmt.Println()
		}
	}

	//log.Printf("Received: %v", in.GetMessage())
	return &pb.Response{Message: result}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPlaylistServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

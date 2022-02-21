package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	pb "test1/proto"
)

const (
	defaultName = "getPlayList"
)

var (
	// TO-DO тут нужно подставить адрес докера
	addr = flag.String("addr", "localhost:5555", "the address to connect to")
	name = flag.String("name", defaultName, "Name of client")
)

func main() {
	http.HandleFunc("/getPlayList", getPlayListHandler)
	http.Handle("/", http.FileServer(http.Dir("./assets")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func getPlayListHandler(w http.ResponseWriter, r *http.Request) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer conn.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	client := pb.NewPlaylistClient(conn)
	response, err := client.Playlist(context.Background(), &pb.Request{Message: string(body)})

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, response.Message)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"thumbnail/internal/app"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "thumbnail/internal/proto"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %s\n", err.Error())
	}
}

// func main() {
// 	a, err := app.New()
// 	if err != nil {
// 		log.Fatalf("error creating app instance: %s\n", err.Error())
// 	}

// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		if err := a.Run(); err != nil {
// 			log.Fatalf("Application startup error: %s\n", err.Error())
// 		}
// 	}()

// 	<-sigChan
// 	log.Println("Shutting down...")
// 	a.Stop()
// }

func main() {
	isServer := flag.Bool("server", false, "Run in server mode")
	outputDir := flag.String("output", "thumbnails", "Output directory for thumbnails")
	async := flag.Bool("async", false, "Download thumbnails asynchronously")
	flag.Parse()

	urls := flag.Args()

	if *isServer {
		runServer()
	} else {
		if len(urls) == 0 {
			log.Fatal("At least one YouTube URL is required")
		}
		runClient(*outputDir, *async, urls)
	}
}

func runServer() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("error creating app instance: %s\n", err.Error())
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.Run(); err != nil {
			log.Fatalf("application startup error: %s\n", err.Error())
		}
	}()

	<-sigChan
	log.Println("shutting down...")
	a.Stop()
}

func runClient(outputDir string, async bool, urls []string) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("failed to create output directory %s: %v", outputDir, err)
	}

	con, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer con.Close()

	client := pb.NewThumbnailServiceClient(con)

	if async {
		downloadThumbnailsAsync(client, urls, outputDir)
	} else {
		downloadThumbnailsSync(client, urls, outputDir)
	}
}

func downloadThumbnailsSync(client pb.ThumbnailServiceClient, urls []string, outputDir string) {
	for _, url := range urls {
		resp, err := client.GetThumbnail(context.Background(), &pb.GetThumbnailRequest{Url: url})
		if err != nil {
			log.Printf("error getting thumbnail for %s: %v", url, err)
			continue
		}

		if resp.Error != "" {
			log.Printf("server error for %s: %s", url, resp.Error)
			continue
		}

		filename := filepath.Join(outputDir, fmt.Sprintf("%x.jpg", url))
		if err := os.WriteFile(filename, resp.Thumbnail, 0644); err != nil {
			log.Printf("error saving thumbnail %s: %v", filename, err)
			continue
		}

		log.Printf("saved: %s", filename)
	}
}

func downloadThumbnailsAsync(client pb.ThumbnailServiceClient, urls []string, outputDir string) {
	stream, err := client.GetThumbnailAsync(context.Background(), &pb.GetThumbnailsRequestAsync{Urls: urls})
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		if resp.Error != "" {
			log.Printf("server error for %s: %s", resp.Url, resp.Error)
			continue
		}

		filename := filepath.Join(outputDir, fmt.Sprintf("%x.jpg", resp.Url))
		if err := os.WriteFile(filename, resp.Thumbnail, 0644); err != nil {
			log.Printf("error saving thumbnail %s: %v", filename, err)
			continue
		}

		log.Printf("saved: %s", filename)
	}
}

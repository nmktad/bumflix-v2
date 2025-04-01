// package main
//
// import (
// 	"fmt"
//
// 	"github.com/nmktad/bumflix/pkg/ffmpeg"
// )
//
// func main() {
// 	err := ffmpeg.ProcessVideoForStream(
// 		&ffmpeg.Video{
// 			Filename:       "Its.A.Wonderful.Life.1946.2160p.4K.BluRay.x265.10bit.AAC5.1-[YTS.MX].mkv",
// 			FileType:       "mkv",
// 			BucketName:     "movies",
// 			DestBucketName: "movies-processed",
// 			SignedURL:      "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL21vdmllcy9JdHMuQS5Xb25kZXJmdWwuTGlmZS4xOTQ2LjIxNjBwLjRLLkJsdVJheS54MjY1LjEwYml0LkFBQzUuMS0lNUJZVFMuTVglNUQubWt2P1gtQW16LUFsZ29yaXRobT1BV1M0LUhNQUMtU0hBMjU2JlgtQW16LUNyZWRlbnRpYWw9MlVWN1NVMlNZME1RSzNKVUJNVkolMkYyMDI1MDQwMSUyRnVzLWVhc3QtMSUyRnMzJTJGYXdzNF9yZXF1ZXN0JlgtQW16LURhdGU9MjAyNTA0MDFUMDkyNzQ5WiZYLUFtei1FeHBpcmVzPTQzMTk5JlgtQW16LVNlY3VyaXR5LVRva2VuPWV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpoWTJObGMzTkxaWGtpT2lJeVZWWTNVMVV5VTFrd1RWRkxNMHBWUWsxV1NpSXNJbVY0Y0NJNk1UYzBNelUwTWpBeU55d2ljR0Z5Wlc1MElqb2lZblZ0Wm14cGVDSjkuZnMtWlRWT2RIZkFRa0h4ZHpyUUdhN3B1VEZXdnd0VzE3R0ZqQWtuLWd2ZzdCV2xkWFFHV2FKRGc1QTlBektuUEdRc21ZSTlOMGZTenZMWEZxV1dJY0EmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JnZlcnNpb25JZD1udWxsJlgtQW16LVNpZ25hdHVyZT01ZWNlMTg3YzMxY2EzNDgxOTNjMjA5MTgwZGFjMjE0ZGZiZmQyMjIwZmYyZTEwNDkzNzhmNWFjOWU4NzEwOTRk",
// 		},
// 	)
// 	if err != nil {
// 		panic(fmt.Sprintf("Error processing video: %v", err))
// 	}
//
// 	fmt.Println("Video processed successfully")
// }

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/nmktad/bumflix/pkg/env"
)

var (
	minioClient *minio.Client
	envVars     *env.Env
)

func main() {
	http.HandleFunc("/", serveHLS)

	err := env.LoadEnv()
	if err != nil {
		log.Fatal("couldn't load env variables")
	}

	envVars = env.EnvInstance

	minioClient, err = minio.New(string(envVars.Endpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(string(envVars.AccessKeyID), string(envVars.SecretAccessKey), ""),
		Secure: envVars.UseSSL,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Error creating minio client: %v", err))
	}

	// Route to serve video content
	http.HandleFunc("/video/playlist", serveHLS)

	fmt.Printf("Starting server on port %d...\n", 8080)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil))
}

// serveHLS fetches HLS files from MinIO and serves them with correct headers
func serveHLS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", string(envVars.FrontendURL))
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// If it's a preflight OPTIONS request, respond with a 200
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract file path from the request
	objectName := strings.TrimPrefix(r.URL.Path, "/video/")

	// Generate a presigned URL for secure access
	presignedURL, err := getPresignedURL(objectName)
	if err != nil {
		http.Error(w, "Error generating presigned URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL, http.StatusFound)
}

// getPresignedURL generates a temporary URL for an object in MinIO
func getPresignedURL(objectName string) (string, error) {
	// Define expiration for the signed URL

	params := url.Values{}

	if strings.HasSuffix(objectName, ".m3u8") {
		params.Set("Content-Type", "application/vnd.apple.mpegurl")
	} else if strings.HasSuffix(objectName, ".ts") {
		params.Set("Content-Type", "video/mp2t")
	} else {
		params.Set("Content-Type", "application/octet-stream")
	}

	// Generate presigned URL
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), "movies-processed", objectName, 15*time.Minute, params)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

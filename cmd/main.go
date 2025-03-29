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

	"github.com/nmktad/bumflix/pkg/ffmpeg"
)

var minioClient *minio.Client

func main() {
	http.HandleFunc("/", serveHLS)

	env, err := ffmpeg.LoadMinioEnv()
	if err != nil {
		log.Fatal("error loading minio env")
	}

	minioClient, err = minio.New(string(env.Endpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(string(env.AccessKeyID), string(env.SecretAccessKey), ""),
		Secure: env.UseSSL,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Error creating minio client: %v", err))
	}

	// Route to serve video content
	http.HandleFunc("/video/", serveHLS)

	fmt.Printf("Starting server on port %d...\n", 8080)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil))
}

// serveHLS fetches HLS files from MinIO and serves them with correct headers
func serveHLS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
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

// package main
//
// import (
//
//	"fmt"
//
//	"github.com/nmktad/bumflix/pkg/ffmpeg"
//
// )
//
//	func main() {
//		err := ffmpeg.ProcessVideoForStream(
//			&ffmpeg.Video{
//				Filename:       "It.Happened.One.Night.1934.2160p.4K.BluRay.x265.10bit.AAC5.1-[YTS.MX].mkv",
//				FileType:       "mkv",
//				BucketName:     "movies",
//				DestBucketName: "movies-processed",
//				SignedURL:      "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL21vdmllcy9JdC5IYXBwZW5lZC5PbmUuTmlnaHQuMTkzNC4yMTYwcC40Sy5CbHVSYXkueDI2NS4xMGJpdC5BQUM1LjEtJTVCWVRTLk1YJTVELm1rdj9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPVVWSVU3VFU2WFI4WENGSUpCWFVUJTJGMjAyNTAzMjglMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjUwMzI4VDE1MDk1N1omWC1BbXotRXhwaXJlcz00MzIwMCZYLUFtei1TZWN1cml0eS1Ub2tlbj1leUpoYkdjaU9pSklVelV4TWlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05MWlhraU9pSlZWa2xWTjFSVk5saFNPRmhEUmtsS1FsaFZWQ0lzSW1WNGNDSTZNVGMwTXpJeE56VTNOQ3dpY0dGeVpXNTBJam9pWW5WdFpteHBlQ0o5Lkh1NXpyWmdOTXVONDI3ZHlwZmV3ZUhkM0xwdmFTNk1malMtWm1fYWJWZ1NNQnlURERYcnVsSFhNVHR4cVFmbW8xb04yUFl2VXBUcmRIWldPR1dPX2tBJlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCZ2ZXJzaW9uSWQ9bnVsbCZYLUFtei1TaWduYXR1cmU9MmE4NTFjYjYxMGE2NDhhNWE1YjQ5N2I2MjlkOTczY2EwN2E2MzA2OTE0YjRmNjk1MDgyYjgwMDliZjJhN2U0OA",
//			},
//		)
//		if err != nil {
//			panic(fmt.Sprintf("Error processing video: %v", err))
//		}
//
//		fmt.Println("Video processed successfully")
//	}

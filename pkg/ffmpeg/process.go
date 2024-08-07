package ffmpeg

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	goffmpeg "github.com/u2takey/ffmpeg-go"
)

type env struct {
	endpoint        []byte
	accessKeyID     []byte
	secretAccessKey []byte
	useSSL          bool
}

func loadMinioEnv() (*env, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	endpoint := []byte(os.Getenv("MINIO_ENDPOINT"))
	accessKeyID := []byte(os.Getenv("MINIO_ACCESS_KEY"))
	secretAccessKey := []byte(os.Getenv("MINIO_SECRET_KEY"))
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	return &env{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		useSSL:          useSSL,
	}, nil
}

type Video struct {
	Filename       string
	FileType       string
	BucketName     string
	DestBucketName string
	SignedURL      string
}

func ProcessVideoForStream(videoinfo *Video) error {
	env, err := loadMinioEnv()
	if err != nil {
		return err
	}

	minioClient, err := minio.New(string(env.endpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(string(env.accessKeyID), string(env.secretAccessKey), ""),
		Secure: env.useSSL,
	})
	if err != nil {
		return fmt.Errorf("Error creating minio client: %v", err)
	}

	ctx := context.Background()

	// check if bucket exists
	if found, err := minioClient.BucketExists(ctx, videoinfo.BucketName); err != nil {
		return fmt.Errorf("Error checking if bucket exists: %v", err)
	} else if !found {
		fmt.Printf("Bucket %s does not exist\n", videoinfo.BucketName)
	}

	// check if video exists
	if _, err := minioClient.StatObject(
		ctx,
		videoinfo.BucketName,
		videoinfo.Filename,
		minio.StatObjectOptions{},
	); err != nil {
		return fmt.Errorf("Error checking if video exists: %v", err)
	}

	// download video
	object, err := minioClient.GetObject(
		ctx,
		videoinfo.BucketName,
		videoinfo.Filename,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("Error downloading video: %v", err)
	}
	defer object.Close()

	err = os.MkdirAll(fmt.Sprintf("/tmp/bumflix/%s", videoinfo.Filename), os.ModePerm)
	if err != nil {
		return err
	}

	localFile, err := os.Create(fmt.Sprintf("/tmp/bumflix/%s/%s", videoinfo.Filename, videoinfo.Filename))
	if err != nil {
		return err
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, object); err != nil {
		return err
	}

	// process video
	err = goffmpeg.
		Input(fmt.Sprintf("/tmp/bumflix/%s/%s", videoinfo.Filename, videoinfo.Filename)).
		Output(
			fmt.Sprintf("/tmp/bumflix/%s/%s.m3u8", videoinfo.Filename, videoinfo.Filename),
			goffmpeg.KwArgs{
				"hls_time":      10,
				"hls_list_size": 0,
				"start_number":  0,
				"codec":         "copy",
				"f":             "hls",
			},
		).Run()

	if err != nil {
		return fmt.Errorf("Error with ffmpeg: %v", err)
	}

	err = filepath.Walk(fmt.Sprintf("/tmp/bumflix/%s", videoinfo.Filename), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Error walking path at file %s: %v", info.Name(), err)
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".m3u8" || filepath.Ext(info.Name()) == ".ts" {
			uploadInfo, err := minioClient.FPutObject(
				ctx,
				videoinfo.DestBucketName,
				info.Name(),
				path,
				minio.PutObjectOptions{},
			)
			if err != nil {
				return fmt.Errorf("Error uploading file %s: %v", info.Name(), err)
			}

			fmt.Printf("Successfully uploaded %s to %s\n", info.Name(), uploadInfo.Location)
			fmt.Printf("Uploaded information: %v\n", uploadInfo)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Error walking path: %v", err)
	}

	return nil
}

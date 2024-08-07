package main

import (
	"fmt"

	"github.com/nmktad/bumflix/pkg/ffmpeg"
)

func main() {
	err := ffmpeg.ProcessVideoForStream(
		&ffmpeg.Video{
			Filename:       "It.Happened.One.Night.1934.2160p.4K.BluRay.x265.10bit.AAC5.1-[YTS.MX].mkv",
			FileType:       "mkv",
			BucketName:     "movies",
			DestBucketName: "movies-processed",
			SignedURL:      "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL21vdmllcy9JdC5IYXBwZW5lZC5PbmUuTmlnaHQuMTkzNC4yMTYwcC40Sy5CbHVSYXkueDI2NS4xMGJpdC5BQUM1LjEtJTVCWVRTLk1YJTVELm1rdj9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPThDT0FQRUVPRE9YWlk4Wk1OR1RHJTJGMjAyNDA4MDUlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjQwODA1VDA5MzMzMVomWC1BbXotRXhwaXJlcz00MzE5OSZYLUFtei1TZWN1cml0eS1Ub2tlbj1leUpoYkdjaU9pSklVelV4TWlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05MWlhraU9pSTRRMDlCVUVWRlQwUlBXRnBaT0ZwTlRrZFVSeUlzSW1WNGNDSTZNVGN5TWpnNU16STBOeXdpY0dGeVpXNTBJam9pWVdSdGFXNGlmUS5HdXhzNC1YbzlRQW5tMUdYeU9OQVduMFltZkhfb1hwXzhiWnZicmF6aGRJQUZmSDhtbk55LU1YVEJadFNFeVZLakJmckNZRjdQOW1LdGtpYmtYdW9zZyZYLUFtei1TaWduZWRIZWFkZXJzPWhvc3QmdmVyc2lvbklkPW51bGwmWC1BbXotU2lnbmF0dXJlPWY1ZTA3MGZiNDQxODFhM2ZlMDdhNzAyZjZiZDBkODBhYzQyZGI3NmNjMWQ2OWY0NWU4ZDdiNjYxMmMyN2EwNmI",
		},
	)
	if err != nil {
		panic(fmt.Sprintf("Error processing video: %v", err))
	}

	fmt.Println("Video processed successfully")
}

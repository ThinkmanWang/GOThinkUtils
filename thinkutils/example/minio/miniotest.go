package main

import (
	"context"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"runtime"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    log.Info("Hello World")

	endpoint := "10.0.0.3:9000"
	accessKeyID := "thinkman"
	secretAccessKey := "Ab123145"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Error("%s", err.Error())
	}

	log.Info("%#v\n", minioClient) // minioClient is now set up

	objectName := "商户SSO接入.zip"
	filePath := "/home/thinkman/Desktop/商户SSO接入.zip"
	contentType := "application/zip"

	ctx := context.Background()
	// Upload the zip file with FPutObject
	info, err := minioClient.FPutObject(ctx, "think-bucket", objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

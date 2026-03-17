package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path"
	"strings"

	"cloud.google.com/go/storage"
)

type GCImageRepository struct {
	Storage    *storage.Client
	BucketName string
}

func NewGCImageRepository(storageClient *storage.Client, bucketName string) *GCImageRepository {
	return &GCImageRepository{
		Storage:    storageClient,
		BucketName: bucketName,
	}
}

func (r *GCImageRepository) UpdateProfile(
	ctx context.Context,
	objName string,
	imageFile multipart.File,
) (string, error) {
	bucket := r.Storage.Bucket(r.BucketName)
	object := bucket.Object(objName)

	writer := object.NewWriter(ctx)

	// 讓瀏覽器盡量不要把舊頭像 cache 住
	writer.CacheControl = "no-cache, max-age=0"

	// 根據副檔名設定基本 Content-Type
	switch strings.ToLower(path.Ext(objName)) {
	case ".png":
		writer.ContentType = "image/png"
	case ".gif":
		writer.ContentType = "image/gif"
	case ".webp":
		writer.ContentType = "image/webp"
	case ".jpg", ".jpeg":
		writer.ContentType = "image/jpeg"
	default:
		writer.ContentType = "application/octet-stream"
	}

	// 把前端上傳的圖片內容寫進 GCS
	if _, err := io.Copy(writer, imageFile); err != nil {
		_ = writer.Close()
		log.Printf("failed to write image to GCS: %v\n", err)
		return "", err
	}

	// 一定要 Close，GCS 寫入才算真正完成
	if err := writer.Close(); err != nil {
		log.Printf("failed to close GCS writer: %v\n", err)
		return "", err
	}

	imageURL := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		r.BucketName,
		objName,
	)

	return imageURL, nil
}

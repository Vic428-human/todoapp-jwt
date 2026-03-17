package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
)

type gcImageRepository struct {
	Storage    *storage.Client
	BucketName string
}

/*
怎樣把圖片存到 GCS，只專注在圖片儲存這件事
處理「跟外部資源怎麼互動」這件事，這裡的外部資源不是資料庫，而是 Google Cloud Storage
*/
func (r *gcImageRepository) UpdateProfile(
	ctx context.Context, // Context：用於取消操作或超時控制
	objName string, // objName：GCS 中的檔案名稱（含路徑，如 "users/123/avatar.jpg"）
	imageFile multipart.File, // imageFile：前端 multipart form 上傳的檔案
) (string, error) {
	bckt := r.Storage.Bucket(r.BucketName)                            // 從 storage client 取得 bucket handle
	object := bckt.Object(objName)                                    //  指定要寫入的 object（檔案
	wc := object.NewWriter(ctx)                                       // 建立 Writer：這是上傳的起點，GCS 準備接收資料
	wc.ObjectAttrs.CacheControl = "Cache-Control:no-cache, max-age=0" // 設定瀏覽器快取：每次都重新載入（適合 profile 圖

	if _, err := io.Copy(wc, imageFile); err != nil { // 串流上傳（邊讀邊寫）
		// io.Copy 自動從 imageFile 讀取所有 bytes 寫入 GCS
		// multipart.File 內建 Read() 方法，所以可以直接當 reader 用
		log.Printf("Unable to write file to Google Cloud Storage: %v\n", err)
		return "", apperrors.NewInternal()
	}

	if err := wc.Close(); err != nil { // GCS 儲存檔案 + 產生 URL
		return "", fmt.Errorf("Writer.Close: %v", err) // Close() 會 flush 所有資料到 GCS 並 finalize 上傳
	}

	imageURL := fmt.Sprintf( // 產生公開可存取的 URL
		"https://storage.googleapis.com/%s/%s",
		r.BucketName, objName,
	)
	return imageURL, nil
}

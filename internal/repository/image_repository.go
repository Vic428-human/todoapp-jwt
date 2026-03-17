/*
凡是跟 GCS 圖片存取有關的，都放這裡
*/

package repository

import (
	"context"
	"mime/multipart"
)

type ImageRepository interface {
	UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error)
}

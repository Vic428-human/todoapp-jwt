package models

import (
	"context"
	"mime/multipart"
)

// ImageRepository defines methods it  expoexts a repository
// it interacts with to implement
type ImageRepository interface {
	UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error)
}

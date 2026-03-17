package service

import (
	"context"
	"log"
	"mime/multipart"

	"todo_api/internal/models"
	"todo_api/internal/repository"
	"todo_api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	DB              *pgxpool.Pool
	ImageRepository repository.ImageRepository
}

func NewUserService(
	db *pgxpool.Pool,
	imageRepository repository.ImageRepository,
) *UserService {
	return &UserService{
		DB:              db,
		ImageRepository: imageRepository,
	}
}

/*
SetProfileImage 處理使用者更新頭像的流程。

整體流程：
1. 查 user
2. 從舊 imageURL 取出 objName；如果原本沒有圖，就建立新的 objName
3. 開檔
4. 上傳 GCS
5. 更新 DB 裡的 image_url
6. 回傳更新後的 user
*/
func (s *UserService) SetProfileImage(
	ctx context.Context,
	id string,
	imageFileHeader *multipart.FileHeader,
) (*models.User, error) {
	log.Printf("start SetProfileImage: userID=%s\n", id)

	user, err := repository.GetUserByID(s.DB, id)
	if err != nil {
		log.Printf("failed to find user by id: %v\n", err)
		return nil, err
	}
	log.Printf("found user: id=%s email=%s oldImageURL=%s\n", user.ID, user.Email, user.ImageURL)

	objName, err := utils.ObjNameFromURL(user.ImageURL, imageFileHeader.Filename)
	if err != nil {
		log.Printf("failed to get object name from image url: %v\n", err)
		return nil, err
	}
	log.Printf("resolved objName=%s\n", objName)

	imageFile, err := imageFileHeader.Open()
	if err != nil {
		log.Printf("failed to open image file: %v\n", err)
		return nil, err
	}
	defer imageFile.Close()

	log.Printf("uploading image to storage...\n")
	imageURL, err := s.ImageRepository.UpdateProfile(ctx, objName, imageFile)
	if err != nil {
		log.Printf("failed to upload image to storage: %v\n", err)
		return nil, err
	}
	log.Printf("uploaded imageURL=%s\n", imageURL)

	updatedUser, err := repository.UpdateUserImage(ctx, s.DB, id, imageURL)
	if err != nil {
		log.Printf("failed to update user image url in db: %v\n", err)
		return nil, err
	}
	log.Printf("updated user image in db: id=%s imageURL=%s\n", updatedUser.ID, updatedUser.ImageURL)

	return updatedUser, nil
}

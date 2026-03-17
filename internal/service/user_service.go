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

	// === 步驟 1: 查詢使用者資料 ===
	log.Printf("start SetProfileImage: userID=%s\n", id)
	user, err := repository.GetUserByID(s.DB, id)
	if err != nil {
		log.Printf("❌ 步驟1失敗: 找不到 user(id=%s): %v\n", id, err)
		return nil, err
	}
	log.Printf("✅ 步驟1完成: 找到 user(id=%s, email=%s, oldImageURL=%s)\n",
		user.ID, user.Email, user.ImageURL)

	// === 步驟 2: 解析/生成 GCS object 名稱 ===
	objName, err := utils.ObjNameFromURL(user.ImageURL, imageFileHeader.Filename)
	if err != nil {
		log.Printf("❌ 步驟2失敗: 無法解析 objName (oldURL=%s, filename=%s): %v\n",
			user.ImageURL, imageFileHeader.Filename, err)
		return nil, err
	}
	log.Printf("✅ 步驟2完成: objName=%s\n", objName)

	// === 步驟 3: 開啟上傳檔案 ===
	imageFile, err := imageFileHeader.Open()
	if err != nil {
		log.Printf("❌ 步驟3失敗: 無法開啟檔案 %s: %v\n", imageFileHeader.Filename, err)
		return nil, err
	}
	defer imageFile.Close()
	log.Printf("✅ 步驟3完成: 成功開啟檔案 %s\n", imageFileHeader.Filename)

	// === 步驟 4: 上傳到 Google Cloud Storage ===
	log.Printf("🔄 步驟4進行中: 上傳到 GCS (objName=%s)...\n", objName)
	imageURL, err := s.ImageRepository.UpdateProfile(ctx, objName, imageFile)
	if err != nil {
		log.Printf("❌ 步驟4失敗: GCS 上傳失敗 (objName=%s): %v\n", objName, err)
		return nil, err
	}
	log.Printf("✅ 步驟4完成: 新 imageURL=%s\n", imageURL)

	// === 步驟 5: 更新資料庫中的 image_url ===
	log.Printf("🔄 步驟5進行中: 更新 DB...\n")
	updatedUser, err := repository.UpdateUserImage(ctx, s.DB, id, imageURL)
	if err != nil {
		log.Printf("❌ 步驟5失敗: DB 更新失敗: %v\n", err)
		return nil, err
	}
	log.Printf("✅ 步驟5完成: DB 已更新 (id=%s, newImageURL=%s)\n",
		updatedUser.ID, updatedUser.ImageURL)

	// === 步驟 6: 回傳成功結果 ===
	log.Printf("🎉 全部流程完成! 返回更新後的 user\n")
	return updatedUser, nil
}

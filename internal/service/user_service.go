package service

/*
先用 uid 找使用者
從原本 ImageURL 推出 objName
打開上傳檔案
呼叫 ImageRepository.UpdateProfile(...) 上傳圖片
再呼叫 UserRepository.UpdateImage(...) 更新資料庫中的 imageURL
最後回傳更新後的 user
這種「串接多個 repository / 多個步驟 / 有流程判斷」的內容，就是很典型的 service。
*/

func (s *userService) SetProfileImage(
	ctx context.Context,
	uid uuid.UUID,
	imageFileHeader *multipart.FileHeader,
) (*model.User, error) {
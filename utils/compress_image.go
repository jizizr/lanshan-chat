package utils

import "mime/multipart"

// CompressImage 从 *multipart.FileHeader 读取并压缩图片
func CompressImage(img *multipart.FileHeader) ([]byte, error) {
	f, err := img.Open()
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()
	return []byte{}, nil
}
